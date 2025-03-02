package pdu

import (
	"encoding/binary"
	"errors"
)

// BindTransceiver represents an SMPP bind_transceiver PDU
type BindTransceiver struct {
	Header           *Header
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion uint8
	AddrTON          uint8
	AddrNPI          uint8
	AddressRange     string
	TLVParams        map[uint16]*TLVParam
}

// NewBindTransceiver creates a new BindTransceiver PDU
func NewBindTransceiver() *BindTransceiver {
	return &BindTransceiver{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (bt *BindTransceiver) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(bt.SystemID) + 1
	length += len(bt.Password) + 1
	length += len(bt.SystemType) + 1
	length += 3 // Interface version + addr_ton + addr_npi
	length += len(bt.AddressRange) + 1

	// Add TLV parameters length
	for _, tlv := range bt.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	bt.Header.CommandLength = uint32(length)
	bt.Header.CommandID = BIND_TRANSCEIVER

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := bt.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write system_id
	copy(buf[offset:], bt.SystemID)
	offset += len(bt.SystemID) + 1

	// Write password
	copy(buf[offset:], bt.Password)
	offset += len(bt.Password) + 1

	// Write system_type
	copy(buf[offset:], bt.SystemType)
	offset += len(bt.SystemType) + 1

	// Write interface_version
	buf[offset] = bt.InterfaceVersion
	offset++

	// Write addr_ton
	buf[offset] = bt.AddrTON
	offset++

	// Write addr_npi
	buf[offset] = bt.AddrNPI
	offset++

	// Write address_range
	copy(buf[offset:], bt.AddressRange)
	offset += len(bt.AddressRange) + 1

	// Write TLV parameters
	for _, tlv := range bt.TLVParams {
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
func (bt *BindTransceiver) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	bt.Header = &Header{}
	if err = bt.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read system_id
	if bt.SystemID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read password
	if bt.Password, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read system_type
	if bt.SystemType, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read interface_version
	bt.InterfaceVersion = data[offset]
	offset++

	// Read addr_ton
	bt.AddrTON = data[offset]
	offset++

	// Read addr_npi
	bt.AddrNPI = data[offset]
	offset++

	// Read address_range
	if bt.AddressRange, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read TLV parameters if any remain
	remaining := int(bt.Header.CommandLength) - offset
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

		bt.TLVParams[tag] = tlv
		currentOffset += 4 + int(length)
		remaining -= 4 + int(length)
	}

	return nil
}
