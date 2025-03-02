package pdu

// SMPP Versions
const (
	SMPP_V33 uint32 = 0x33 // SMPP version 3.3
	SMPP_V34 uint32 = 0x34 // SMPP version 3.4
	SMPP_V50 uint32 = 0x50 // SMPP version 5.0
)

// Interface Versions
const (
	IF_VERSION_33 string = "3.3"
	IF_VERSION_34 string = "3.4"
	IF_VERSION_50 string = "5.0"
)

// SMPP Command IDs
const (
	GENERIC_NACK          uint32 = 0x80000000
	BIND_RECEIVER         uint32 = 0x00000001
	BIND_RECEIVER_RESP    uint32 = 0x80000001
	BIND_TRANSMITTER      uint32 = 0x00000002
	BIND_TRANSMITTER_RESP uint32 = 0x80000002
	QUERY_SM              uint32 = 0x00000003
	QUERY_SM_RESP         uint32 = 0x80000003
	SUBMIT_SM             uint32 = 0x00000004
	SUBMIT_SM_RESP        uint32 = 0x80000004
	DELIVER_SM            uint32 = 0x00000005
	DELIVER_SM_RESP       uint32 = 0x80000005
	UNBIND                uint32 = 0x00000006
	UNBIND_RESP           uint32 = 0x80000006
	REPLACE_SM            uint32 = 0x00000007
	REPLACE_SM_RESP       uint32 = 0x80000007
	CANCEL_SM             uint32 = 0x00000008
	CANCEL_SM_RESP        uint32 = 0x80000008
	BIND_TRANSCEIVER      uint32 = 0x00000009
	BIND_TRANSCEIVER_RESP uint32 = 0x80000009
	OUTBIND               uint32 = 0x0000000B
	ENQUIRE_LINK          uint32 = 0x00000015
	ENQUIRE_LINK_RESP     uint32 = 0x80000015
	SUBMIT_MULTI          uint32 = 0x00000021
	SUBMIT_MULTI_RESP     uint32 = 0x80000021
	ALERT_NOTIFICATION    uint32 = 0x00000102
	DATA_SM               uint32 = 0x00000103
	DATA_SM_RESP          uint32 = 0x80000103
	// SMPP v5.0 specific commands
	BROADCAST_SM             uint32 = 0x00000111 // v5.0
	BROADCAST_SM_RESP        uint32 = 0x80000111 // v5.0
	QUERY_BROADCAST_SM       uint32 = 0x00000112 // v5.0
	QUERY_BROADCAST_SM_RESP  uint32 = 0x80000112 // v5.0
	CANCEL_BROADCAST_SM      uint32 = 0x00000113 // v5.0
	CANCEL_BROADCAST_SM_RESP uint32 = 0x80000113 // v5.0
)

