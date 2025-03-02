package pdu

// EnquireLink represents an SMPP enquire_link PDU
type EnquireLink struct {
	Header *Header
}

// NewEnquireLink creates a new EnquireLink PDU
func NewEnquireLink() *EnquireLink {
	return &EnquireLink{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (e *EnquireLink) Marshal() ([]byte, error) {
	// Set command length in header
	e.Header.CommandLength = 16 // Only header, no body
	e.Header.CommandID = ENQUIRE_LINK

	// Marshal header
	return e.Header.Marshal()
}

// Unmarshal deserializes the PDU from bytes
func (e *EnquireLink) Unmarshal(data []byte) error {
	// Unmarshal header
	e.Header = &Header{}
	return e.Header.Unmarshal(data[:16])
}
