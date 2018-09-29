package types

// Schema is a representation of a stubbed endpoint.
type Schema struct {
	Method   string
	Path     string
	Status   int
	Response []byte
}