// SMPP Command Statuses
const (
	ESME_ROK              uint32 = 0x00000000 // No Error
	ESME_RINVMSGLEN       uint32 = 0x00000001 // Message Length is invalid
	ESME_RINVCMDID        uint32 = 0x00000002 // Invalid Command ID
	ESME_RINVPASWD        uint32 = 0x0000000E // Invalid Password
	ESME_RINVSYSID        uint32 = 0x0000000F // Invalid System ID
	ESME_ROK              uint32 = 0x00000000 // No Error
	ESME_RINVMSGLEN       uint32 = 0x00000001 // Message Length is invalid
	ESME_RINVCMDLEN       uint32 = 0x00000002 // Command Length is invalid
	ESME_RINVCMDID        uint32 = 0x00000003 // Invalid Command ID
	ESME_RINVBNDSTS       uint32 = 0x00000004 // Incorrect BIND Status for given command
	ESME_RALYBND          uint32 = 0x00000005 // ESME Already in Bound State
	ESME_RINVPRTFLG       uint32 = 0x00000006 // Invalid Priority Flag
	ESME_RINVREGDLVFLG    uint32 = 0x00000007 // Invalid Registered Delivery Flag
	ESME_RSYSERR          uint32 = 0x00000008 // System Error
	ESME_RINVSRCADR       uint32 = 0x0000000A // Invalid Source Address
	ESME_RINVDSTADR       uint32 = 0x0000000B // Invalid Dest Addr
	ESME_RINVMSGID        uint32 = 0x0000000C // Message ID is invalid
	ESME_RBINDFAIL        uint32 = 0x0000000D // Bind Failed
	ESME_RINVPASWD        uint32 = 0x0000000E // Invalid Password
	ESME_RINVSYSID        uint32 = 0x0000000F // Invalid System ID
	ESME_RCANCELFAIL      uint32 = 0x00000011 // Cancel SM Failed
	ESME_RREPLACEFAIL     uint32 = 0x00000013 // Replace SM Failed
	ESME_RMSGQFUL         uint32 = 0x00000014 // Message Queue Full
	ESME_RINVSERTYP       uint32 = 0x00000015 // Invalid Service Type
	ESME_RINVNUMDESTS     uint32 = 0x00000033 // Invalid number of destinations
	ESME_RINVDLNAME       uint32 = 0x00000034 // Invalid Distribution List name
	ESME_RINVDESTFLAG     uint32 = 0x00000040 // Destination flag is invalid
	ESME_RINVSUBREP       uint32 = 0x00000042 // Invalid 'submit with replace' request
	ESME_RINVESMCLASS     uint32 = 0x00000043 // Invalid esm_class field data
	ESME_RCNTSUBDL        uint32 = 0x00000044 // Cannot Submit to Distribution List
	ESME_RSUBMITFAIL      uint32 = 0x00000045 // submit_sm or submit_multi failed
	ESME_RINVSRCTON       uint32 = 0x00000048 // Invalid Source address TON
	ESME_RINVSRCNPI       uint32 = 0x00000049 // Invalid Source address NPI
	ESME_RINVDSTTON       uint32 = 0x00000050 // Invalid Destination address TON
	ESME_RINVDSTNPI       uint32 = 0x00000051 // Invalid Destination address NPI
	ESME_RINVSYSTYP       uint32 = 0x00000053 // Invalid system_type field
	ESME_RINVREPFLAG      uint32 = 0x00000054 // Invalid replace_if_present flag
	ESME_RINVNUMMSGS      uint32 = 0x00000055 // Invalid number of messages
	ESME_RTHROTTLED       uint32 = 0x00000058 // Throttling error
	ESME_RINVSCHED        uint32 = 0x00000061 // Invalid Scheduled Delivery Time
	ESME_RINVEXPIRY       uint32 = 0x00000062 // Invalid message validity period (Expiry time)
	ESME_RINVDFTMSGID     uint32 = 0x00000063 // Predefined Message Invalid or Not Found
	ESME_RX_T_APPN        uint32 = 0x00000064 // ESME Receiver Temporary App Error Code
	ESME_RX_P_APPN        uint32 = 0x00000065 // ESME Receiver Permanent App Error Code
	ESME_RX_R_APPN        uint32 = 0x00000066 // ESME Receiver Reject Message Error Code
	ESME_RQUERYFAIL       uint32 = 0x00000067 // query_sm request failed
	ESME_RINVOPTPARSTREAM uint32 = 0x000000C0 // Error in the optional part of the PDU Body
	ESME_ROPTPARNOTALLWD  uint32 = 0x000000C1 // Optional Parameter not allowed
	ESME_RINVPARLEN       uint32 = 0x000000C2 // Invalid Parameter Length
	ESME_RMISSINGOPTPARAM uint32 = 0x000000C3 // Expected Optional Parameter missing
	ESME_RINVOPTPARAMVAL  uint32 = 0x000000C4 // Invalid Optional Parameter Value
	ESME_RDELIVERYFAILURE uint32 = 0x000000FE // Delivery Failure (used for data_sm_resp)
	ESME_RUNKNOWNERR      uint32 = 0x000000FF // Unknown Error

	// Network specific status codes (0x0100 - 0x03FF)
	ESME_NETWORK_ERROR_START     uint32 = 0x00000100
	ESME_RMESSAGE_TOO_LONG       uint32 = 0x00000101 // Message too long
	ESME_RSERVICE_TYPE_NOT_FOUND uint32 = 0x00000102 // Service type not found
	ESME_ROPER_NOT_ALLOWED       uint32 = 0x00000103 // Operation not allowed
	ESME_RSERVICE_UNAVAILABLE    uint32 = 0x00000104 // Service unavailable
	ESME_RSERVICE_DENIED         uint32 = 0x00000105 // Service denied
	ESME_RINVALID_REFERENCE      uint32 = 0x00000106 // Invalid reference number
	ESME_RINVALID_DELIVERY_TIME  uint32 = 0x00000107 // Invalid delivery time
	ESME_RINVALID_DESTS          uint32 = 0x00000108 // Invalid destinations
	ESME_RUNKNOWN_DEST           uint32 = 0x00000109 // Unknown destination
	ESME_RDEST_UNAVAILABLE       uint32 = 0x0000010A // Destination unavailable
	ESME_RDEST_FLAGGED           uint32 = 0x0000010B // Destination flagged as invalid
	ESME_RDUPLICATE_MSGID        uint32 = 0x0000010C // Duplicate message ID
	ESME_RPAYLOAD_NOT_SUPPORTED  uint32 = 0x0000010D // Payload not supported

	// Application specific status codes (0x0400 - 0x04FF)
	ESME_APP_ERROR_START      uint32 = 0x00000400
	ESME_RAPP_BUSY            uint32 = 0x00000401 // Application busy
	ESME_RAPP_QUEUE_FULL      uint32 = 0x00000402 // Application queue full
	ESME_RAPP_NOT_AVAILABLE   uint32 = 0x00000403 // Application not available
	ESME_RAPP_INVALID_REQUEST uint32 = 0x00000404 // Invalid application request

	// Security related status codes (0x0500 - 0x05FF)
	ESME_SECURITY_ERROR_START    uint32 = 0x00000500
	ESME_RAUTHENTICATION_FAILURE uint32 = 0x00000501 // Authentication failure
	ESME_RSECURITY_VIOLATION     uint32 = 0x00000502 // Security violation
	ESME_RPROHIBITED_BY_SECURITY uint32 = 0x00000503 // Prohibited by security settings

	// Billing related status codes (0x0600 - 0x06FF)
	ESME_BILLING_ERROR_START    uint32 = 0x00000600
	ESME_RINSUFFICIENT_CREDITS  uint32 = 0x00000601 // Insufficient credits
	ESME_RBILLING_FAILED        uint32 = 0x00000602 // Billing operation failed
	ESME_RBILLING_NOT_SUPPORTED uint32 = 0x00000603 // Billing not supported

	// Batch submission related status codes (0x0700 - 0x07FF)
	ESME_BATCH_ERROR_START        uint32 = 0x00000700
	ESME_RBATCH_SUBMISSION_FAILED uint32 = 0x00000701 // Batch submission failed
	ESME_RBATCH_SIZE_EXCEEDED     uint32 = 0x00000702 // Batch size exceeded

	// Message format related status codes (0x0800 - 0x08FF)
	ESME_FORMAT_ERROR_START   uint32 = 0x00000800
	ESME_RINVALID_MSG_FORMAT  uint32 = 0x00000801 // Invalid message format
	ESME_RUNSUPPORTED_CHARSET uint32 = 0x00000802 // Unsupported character set
	ESME_RINVALID_ENCODING    uint32 = 0x00000803 // Invalid encoding

	// Protocol related status codes (0x0900 - 0x09FF)
	ESME_PROTOCOL_ERROR_START uint32 = 0x00000900
	ESME_RPROTOCOL_VERSION    uint32 = 0x00000901 // Protocol version mismatch
	ESME_RSEQUENCE_ERROR      uint32 = 0x00000902 // Sequence number error
	ESME_RINVALID_TLV         uint32 = 0x00000903 // Invalid TLV parameter

	// Temporary errors (0x0A00 - 0x0AFF)
	ESME_TEMP_ERROR_START    uint32 = 0x00000A00
	ESME_RTEMP_NETWORK_ERROR uint32 = 0x00000A01 // Temporary network error
	ESME_RTEMP_SYSTEM_ERROR  uint32 = 0x00000A02 // Temporary system error
	ESME_RTEMP_APP_ERROR     uint32 = 0x00000A03 // Temporary application error

	// Permanent errors (0x0B00 - 0x0BFF)
	ESME_PERM_ERROR_START    uint32 = 0x00000B00
	ESME_RPERM_NETWORK_ERROR uint32 = 0x00000B01 // Permanent network error
	ESME_RPERM_SYSTEM_ERROR  uint32 = 0x00000B02 // Permanent system error
	ESME_RPERM_APP_ERROR     uint32 = 0x00000B03 // Permanent application error

	// Reserved for future use (0x0C00 - 0xFFFF)
	ESME_RESERVED_START uint32 = 0x00000C00
	ESME_RESERVED_END   uint32 = 0x0000FFFF
)

