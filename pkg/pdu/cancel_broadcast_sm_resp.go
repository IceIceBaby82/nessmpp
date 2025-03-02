package pdu

// CancelBroadcastSMResp represents an SMPP cancel_broadcast_sm_resp PDU (SMPP v5.0)
type CancelBroadcastSMResp struct {
	Header *Header
}

// NewCancelBroadcastSMResp creates a new CancelBroadcastSMResp PDU
func NewCancelBroadcastSMResp() *CancelBroadcastSMResp {
	return &CancelBroadcastSMResp{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (c *CancelBroadcastSMResp) Marshal() ([]byte, error) {
	// Set command length in header
	c.Header.CommandLength = 16 // Only header, no body
	c.Header.CommandID = CANCEL_BROADCAST_SM_RESP

	// Marshal header
	return c.Header.Marshal()
}

// Unmarshal deserializes the PDU from bytes
func (c *CancelBroadcastSMResp) Unmarshal(data []byte) error {
	// Unmarshal header
	c.Header = &Header{}
	return c.Header.Unmarshal(data[:16])
}
