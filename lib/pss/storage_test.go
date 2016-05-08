package pss

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func genTestStorage(n int) (s Storage) {
	samples := Storage{
		&Record{"google.com", "Bob Smith", "qwerty"},
		&Record{"yandex.ru", "Vladimir Koroviev", "12345678"},
		&Record{"microsoft.com", "bill_gates", "applesucks000"},
		&Record{"apple.com", "steve_jobs", "MSeatsCHIT"},
		&Record{"google.com", "SergeyBrin", "shutup_you_both"},
	}
	for n > 0 {
		if n >= 3 {
			s = append(s, samples...)
		} else {
			s = append(s, samples[:n]...)
		}
		n -= 3
	}
	return
}

func TestEncryptDecrypt(t *testing.T) {
	storage := genTestStorage(10)
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
	storage := genTestStorage(10)
	buffer := &bytes.Buffer{}

	if err := Encrypt(storage, buffer, "First password"); err != nil {
		t.Fatal(err)
	}
	if _, err := Decrypt(buffer, "Second password"); err != WrongPwdErr {
		t.Fatal(err)
	}
}

func TestGoldenStorage(t *testing.T) {
	const password = "golden"

	genc, err := os.Open("golden.enc")
	if err != nil {
		t.Fatal(err)
	}
	defer genc.Close()

	gjson, err := os.Open("golden.json")
	if err != nil {
		t.Fatal(err)
	}
	defer gjson.Close()

	var (
		expected Storage
		actual   Storage
	)

	expBuf, err := ioutil.ReadAll(gjson)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(expBuf, &expected); err != nil {
		t.Fatal(err)
	}
	actual, err = Decrypt(genc, password)
	if err != nil {
		t.Fatal(err)
	}

	if len(expected) != len(actual) {
		t.Fatalf("Длины не совпадают. Ожидалось %v, получено %v", len(expected), len(actual))
	}
	for i := range expected {
		if *expected[i] != *actual[i] {
			t.Errorf("Записи не совпадают. Ожидалось %v, получено %v", expected[i], actual[i])
		}
	}
}
