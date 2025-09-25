package services

import (
    "errors"
    "fmt"
    "testing"
    "time"

    "github.com/apkatsikas/artist-entities/interfaces/mocks"
    sc "github.com/apkatsikas/artist-entities/storageclient"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

const ss = "string"

func happyPathFiles() []sc.BackupFile {
    // Data
    now := time.Now()
    yesterday := now.AddDate(0, 0, -1)
    twoDaysAgo := now.AddDate(0, 0, -2)
    oldestName := "twodaysago.sqlite"

    return []sc.BackupFile{
        {
            Name:    "yesterday.sqlite",
            Updated: yesterday,
        },
        {
            Name:    "coolfile.sqlite",
            Updated: now,
        },
        {
            Name:    oldestName,
            Updated: twoDaysAgo,
        },
    }
}

type adminServiceTestMocks struct {
    *mocks.IAdminRepository
    *mocks.IStorageClient
    *mocks.IFileUtil
    *mocks.IAdminRules
}

func adminServiceReqMocks(t *testing.T) adminServiceTestMocks {
    return adminServiceTestMocks{
        IAdminRepository: mocks.NewIAdminRepository(t),
        IStorageClient:   mocks.NewIStorageClient(t),
        IFileUtil:        mocks.NewIFileUtil(t),
        IAdminRules:      mocks.NewIAdminRules(t),
    }
}

func injectedAdminService(mocks adminServiceTestMocks) AdminService {
    return AdminService{
        AdminRepository: mocks.IAdminRepository,
        StorageClient:   mocks.IStorageClient,
        FileUtil:        mocks.IFileUtil,
        Rules:           mocks.IAdminRules,
    }
}

func TestBackup(t *testing.T) {
    // Data
    bucketFiles := happyPathFiles()
    oldest := bucketFiles[len(bucketFiles)-1]
    oldestName := oldest.Name

    file := vacuumFileName

    // Setup mocks
    mocks := adminServiceReqMocks(t)

    // Happy path normal scenario
    // Local backup file already exists and is deleted successfully
    // New backup is successful
    // Upload of the backup is successful
    // 2 files exist on the bucket prior to upload
    // We delete the oldest one successfully
    mocks.IFileUtil.EXPECT().DeleteIfExists(file).Return(nil)
    mocks.IAdminRepository.EXPECT().CreateBackup(vacuumFileName).Return(nil)
    mocks.IStorageClient.EXPECT().UploadFile(vacuumFileName, mock.AnythingOfType(ss)).Return(nil)
    mocks.IStorageClient.EXPECT().ListFiles().Return(bucketFiles, nil)
    mocks.IAdminRules.EXPECT().FileToDelete(bucketFiles).Return(oldestName)
    mocks.IStorageClient.EXPECT().DeleteFile(oldestName).Return(nil)

    // Inject service
    adminService := injectedAdminService(mocks)

    // Backup
    err := adminService.Backup()

    // Check that there is no error
    assert.Nil(t, err)
}

func TestBackupNoDelete(t *testing.T) {
    // Data
    bucketFiles := happyPathFiles()
    fileToDelete := ""

    file := vacuumFileName

    // Setup mocks
    mocks := adminServiceReqMocks(t)

    // Happy path normal scenario
    // Local backup file already exists and is deleted successfully
    // New backup is successful
    // Upload of the backup is successful
    // 2 files exist on the bucket prior to upload
    // We delete the oldest one successfully
    mocks.IFileUtil.EXPECT().DeleteIfExists(file).Return(nil)
    mocks.IAdminRepository.EXPECT().CreateBackup(vacuumFileName).Return(nil)
    mocks.IStorageClient.EXPECT().UploadFile(vacuumFileName, mock.AnythingOfType(ss)).Return(nil)
    mocks.IStorageClient.EXPECT().ListFiles().Return(bucketFiles, nil)
    mocks.IAdminRules.EXPECT().FileToDelete(bucketFiles).Return(fileToDelete)

    // Inject service
    adminService := injectedAdminService(mocks)

    // Backup
    err := adminService.Backup()

    // Check that there is no error
    assert.Nil(t, err)
}

func TestBackupFileUtilFails(t *testing.T) {
    // Test data
    file := vacuumFileName
    expectedError := fmt.Errorf("??")

    // Setup mocks
    mocks := adminServiceReqMocks(t)
    mocks.IFileUtil.EXPECT().DeleteIfExists(file).Return(expectedError)

    // Inject service
    adminService := injectedAdminService(mocks)

    // Backup
    err := adminService.Backup()

    // Check error
    assert.True(t, errors.Is(err, expectedError))
}

func TestBackupAdminFails(t *testing.T) {
    // Test data
    expectedError := fmt.Errorf("??")
    file := vacuumFileName

    // Setup mocks
    mocks := adminServiceReqMocks(t)
    mocks.IFileUtil.EXPECT().DeleteIfExists(file).Return(nil)
    // Fail
    mocks.IAdminRepository.EXPECT().CreateBackup(vacuumFileName).Return(expectedError)

    // Inject service
    adminService := injectedAdminService(mocks)

    // Backup
    err := adminService.Backup()

    // Check error
    assert.True(t, errors.Is(err, expectedError))
}

func TestBackupUploadStorageFails(t *testing.T) {
    // Test data
    expectedError := fmt.Errorf("??")
    file := vacuumFileName

    // Setup mocks
    mocks := adminServiceReqMocks(t)
    mocks.IFileUtil.EXPECT().DeleteIfExists(file).Return(nil)
    mocks.IAdminRepository.EXPECT().CreateBackup(vacuumFileName).Return(nil)
    // Fail
    mocks.IStorageClient.EXPECT().UploadFile(vacuumFileName, mock.AnythingOfType(ss)).Return(expectedError)

    // Inject service
    adminService := injectedAdminService(mocks)

    // Backup
    err := adminService.Backup()

    // Check error
    assert.True(t, errors.Is(err, expectedError))
}

func TestBackupListStorageFails(t *testing.T) {
    // Test data
    expectedError := fmt.Errorf("??")
    file := vacuumFileName

    // Setup mocks
    mocks := adminServiceReqMocks(t)
    mocks.IFileUtil.EXPECT().DeleteIfExists(file).Return(nil)
    mocks.IAdminRepository.EXPECT().CreateBackup(vacuumFileName).Return(nil)
    mocks.IStorageClient.EXPECT().UploadFile(vacuumFileName, mock.AnythingOfType(ss)).Return(nil)
    // Fail
    mocks.IStorageClient.EXPECT().ListFiles().Return(nil, expectedError)

    // Inject service
    adminService := injectedAdminService(mocks)

    // Backup
    err := adminService.Backup()

    // Check error
    assert.True(t, errors.Is(err, expectedError))
}

func TestBackupDeleteStorageFails(t *testing.T) {
    // Test data
    bucketFiles := happyPathFiles()
    expectedError := fmt.Errorf("??")
    oldest := bucketFiles[len(bucketFiles)-1].Name
    file := vacuumFileName

    // Setup mocks
    mocks := adminServiceReqMocks(t)
    mocks.IFileUtil.EXPECT().DeleteIfExists(file).Return(nil)
    mocks.IAdminRepository.EXPECT().CreateBackup(vacuumFileName).Return(nil)
    mocks.IStorageClient.EXPECT().UploadFile(vacuumFileName, mock.AnythingOfType(ss)).Return(nil)
    mocks.IStorageClient.EXPECT().ListFiles().Return(bucketFiles, nil)
    mocks.IAdminRules.EXPECT().FileToDelete(bucketFiles).Return(oldest)
    // Fail
    mocks.IStorageClient.EXPECT().DeleteFile(oldest).Return(expectedError)

    // Inject service
    adminService := injectedAdminService(mocks)

    // Backup
    err := adminService.Backup()

    // Check error
    assert.True(t, errors.Is(err, expectedError))
}
