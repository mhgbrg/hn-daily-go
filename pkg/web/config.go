package web

type Config struct {
	Hostname    string
	Port        int
	DatabaseURL string
	CryptoKeys  CryptoKeys
}

type CryptoKeys struct {
	HashKey       []byte
	EncryptionKey []byte
}
