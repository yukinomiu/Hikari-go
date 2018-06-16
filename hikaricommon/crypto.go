package hikaricommon

type Crypto interface {
	Encrypt([]byte)
	Decrypt([]byte)
	GetIV() []byte
}
