package signatures

type SingatureVerifier interface {
	Verify(signature string, accountAddress string) (bool, error)
}
