package pdu

import "errors"

// Outbind represents an SMPP outbind PDU
type Outbind struct {
	Header   *Header
	SystemID string
	Password string
}

// NewOutbind creates a new Outbind PDU
func NewOutbind() *Outbind {
	return &Outbind{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (o *Outbind) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(o.SystemID) + 1
	length += len(o.Password) + 1

	// Set command length in header
	o.Header.CommandLength = uint32(length)
	o.Header.CommandID = OUTBIND

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := o.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write system_id
	copy(buf[offset:], o.SystemID)
	offset += len(o.SystemID) + 1

	// Write password
	copy(buf[offset:], o.Password)
	offset += len(o.Password) + 1

	return buf, nil
}

// Unmarshal deserializes the PDU from bytes
func (o *Outbind) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	o.Header = &Header{}
	if err = o.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read system_id
	if o.SystemID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read password
	if o.Password, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Verify we've read all the data
	if offset != int(o.Header.CommandLength) {
		return errors.New("invalid PDU length")
	}

	return nil
}
