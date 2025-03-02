package pdu

// AlertNotification represents an SMPP Alert Notification PDU
type AlertNotification struct {
	Header         *Header
	SourceAddr     string
	SourceAddrTON  uint8
	SourceAddrNPI  uint8
	EsmeAddr       string
	EsmeAddrTON    uint8
	EsmeAddrNPI    uint8
	OptionalParams map[string]interface{}
}

// NewAlertNotification creates a new Alert Notification PDU
func NewAlertNotification() *AlertNotification {
	return &AlertNotification{
		Header:         NewHeader(),
		OptionalParams: make(map[string]interface{}),
	}
}

// SetSourceAddr sets the source address parameters
func (an *AlertNotification) SetSourceAddr(addr string, ton, npi uint8) {
	an.SourceAddr = addr
	an.SourceAddrTON = ton
	an.SourceAddrNPI = npi
}

// SetEsmeAddr sets the ESME address parameters
func (an *AlertNotification) SetEsmeAddr(addr string, ton, npi uint8) {
	an.EsmeAddr = addr
	an.EsmeAddrTON = ton
	an.EsmeAddrNPI = npi
}

// Marshal serializes the PDU into bytes
func (an *AlertNotification) Marshal() ([]byte, error) {
	// Calculate the total length
	length := 16                             // Header length
	length += 1 + 1 + len(an.SourceAddr) + 1 // Source address params
	length += 1 + 1 + len(an.EsmeAddr) + 1   // ESME address params

	// Set the command length in the header
	an.Header.CommandLength = uint32(length)
	an.Header.CommandID = ALERT_NOTIFICATION

	// Create the byte slice
	buf := make([]byte, length)

	// Marshal header
	headerBytes, err := an.Header.Marshal()
	if err != nil {
		return nil, err
	}
	copy(buf[0:], headerBytes)

	offset := 16

	// Write source address parameters
	buf[offset] = an.SourceAddrTON
	buf[offset+1] = an.SourceAddrNPI
	copy(buf[offset+2:], []byte(an.SourceAddr))
	offset += 2 + len(an.SourceAddr) + 1

	// Write ESME address parameters
	buf[offset] = an.EsmeAddrTON
	buf[offset+1] = an.EsmeAddrNPI
	copy(buf[offset+2:], []byte(an.EsmeAddr))

	return buf, nil
}

// Unmarshal deserializes the PDU from bytes
func (an *AlertNotification) Unmarshal(data []byte) error {
	// Unmarshal header
	an.Header = &Header{}
	if err := an.Header.Unmarshal(data[:16]); err != nil {
		return err
	}

	offset := 16

	// Read source address parameters
	an.SourceAddrTON = data[offset]
	an.SourceAddrNPI = data[offset+1]
	offset += 2

	// Read source address
	i := offset
	for ; data[i] != 0; i++ {
	}
	an.SourceAddr = string(data[offset:i])
	offset = i + 1

	// Read ESME address parameters
	an.EsmeAddrTON = data[offset]
	an.EsmeAddrNPI = data[offset+1]
	offset += 2

	// Read ESME address
	i = offset
	for ; data[i] != 0; i++ {
	}
	an.EsmeAddr = string(data[offset:i])

	return nil
}
