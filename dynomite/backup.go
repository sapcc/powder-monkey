package dynomite

import (
	"fmt"
	"os"
	"time"

	"github.com/majewsky/schwift"
	"github.com/sapcc/go-bits/logg"
)

// Backup triggers a BGSAVE and uploads the dumpfile to Swift
func (dyno Dynomite) Backup(containerName, prefix string) error {
	currentTime := time.Now()

	size, err := dyno.Backend.DBSize()
	if err != nil {
		return err
	}

	if size <= 0 {
		return fmt.Errorf("Skipping this backup. No Keys in backend")
	}

	_, err = dyno.Backend.BGSave(5 * time.Minute)
	if err != nil {
		return err
	}

	objectName := prefix + "/" + currentTime.Format("2006-01-02_1504") + "/dump.rdb"
	err = uploadDump("/data/dump.rdb", containerName, objectName)
	if err != nil {
		return err
	}

	return nil
}

// BackupEvery calls Backup() periodicly
func (dyno Dynomite) BackupEvery(every time.Duration, containerName, prefix string) error {
	backup := func() {
		err := dyno.Backup(containerName, prefix)
		if err != nil {
			logg.Error("Backup failure: %s", err.Error())
		}
	}

	backup()
	ticker := time.NewTicker(every)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			backup()
		}
	}
}

func uploadDump(dumpFileName, containerName, objectName string) error {
	currentTime := time.Now()

	account, err := newObjectStoreAccount()
	if err != nil {
		return err
	}

	_, err = account.Container(containerName).EnsureExists()
	if err != nil {
		return err
	}

	f, err := os.Open(dumpFileName)
	if err != nil {
		return err
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	segmentBytes := int64(1 << 30) // 1<30 bytes = 1 GiB per segment
	object := account.Container(containerName).Object(objectName)

	logg.Info("Upload dump to %s", object.FullName())

	header := schwift.NewObjectHeaders()
	header.ExpiresAt().Set(time.Now().Add(14 * 24 * time.Hour)) // Keep backups for 2 weeks

	if fi.Size() > segmentBytes {
		segmentContainer, err := account.Container(containerName + "_segments").EnsureExists()
		if err != nil {
			return err
		}

		largeObject, err := object.AsNewLargeObject(schwift.SegmentingOptions{
			SegmentContainer: segmentContainer,
			//use defaults for everything else
		}, &schwift.TruncateOptions{
			//if there's already a large object here, clean it up
			DeleteSegments: true,
		})

		err = largeObject.Append(f, segmentBytes, header.ToOpts())
		if err != nil {
			return err
		}

		err = largeObject.WriteManifest(header.ToOpts())
		if err != nil {
			return err
		}
	} else {
		err = object.Upload(f, &schwift.UploadOptions{
			//if there's already a large object here, clean it up
			DeleteSegments: true,
		}, header.ToOpts())
		if err != nil {
			return err
		}
	}

	logg.Info("Upload finished in %s", time.Now().Sub(currentTime))

	return nil
}
