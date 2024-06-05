# NetworkResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BondRewardRune** | **string** | total amount of RUNE awarded to node operators | 
**BurnedBep2Rune** | **string** | total of burned BEP2 RUNE | 
**BurnedErc20Rune** | **string** | total of burned ERC20 RUNE | 
**TotalBondUnits** | **string** | total bonded RUNE | 
**EffectiveSecurityBond** | **string** | effective security bond used to determine maximum pooled RUNE | 
**TotalReserve** | **string** | total reserve RUNE | 
**VaultsMigrating** | **bool** | Returns true if there exist RetiringVaults which have not finished migrating funds to new ActiveVaults | 
**GasSpentRune** | **string** | Sum of the gas the network has spent to send outbounds | 
**GasWithheldRune** | **string** | Sum of the gas withheld from users to cover outbound gas | 
**OutboundFeeMultiplier** | Pointer to **string** | Current outbound fee multiplier, in basis points | [optional] 
**NativeOutboundFeeRune** | **string** | the outbound transaction fee in rune, converted from the NativeOutboundFeeUSD mimir (after USD fees are enabled) | 
**NativeTxFeeRune** | **string** | the native transaction fee in rune, converted from the NativeTransactionFeeUSD mimir (after USD fees are enabled) | 
**TnsRegisterFeeRune** | **string** | the thorname register fee in rune, converted from the TNSRegisterFeeUSD mimir (after USD fees are enabled) | 
**TnsFeePerBlockRune** | **string** | the thorname fee per block in rune, converted from the TNSFeePerBlockUSD mimir (after USD fees are enabled) | 
**RunePriceInTor** | **string** | the rune price in tor | 
**TorPriceInRune** | **string** | the tor price in rune | 

## Methods

### NewNetworkResponse

`func NewNetworkResponse(bondRewardRune string, burnedBep2Rune string, burnedErc20Rune string, totalBondUnits string, effectiveSecurityBond string, totalReserve string, vaultsMigrating bool, gasSpentRune string, gasWithheldRune string, nativeOutboundFeeRune string, nativeTxFeeRune string, tnsRegisterFeeRune string, tnsFeePerBlockRune string, runePriceInTor string, torPriceInRune string, ) *NetworkResponse`

NewNetworkResponse instantiates a new NetworkResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNetworkResponseWithDefaults

`func NewNetworkResponseWithDefaults() *NetworkResponse`

NewNetworkResponseWithDefaults instantiates a new NetworkResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBondRewardRune

`func (o *NetworkResponse) GetBondRewardRune() string`

GetBondRewardRune returns the BondRewardRune field if non-nil, zero value otherwise.

### GetBondRewardRuneOk

`func (o *NetworkResponse) GetBondRewardRuneOk() (*string, bool)`

GetBondRewardRuneOk returns a tuple with the BondRewardRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBondRewardRune

`func (o *NetworkResponse) SetBondRewardRune(v string)`

SetBondRewardRune sets BondRewardRune field to given value.


### GetBurnedBep2Rune

`func (o *NetworkResponse) GetBurnedBep2Rune() string`

GetBurnedBep2Rune returns the BurnedBep2Rune field if non-nil, zero value otherwise.

### GetBurnedBep2RuneOk

`func (o *NetworkResponse) GetBurnedBep2RuneOk() (*string, bool)`

GetBurnedBep2RuneOk returns a tuple with the BurnedBep2Rune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBurnedBep2Rune

`func (o *NetworkResponse) SetBurnedBep2Rune(v string)`

SetBurnedBep2Rune sets BurnedBep2Rune field to given value.


### GetBurnedErc20Rune

`func (o *NetworkResponse) GetBurnedErc20Rune() string`

GetBurnedErc20Rune returns the BurnedErc20Rune field if non-nil, zero value otherwise.

### GetBurnedErc20RuneOk

`func (o *NetworkResponse) GetBurnedErc20RuneOk() (*string, bool)`

GetBurnedErc20RuneOk returns a tuple with the BurnedErc20Rune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBurnedErc20Rune

