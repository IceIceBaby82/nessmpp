package pdu

import (
	"errors"
)

// SubmitSM represents an SMPP submit_sm PDU
type SubmitSM struct {
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

// NewSubmitSM creates a new SubmitSM PDU
func NewSubmitSM() *SubmitSM {
	return &SubmitSM{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// SetMessageText sets the message text with the specified data coding
func (s *SubmitSM) SetMessageText(text string, coding uint8) error {
	s.DataCoding = coding
	s.ShortMessage = []byte(text)
	if len(s.ShortMessage) > 255 {
		return errors.New("message text too long")
	}
	s.SMLength = uint8(len(s.ShortMessage))
	return nil
}

// SetUDH sets User Data Header for concatenated messages
func (s *SubmitSM) SetUDH(refNum uint16, total, seqNum uint8) {
	// Set ESM class to indicate UDH presence
	s.ESMClass |= 0x40

	// Create UDH for 16-bit reference
	udh := make([]byte, 6)
	udh[0] = 5                 // UDH Length
	udh[1] = 0x08              // IE Identifier (16-bit reference)
	udh[2] = 0x04              // IE Length
	udh[3] = byte(refNum >> 8) // Reference high byte
	udh[4] = byte(refNum)      // Reference low byte
	udh[5] = total             // Total segments
	udh[6] = seqNum            // Sequence number

	// Prepend UDH to message
	newMsg := make([]byte, len(udh)+len(s.ShortMessage))
	copy(newMsg, udh)
	copy(newMsg[len(udh):], s.ShortMessage)
	s.ShortMessage = newMsg
	s.SMLength = uint8(len(s.ShortMessage))
}

// Marshal serializes the PDU into bytes
func (s *SubmitSM) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(s.ServiceType) + 1
	length += 2 // Source addr ton and npi
	length += len(s.SourceAddr) + 1
	length += 2 // Dest addr ton and npi
	length += len(s.DestinationAddr) + 1
	length += 5 // ESM class, protocol id, priority flag, registered delivery, replace if present
	length += len(s.ScheduleDeliveryTime) + 1
	length += len(s.ValidityPeriod) + 1
	length += 3 // Data coding, sm default msg id, sm length
	length += len(s.ShortMessage)

	// Add TLV parameters length
	for _, tlv := range s.TLVParams {
		length += 4 + len(tlv.Value)
	}

	// Set command length in header
	s.Header.CommandLength = uint32(length)
	s.Header.CommandID = SUBMIT_SM

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := s.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write service_type
	copy(buf[offset:], s.ServiceType)
	offset += len(s.ServiceType) + 1

	// Write source_addr_ton
	buf[offset] = s.SourceAddrTON
	offset++

	// Write source_addr_npi
	buf[offset] = s.SourceAddrNPI
	offset++

	// Write source_addr
	copy(buf[offset:], s.SourceAddr)
	offset += len(s.SourceAddr) + 1

	// Write dest_addr_ton
	buf[offset] = s.DestAddrTON
	offset++

	// Write dest_addr_npi
	buf[offset] = s.DestAddrNPI
	offset++

	// Write destination_addr
	copy(buf[offset:], s.DestinationAddr)
	offset += len(s.DestinationAddr) + 1

	// Write esm_class
	buf[offset] = s.ESMClass
	offset++

	// Write protocol_id
	buf[offset] = s.ProtocolID
	offset++

	// Write priority_flag
	buf[offset] = s.PriorityFlag
	offset++

	// Write schedule_delivery_time
	copy(buf[offset:], s.ScheduleDeliveryTime)
	offset += len(s.ScheduleDeliveryTime) + 1

	// Write validity_period
	copy(buf[offset:], s.ValidityPeriod)
	offset += len(s.ValidityPeriod) + 1

	// Write registered_delivery
	buf[offset] = s.RegisteredDelivery
	offset++

	// Write replace_if_present
	buf[offset] = s.ReplaceIfPresent
	offset++

	// Write data_coding
	buf[offset] = s.DataCoding
	offset++

	// Write sm_default_msg_id
	buf[offset] = s.SMDefaultMsgID
	offset++

	// Write sm_length
	buf[offset] = s.SMLength
	offset++

	// Write short_message
	copy(buf[offset:], s.ShortMessage)
	offset += len(s.ShortMessage)

	// Write TLV parameters
	for _, tlv := range s.TLVParams {
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
func (s *SubmitSM) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	s.Header = &Header{}
	if err = s.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read service_type
	if s.ServiceType, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read source_addr_ton
	s.SourceAddrTON = data[offset]
	offset++

	// Read source_addr_npi
	s.SourceAddrNPI = data[offset]
	offset++

	// Read source_addr
	if s.SourceAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read dest_addr_ton
	s.DestAddrTON = data[offset]
	offset++

	// Read dest_addr_npi
	s.DestAddrNPI = data[offset]
	offset++

	// Read destination_addr
	if s.DestinationAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read esm_class
	s.ESMClass = data[offset]
	offset++

	// Read protocol_id
	s.ProtocolID = data[offset]
	offset++

	// Read priority_flag
	s.PriorityFlag = data[offset]
	offset++

	// Read schedule_delivery_time
	if s.ScheduleDeliveryTime, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read validity_period
	if s.ValidityPeriod, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read registered_delivery
	s.RegisteredDelivery = data[offset]
	offset++

	// Read replace_if_present
	s.ReplaceIfPresent = data[offset]
	offset++

	// Read data_coding
	s.DataCoding = data[offset]
	offset++

	// Read sm_default_msg_id
	s.SMDefaultMsgID = data[offset]
	offset++

	// Read sm_length
	s.SMLength = data[offset]
	offset++

	// Read short_message
	if int(s.SMLength) > len(data[offset:]) {
		return errors.New("invalid short message length")
	}
	s.ShortMessage = make([]byte, s.SMLength)
	copy(s.ShortMessage, data[offset:offset+int(s.SMLength)])
	offset += int(s.SMLength)

	// Read TLV parameters if any remain
	remaining := int(s.Header.CommandLength) - offset
	currentOffset := offset

	for remaining > 0 {
		if remaining < 4 {
			return errors.New("invalid TLV parameter: insufficient data")
		}

		tlv := &TLVParam{}
		if err := tlv.Unmarshal(data[currentOffset:]); err != nil {
			return err
		}

		s.TLVParams[tlv.Tag] = tlv
		currentOffset += 4 + int(tlv.Length)
		remaining -= 4 + int(tlv.Length)
	}

	return nil
}
