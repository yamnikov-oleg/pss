package pss

import (
	"os"
	"os/user"
	"path"
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
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user.HomeDir
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
