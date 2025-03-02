package pdu

// CancelSMResp represents an SMPP cancel_sm_resp PDU
type CancelSMResp struct {
	Header *Header
}

// NewCancelSMResp creates a new CancelSMResp PDU
func NewCancelSMResp() *CancelSMResp {
	return &CancelSMResp{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (c *CancelSMResp) Marshal() ([]byte, error) {
	// Set command length in header
	c.Header.CommandLength = 16 // Only header, no body
	c.Header.CommandID = CANCEL_SM_RESP

	// Marshal header
	return c.Header.Marshal()
}

// Unmarshal deserializes the PDU from bytes
func (c *CancelSMResp) Unmarshal(data []byte) error {
	// Unmarshal header
	c.Header = &Header{}
	return c.Header.Unmarshal(data[:16])
}
