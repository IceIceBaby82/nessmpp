package pdu

import "errors"

// SubmitSMResp represents an SMPP submit_sm_resp PDU
type SubmitSMResp struct {
	Header    *Header
	MessageID string
	TLVParams map[uint16]*TLVParam
}

// NewSubmitSMResp creates a new SubmitSMResp PDU
func NewSubmitSMResp() *SubmitSMResp {
	return &SubmitSMResp{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (s *SubmitSMResp) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(s.MessageID) + 1

	// Add TLV parameters length
	for _, tlv := range s.TLVParams {
		length += 4 + len(tlv.Value)
	}

	// Set command length in header
	s.Header.CommandLength = uint32(length)
	s.Header.CommandID = SUBMIT_SM_RESP

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := s.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write message_id
	copy(buf[offset:], s.MessageID)
	offset += len(s.MessageID) + 1

	// Write TLV parameters
	for _, tlv := range s.TLVParams {
		tlvBytes, err := tlv.Marshal()
		if err != nil {
			return nil, err
		}
		copy(buf[offset:], tlvBytes)
		offset += 4 + len(tlv.Value)
	}

	return buf, nil
}

// Unmarshal deserializes the PDU from bytes
func (s *SubmitSMResp) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	s.Header = &Header{}
	if err = s.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read message_id
	if s.MessageID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read TLV parameters if any remain
	remaining := int(s.Header.CommandLength) - offset
	currentOffset := offset

	for remaining > 0 {
		if remaining < 4 {
			return errors.New("invalid TLV parameter: insufficient data")
		}

		tlv := &TLVParam{}
		if err := tlv.Unmarshal(data[currentOffset:]); err != nil {
			return err
		}

		s.TLVParams[tlv.Tag] = tlv
		currentOffset += 4 + int(tlv.Length)
		remaining -= 4 + int(tlv.Length)
	}

	return nil
}
