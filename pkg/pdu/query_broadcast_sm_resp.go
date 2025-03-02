package pdu

import (
	"encoding/binary"
	"errors"
)

// QueryBroadcastSMResp represents an SMPP query_broadcast_sm_resp PDU (SMPP v5.0)
type QueryBroadcastSMResp struct {
	Header    *Header
	MessageID string
	TLVParams map[uint16]*TLVParam
}

// NewQueryBroadcastSMResp creates a new QueryBroadcastSMResp PDU
func NewQueryBroadcastSMResp() *QueryBroadcastSMResp {
	return &QueryBroadcastSMResp{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (q *QueryBroadcastSMResp) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(q.MessageID) + 1

	// Add TLV parameters length
	for _, tlv := range q.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	q.Header.CommandLength = uint32(length)
	q.Header.CommandID = QUERY_BROADCAST_SM_RESP

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := q.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write message_id
	copy(buf[offset:], q.MessageID)
	offset += len(q.MessageID) + 1

	// Write TLV parameters
	for _, tlv := range q.TLVParams {
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
func (q *QueryBroadcastSMResp) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	q.Header = &Header{}
	if err = q.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read message_id
	if q.MessageID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read TLV parameters if any remain
	remaining := int(q.Header.CommandLength) - offset
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

		q.TLVParams[tag] = tlv
		currentOffset += 4 + int(length)
		remaining -= 4 + int(length)
	}

	return nil
}
