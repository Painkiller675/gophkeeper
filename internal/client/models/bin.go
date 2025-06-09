package models

var _ Secret = (*Bin)(nil)

// Bin random binary data
type Bin struct {
	Data []byte
}

// Type returns the saved info type
func (b Bin) Type() SecretType {
	return secretTypeBin
}

// String the func shows private info
func (b Bin) String() string {
	return "BINARY DATA"
}
