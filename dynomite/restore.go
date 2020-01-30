package dynomite

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/majewsky/schwift"
	"github.com/sapcc/go-bits/logg"
)

// Restore downloads the latest backup to disk
func Restore(containerName, prefix string) error {
	currentTime := time.Now()

	iter, err := newIterator(containerName, prefix)
	if err != nil {
		return err
	}

	// Get the last backup
	objects, err := iter.NextPage(1)
	if err != nil {
		return err
	}

	if len(objects) == 0 {
		return fmt.Errorf("Restore failed - No Backups found in: %s", containerName+prefix)
	}

	logg.Info("Downloading %s", objects[0].FullName())

	data, err := objects[0].Download(nil).AsReadCloser()
	if err != nil {
		return err
	}
	defer data.Close()

	file, err := os.Create("/data/dump.rdb")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}

	logg.Info("Restore finished in %s", time.Now().Sub(currentTime))

	return nil
}

// ListBackups shows the possible restore candidates
func ListBackups(containerName, prefix string, limit int) error {
	iter, err := newIterator(containerName, prefix)
	if err != nil {
		return err
	}

	// Get the last backup
	objectInfos, err := iter.NextPageDetailed(limit)
	if err != nil {
		return err
	}

	if len(objectInfos) == 0 {
		return fmt.Errorf("Listing failed - No Backups found in: %s", containerName+prefix)
	}

	var backupItems string
	for _, objectInfo := range objectInfos {
		backupItems = backupItems + fmt.Sprintf("\n%s - %d bytes", objectInfo.Object.FullName(), objectInfo.SizeBytes)
	}

	logg.Info("Backup Candidates: %s", backupItems)

	return nil
}

func newIterator(containerName, prefix string) (*schwift.ObjectIterator, error) {
	var iter *schwift.ObjectIterator

	account, err := newObjectStoreAccount()
	if err != nil {
		return iter, err
	}

	// Use reverse listing, means the newest backup
	opts := &schwift.RequestOptions{
		Values: url.Values{
			"reverse": []string{"on"},
		},
	}

	container := account.Container(containerName)
	iter = container.Objects()
	iter.Prefix = prefix
	iter.Options = opts

	return iter, nil
}
