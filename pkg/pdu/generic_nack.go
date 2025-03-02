package pdu

// GenericNack represents an SMPP generic_nack PDU
type GenericNack struct {
	Header *Header
}

// NewGenericNack creates a new GenericNack PDU
func NewGenericNack() *GenericNack {
	return &GenericNack{
		Header: NewHeader(),
	}
}

// SetErrorCode sets the command status in the header
func (g *GenericNack) SetErrorCode(code uint32) {
	g.Header.CommandStatus = code
}

// Marshal serializes the PDU into bytes
func (g *GenericNack) Marshal() ([]byte, error) {
	// Set command length in header
	g.Header.CommandLength = 16 // Only header, no body
	g.Header.CommandID = GENERIC_NACK

	// Marshal header
	return g.Header.Marshal()
}

// Unmarshal deserializes the PDU from bytes
func (g *GenericNack) Unmarshal(data []byte) error {
	// Unmarshal header
	g.Header = &Header{}
	return g.Header.Unmarshal(data[:16])
}
