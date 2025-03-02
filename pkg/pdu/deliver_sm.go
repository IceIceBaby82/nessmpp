package pdu

import (
	"errors"
)

// DeliverSM represents an SMPP deliver_sm PDU
type DeliverSM struct {
	Header               *Header
	ServiceType          string
	SourceAddrTON        uint8
	SourceAddrNPI        uint8
	SourceAddr           string
	DestAddrTON          uint8
	DestAddrNPI          uint8
	DestinationAddr      string
	ESMClass             uint8
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   uint8
	ReplaceIfPresent     uint8
	DataCoding           uint8
	SMDefaultMsgID       uint8
	SMLength             uint8
	ShortMessage         []byte
	TLVParams            map[uint16]*TLVParam
}

// NewDeliverSM creates a new DeliverSM PDU
func NewDeliverSM() *DeliverSM {
	return &DeliverSM{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// IsDeliveryReceipt checks if this PDU is a delivery receipt
func (d *DeliverSM) IsDeliveryReceipt() bool {
	return (d.ESMClass & 0x04) != 0
}

// SetMessageText sets the message text with the specified data coding
func (d *DeliverSM) SetMessageText(text string, coding uint8) error {
	d.DataCoding = coding
	d.ShortMessage = []byte(text)
	if len(d.ShortMessage) > 255 {
		return errors.New("message text too long")
	}
	d.SMLength = uint8(len(d.ShortMessage))
	return nil
}

// Marshal serializes the PDU into bytes
func (d *DeliverSM) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(d.ServiceType) + 1
	length += 2 // Source addr ton and npi
	length += len(d.SourceAddr) + 1
	length += 2 // Dest addr ton and npi
	length += len(d.DestinationAddr) + 1
	length += 5 // ESM class, protocol id, priority flag, registered delivery, replace if present
	length += len(d.ScheduleDeliveryTime) + 1
	length += len(d.ValidityPeriod) + 1
	length += 3 // Data coding, sm default msg id, sm length
	length += len(d.ShortMessage)

	// Add TLV parameters length
	for _, tlv := range d.TLVParams {
		length += 4 + len(tlv.Value)
	}

	// Set command length in header
	d.Header.CommandLength = uint32(length)
	d.Header.CommandID = DELIVER_SM

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

	// Write protocol_id
	buf[offset] = d.ProtocolID
	offset++

	// Write priority_flag
	buf[offset] = d.PriorityFlag
	offset++

	// Write schedule_delivery_time
	copy(buf[offset:], d.ScheduleDeliveryTime)
	offset += len(d.ScheduleDeliveryTime) + 1

	// Write validity_period
	copy(buf[offset:], d.ValidityPeriod)
	offset += len(d.ValidityPeriod) + 1

	// Write registered_delivery
	buf[offset] = d.RegisteredDelivery
	offset++

	// Write replace_if_present
	buf[offset] = d.ReplaceIfPresent
	offset++

	// Write data_coding
	buf[offset] = d.DataCoding
	offset++

	// Write sm_default_msg_id
	buf[offset] = d.SMDefaultMsgID
	offset++

	// Write sm_length
	buf[offset] = d.SMLength
	offset++

	// Write short_message
	copy(buf[offset:], d.ShortMessage)
	offset += len(d.ShortMessage)

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
func (d *DeliverSM) Unmarshal(data []byte) error {
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

	// Read protocol_id
	d.ProtocolID = data[offset]
	offset++

	// Read priority_flag
	d.PriorityFlag = data[offset]
	offset++

	// Read schedule_delivery_time
	if d.ScheduleDeliveryTime, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read validity_period
	if d.ValidityPeriod, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read registered_delivery
	d.RegisteredDelivery = data[offset]
	offset++

	// Read replace_if_present
	d.ReplaceIfPresent = data[offset]
	offset++

	// Read data_coding
	d.DataCoding = data[offset]
	offset++

	// Read sm_default_msg_id
	d.SMDefaultMsgID = data[offset]
	offset++

	// Read sm_length
	d.SMLength = data[offset]
	offset++

	// Read short_message
	if int(d.SMLength) > len(data[offset:]) {
		return errors.New("invalid short message length")
	}
	d.ShortMessage = make([]byte, d.SMLength)
	copy(d.ShortMessage, data[offset:offset+int(d.SMLength)])
	offset += int(d.SMLength)

	// Read TLV parameters if any remain
	remaining := int(d.Header.CommandLength) - offset
	currentOffset := offset

	for remaining > 0 {
		if remaining < 4 {
			return errors.New("invalid TLV parameter: insufficient data")
		}

		tlv := &TLVParam{}
		if err := tlv.Unmarshal(data[currentOffset:]); err != nil {
			return err
		}

		d.TLVParams[tlv.Tag] = tlv
		currentOffset += 4 + int(tlv.Length)
		remaining -= 4 + int(tlv.Length)
	}

	return nil
}
