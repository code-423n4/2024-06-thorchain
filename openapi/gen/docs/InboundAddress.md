# InboundAddress

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Chain** | Pointer to **string** |  | [optional] 
**PubKey** | Pointer to **string** |  | [optional] 
**Address** | Pointer to **string** |  | [optional] 
**Router** | Pointer to **string** |  | [optional] 
**Halted** | **bool** | Returns true if trading is unavailable for this chain, either because trading is halted globally or specifically for this chain | 
**GlobalTradingPaused** | Pointer to **bool** | Returns true if trading is paused globally | [optional] 
**ChainTradingPaused** | Pointer to **bool** | Returns true if trading is paused for this chain | [optional] 
**ChainLpActionsPaused** | Pointer to **bool** | Returns true if LP actions are paused for this chain | [optional] 
**GasRate** | Pointer to **string** | The minimum fee rate used by vaults to send outbound TXs. The actual fee rate may be higher. For EVM chains this is returned in gwei (1e9). | [optional] 
**GasRateUnits** | Pointer to **string** | Units of the gas_rate. | [optional] 
**OutboundTxSize** | Pointer to **string** | Avg size of outbound TXs on each chain. For UTXO chains it may be larger than average, as it takes into account vault consolidation txs, which can have many vouts | [optional] 
**OutboundFee** | Pointer to **string** | The total outbound fee charged to the user for outbound txs in the gas asset of the chain. | [optional] 
**DustThreshold** | Pointer to **string** | Defines the minimum transaction size for the chain in base units (sats, wei, uatom). Transactions with asset amounts lower than the dust_threshold are ignored. | [optional] 

## Methods

### NewInboundAddress

`func NewInboundAddress(halted bool, ) *InboundAddress`

NewInboundAddress instantiates a new InboundAddress object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInboundAddressWithDefaults

`func NewInboundAddressWithDefaults() *InboundAddress`

NewInboundAddressWithDefaults instantiates a new InboundAddress object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChain

`func (o *InboundAddress) GetChain() string`

GetChain returns the Chain field if non-nil, zero value otherwise.

### GetChainOk

`func (o *InboundAddress) GetChainOk() (*string, bool)`

GetChainOk returns a tuple with the Chain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChain

`func (o *InboundAddress) SetChain(v string)`

SetChain sets Chain field to given value.

### HasChain

`func (o *InboundAddress) HasChain() bool`

HasChain returns a boolean if a field has been set.

### GetPubKey

`func (o *InboundAddress) GetPubKey() string`

GetPubKey returns the PubKey field if non-nil, zero value otherwise.

### GetPubKeyOk

`func (o *InboundAddress) GetPubKeyOk() (*string, bool)`

GetPubKeyOk returns a tuple with the PubKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPubKey

`func (o *InboundAddress) SetPubKey(v string)`

SetPubKey sets PubKey field to given value.

### HasPubKey

`func (o *InboundAddress) HasPubKey() bool`

HasPubKey returns a boolean if a field has been set.

### GetAddress

`func (o *InboundAddress) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *InboundAddress) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *InboundAddress) SetAddress(v string)`

SetAddress sets Address field to given value.

### HasAddress

`func (o *InboundAddress) HasAddress() bool`

HasAddress returns a boolean if a field has been set.

### GetRouter

`func (o *InboundAddress) GetRouter() string`

GetRouter returns the Router field if non-nil, zero value otherwise.

### GetRouterOk

`func (o *InboundAddress) GetRouterOk() (*string, bool)`

GetRouterOk returns a tuple with the Router field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouter

`func (o *InboundAddress) SetRouter(v string)`

SetRouter sets Router field to given value.

### HasRouter

`func (o *InboundAddress) HasRouter() bool`

HasRouter returns a boolean if a field has been set.

### GetHalted

`func (o *InboundAddress) GetHalted() bool`

GetHalted returns the Halted field if non-nil, zero value otherwise.

### GetHaltedOk

`func (o *InboundAddress) GetHaltedOk() (*bool, bool)`

GetHaltedOk returns a tuple with the Halted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHalted

`func (o *InboundAddress) SetHalted(v bool)`

SetHalted sets Halted field to given value.


### GetGlobalTradingPaused

`func (o *InboundAddress) GetGlobalTradingPaused() bool`

GetGlobalTradingPaused returns the GlobalTradingPaused field if non-nil, zero value otherwise.

### GetGlobalTradingPausedOk

`func (o *InboundAddress) GetGlobalTradingPausedOk() (*bool, bool)`

GetGlobalTradingPausedOk returns a tuple with the GlobalTradingPaused field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalTradingPaused

`func (o *InboundAddress) SetGlobalTradingPaused(v bool)`

SetGlobalTradingPaused sets GlobalTradingPaused field to given value.

### HasGlobalTradingPaused

