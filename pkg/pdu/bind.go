package pdu

import (
	"errors"
)

// BindTransmitter represents a bind_transmitter PDU
type BindTransmitter struct {
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

// BindReceiver represents a bind_receiver PDU
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

// BindTransceiver represents a bind_transceiver PDU
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

// NewBindTransmitter creates a new BindTransmitter PDU
func NewBindTransmitter() *BindTransmitter {
	return &BindTransmitter{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// NewBindReceiver creates a new BindReceiver PDU
func NewBindReceiver() *BindReceiver {
	return &BindReceiver{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// NewBindTransceiver creates a new BindTransceiver PDU
func NewBindTransceiver() *BindTransceiver {
	return &BindTransceiver{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (b *BindTransmitter) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(b.SystemID) + 1
	length += len(b.Password) + 1
	length += len(b.SystemType) + 1
	length += 3 // Interface version, TON, NPI
	length += len(b.AddressRange) + 1

	// Set command length in header
	b.Header.CommandLength = uint32(length)
	b.Header.CommandID = BIND_TRANSMITTER

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := b.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write system_id
	copy(buf[offset:], b.SystemID)
	offset += len(b.SystemID) + 1

	// Write password
	copy(buf[offset:], b.Password)
	offset += len(b.Password) + 1

	// Write system_type
	copy(buf[offset:], b.SystemType)
	offset += len(b.SystemType) + 1

	// Write interface version and address params
	buf[offset] = b.InterfaceVersion
	buf[offset+1] = b.AddrTON
	buf[offset+2] = b.AddrNPI
	offset += 3

	// Write address_range
	copy(buf[offset:], b.AddressRange)

	return buf, nil
}

// Marshal serializes the PDU into bytes
func (b *BindReceiver) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(b.SystemID) + 1
	length += len(b.Password) + 1
	length += len(b.SystemType) + 1
	length += 3 // Interface version, TON, NPI
	length += len(b.AddressRange) + 1

	// Set command length in header
	b.Header.CommandLength = uint32(length)
	b.Header.CommandID = BIND_RECEIVER

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := b.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write system_id
	copy(buf[offset:], b.SystemID)
	offset += len(b.SystemID) + 1

	// Write password
	copy(buf[offset:], b.Password)
	offset += len(b.Password) + 1

	// Write system_type
	copy(buf[offset:], b.SystemType)
	offset += len(b.SystemType) + 1

	// Write interface version and address params
	buf[offset] = b.InterfaceVersion
	buf[offset+1] = b.AddrTON
	buf[offset+2] = b.AddrNPI
	offset += 3

	// Write address_range
	copy(buf[offset:], b.AddressRange)

	return buf, nil
}

// Marshal serializes the PDU into bytes
func (b *BindTransceiver) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(b.SystemID) + 1
	length += len(b.Password) + 1
	length += len(b.SystemType) + 1
	length += 3 // Interface version, TON, NPI
	length += len(b.AddressRange) + 1

	// Set command length in header
	b.Header.CommandLength = uint32(length)
	b.Header.CommandID = BIND_TRANSCEIVER

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := b.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write system_id
	copy(buf[offset:], b.SystemID)
	offset += len(b.SystemID) + 1

	// Write password
	copy(buf[offset:], b.Password)
	offset += len(b.Password) + 1

	// Write system_type
	copy(buf[offset:], b.SystemType)
	offset += len(b.SystemType) + 1

	// Write interface version and address params
	buf[offset] = b.InterfaceVersion
	buf[offset+1] = b.AddrTON
	buf[offset+2] = b.AddrNPI
	offset += 3

	// Write address_range
	copy(buf[offset:], b.AddressRange)

	return buf, nil
}

// Unmarshal deserializes bytes into a BindTransmitter PDU
func (b *BindTransmitter) Unmarshal(data []byte) error {
	if len(data) < 16 {
		return errors.New("data too short for bind PDU")
	}

	// Read null-terminated strings
	var pos int
	var err error

	b.SystemID, pos, err = readCString(data, 0)
	if err != nil {
		return err
	}

	b.Password, pos, err = readCString(data, pos)
	if err != nil {
		return err
	}

	b.SystemType, pos, err = readCString(data, pos)
	if err != nil {
		return err
	}

	if len(data[pos:]) < 3 {
		return errors.New("data too short for bind parameters")
	}

	b.InterfaceVersion = data[pos]
	b.AddrTON = data[pos+1]
	b.AddrNPI = data[pos+2]
	pos += 3

	b.AddressRange, _, err = readCString(data, pos)
	if err != nil {
		return err
	}

	return nil
}

// Unmarshal deserializes bytes into a BindReceiver PDU
func (b *BindReceiver) Unmarshal(data []byte) error {
	if len(data) < 16 {
		return errors.New("data too short for bind PDU")
	}

	// Read null-terminated strings
	var pos int
	var err error

	b.SystemID, pos, err = readCString(data, 0)
	if err != nil {
		return err
	}

	b.Password, pos, err = readCString(data, pos)
	if err != nil {
		return err
	}

	b.SystemType, pos, err = readCString(data, pos)
	if err != nil {
		return err
	}

	if len(data[pos:]) < 3 {
		return errors.New("data too short for bind parameters")
	}

	b.InterfaceVersion = data[pos]
	b.AddrTON = data[pos+1]
	b.AddrNPI = data[pos+2]
	pos += 3

	b.AddressRange, _, err = readCString(data, pos)
	if err != nil {
		return err
	}

	return nil
}

// Unmarshal deserializes bytes into a BindTransceiver PDU
func (b *BindTransceiver) Unmarshal(data []byte) error {
	if len(data) < 16 {
		return errors.New("data too short for bind PDU")
	}

	// Read null-terminated strings
	var pos int
	var err error

	b.SystemID, pos, err = readCString(data, 0)
	if err != nil {
		return err
	}

	b.Password, pos, err = readCString(data, pos)
	if err != nil {
		return err
	}

	b.SystemType, pos, err = readCString(data, pos)
	if err != nil {
		return err
	}

	if len(data[pos:]) < 3 {
		return errors.New("data too short for bind parameters")
	}

	b.InterfaceVersion = data[pos]
	b.AddrTON = data[pos+1]
	b.AddrNPI = data[pos+2]
	pos += 3

	b.AddressRange, _, err = readCString(data, pos)
	if err != nil {
		return err
	}

	return nil
}

// readCString reads a null-terminated string from data starting at offset
// returns the string, the new offset, and any error
func readCString(data []byte, offset int) (string, int, error) {
	if offset >= len(data) {
		return "", offset, errors.New("offset beyond data length")
	}

	end := offset
	for end < len(data) && data[end] != 0 {
		end++
	}

	if end >= len(data) {
		return "", offset, errors.New("no null terminator found")
	}

	return string(data[offset:end]), end + 1, nil
}
