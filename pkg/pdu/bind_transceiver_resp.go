package pdu

import "errors"

// BindTransceiverResp represents an SMPP bind_transceiver_resp PDU
type BindTransceiverResp struct {
	Header    *Header
	SystemID  string
	TLVParams map[uint16]*TLVParam
}

// NewBindTransceiverResp creates a new BindTransceiverResp PDU
func NewBindTransceiverResp() *BindTransceiverResp {
	return &BindTransceiverResp{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (btr *BindTransceiverResp) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(btr.SystemID) + 1

	// Add TLV parameters length
	for _, tlv := range btr.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	btr.Header.CommandLength = uint32(length)
	btr.Header.CommandID = BIND_TRANSCEIVER_RESP

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := btr.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write system_id
	copy(buf[offset:], btr.SystemID)
	offset += len(btr.SystemID) + 1

	// Write TLV parameters
	for _, tlv := range btr.TLVParams {
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
func (btr *BindTransceiverResp) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	btr.Header = &Header{}
	if err = btr.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read system_id
	if btr.SystemID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read TLV parameters if any remain
	remaining := int(btr.Header.CommandLength) - offset
	currentOffset := offset

	for remaining > 0 {
		if remaining < 4 {
			return errors.New("invalid TLV parameter: insufficient data")
		}

		tlv := &TLVParam{}
		if err := tlv.Unmarshal(data[currentOffset:]); err != nil {
			return err
		}

		btr.TLVParams[tlv.Tag] = tlv
		currentOffset += 4 + int(tlv.Length)
		remaining -= 4 + int(tlv.Length)
	}

	return nil
}
