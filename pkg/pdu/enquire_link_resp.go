package pdu

// EnquireLinkResp represents an SMPP enquire_link_resp PDU
type EnquireLinkResp struct {
	Header *Header
}

// NewEnquireLinkResp creates a new EnquireLinkResp PDU
func NewEnquireLinkResp() *EnquireLinkResp {
	return &EnquireLinkResp{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (e *EnquireLinkResp) Marshal() ([]byte, error) {
	// Set command length in header
	e.Header.CommandLength = 16 // Only header, no body
	e.Header.CommandID = ENQUIRE_LINK_RESP

	// Marshal header
	return e.Header.Marshal()
}

// Unmarshal deserializes the PDU from bytes
func (e *EnquireLinkResp) Unmarshal(data []byte) error {
	// Unmarshal header
	e.Header = &Header{}
	return e.Header.Unmarshal(data[:16])
}
