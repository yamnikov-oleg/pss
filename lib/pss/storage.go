package pss

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Record - запись в хранилище паролей.
type Record struct {
	Website  string
	Username string
	Password string
}

// Storage - хранилище паролей в памяти.
type Storage []*Record

// keySize - используемый размер ключа.
const keySize = sha256.Size // 32

// WrongPwdErr - ошибка о некорректности пароля при дешифровании.
const WrongPwdErr = Error("wrong pwd")

// Decrypt производит дешифровку хранилища по паролю. Если пароль неверный,
// вернет WrongPwdErr.
func Decrypt(r io.Reader, pwd string) (s Storage, err error) {
	var (
		key     [keySize]byte
		block   cipher.Block
		iv      []byte
		payload []byte
	)

	// Prepare decryptor
	key = sha256.Sum256([]byte(pwd))
	if block, err = aes.NewCipher(key[:]); err != nil {
		return nil, err
	}
	iv = make([]byte, block.BlockSize())
	if _, err = r.Read(iv); err != nil {
		return nil, err
	}
	r = cipher.StreamReader{S: cipher.NewCFBDecrypter(block, iv), R: r}

	// Check for consistensy
	head := make([]byte, keySize)
	if _, err = r.Read(head); err != nil {
		return nil, err
	}
	if !bytes.Equal(head, key[:]) {
		return nil, WrongPwdErr
	}

	// Read and decode
	if payload, err = ioutil.ReadAll(r); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(payload, &s); err != nil {
		return nil, err
	}

	return
}

// Encrypt производит шифрования хранилища по паролю.
func Encrypt(s Storage, w io.Writer, pwd string) (err error) {
	var (
		key     [keySize]byte
		block   cipher.Block
		iv      []byte
		payload []byte
	)

	// Prepare encryptor
	key = sha256.Sum256([]byte(pwd))
	if block, err = aes.NewCipher(key[:]); err != nil {
		return err
	}
	iv = make([]byte, block.BlockSize())
	if _, err = rand.Read(iv); err != nil {
		return err
	}
	if _, err = w.Write(iv); err != nil {
		return err
	}
	w = cipher.StreamWriter{S: cipher.NewCFBEncrypter(block, iv), W: w}

	// Encode and write
	if _, err = w.Write(key[:]); err != nil {
		return err
	}
	if payload, err = json.Marshal(s); err != nil {
		return err
	}
	if _, err := w.Write(payload); err != nil {
		return err
	}
	return nil
}
