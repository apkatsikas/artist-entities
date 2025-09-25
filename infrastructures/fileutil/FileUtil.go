package fileutil

import (
    "os"
)

type FileUtil struct {
}

func (fu *FileUtil) fileExists(file string) (bool, error) {
    info, err := os.Stat(file)
    // If error is not exist, return false
    if os.IsNotExist(err) {
        return false, nil
    } else if err != nil {
        return false, err
    }
    // Return if this is not a directory
    return !info.IsDir(), nil
}

func (fu *FileUtil) DeleteIfExists(file string) error {
    exists, err := fu.fileExists(file)
    if err != nil {
        return err
    }

    if exists {
        err = os.Remove(file)

        if err != nil {
            return err
        }
    }
    return nil
}