`func (o *NetworkResponse) SetBurnedErc20Rune(v string)`

SetBurnedErc20Rune sets BurnedErc20Rune field to given value.


### GetTotalBondUnits

`func (o *NetworkResponse) GetTotalBondUnits() string`

GetTotalBondUnits returns the TotalBondUnits field if non-nil, zero value otherwise.

### GetTotalBondUnitsOk

`func (o *NetworkResponse) GetTotalBondUnitsOk() (*string, bool)`

GetTotalBondUnitsOk returns a tuple with the TotalBondUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalBondUnits

`func (o *NetworkResponse) SetTotalBondUnits(v string)`

SetTotalBondUnits sets TotalBondUnits field to given value.


### GetEffectiveSecurityBond

`func (o *NetworkResponse) GetEffectiveSecurityBond() string`

GetEffectiveSecurityBond returns the EffectiveSecurityBond field if non-nil, zero value otherwise.

### GetEffectiveSecurityBondOk

`func (o *NetworkResponse) GetEffectiveSecurityBondOk() (*string, bool)`

GetEffectiveSecurityBondOk returns a tuple with the EffectiveSecurityBond field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEffectiveSecurityBond

`func (o *NetworkResponse) SetEffectiveSecurityBond(v string)`

SetEffectiveSecurityBond sets EffectiveSecurityBond field to given value.


### GetTotalReserve

`func (o *NetworkResponse) GetTotalReserve() string`

GetTotalReserve returns the TotalReserve field if non-nil, zero value otherwise.

### GetTotalReserveOk

`func (o *NetworkResponse) GetTotalReserveOk() (*string, bool)`

GetTotalReserveOk returns a tuple with the TotalReserve field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalReserve

`func (o *NetworkResponse) SetTotalReserve(v string)`

SetTotalReserve sets TotalReserve field to given value.


### GetVaultsMigrating

`func (o *NetworkResponse) GetVaultsMigrating() bool`

GetVaultsMigrating returns the VaultsMigrating field if non-nil, zero value otherwise.

### GetVaultsMigratingOk

`func (o *NetworkResponse) GetVaultsMigratingOk() (*bool, bool)`

GetVaultsMigratingOk returns a tuple with the VaultsMigrating field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVaultsMigrating

`func (o *NetworkResponse) SetVaultsMigrating(v bool)`

SetVaultsMigrating sets VaultsMigrating field to given value.


### GetGasSpentRune

`func (o *NetworkResponse) GetGasSpentRune() string`

GetGasSpentRune returns the GasSpentRune field if non-nil, zero value otherwise.

### GetGasSpentRuneOk

`func (o *NetworkResponse) GetGasSpentRuneOk() (*string, bool)`

GetGasSpentRuneOk returns a tuple with the GasSpentRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasSpentRune

`func (o *NetworkResponse) SetGasSpentRune(v string)`

SetGasSpentRune sets GasSpentRune field to given value.


### GetGasWithheldRune

`func (o *NetworkResponse) GetGasWithheldRune() string`

GetGasWithheldRune returns the GasWithheldRune field if non-nil, zero value otherwise.

### GetGasWithheldRuneOk

`func (o *NetworkResponse) GetGasWithheldRuneOk() (*string, bool)`

GetGasWithheldRuneOk returns a tuple with the GasWithheldRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasWithheldRune

`func (o *NetworkResponse) SetGasWithheldRune(v string)`

SetGasWithheldRune sets GasWithheldRune field to given value.


### GetOutboundFeeMultiplier

`func (o *NetworkResponse) GetOutboundFeeMultiplier() string`

GetOutboundFeeMultiplier returns the OutboundFeeMultiplier field if non-nil, zero value otherwise.

### GetOutboundFeeMultiplierOk

`func (o *NetworkResponse) GetOutboundFeeMultiplierOk() (*string, bool)`

GetOutboundFeeMultiplierOk returns a tuple with the OutboundFeeMultiplier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundFeeMultiplier

`func (o *NetworkResponse) SetOutboundFeeMultiplier(v string)`

