package pdu

// ReplaceSMResp represents an SMPP replace_sm_resp PDU
type ReplaceSMResp struct {
	Header *Header
}

// NewReplaceSMResp creates a new ReplaceSMResp PDU
func NewReplaceSMResp() *ReplaceSMResp {
	return &ReplaceSMResp{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (r *ReplaceSMResp) Marshal() ([]byte, error) {
	// Set command length in header
	r.Header.CommandLength = 16 // Only header, no body
	r.Header.CommandID = REPLACE_SM_RESP

	// Marshal header
	return r.Header.Marshal()
}

// Unmarshal deserializes the PDU from bytes
func (r *ReplaceSMResp) Unmarshal(data []byte) error {
	// Unmarshal header
	r.Header = &Header{}
	return r.Header.Unmarshal(data[:16])
}
