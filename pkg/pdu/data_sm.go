package pdu

import (
	"encoding/binary"
	"errors"
)

// DataSM represents an SMPP data_sm PDU
type DataSM struct {
	Header             *Header
	ServiceType        string
	SourceAddrTON      uint8
	SourceAddrNPI      uint8
	SourceAddr         string
	DestAddrTON        uint8
	DestAddrNPI        uint8
	DestinationAddr    string
	ESMClass           uint8
	RegisteredDelivery uint8
	DataCoding         uint8
	TLVParams          map[uint16]*TLVParam
}

// NewDataSM creates a new DataSM PDU
func NewDataSM() *DataSM {
	return &DataSM{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (d *DataSM) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(d.ServiceType) + 1
	length += 2 // Source addr ton and npi
	length += len(d.SourceAddr) + 1
	length += 2 // Dest addr ton and npi
	length += len(d.DestinationAddr) + 1
	length += 3 // ESM class, registered delivery, data coding

	// Add TLV parameters length
	for _, tlv := range d.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	d.Header.CommandLength = uint32(length)
	d.Header.CommandID = DATA_SM

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := d.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write service_type
	copy(buf[offset:], d.ServiceType)
	offset += len(d.ServiceType) + 1

	// Write source_addr_ton
	buf[offset] = d.SourceAddrTON
	offset++

	// Write source_addr_npi
	buf[offset] = d.SourceAddrNPI
	offset++

	// Write source_addr
	copy(buf[offset:], d.SourceAddr)
	offset += len(d.SourceAddr) + 1

	// Write dest_addr_ton
	buf[offset] = d.DestAddrTON
	offset++

	// Write dest_addr_npi
	buf[offset] = d.DestAddrNPI
	offset++

	// Write destination_addr
	copy(buf[offset:], d.DestinationAddr)
	offset += len(d.DestinationAddr) + 1

	// Write esm_class
	buf[offset] = d.ESMClass
	offset++

	// Write registered_delivery
	buf[offset] = d.RegisteredDelivery
	offset++

	// Write data_coding
	buf[offset] = d.DataCoding
	offset++

	// Write TLV parameters
	for _, tlv := range d.TLVParams {
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
func (d *DataSM) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	d.Header = &Header{}
	if err = d.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read service_type
	if d.ServiceType, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read source_addr_ton
	d.SourceAddrTON = data[offset]
	offset++

	// Read source_addr_npi
	d.SourceAddrNPI = data[offset]
	offset++

	// Read source_addr
	if d.SourceAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read dest_addr_ton
	d.DestAddrTON = data[offset]
	offset++

	// Read dest_addr_npi
	d.DestAddrNPI = data[offset]
	offset++

	// Read destination_addr
	if d.DestinationAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read esm_class
	d.ESMClass = data[offset]
	offset++

	// Read registered_delivery
	d.RegisteredDelivery = data[offset]
	offset++

	// Read data_coding
	d.DataCoding = data[offset]
	offset++

	// Read TLV parameters if any remain
	remaining := int(d.Header.CommandLength) - offset
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

		d.TLVParams[tag] = tlv
		currentOffset += 4 + int(length)
		remaining -= 4 + int(length)
	}

	return nil
}