SetOutboundFeeMultiplier sets OutboundFeeMultiplier field to given value.

### HasOutboundFeeMultiplier

`func (o *NetworkResponse) HasOutboundFeeMultiplier() bool`

HasOutboundFeeMultiplier returns a boolean if a field has been set.

### GetNativeOutboundFeeRune

`func (o *NetworkResponse) GetNativeOutboundFeeRune() string`

GetNativeOutboundFeeRune returns the NativeOutboundFeeRune field if non-nil, zero value otherwise.

### GetNativeOutboundFeeRuneOk

`func (o *NetworkResponse) GetNativeOutboundFeeRuneOk() (*string, bool)`

GetNativeOutboundFeeRuneOk returns a tuple with the NativeOutboundFeeRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNativeOutboundFeeRune

`func (o *NetworkResponse) SetNativeOutboundFeeRune(v string)`

SetNativeOutboundFeeRune sets NativeOutboundFeeRune field to given value.


### GetNativeTxFeeRune

`func (o *NetworkResponse) GetNativeTxFeeRune() string`

GetNativeTxFeeRune returns the NativeTxFeeRune field if non-nil, zero value otherwise.

### GetNativeTxFeeRuneOk

`func (o *NetworkResponse) GetNativeTxFeeRuneOk() (*string, bool)`

GetNativeTxFeeRuneOk returns a tuple with the NativeTxFeeRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNativeTxFeeRune

`func (o *NetworkResponse) SetNativeTxFeeRune(v string)`

SetNativeTxFeeRune sets NativeTxFeeRune field to given value.


### GetTnsRegisterFeeRune

`func (o *NetworkResponse) GetTnsRegisterFeeRune() string`

GetTnsRegisterFeeRune returns the TnsRegisterFeeRune field if non-nil, zero value otherwise.

### GetTnsRegisterFeeRuneOk

`func (o *NetworkResponse) GetTnsRegisterFeeRuneOk() (*string, bool)`

GetTnsRegisterFeeRuneOk returns a tuple with the TnsRegisterFeeRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTnsRegisterFeeRune

`func (o *NetworkResponse) SetTnsRegisterFeeRune(v string)`

SetTnsRegisterFeeRune sets TnsRegisterFeeRune field to given value.


### GetTnsFeePerBlockRune

`func (o *NetworkResponse) GetTnsFeePerBlockRune() string`

GetTnsFeePerBlockRune returns the TnsFeePerBlockRune field if non-nil, zero value otherwise.

### GetTnsFeePerBlockRuneOk

`func (o *NetworkResponse) GetTnsFeePerBlockRuneOk() (*string, bool)`

GetTnsFeePerBlockRuneOk returns a tuple with the TnsFeePerBlockRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTnsFeePerBlockRune

`func (o *NetworkResponse) SetTnsFeePerBlockRune(v string)`

SetTnsFeePerBlockRune sets TnsFeePerBlockRune field to given value.


### GetRunePriceInTor

`func (o *NetworkResponse) GetRunePriceInTor() string`

GetRunePriceInTor returns the RunePriceInTor field if non-nil, zero value otherwise.

### GetRunePriceInTorOk

`func (o *NetworkResponse) GetRunePriceInTorOk() (*string, bool)`

GetRunePriceInTorOk returns a tuple with the RunePriceInTor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunePriceInTor

`func (o *NetworkResponse) SetRunePriceInTor(v string)`

SetRunePriceInTor sets RunePriceInTor field to given value.


### GetTorPriceInRune

`func (o *NetworkResponse) GetTorPriceInRune() string`

GetTorPriceInRune returns the TorPriceInRune field if non-nil, zero value otherwise.

### GetTorPriceInRuneOk

`func (o *NetworkResponse) GetTorPriceInRuneOk() (*string, bool)`

GetTorPriceInRuneOk returns a tuple with the TorPriceInRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTorPriceInRune

`func (o *NetworkResponse) SetTorPriceInRune(v string)`

SetTorPriceInRune sets TorPriceInRune field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


