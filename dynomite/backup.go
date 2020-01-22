package dynomite

import (
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/majewsky/schwift"
	"github.com/majewsky/schwift/gopherschwift"
	"github.com/sapcc/go-bits/logg"
)

// Backup triggers a BGSAVE and uploads the dumpfile to Swift
func (dyno Dynomite) Backup(containerName, prefix string) error {
	currentTime := time.Now()

	_, err := dyno.Backend.BGSave(5 * time.Minute)
	if err != nil {
		return err
	}

	objectName := prefix + "/" + currentTime.Format("2006-01-02_1504") + "/dump.rdb"
	err = uploadDump("/Users/d044166/upload_test/upload_test_1K", containerName, objectName)
	//err = uploadDump("/data/dump.rdb", containerName, objectName)
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

	ao, err := clientconfig.AuthOptions(nil)
	if err != nil {
		return err
	}
	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}
	err = openstack.Authenticate(provider, *ao)
	client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}
	account, err := gopherschwift.Wrap(client, nil)
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

	logg.Info("Upload dump to in %s", object.FullName())

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
