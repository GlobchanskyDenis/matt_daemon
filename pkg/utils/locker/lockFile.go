package locker

import (
	"os"
	"io"
)

var gLockFile io.Closer
var gLockFilePath *string

func Lock(pathfile string) error {
	file, err := os.OpenFile(pathfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	gLockFile = file
	gLockFilePath = &pathfile
	return nil
}

func Unlock() error {
	if gLockFile != nil && gLockFilePath != nil {
		if err := gLockFile.Close(); err != nil {
			return err
		}
		gLockFile = nil
		_ = os.Remove(*gLockFilePath)
		gLockFilePath = nil
	}
	return nil
}

func IsLocked(pathfile string) bool {
	file, err := os.OpenFile(pathfile, os.O_RDONLY, 0755)
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}