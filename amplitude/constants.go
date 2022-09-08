package amplitude

import (
	"time"
)

const (
	SdkLibrary = "amplitude-go"
	SdkVersion = "0.0.2"

	IdentifyEventType      = "$identify"
	GroupIdentifyEventType = "$groupidentify"
	RevenueEventType       = "revenue_amount"

	LoggerName = "amplitude"

	RevenueProductID  = "$productId"
	RevenueQuantity   = "$quantity"
	RevenuePrice      = "$price"
	RevenueType       = "$revenueType"
	RevenueReceipt    = "$receipt"
	RevenueReceiptSig = "$receiptSig"
	DefaultRevenue    = "$revenue"

	MaxPropertyKeys = 1024
	MaxStringLength = 1024

	MaxBufferCapacity = 20000
)

type (
	IdentityOp string
)

const (
	IdentityOpAdd        IdentityOp = "$add"
	IdentityOpAppend     IdentityOp = "$append"
	IdentityOpClearAll   IdentityOp = "$clearAll"
	IdentityOpPrepend    IdentityOp = "$prepend"
	IdentityOpSet        IdentityOp = "$set"
	IdentityOpSetOnce    IdentityOp = "$setOnce"
	IdentityOpUnset      IdentityOp = "$unset"
	IdentityOpPreInsert  IdentityOp = "$preInsert"
	IdentityOpPostInsert IdentityOp = "$postInsert"
	IdentityOpRemove     IdentityOp = "$remove"
	UnsetValue           string     = "-"
)

type (
	PluginType int
)

const (
	PluginTypeBefore PluginType = iota
	PluginTypeEnrichment
	PluginTypeDestination
	PluginTypeObserve
)

type ServerZone string

const (
	ServerZoneUS ServerZone = "US"
	ServerZoneEU ServerZone = "EU"
)

var ServerURLs = map[ServerZone]string{
	ServerZoneUS: "https://api2.amplitude.com/2/httpapi",
	ServerZoneEU: "https://api.eu.amplitude.com/2/httpapi",
}

var ServerBatchURLs = map[ServerZone]string{
	ServerZoneUS: "https://api2.amplitude.com/batch",
	ServerZoneEU: "https://api.eu.amplitude.com/batch",
}

var DefaultConfig = Config{
	FlushInterval:     time.Second * 10,
	FlushQueueSize:    200,
	FlushMaxRetries:   12,
	ServerZone:        ServerZoneUS,
	ConnectionTimeout: time.Second * 10,
}
