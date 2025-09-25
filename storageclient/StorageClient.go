package storageclient

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/apkatsikas/artist-entities/infrastructures/logutil"
	"google.golang.org/api/iterator"
)

const (
	timeoutSeconds = 60
	credsEnvVar    = "GOOGLE_APPLICATION_CREDENTIALS"
)

type StorageClient struct {
	client     *storage.Client
	projectID  string
	bucketName string
}

type BackupFile struct {
	Name    string
	Updated time.Time
}

func New() *StorageClient {
	os.Setenv(credsEnvVar, os.Getenv("GCS_CREDS_FILE"))

	client, err := storage.NewClient(context.Background())
	if err != nil {
		logutil.Fatal("Failed to create StorageClient: %v", err)
	}

	return &StorageClient{
		client:     client,
		bucketName: os.Getenv("GCS_BUCKET_NAME"),
		projectID:  os.Getenv("GCS_PROJECT_ID"),
	}

}

func (sc *StorageClient) DeleteFile(object string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*timeoutSeconds)

	defer cancel()

	objForDeletion := sc.client.Bucket(sc.bucketName).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to delete the file is aborted
	// if the object's generation number does not match your precondition.
	attrs, err := objForDeletion.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("error getting object attributes: %v", err)
	}
	objForDeletion = objForDeletion.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := objForDeletion.Delete(ctx); err != nil {
		return fmt.Errorf("error during delete: %v", err)
	}
	return nil
}

// ListFiles lists files
// Returns value because we need to sort later
func (sc *StorageClient) ListFiles() ([]BackupFile, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*timeoutSeconds)
	defer cancel()

	it := sc.client.Bucket(sc.bucketName).Objects(ctx, nil)

	var files []BackupFile

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		file := BackupFile{Name: attrs.Name, Updated: attrs.Updated}

		files = append(files, file)
	}

	return files, nil
}

// UploadFile uploads an object
func (sc *StorageClient) UploadFile(path string, destObject string) error {
	blobFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer blobFile.Close()

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*timeoutSeconds)
	defer cancel()

	// Upload an object with storage.Writer.
	obj := sc.client.Bucket(sc.bucketName).Object(destObject)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	obj = obj.If(storage.Conditions{DoesNotExist: true})

	// Upload an object with storage.Writer.
	wc := obj.NewWriter(ctx)

	if _, err := io.Copy(wc, blobFile); err != nil {
		return fmt.Errorf("error on Copy to bucket %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("error on Close during bucket upload: %v", err)
	}

	return nil
}
