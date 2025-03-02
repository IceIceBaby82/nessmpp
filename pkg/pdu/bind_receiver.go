package pdu

import (
	"encoding/binary"
	"errors"
)

// BindReceiver represents an SMPP bind_receiver PDU
type BindReceiver struct {
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

// NewBindReceiver creates a new BindReceiver PDU
func NewBindReceiver() *BindReceiver {
	return &BindReceiver{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (br *BindReceiver) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(br.SystemID) + 1
	length += len(br.Password) + 1
	length += len(br.SystemType) + 1
	length += 3 // Interface version + addr_ton + addr_npi
	length += len(br.AddressRange) + 1

	// Add TLV parameters length
	for _, tlv := range br.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	br.Header.CommandLength = uint32(length)
	br.Header.CommandID = BIND_RECEIVER

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := br.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write system_id
	copy(buf[offset:], br.SystemID)
	offset += len(br.SystemID) + 1

	// Write password
	copy(buf[offset:], br.Password)
	offset += len(br.Password) + 1

	// Write system_type
	copy(buf[offset:], br.SystemType)
	offset += len(br.SystemType) + 1

	// Write interface_version
	buf[offset] = br.InterfaceVersion
	offset++

	// Write addr_ton
	buf[offset] = br.AddrTON
	offset++

	// Write addr_npi
	buf[offset] = br.AddrNPI
	offset++

	// Write address_range
	copy(buf[offset:], br.AddressRange)
	offset += len(br.AddressRange) + 1

	// Write TLV parameters
	for _, tlv := range br.TLVParams {
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
func (br *BindReceiver) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	br.Header = &Header{}
	if err = br.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read system_id
	if br.SystemID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read password
	if br.Password, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read system_type
	if br.SystemType, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read interface_version
	br.InterfaceVersion = data[offset]
	offset++

	// Read addr_ton
	br.AddrTON = data[offset]
	offset++

	// Read addr_npi
	br.AddrNPI = data[offset]
	offset++

	// Read address_range
	if br.AddressRange, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read TLV parameters if any remain
	remaining := int(br.Header.CommandLength) - offset
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

		br.TLVParams[tag] = tlv
		currentOffset += 4 + int(length)
		remaining -= 4 + int(length)
	}

	return nil
}
