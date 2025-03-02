package pdu

// Unbind represents an SMPP unbind PDU
type Unbind struct {
	Header *Header
}

// NewUnbind creates a new Unbind PDU
func NewUnbind() *Unbind {
	return &Unbind{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (u *Unbind) Marshal() ([]byte, error) {
	// Set command length in header
	u.Header.CommandLength = 16 // Only header, no body
	u.Header.CommandID = UNBIND

	// Marshal header
	return u.Header.Marshal()
}

// Unmarshal deserializes the PDU from bytes
func (u *Unbind) Unmarshal(data []byte) error {
	// Unmarshal header
	u.Header = &Header{}
	return u.Header.Unmarshal(data[:16])
}
