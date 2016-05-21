package pss

import (
	"os"
	"path"

	"github.com/yamnikov-oleg/pss/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
)

const (
	dir         = ".pss_storage"
	storageFile = "storage"
)

// Пути стандратного расположения хранилища паролей.
var (
	StorageDir  = path.Join(mustHomeDir(), dir)
	StoragePath = path.Join(StorageDir, storageFile)
)

func mustHomeDir() string {
	path, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return path
}

// EncryptDefault записывает хранилище в стандартный файл.
func EncryptDefault(s Storage, pwd string) error {
	if err := os.MkdirAll(StorageDir, 0700); err != nil {
		return err
	}
	file, err := os.Create(StoragePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return Encrypt(s, file, pwd)
}

// DecryptDefault считывает хранилище из стандартного файла.
func DecryptDefault(pwd string) (Storage, error) {
	file, err := os.Open(StoragePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Decrypt(file, pwd)
}