// SMPP v3.4 specific status codes
const (
	// Message State
	SMPP_34_MESSAGE_STATE_ENROUTE       uint32 = 0x00000001
	SMPP_34_MESSAGE_STATE_DELIVERED     uint32 = 0x00000002
	SMPP_34_MESSAGE_STATE_EXPIRED       uint32 = 0x00000003
	SMPP_34_MESSAGE_STATE_DELETED       uint32 = 0x00000004
	SMPP_34_MESSAGE_STATE_UNDELIVERABLE uint32 = 0x00000005
	SMPP_34_MESSAGE_STATE_ACCEPTED      uint32 = 0x00000006
	SMPP_34_MESSAGE_STATE_UNKNOWN       uint32 = 0x00000007
	SMPP_34_MESSAGE_STATE_REJECTED      uint32 = 0x00000008

	// TON (Type of Number) Values
	SMPP_34_TON_UNKNOWN           uint8 = 0x00
	SMPP_34_TON_INTERNATIONAL     uint8 = 0x01
	SMPP_34_TON_NATIONAL          uint8 = 0x02
	SMPP_34_TON_NETWORK_SPECIFIC  uint8 = 0x03
	SMPP_34_TON_SUBSCRIBER_NUMBER uint8 = 0x04
	SMPP_34_TON_ALPHANUMERIC      uint8 = 0x05
	SMPP_34_TON_ABBREVIATED       uint8 = 0x06

	// NPI (Numbering Plan Indicator) Values
	SMPP_34_NPI_UNKNOWN       uint8 = 0x00
	SMPP_34_NPI_ISDN          uint8 = 0x01
	SMPP_34_NPI_DATA          uint8 = 0x03
	SMPP_34_NPI_TELEX         uint8 = 0x04
	SMPP_34_NPI_LAND_MOBILE   uint8 = 0x06
	SMPP_34_NPI_NATIONAL      uint8 = 0x08
	SMPP_34_NPI_PRIVATE       uint8 = 0x09
	SMPP_34_NPI_ERMES         uint8 = 0x0A
	SMPP_34_NPI_INTERNET      uint8 = 0x0E
	SMPP_34_NPI_WAP_CLIENT_ID uint8 = 0x12
)

