package pdu

import "errors"

// BindReceiverResp represents an SMPP bind_receiver_resp PDU
type BindReceiverResp struct {
	Header    *Header
	SystemID  string
	TLVParams map[uint16]*TLVParam
}

// NewBindReceiverResp creates a new BindReceiverResp PDU
func NewBindReceiverResp() *BindReceiverResp {
	return &BindReceiverResp{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (brr *BindReceiverResp) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(brr.SystemID) + 1

	// Add TLV parameters length
	for _, tlv := range brr.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	brr.Header.CommandLength = uint32(length)
	brr.Header.CommandID = BIND_RECEIVER_RESP

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := brr.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write system_id
	copy(buf[offset:], brr.SystemID)
	offset += len(brr.SystemID) + 1

	// Write TLV parameters
	for _, tlv := range brr.TLVParams {
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
func (brr *BindReceiverResp) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	brr.Header = &Header{}
	if err = brr.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read system_id
	if brr.SystemID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read TLV parameters if any remain
	remaining := int(brr.Header.CommandLength) - offset
	currentOffset := offset

	for remaining > 0 {
		if remaining < 4 {
			return errors.New("invalid TLV parameter: insufficient data")
		}

		tlv := &TLVParam{}
		if err := tlv.Unmarshal(data[currentOffset:]); err != nil {
			return err
		}

		brr.TLVParams[tlv.Tag] = tlv
		currentOffset += 4 + int(tlv.Length)
		remaining -= 4 + int(tlv.Length)
	}

	return nil
}
