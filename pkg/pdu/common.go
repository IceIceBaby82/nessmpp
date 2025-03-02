package pdu

import "errors"

// Data Coding Scheme Constants
const (
	// GSM 03.38 Coding Schemes
	DATA_CODING_DEFAULT    uint8 = 0x00 // Default GSM 7-bit alphabet
	DATA_CODING_IA5        uint8 = 0x01 // IA5 (CCITT T.50)/ASCII (ANSI X3.4)
	DATA_CODING_BINARY     uint8 = 0x02 // 8-bit binary
	DATA_CODING_ISO8859_1  uint8 = 0x03 // ISO-8859-1 (Latin-1)
	DATA_CODING_UCS2       uint8 = 0x08 // UCS2 (ISO/IEC-10646)
	DATA_CODING_PICTOGRAM  uint8 = 0x09 // Pictogram Encoding
	DATA_CODING_ISO2022_JP uint8 = 0x0A // ISO-2022-JP (Music Codes)
	DATA_CODING_KANJI      uint8 = 0x0D // Extended Kanji JIS
	DATA_CODING_KSC5601    uint8 = 0x0E // KS C 5601

	// Message Class Constants
	MSG_CLASS_0 uint8 = 0x00 // Class 0 (Immediate Display)
	MSG_CLASS_1 uint8 = 0x01 // Class 1 (ME Specific)
	MSG_CLASS_2 uint8 = 0x02 // Class 2 (SIM Specific)
	MSG_CLASS_3 uint8 = 0x03 // Class 3 (TE Specific)
)

// National Language Identifier Constants
const (
	NLI_DEFAULT    uint8 = 0x00
	NLI_TURKISH    uint8 = 0x01
	NLI_SPANISH    uint8 = 0x02
	NLI_PORTUGUESE uint8 = 0x03
	NLI_BENGALI    uint8 = 0x04
	NLI_GUJARATI   uint8 = 0x05
	NLI_HINDI      uint8 = 0x06
	NLI_KANNADA    uint8 = 0x07
	NLI_MALAYALAM  uint8 = 0x08
	NLI_ORIYA      uint8 = 0x09
	NLI_PUNJABI    uint8 = 0x0A
	NLI_TAMIL      uint8 = 0x0B
	NLI_TELUGU     uint8 = 0x0C
	NLI_URDU       uint8 = 0x0D
)

// SAR (Segmentation And Reassembly) Parameters
type SARParams struct {
	RefNum uint16 // Reference number
	Total  uint8  // Total number of segments
	SeqNum uint8  // Sequence number of this segment
}

// UDH (User Data Header) Information Element Identifiers
const (
	UDH_IE_CONCAT_8BIT      uint8 = 0x00 // Concatenated messages, 8-bit reference
	UDH_IE_SPECIAL_SMS      uint8 = 0x01 // Special SMS Message Indication
	UDH_IE_PORT_8BIT        uint8 = 0x04 // Application port addressing scheme, 8 bit
	UDH_IE_PORT_16BIT       uint8 = 0x05 // Application port addressing scheme, 16 bit
	UDH_IE_CONCAT_16BIT     uint8 = 0x08 // Concatenated messages, 16-bit reference
	UDH_IE_WIRELESS_CTRL    uint8 = 0x09 // Wireless Control Message Protocol
	UDH_IE_TEXT_FORMAT      uint8 = 0x0A // Text Formatting
	UDH_IE_PREDEFINED_SOUND uint8 = 0x0B // Predefined Sound
	UDH_IE_USER_PROMPT      uint8 = 0x0C // User Prompt Indicator
	UDH_IE_EMS_VAR_PIC      uint8 = 0x0D // Extended Object
)