// SMPP v5.0 specific status codes
const (
	// Broadcast Message State
	SMPP_50_BCAST_STATE_SCHEDULED  uint32 = 0x00000001
	SMPP_50_BCAST_STATE_COMPLETE   uint32 = 0x00000002
	SMPP_50_BCAST_STATE_INCOMPLETE uint32 = 0x00000003
	SMPP_50_BCAST_STATE_CANCELLED  uint32 = 0x00000004

	// Broadcast Area Formats
	SMPP_50_BCAST_AREA_FORMAT_NAME  uint8 = 0x00
	SMPP_50_BCAST_AREA_FORMAT_ALIAS uint8 = 0x01
	SMPP_50_BCAST_AREA_FORMAT_MSC   uint8 = 0x02
	SMPP_50_BCAST_AREA_FORMAT_LAC   uint8 = 0x03
	SMPP_50_BCAST_AREA_FORMAT_CELL  uint8 = 0x04
	SMPP_50_BCAST_AREA_FORMAT_HLR   uint8 = 0x05
	SMPP_50_BCAST_AREA_FORMAT_ALL   uint8 = 0x06

	// Additional v5.0 Error Codes
	ESME_RBCAST_QUERY_FAIL          uint32 = 0x00000425 // Broadcast query operation failed
	ESME_RBCAST_CANCEL_FAIL         uint32 = 0x00000426 // Broadcast cancel operation failed
	ESME_RBCAST_REPLACE_FAIL        uint32 = 0x00000427 // Broadcast replace operation failed
	ESME_RBCAST_MSG_NOT_FOUND       uint32 = 0x00000428 // Broadcast message not found
	ESME_RBCAST_AREA_FORMAT_INVALID uint32 = 0x00000429 // Broadcast area format invalid
	ESME_RBCAST_AREA_NOT_SUPPORTED  uint32 = 0x0000042A // Broadcast area not supported
	ESME_RBCAST_PRIORITY_INVALID    uint32 = 0x0000042B // Broadcast priority invalid
	ESME_RBCAST_CHANNEL_INVALID     uint32 = 0x0000042C // Broadcast channel invalid
	ESME_RBCAST_CHANNEL_NOT_AVAIL   uint32 = 0x0000042D // Broadcast channel not available
	ESME_RBCAST_EMERGENCY_NOT_SUPP  uint32 = 0x0000042E // Emergency broadcast not supported
	ESME_RBCAST_RESPONSE_TIMEOUT    uint32 = 0x0000042F // Broadcast response timeout

	// v5.0 Message Payload Types
	SMPP_50_PAYLOAD_TYPE_DEFAULT uint8 = 0x00
	SMPP_50_PAYLOAD_TYPE_WCMP    uint8 = 0x01
	SMPP_50_PAYLOAD_TYPE_WAP     uint8 = 0x02
	SMPP_50_PAYLOAD_TYPE_JAVA    uint8 = 0x03
	SMPP_50_PAYLOAD_TYPE_BINARY  uint8 = 0x04
)

