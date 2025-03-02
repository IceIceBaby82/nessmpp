package pdu

import (
	"encoding/binary"
	"errors"
)

// BroadcastSM represents an SMPP broadcast_sm PDU (SMPP v5.0)
type BroadcastSM struct {
	Header               *Header
	ServiceType          string
	SourceAddrTON        uint8
	SourceAddrNPI        uint8
	SourceAddr           string
	MessageID            string
	PriorityFlag         uint8
	ScheduleDeliveryTime string
	ValidityPeriod       string
	ReplaceIfPresent     uint8
	DataCoding           uint8
	SMDefaultMsgID       uint8
	TLVParams            map[uint16]*TLVParam
}

// NewBroadcastSM creates a new BroadcastSM PDU
func NewBroadcastSM() *BroadcastSM {
	return &BroadcastSM{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// Marshal serializes the PDU into bytes
func (b *BroadcastSM) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(b.ServiceType) + 1
	length += 2 // Source addr ton and npi
	length += len(b.SourceAddr) + 1
	length += len(b.MessageID) + 1
	length += 1 // Priority flag
	length += len(b.ScheduleDeliveryTime) + 1
	length += len(b.ValidityPeriod) + 1
	length += 3 // Replace if present, data coding, sm default msg id

	// Add TLV parameters length
	for _, tlv := range b.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	b.Header.CommandLength = uint32(length)
	b.Header.CommandID = BROADCAST_SM

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := b.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write service_type
	copy(buf[offset:], b.ServiceType)
	offset += len(b.ServiceType) + 1

	// Write source_addr_ton
	buf[offset] = b.SourceAddrTON
	offset++

	// Write source_addr_npi
	buf[offset] = b.SourceAddrNPI
	offset++

	// Write source_addr
	copy(buf[offset:], b.SourceAddr)
	offset += len(b.SourceAddr) + 1

	// Write message_id
	copy(buf[offset:], b.MessageID)
	offset += len(b.MessageID) + 1

	// Write priority_flag
	buf[offset] = b.PriorityFlag
	offset++

	// Write schedule_delivery_time
	copy(buf[offset:], b.ScheduleDeliveryTime)
	offset += len(b.ScheduleDeliveryTime) + 1

	// Write validity_period
	copy(buf[offset:], b.ValidityPeriod)
	offset += len(b.ValidityPeriod) + 1

	// Write replace_if_present
	buf[offset] = b.ReplaceIfPresent
	offset++

	// Write data_coding
	buf[offset] = b.DataCoding
	offset++

	// Write sm_default_msg_id
	buf[offset] = b.SMDefaultMsgID
	offset++

	// Write TLV parameters
	for _, tlv := range b.TLVParams {
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
func (b *BroadcastSM) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	b.Header = &Header{}
	if err = b.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read service_type
	if b.ServiceType, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read source_addr_ton
	b.SourceAddrTON = data[offset]
	offset++

	// Read source_addr_npi
	b.SourceAddrNPI = data[offset]
	offset++

	// Read source_addr
	if b.SourceAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read message_id
	if b.MessageID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read priority_flag
	b.PriorityFlag = data[offset]
	offset++

	// Read schedule_delivery_time
	if b.ScheduleDeliveryTime, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read validity_period
	if b.ValidityPeriod, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read replace_if_present
	b.ReplaceIfPresent = data[offset]
	offset++

	// Read data_coding
	b.DataCoding = data[offset]
	offset++

	// Read sm_default_msg_id
	b.SMDefaultMsgID = data[offset]
	offset++

	// Read TLV parameters if any remain
	remaining := int(b.Header.CommandLength) - offset
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

		b.TLVParams[tag] = tlv
		currentOffset += 4 + int(length)
		remaining -= 4 + int(length)
	}

	return nil
}
