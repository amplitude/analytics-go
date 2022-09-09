package amplitude

import (
	"github.com/amplitude/analytics-go/amplitude/constants"
	"github.com/amplitude/analytics-go/amplitude/types"
)

type (
	Config     = types.Config
	Plan       = types.Plan
	ServerZone = types.ServerZone

	EventOptions = types.EventOptions
	Event        = types.Event
	IdentityOp   = types.IdentityOp
	Identify     = types.Identify
	Revenue      = types.Revenue

	PluginType                = types.PluginType
	Plugin                    = types.Plugin
	EnrichmentPlugin          = types.EnrichmentPlugin
	DestinationPlugin         = types.DestinationPlugin
	ExtendedDestinationPlugin = types.ExtendedDestinationPlugin
	ExecuteResult             = types.ExecuteResult

	EventStorage = types.EventStorage
)

const (
	ServerZoneUS = types.ServerZoneUS
	ServerZoneEU = types.ServerZoneEU

	PluginTypeBefore      = types.PluginTypeBefore
	PluginTypeEnrichment  = types.PluginTypeEnrichment
	PluginTypeDestination = types.PluginTypeDestination
	PluginTypeObserve     = types.PluginTypeObserve

	IdentifyEventType      = constants.IdentifyEventType
	GroupIdentifyEventType = constants.GroupIdentifyEventType
	RevenueEventType       = constants.RevenueEventType
)

var NewConfig = types.NewConfig
