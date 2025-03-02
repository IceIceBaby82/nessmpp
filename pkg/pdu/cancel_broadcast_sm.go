package pdu

import (
	"encoding/binary"
	"errors"
)

// CancelBroadcastSM represents an SMPP cancel_broadcast_sm PDU (SMPP v5.0)
type CancelBroadcastSM struct {
	Header        *Header
	ServiceType   string
	MessageID     string
	SourceAddrTON uint8
	SourceAddrNPI uint8
	SourceAddr    string
	TLVParams     map[uint16]*TLVParam
}

// NewCancelBroadcastSM creates a new CancelBroadcastSM PDU
func NewCancelBroadcastSM() *CancelBroadcastSM {
	return &CancelBroadcastSM{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (c *CancelBroadcastSM) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(c.ServiceType) + 1
	length += len(c.MessageID) + 1
	length += 2 // Source addr ton and npi
	length += len(c.SourceAddr) + 1

	// Add TLV parameters length
	for _, tlv := range c.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	c.Header.CommandLength = uint32(length)
	c.Header.CommandID = CANCEL_BROADCAST_SM

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := c.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write service_type
	copy(buf[offset:], c.ServiceType)
	offset += len(c.ServiceType) + 1

	// Write message_id
	copy(buf[offset:], c.MessageID)
	offset += len(c.MessageID) + 1

	// Write source_addr_ton
	buf[offset] = c.SourceAddrTON
	offset++

	// Write source_addr_npi
	buf[offset] = c.SourceAddrNPI
	offset++

	// Write source_addr
	copy(buf[offset:], c.SourceAddr)
	offset += len(c.SourceAddr) + 1

	// Write TLV parameters
	for _, tlv := range c.TLVParams {
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
func (c *CancelBroadcastSM) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	c.Header = &Header{}
	if err = c.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read service_type
	if c.ServiceType, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read message_id
	if c.MessageID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read source_addr_ton
	c.SourceAddrTON = data[offset]
	offset++

	// Read source_addr_npi
	c.SourceAddrNPI = data[offset]
	offset++

	// Read source_addr
	if c.SourceAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read TLV parameters if any remain
	remaining := int(c.Header.CommandLength) - offset
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

		c.TLVParams[tag] = tlv
		currentOffset += 4 + int(length)
		remaining -= 4 + int(length)
	}

	return nil
}
