package abstract

type Encryption interface {
	Encrypt(val interface{}) (string, error)
	Decrypt(val string) (string, error)
}
