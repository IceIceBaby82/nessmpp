package pdu

import (
	"encoding/binary"
	"errors"
)

// DataSMResp represents an SMPP data_sm_resp PDU
type DataSMResp struct {
	Header    *Header
	MessageID string
	TLVParams map[uint16]*TLVParam
}

// NewDataSMResp creates a new DataSMResp PDU
func NewDataSMResp() *DataSMResp {
	return &DataSMResp{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (d *DataSMResp) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(d.MessageID) + 1

	// Add TLV parameters length
	for _, tlv := range d.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	d.Header.CommandLength = uint32(length)
	d.Header.CommandID = DATA_SM_RESP

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
func (d *DataSMResp) Unmarshal(data []byte) error {
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

		tag := binary.BigEndian.Uint16(data[currentOffset:])
		length := binary.BigEndian.Uint16(data[currentOffset+2:])

		if remaining < int(4+length) {
			return errors.New("invalid TLV parameter: insufficient data for value")
		}

		tlv := &TLVParam{
			Tag:    tag,
			Length: length,
			Value:  make([]byte, length),
		}
		copy(tlv.Value, data[currentOffset+4:currentOffset+4+int(length)])

		d.TLVParams[tag] = tlv
		currentOffset += 4 + int(length)
		remaining -= 4 + int(length)
	}

	return nil
}
