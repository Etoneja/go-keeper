package crypto

type Crypter interface {
	GenerateHash(data []byte) string
	EncryptData(data []byte) ([]byte, error)
	DecryptData(encryptedData []byte) ([]byte, error)
}