// TLV (Tag Length Value) Tag Definitions
const (
	// SMPP v3.4 TLV Tags
	TLV_DEST_ADDR_SUBUNIT            uint16 = 0x0005
	TLV_DEST_NETWORK_TYPE            uint16 = 0x0006
	TLV_DEST_BEARER_TYPE             uint16 = 0x0007
	TLV_DEST_TELEMATICS_ID           uint16 = 0x0008
	TLV_SOURCE_ADDR_SUBUNIT          uint16 = 0x000D
	TLV_SOURCE_NETWORK_TYPE          uint16 = 0x000E
	TLV_SOURCE_BEARER_TYPE           uint16 = 0x000F
	TLV_SOURCE_TELEMATICS_ID         uint16 = 0x0010
	TLV_QOS_TIME_TO_LIVE             uint16 = 0x0017
	TLV_PAYLOAD_TYPE                 uint16 = 0x0019
	TLV_ADDITIONAL_STATUS_INFO_TEXT  uint16 = 0x001D
	TLV_RECEIPTED_MESSAGE_ID         uint16 = 0x001E
	TLV_MS_MSG_WAIT_FACILITIES       uint16 = 0x0030
	TLV_PRIVACY_INDICATOR            uint16 = 0x0201
	TLV_SOURCE_SUBADDRESS            uint16 = 0x0202
	TLV_DEST_SUBADDRESS              uint16 = 0x0203
	TLV_USER_MESSAGE_REFERENCE       uint16 = 0x0204
	TLV_USER_RESPONSE_CODE           uint16 = 0x0205
	TLV_SOURCE_PORT                  uint16 = 0x020A
	TLV_DESTINATION_PORT             uint16 = 0x020B
	TLV_SAR_MSG_REF_NUM              uint16 = 0x020C
	TLV_LANGUAGE_INDICATOR           uint16 = 0x020D
	TLV_SAR_TOTAL_SEGMENTS           uint16 = 0x020E
	TLV_SAR_SEGMENT_SEQNUM           uint16 = 0x020F
	TLV_SC_INTERFACE_VERSION         uint16 = 0x0210
	TLV_CALLBACK_NUM_PRES_IND        uint16 = 0x0302
	TLV_CALLBACK_NUM_ATAG            uint16 = 0x0303
	TLV_NUMBER_OF_MESSAGES           uint16 = 0x0304
	TLV_CALLBACK_NUM                 uint16 = 0x0381
	TLV_DPF_RESULT                   uint16 = 0x0420
	TLV_SET_DPF                      uint16 = 0x0421
	TLV_MS_AVAILABILITY_STATUS       uint16 = 0x0422
	TLV_NETWORK_ERROR_CODE           uint16 = 0x0423
	TLV_MESSAGE_PAYLOAD              uint16 = 0x0424
	TLV_DELIVERY_FAILURE_REASON      uint16 = 0x0425
	TLV_MORE_MESSAGES_TO_SEND        uint16 = 0x0426
	TLV_MESSAGE_STATE                uint16 = 0x0427
	TLV_CONGESTION_STATE             uint16 = 0x0428
	TLV_USSD_SERVICE_OP              uint16 = 0x0501
	TLV_BROADCAST_CHANNEL_INDICATOR  uint16 = 0x0600
	TLV_BROADCAST_CONTENT_TYPE       uint16 = 0x0601
	TLV_BROADCAST_CONTENT_TYPE_INFO  uint16 = 0x0602
	TLV_BROADCAST_MESSAGE_CLASS      uint16 = 0x0603
	TLV_BROADCAST_REP_NUM            uint16 = 0x0604
	TLV_BROADCAST_FREQUENCY_INTERVAL uint16 = 0x0605
	TLV_BROADCAST_AREA_IDENTIFIER    uint16 = 0x0606
	TLV_BROADCAST_ERROR_STATUS       uint16 = 0x0607
	TLV_BROADCAST_AREA_SUCCESS       uint16 = 0x0608
	TLV_BROADCAST_END_TIME           uint16 = 0x0609
	TLV_BROADCAST_SERVICE_GROUP      uint16 = 0x060A
	TLV_BILLING_IDENTIFICATION       uint16 = 0x060B
	TLV_SOURCE_NETWORK_ID            uint16 = 0x060D
	TLV_DEST_NETWORK_ID              uint16 = 0x060E
	TLV_SOURCE_NODE_ID               uint16 = 0x060F
	TLV_DEST_NODE_ID                 uint16 = 0x0610
	TLV_DEST_ADDR_NP_RESOLUTION      uint16 = 0x0611
	TLV_DEST_ADDR_NP_INFORMATION     uint16 = 0x0612
	TLV_DEST_ADDR_NP_COUNTRY         uint16 = 0x0613
)

// TLV Parameter
type TLVParam struct {
	Tag    uint16
	Length uint16
	Value  []byte
}

// TLV Methods
func NewTLVParam(tag uint16, value []byte) *TLVParam {
	return &TLVParam{
		Tag:    tag,
		Length: uint16(len(value)),
		Value:  value,
	}
}

func (tlv *TLVParam) Marshal() ([]byte, error) {
	buf := make([]byte, 4+len(tlv.Value))

	// Write Tag (2 bytes)
	buf[0] = byte(tlv.Tag >> 8)
	buf[1] = byte(tlv.Tag)

	// Write Length (2 bytes)
	buf[2] = byte(tlv.Length >> 8)
	buf[3] = byte(tlv.Length)

	// Write Value
	copy(buf[4:], tlv.Value)

	return buf, nil
}

func (tlv *TLVParam) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return errors.New("invalid TLV data: too short")
	}

	tlv.Tag = uint16(data[0])<<8 | uint16(data[1])
	tlv.Length = uint16(data[2])<<8 | uint16(data[3])

	if len(data) < int(4+tlv.Length) {
		return errors.New("invalid TLV data: value length mismatch")
	}

	tlv.Value = make([]byte, tlv.Length)
	copy(tlv.Value, data[4:4+tlv.Length])

	return nil
}
