package pdu

import "errors"

// QuerySM represents an SMPP query_sm PDU
type QuerySM struct {
	Header        *Header
	MessageID     string
	SourceAddrTON uint8
	SourceAddrNPI uint8
	SourceAddr    string
}

// NewQuerySM creates a new QuerySM PDU
func NewQuerySM() *QuerySM {
	return &QuerySM{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (q *QuerySM) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(q.MessageID) + 1
	length += 2 // Source addr ton and npi
	length += len(q.SourceAddr) + 1

	// Set command length in header
	q.Header.CommandLength = uint32(length)
	q.Header.CommandID = QUERY_SM

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

	// Write source_addr_ton
	buf[offset] = q.SourceAddrTON
	offset++

	// Write source_addr_npi
	buf[offset] = q.SourceAddrNPI
	offset++

	// Write source_addr
	copy(buf[offset:], q.SourceAddr)
	offset += len(q.SourceAddr) + 1

	return buf, nil
}

// Unmarshal deserializes the PDU from bytes
func (q *QuerySM) Unmarshal(data []byte) error {
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

	// Read source_addr_ton
	q.SourceAddrTON = data[offset]
	offset++

	// Read source_addr_npi
	q.SourceAddrNPI = data[offset]
	offset++

	// Read source_addr
	if q.SourceAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Verify we've read all the data
	if offset != int(q.Header.CommandLength) {
		return errors.New("invalid PDU length")
	}

	return nil
}
