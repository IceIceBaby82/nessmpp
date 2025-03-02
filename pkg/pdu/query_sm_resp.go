package pdu

import "errors"

// QuerySMResp represents an SMPP query_sm_resp PDU
type QuerySMResp struct {
	Header       *Header
	MessageID    string
	FinalDate    string
	MessageState uint8
	ErrorCode    uint8
	TLVParams    map[uint16]*TLVParam
}

// NewQuerySMResp creates a new QuerySMResp PDU
func NewQuerySMResp() *QuerySMResp {
	return &QuerySMResp{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (q *QuerySMResp) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(q.MessageID) + 1
	length += len(q.FinalDate) + 1
	length += 2 // Message state and error code

	// Add TLV parameters length
	for _, tlv := range q.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	q.Header.CommandLength = uint32(length)
	q.Header.CommandID = QUERY_SM_RESP

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

	// Write final_date
	copy(buf[offset:], q.FinalDate)
	offset += len(q.FinalDate) + 1

	// Write message_state
	buf[offset] = q.MessageState
	offset++

	// Write error_code
	buf[offset] = q.ErrorCode
	offset++

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
func (q *QuerySMResp) Unmarshal(data []byte) error {
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

	// Read final_date
	if q.FinalDate, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read message_state
	q.MessageState = data[offset]
	offset++

	// Read error_code
	q.ErrorCode = data[offset]
	offset++

	// Read TLV parameters if any remain
	remaining := int(q.Header.CommandLength) - offset
	currentOffset := offset

	for remaining > 0 {
		if remaining < 4 {
			return errors.New("invalid TLV parameter: insufficient data")
		}

		tlv := &TLVParam{}
		if err := tlv.Unmarshal(data[currentOffset:]); err != nil {
			return err
		}

		q.TLVParams[tlv.Tag] = tlv
		currentOffset += 4 + int(tlv.Length)
		remaining -= 4 + int(tlv.Length)
	}

	return nil
}
