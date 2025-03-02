package pdu

import "errors"

// CancelSM represents an SMPP cancel_sm PDU
type CancelSM struct {
	Header          *Header
	ServiceType     string
	MessageID       string
	SourceAddrTON   uint8
	SourceAddrNPI   uint8
	SourceAddr      string
	DestAddrTON     uint8
	DestAddrNPI     uint8
	DestinationAddr string
}

// NewCancelSM creates a new CancelSM PDU
func NewCancelSM() *CancelSM {
	return &CancelSM{
		Header: NewHeader(),
	}
}

// Marshal serializes the PDU into bytes
func (c *CancelSM) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(c.ServiceType) + 1
	length += len(c.MessageID) + 1
	length += 2 // Source addr ton and npi
	length += len(c.SourceAddr) + 1
	length += 2 // Dest addr ton and npi
	length += len(c.DestinationAddr) + 1

	// Set command length in header
	c.Header.CommandLength = uint32(length)
	c.Header.CommandID = CANCEL_SM

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

	// Write dest_addr_ton
	buf[offset] = c.DestAddrTON
	offset++

	// Write dest_addr_npi
	buf[offset] = c.DestAddrNPI
	offset++

	// Write destination_addr
	copy(buf[offset:], c.DestinationAddr)
	offset += len(c.DestinationAddr) + 1

	return buf, nil
}

// Unmarshal deserializes the PDU from bytes
func (c *CancelSM) Unmarshal(data []byte) error {
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

	// Read dest_addr_ton
	c.DestAddrTON = data[offset]
	offset++

	// Read dest_addr_npi
	c.DestAddrNPI = data[offset]
	offset++

	// Read destination_addr
	if c.DestinationAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Verify we've read all the data
	if offset != int(c.Header.CommandLength) {
		return errors.New("invalid PDU length")
	}

	return nil
}
