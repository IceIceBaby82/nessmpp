package pdu

import "errors"

// DeliverSMResp represents an SMPP deliver_sm_resp PDU
type DeliverSMResp struct {
	Header    *Header
	MessageID string
	TLVParams map[uint16]*TLVParam
}

// NewDeliverSMResp creates a new DeliverSMResp PDU
func NewDeliverSMResp() *DeliverSMResp {
	return &DeliverSMResp{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (d *DeliverSMResp) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(d.MessageID) + 1

	// Add TLV parameters length
	for _, tlv := range d.TLVParams {
		length += 4 + len(tlv.Value)
	}

	// Set command length in header
	d.Header.CommandLength = uint32(length)
	d.Header.CommandID = DELIVER_SM_RESP

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := d.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write message_id
	copy(buf[offset:], d.MessageID)
	offset += len(d.MessageID) + 1

	// Write TLV parameters
	for _, tlv := range d.TLVParams {
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
func (d *DeliverSMResp) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	d.Header = &Header{}
	if err = d.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read message_id
	if d.MessageID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read TLV parameters if any remain
	remaining := int(d.Header.CommandLength) - offset
	currentOffset := offset

	for remaining > 0 {
		if remaining < 4 {
			return errors.New("invalid TLV parameter: insufficient data")
		}

		tlv := &TLVParam{}
		if err := tlv.Unmarshal(data[currentOffset:]); err != nil {
			return err
		}

		d.TLVParams[tlv.Tag] = tlv
		currentOffset += 4 + int(tlv.Length)
		remaining -= 4 + int(tlv.Length)
	}

	return nil
}
