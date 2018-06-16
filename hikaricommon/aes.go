package hikaricommon

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type AESCrypto struct {
	iv        []byte
	encStream *cipher.Stream
	decStream *cipher.Stream
}

func NewAESCrypto(key []byte, iv []byte) (*AESCrypto, error) {
	if iv == nil {
		iv = make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	encStream := cipher.NewCFBEncrypter(block, iv)
	decStream := cipher.NewCFBDecrypter(block, iv)

	return &AESCrypto{iv, &encStream, &decStream}, nil
}

func (c *AESCrypto) Encrypt(in []byte) {
	(*c.encStream).XORKeyStream(in, in)
}

func (c *AESCrypto) Decrypt(in []byte) {
	(*c.decStream).XORKeyStream(in, in)
}

func (c *AESCrypto) GetIV() []byte {
	return c.iv
}
