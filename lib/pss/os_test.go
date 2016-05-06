package pss

import (
	"os"
	"path"
	"testing"
)

func TestEncryptDecryptDefault(t *testing.T) {
	StorageDir = path.Join(os.Getenv("HOME"), ".pss_test")
	StoragePath = path.Join(StorageDir, storageFile)

	var (
		storage  = genTestStorage(10)
		password = "masterPassword"
		storage2 Storage

		err error
	)

	storage = genTestStorage(10)
	password = "masterPassword"

	// Write-read consistensy check
	if err = EncryptDefault(storage, password); err != nil {
		t.Fatal(err)
	}
	if storage2, err = DecryptDefault(password); err != nil {
		t.Fatal(err)
	}

	if len(storage) != len(storage2) {
		t.Fatalf("Длины не совпадают: %v и %v", len(storage2), len(storage))
	}
	for i := range storage {
		if *storage[i] != *storage2[i] {
			t.Errorf("Записи не совпадают (%v): %v и %v", i, *storage2[i], *storage[i])
		}
	}

	// Rewriting check
	if err = EncryptDefault(storage, password); err != nil {
		t.Fatal(err)
	}

	// Clean up
	if err := os.RemoveAll(StorageDir); err != nil {
		t.Log(err)
	}
}