// Address TON (Type of Number) Constants
const (
	// Source TON
	SRC_TON_UNKNOWN           uint8 = 0x00 // Unknown
	SRC_TON_INTERNATIONAL     uint8 = 0x01 // International
	SRC_TON_NATIONAL          uint8 = 0x02 // National
	SRC_TON_NETWORK_SPECIFIC  uint8 = 0x03 // Network Specific
	SRC_TON_SUBSCRIBER_NUMBER uint8 = 0x04 // Subscriber Number
	SRC_TON_ALPHANUMERIC      uint8 = 0x05 // Alphanumeric
	SRC_TON_ABBREVIATED       uint8 = 0x06 // Abbreviated

	// Destination TON
	DST_TON_UNKNOWN           uint8 = 0x00 // Unknown
	DST_TON_INTERNATIONAL     uint8 = 0x01 // International
	DST_TON_NATIONAL          uint8 = 0x02 // National
	DST_TON_NETWORK_SPECIFIC  uint8 = 0x03 // Network Specific
	DST_TON_SUBSCRIBER_NUMBER uint8 = 0x04 // Subscriber Number
	DST_TON_ALPHANUMERIC      uint8 = 0x05 // Alphanumeric
	DST_TON_ABBREVIATED       uint8 = 0x06 // Abbreviated
)

// Address NPI (Numbering Plan Indicator) Constants
const (
	// Source NPI
	SRC_NPI_UNKNOWN     uint8 = 0x00 // Unknown
	SRC_NPI_ISDN        uint8 = 0x01 // ISDN (E163/E164)
	SRC_NPI_DATA        uint8 = 0x03 // Data (X.121)
	SRC_NPI_TELEX       uint8 = 0x04 // Telex (F.69)
	SRC_NPI_LAND_MOBILE uint8 = 0x06 // Land Mobile (E.212)
	SRC_NPI_NATIONAL    uint8 = 0x08 // National
	SRC_NPI_PRIVATE     uint8 = 0x09 // Private
	SRC_NPI_ERMES       uint8 = 0x0A // ERMES
	SRC_NPI_INTERNET    uint8 = 0x0E // Internet (IP)
	SRC_NPI_WAP_CLIENT  uint8 = 0x12 // WAP Client ID

	// Destination NPI
	DST_NPI_UNKNOWN     uint8 = 0x00 // Unknown
	DST_NPI_ISDN        uint8 = 0x01 // ISDN (E163/E164)
	DST_NPI_DATA        uint8 = 0x03 // Data (X.121)
	DST_NPI_TELEX       uint8 = 0x04 // Telex (F.69)
	DST_NPI_LAND_MOBILE uint8 = 0x06 // Land Mobile (E.212)
	DST_NPI_NATIONAL    uint8 = 0x08 // National
	DST_NPI_PRIVATE     uint8 = 0x09 // Private
	DST_NPI_ERMES       uint8 = 0x0A // ERMES
	DST_NPI_INTERNET    uint8 = 0x0E // Internet (IP)
	DST_NPI_WAP_CLIENT  uint8 = 0x12 // WAP Client ID
)
