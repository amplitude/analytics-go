package amplitude

import "time"

type PluginType int

type ServerZone string

type IdentityOp string

const (
	SdkLibrary = "amplitude-go"
	SdkVersion = "0.0.0"

	ServerZoneEU = "EU"
	ServerZoneUS = "US"
	Batch        = "batch"
	HTTPV2       = "v2"

	LoggerName = "amplitude"

	IdentifyEventEventType      = "$identify"
	GroupIdentifyEventEventType = "$groupidentify"
	RevenueEventEventType       = "revenue_amount"

	IdentityOpAdd        = "$add"
	IdentityOpAppend     = "$append"
	IdentityOpClearAll   = "$clearAll"
	IdentityOpPrepend    = "$prepend"
	IdentityOpSet        = "$set"
	IdentityOpSetOnce    = "$setOnce"
	IdentityOpUnset      = "$unset"
	IdentityOpPreInsert  = "$preInsert"
	IdentityOpPostInsert = "$postInsert"
	IdentityOpRemove     = "$remove"
	UnsetValue           = "-"

	RevenueProductID  = "$productId"
	RevenueQuantity   = "$quantity"
	RevenuePrice      = "$price"
	RevenueType       = "$revenueType"
	RevenueReceipt    = "$receipt"
	RevenueReceiptSig = "$receiptSig"
	DefaultRevenue    = "$revenue"

	MaxPropertyKeys = 1024
	MaxStringLength = 1024

	DefaultFlushQueueSize  = 200
	DefaultFlushInterval   = time.Second * 10
	DefaultFlushMaxRetries = 12
	DefaultMinIDLength     = 5
	ConnectionTimeout      = 10.0
	MaxBufferCapacity      = 20000

	BEFORE PluginType = iota
	ENRICHMENT
	DESTINATION
	OBSERVE
)
