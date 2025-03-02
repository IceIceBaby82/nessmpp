package pdu

import (
	"encoding/binary"
)

// Header represents the standard SMPP PDU header
type Header struct {
	CommandLength  uint32
	CommandID      uint32
	CommandStatus  uint32
	SequenceNumber uint32
}

// NewHeader creates a new PDU header
func NewHeader() *Header {
	return &Header{}
}

// Marshal serializes the header into bytes
func (h *Header) Marshal() ([]byte, error) {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint32(buf[0:4], h.CommandLength)
	binary.BigEndian.PutUint32(buf[4:8], h.CommandID)
	binary.BigEndian.PutUint32(buf[8:12], h.CommandStatus)
	binary.BigEndian.PutUint32(buf[12:16], h.SequenceNumber)
	return buf, nil
}

// Unmarshal deserializes the header from bytes
func (h *Header) Unmarshal(data []byte) error {
	h.CommandLength = binary.BigEndian.Uint32(data[0:4])
	h.CommandID = binary.BigEndian.Uint32(data[4:8])
	h.CommandStatus = binary.BigEndian.Uint32(data[8:12])
	h.SequenceNumber = binary.BigEndian.Uint32(data[12:16])
	return nil
}
