package pss

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	storage := Storage{
		&Record{"google.com", "Bob Smith", "qwerty"},
		&Record{"yandex.ru", "Vladimir Koroviev", "12345678"},
		&Record{"microsoft.com", "bill_gates", "applesucks000"},
	}
	buffer := &bytes.Buffer{}
	pwd := "masterPassword"

	if err := Encrypt(storage, buffer, pwd); err != nil {
		t.Fatal(err)
	}

	// Simple encryption check
	if bytes.Contains(buffer.Bytes(), []byte(storage[0].Password)) {
		t.Errorf("Вывод незашифрован: %q", buffer.String())
	}

	storage2, err := Decrypt(buffer, pwd)
	if err != nil {
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
}

func TestEncryptDecrypt_WrongPwd(t *testing.T) {
	storage := Storage{
		&Record{"google.com", "Bob Smith", "qwerty"},
	}
	buffer := &bytes.Buffer{}

	if err := Encrypt(storage, buffer, "First password"); err != nil {
		t.Fatal(err)
	}
	if _, err := Decrypt(buffer, "Second password"); err != WrongPwdErr {
		t.Fatal(err)
	}
}
