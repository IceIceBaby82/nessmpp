package pdu

import "errors"

// ReplaceSM represents an SMPP replace_sm PDU
type ReplaceSM struct {
	Header               *Header
	MessageID            string
	SourceAddrTON        uint8
	SourceAddrNPI        uint8
	SourceAddr           string
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   uint8
	SMDefaultMsgID       uint8
	SMLength             uint8
	ShortMessage         []byte
	TLVParams            map[uint16]*TLVParam
}

// NewReplaceSM creates a new ReplaceSM PDU
func NewReplaceSM() *ReplaceSM {
	return &ReplaceSM{
		Header:    NewHeader(),
		TLVParams: make(map[uint16]*TLVParam),
	}
}

// SetMessageText sets the message text
func (r *ReplaceSM) SetMessageText(text string) error {
	if len(text) > 254 {
		return errors.New("message text too long")
	}
	r.ShortMessage = []byte(text)
	r.SMLength = uint8(len(r.ShortMessage))
	return nil
}

// Marshal serializes the PDU into bytes
func (r *ReplaceSM) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16 // Header length
	length += len(r.MessageID) + 1
	length += 2 // Source addr ton and npi
	length += len(r.SourceAddr) + 1
	length += len(r.ScheduleDeliveryTime) + 1
	length += len(r.ValidityPeriod) + 1
	length += 3 // registered_delivery, sm_default_msg_id, sm_length
	length += len(r.ShortMessage)

	// Add TLV parameters length
	for _, tlv := range r.TLVParams {
		length += 4 + len(tlv.Value) // Tag(2) + Length(2) + Value(n)
	}

	// Set command length in header
	r.Header.CommandLength = uint32(length)
	r.Header.CommandID = REPLACE_SM

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := r.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write message_id
	copy(buf[offset:], r.MessageID)
	offset += len(r.MessageID) + 1

	// Write source_addr_ton
	buf[offset] = r.SourceAddrTON
	offset++

	// Write source_addr_npi
	buf[offset] = r.SourceAddrNPI
	offset++

	// Write source_addr
	copy(buf[offset:], r.SourceAddr)
	offset += len(r.SourceAddr) + 1

	// Write schedule_delivery_time
	copy(buf[offset:], r.ScheduleDeliveryTime)
	offset += len(r.ScheduleDeliveryTime) + 1

	// Write validity_period
	copy(buf[offset:], r.ValidityPeriod)
	offset += len(r.ValidityPeriod) + 1

	// Write registered_delivery
	buf[offset] = r.RegisteredDelivery
	offset++

	// Write sm_default_msg_id
	buf[offset] = r.SMDefaultMsgID
	offset++

	// Write sm_length
	buf[offset] = r.SMLength
	offset++

	// Write short_message
	copy(buf[offset:], r.ShortMessage)
	offset += len(r.ShortMessage)

	// Write TLV parameters
	for _, tlv := range r.TLVParams {
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
func (r *ReplaceSM) Unmarshal(data []byte) error {
	var err error

	// Unmarshal header
	r.Header = &Header{}
	if err = r.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read message_id
	if r.MessageID, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read source_addr_ton
	r.SourceAddrTON = data[offset]
	offset++

	// Read source_addr_npi
	r.SourceAddrNPI = data[offset]
	offset++

	// Read source_addr
	if r.SourceAddr, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read schedule_delivery_time
	if r.ScheduleDeliveryTime, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read validity_period
	if r.ValidityPeriod, offset, err = ReadCString(data[offset:]); err != nil {
		return err
	}

	// Read registered_delivery
	r.RegisteredDelivery = data[offset]
	offset++

	// Read sm_default_msg_id
	r.SMDefaultMsgID = data[offset]
	offset++

	// Read sm_length
	r.SMLength = data[offset]
	offset++

	// Read short_message
	if int(r.SMLength) > len(data[offset:]) {
		return errors.New("invalid short message length")
	}
	r.ShortMessage = make([]byte, r.SMLength)
	copy(r.ShortMessage, data[offset:offset+int(r.SMLength)])
	offset += int(r.SMLength)

	// Read TLV parameters if any remain
	remaining := int(r.Header.CommandLength) - offset
	currentOffset := offset

	for remaining > 0 {
		if remaining < 4 {
			return errors.New("invalid TLV parameter: insufficient data")
		}

		tlv := &TLVParam{}
		if err := tlv.Unmarshal(data[currentOffset:]); err != nil {
			return err
		}

		r.TLVParams[tlv.Tag] = tlv
		currentOffset += 4 + int(tlv.Length)
		remaining -= 4 + int(tlv.Length)
	}

	return nil
}
