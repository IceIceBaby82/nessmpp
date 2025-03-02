package pdu

// UnbindResp represents an SMPP unbind response PDU
type UnbindResp struct {
	Header *Header
}

// NewUnbindResp creates a new UnbindResp PDU
func NewUnbindResp() *UnbindResp {
	return &UnbindResp{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (u *UnbindResp) Marshal() ([]byte, error) {
	// Set command length in header
	u.Header.CommandLength = 16 // Only header, no body
	u.Header.CommandID = UNBIND_RESP

	// Marshal header
	return u.Header.Marshal()
}

// Unmarshal deserializes the PDU from bytes
func (u *UnbindResp) Unmarshal(data []byte) error {
	// Unmarshal header
	u.Header = &Header{}
	return u.Header.Unmarshal(data[:16])
}