`func (o *InboundAddress) HasGlobalTradingPaused() bool`

HasGlobalTradingPaused returns a boolean if a field has been set.

### GetChainTradingPaused

`func (o *InboundAddress) GetChainTradingPaused() bool`

GetChainTradingPaused returns the ChainTradingPaused field if non-nil, zero value otherwise.

### GetChainTradingPausedOk

`func (o *InboundAddress) GetChainTradingPausedOk() (*bool, bool)`

GetChainTradingPausedOk returns a tuple with the ChainTradingPaused field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChainTradingPaused

`func (o *InboundAddress) SetChainTradingPaused(v bool)`

SetChainTradingPaused sets ChainTradingPaused field to given value.

### HasChainTradingPaused

`func (o *InboundAddress) HasChainTradingPaused() bool`

HasChainTradingPaused returns a boolean if a field has been set.

### GetChainLpActionsPaused

`func (o *InboundAddress) GetChainLpActionsPaused() bool`

GetChainLpActionsPaused returns the ChainLpActionsPaused field if non-nil, zero value otherwise.

### GetChainLpActionsPausedOk

`func (o *InboundAddress) GetChainLpActionsPausedOk() (*bool, bool)`

GetChainLpActionsPausedOk returns a tuple with the ChainLpActionsPaused field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChainLpActionsPaused

`func (o *InboundAddress) SetChainLpActionsPaused(v bool)`

SetChainLpActionsPaused sets ChainLpActionsPaused field to given value.

### HasChainLpActionsPaused

`func (o *InboundAddress) HasChainLpActionsPaused() bool`

HasChainLpActionsPaused returns a boolean if a field has been set.

### GetGasRate

`func (o *InboundAddress) GetGasRate() string`

GetGasRate returns the GasRate field if non-nil, zero value otherwise.

### GetGasRateOk

`func (o *InboundAddress) GetGasRateOk() (*string, bool)`

GetGasRateOk returns a tuple with the GasRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasRate

`func (o *InboundAddress) SetGasRate(v string)`

SetGasRate sets GasRate field to given value.

### HasGasRate

`func (o *InboundAddress) HasGasRate() bool`

HasGasRate returns a boolean if a field has been set.

### GetGasRateUnits

`func (o *InboundAddress) GetGasRateUnits() string`

GetGasRateUnits returns the GasRateUnits field if non-nil, zero value otherwise.

### GetGasRateUnitsOk

`func (o *InboundAddress) GetGasRateUnitsOk() (*string, bool)`

GetGasRateUnitsOk returns a tuple with the GasRateUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasRateUnits

`func (o *InboundAddress) SetGasRateUnits(v string)`

SetGasRateUnits sets GasRateUnits field to given value.

### HasGasRateUnits

`func (o *InboundAddress) HasGasRateUnits() bool`

HasGasRateUnits returns a boolean if a field has been set.

### GetOutboundTxSize

`func (o *InboundAddress) GetOutboundTxSize() string`

GetOutboundTxSize returns the OutboundTxSize field if non-nil, zero value otherwise.

### GetOutboundTxSizeOk

`func (o *InboundAddress) GetOutboundTxSizeOk() (*string, bool)`

GetOutboundTxSizeOk returns a tuple with the OutboundTxSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundTxSize

`func (o *InboundAddress) SetOutboundTxSize(v string)`

SetOutboundTxSize sets OutboundTxSize field to given value.

### HasOutboundTxSize

`func (o *InboundAddress) HasOutboundTxSize() bool`

HasOutboundTxSize returns a boolean if a field has been set.

### GetOutboundFee

`func (o *InboundAddress) GetOutboundFee() string`

GetOutboundFee returns the OutboundFee field if non-nil, zero value otherwise.

### GetOutboundFeeOk

`func (o *InboundAddress) GetOutboundFeeOk() (*string, bool)`

GetOutboundFeeOk returns a tuple with the OutboundFee field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundFee

`func (o *InboundAddress) SetOutboundFee(v string)`

SetOutboundFee sets OutboundFee field to given value.

### HasOutboundFee

`func (o *InboundAddress) HasOutboundFee() bool`

HasOutboundFee returns a boolean if a field has been set.

### GetDustThreshold

`func (o *InboundAddress) GetDustThreshold() string`

GetDustThreshold returns the DustThreshold field if non-nil, zero value otherwise.

### GetDustThresholdOk

`func (o *InboundAddress) GetDustThresholdOk() (*string, bool)`

GetDustThresholdOk returns a tuple with the DustThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDustThreshold

`func (o *InboundAddress) SetDustThreshold(v string)`

SetDustThreshold sets DustThreshold field to given value.

### HasDustThreshold

`func (o *InboundAddress) HasDustThreshold() bool`

HasDustThreshold returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


