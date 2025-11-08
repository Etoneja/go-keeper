package types

type SecretData interface {
	Validate() error
}
