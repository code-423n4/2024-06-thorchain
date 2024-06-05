# Pool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** |  | 
**ShortCode** | Pointer to **string** |  | [optional] 
**Status** | **string** |  | 
**Decimals** | Pointer to **int64** |  | [optional] 
**PendingInboundAsset** | **string** |  | 
**PendingInboundRune** | **string** |  | 
**BalanceAsset** | **string** |  | 
**BalanceRune** | **string** |  | 
**AssetTorPrice** | **string** | the USD (TOR) price of the asset in 1e8 | 
**PoolUnits** | **string** | the total pool units, this is the sum of LP and synth units | 
**LPUnits** | **string** | the total pool liquidity provider units | 
**SynthUnits** | **string** | the total synth units in the pool | 
**SynthSupply** | **string** | the total supply of synths for the asset | 
**SaversDepth** | **string** | the balance of L1 asset deposited into the Savers Vault | 
**SaversUnits** | **string** | the number of units owned by Savers | 
**SynthMintPaused** | **bool** | whether additional synths cannot be minted | 
**SynthSupplyRemaining** | **string** | the amount of synth supply remaining before the current max supply is reached | 
**LoanCollateral** | **string** | the amount of collateral collects for loans | 
**LoanCollateralRemaining** | **string** | the amount of remaining collateral collects for loans | 
**LoanCr** | **string** | the current loan collateralization ratio | 
**DerivedDepthBps** | **string** | the depth of the derived virtual pool relative to L1 pool (in basis points) | 

## Methods

### NewPool

`func NewPool(asset string, status string, pendingInboundAsset string, pendingInboundRune string, balanceAsset string, balanceRune string, assetTorPrice string, poolUnits string, lPUnits string, synthUnits string, synthSupply string, saversDepth string, saversUnits string, synthMintPaused bool, synthSupplyRemaining string, loanCollateral string, loanCollateralRemaining string, loanCr string, derivedDepthBps string, ) *Pool`

NewPool instantiates a new Pool object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPoolWithDefaults

`func NewPoolWithDefaults() *Pool`

NewPoolWithDefaults instantiates a new Pool object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *Pool) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *Pool) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *Pool) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetShortCode

`func (o *Pool) GetShortCode() string`

GetShortCode returns the ShortCode field if non-nil, zero value otherwise.

### GetShortCodeOk

`func (o *Pool) GetShortCodeOk() (*string, bool)`

GetShortCodeOk returns a tuple with the ShortCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShortCode

`func (o *Pool) SetShortCode(v string)`

SetShortCode sets ShortCode field to given value.

### HasShortCode

`func (o *Pool) HasShortCode() bool`

HasShortCode returns a boolean if a field has been set.

### GetStatus

`func (o *Pool) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Pool) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Pool) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetDecimals

`func (o *Pool) GetDecimals() int64`

GetDecimals returns the Decimals field if non-nil, zero value otherwise.

### GetDecimalsOk

`func (o *Pool) GetDecimalsOk() (*int64, bool)`

GetDecimalsOk returns a tuple with the Decimals field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDecimals

`func (o *Pool) SetDecimals(v int64)`

SetDecimals sets Decimals field to given value.

### HasDecimals

`func (o *Pool) HasDecimals() bool`

HasDecimals returns a boolean if a field has been set.

### GetPendingInboundAsset

`func (o *Pool) GetPendingInboundAsset() string`

GetPendingInboundAsset returns the PendingInboundAsset field if non-nil, zero value otherwise.

### GetPendingInboundAssetOk

`func (o *Pool) GetPendingInboundAssetOk() (*string, bool)`

GetPendingInboundAssetOk returns a tuple with the PendingInboundAsset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPendingInboundAsset

`func (o *Pool) SetPendingInboundAsset(v string)`

SetPendingInboundAsset sets PendingInboundAsset field to given value.


### GetPendingInboundRune

`func (o *Pool) GetPendingInboundRune() string`

GetPendingInboundRune returns the PendingInboundRune field if non-nil, zero value otherwise.

### GetPendingInboundRuneOk

`func (o *Pool) GetPendingInboundRuneOk() (*string, bool)`

GetPendingInboundRuneOk returns a tuple with the PendingInboundRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPendingInboundRune

`func (o *Pool) SetPendingInboundRune(v string)`

SetPendingInboundRune sets PendingInboundRune field to given value.


### GetBalanceAsset

`func (o *Pool) GetBalanceAsset() string`

GetBalanceAsset returns the BalanceAsset field if non-nil, zero value otherwise.

### GetBalanceAssetOk

`func (o *Pool) GetBalanceAssetOk() (*string, bool)`

GetBalanceAssetOk returns a tuple with the BalanceAsset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalanceAsset

`func (o *Pool) SetBalanceAsset(v string)`

SetBalanceAsset sets BalanceAsset field to given value.


### GetBalanceRune

`func (o *Pool) GetBalanceRune() string`

GetBalanceRune returns the BalanceRune field if non-nil, zero value otherwise.

### GetBalanceRuneOk

`func (o *Pool) GetBalanceRuneOk() (*string, bool)`

GetBalanceRuneOk returns a tuple with the BalanceRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalanceRune

`func (o *Pool) SetBalanceRune(v string)`

SetBalanceRune sets BalanceRune field to given value.


### GetAssetTorPrice

`func (o *Pool) GetAssetTorPrice() string`

GetAssetTorPrice returns the AssetTorPrice field if non-nil, zero value otherwise.

### GetAssetTorPriceOk

`func (o *Pool) GetAssetTorPriceOk() (*string, bool)`

GetAssetTorPriceOk returns a tuple with the AssetTorPrice field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssetTorPrice

`func (o *Pool) SetAssetTorPrice(v string)`

SetAssetTorPrice sets AssetTorPrice field to given value.


### GetPoolUnits

`func (o *Pool) GetPoolUnits() string`

GetPoolUnits returns the PoolUnits field if non-nil, zero value otherwise.

### GetPoolUnitsOk

`func (o *Pool) GetPoolUnitsOk() (*string, bool)`

GetPoolUnitsOk returns a tuple with the PoolUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolUnits

`func (o *Pool) SetPoolUnits(v string)`

SetPoolUnits sets PoolUnits field to given value.


### GetLPUnits

`func (o *Pool) GetLPUnits() string`

GetLPUnits returns the LPUnits field if non-nil, zero value otherwise.

### GetLPUnitsOk

`func (o *Pool) GetLPUnitsOk() (*string, bool)`

GetLPUnitsOk returns a tuple with the LPUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLPUnits

`func (o *Pool) SetLPUnits(v string)`

SetLPUnits sets LPUnits field to given value.


### GetSynthUnits

`func (o *Pool) GetSynthUnits() string`

GetSynthUnits returns the SynthUnits field if non-nil, zero value otherwise.

### GetSynthUnitsOk

`func (o *Pool) GetSynthUnitsOk() (*string, bool)`

GetSynthUnitsOk returns a tuple with the SynthUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSynthUnits

`func (o *Pool) SetSynthUnits(v string)`

SetSynthUnits sets SynthUnits field to given value.


### GetSynthSupply

`func (o *Pool) GetSynthSupply() string`

GetSynthSupply returns the SynthSupply field if non-nil, zero value otherwise.

### GetSynthSupplyOk

`func (o *Pool) GetSynthSupplyOk() (*string, bool)`

GetSynthSupplyOk returns a tuple with the SynthSupply field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSynthSupply

`func (o *Pool) SetSynthSupply(v string)`

SetSynthSupply sets SynthSupply field to given value.


### GetSaversDepth

`func (o *Pool) GetSaversDepth() string`

GetSaversDepth returns the SaversDepth field if non-nil, zero value otherwise.

### GetSaversDepthOk

`func (o *Pool) GetSaversDepthOk() (*string, bool)`

GetSaversDepthOk returns a tuple with the SaversDepth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSaversDepth

`func (o *Pool) SetSaversDepth(v string)`

SetSaversDepth sets SaversDepth field to given value.


### GetSaversUnits

`func (o *Pool) GetSaversUnits() string`

GetSaversUnits returns the SaversUnits field if non-nil, zero value otherwise.

### GetSaversUnitsOk

`func (o *Pool) GetSaversUnitsOk() (*string, bool)`

GetSaversUnitsOk returns a tuple with the SaversUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSaversUnits

`func (o *Pool) SetSaversUnits(v string)`

SetSaversUnits sets SaversUnits field to given value.


### GetSynthMintPaused

`func (o *Pool) GetSynthMintPaused() bool`

GetSynthMintPaused returns the SynthMintPaused field if non-nil, zero value otherwise.

### GetSynthMintPausedOk

`func (o *Pool) GetSynthMintPausedOk() (*bool, bool)`

GetSynthMintPausedOk returns a tuple with the SynthMintPaused field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSynthMintPaused

`func (o *Pool) SetSynthMintPaused(v bool)`

SetSynthMintPaused sets SynthMintPaused field to given value.


### GetSynthSupplyRemaining

`func (o *Pool) GetSynthSupplyRemaining() string`

GetSynthSupplyRemaining returns the SynthSupplyRemaining field if non-nil, zero value otherwise.

### GetSynthSupplyRemainingOk

`func (o *Pool) GetSynthSupplyRemainingOk() (*string, bool)`

GetSynthSupplyRemainingOk returns a tuple with the SynthSupplyRemaining field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSynthSupplyRemaining

`func (o *Pool) SetSynthSupplyRemaining(v string)`

SetSynthSupplyRemaining sets SynthSupplyRemaining field to given value.


### GetLoanCollateral

`func (o *Pool) GetLoanCollateral() string`

GetLoanCollateral returns the LoanCollateral field if non-nil, zero value otherwise.

### GetLoanCollateralOk

`func (o *Pool) GetLoanCollateralOk() (*string, bool)`

GetLoanCollateralOk returns a tuple with the LoanCollateral field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoanCollateral

`func (o *Pool) SetLoanCollateral(v string)`

SetLoanCollateral sets LoanCollateral field to given value.


### GetLoanCollateralRemaining

`func (o *Pool) GetLoanCollateralRemaining() string`

GetLoanCollateralRemaining returns the LoanCollateralRemaining field if non-nil, zero value otherwise.

### GetLoanCollateralRemainingOk

`func (o *Pool) GetLoanCollateralRemainingOk() (*string, bool)`

GetLoanCollateralRemainingOk returns a tuple with the LoanCollateralRemaining field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoanCollateralRemaining

`func (o *Pool) SetLoanCollateralRemaining(v string)`

SetLoanCollateralRemaining sets LoanCollateralRemaining field to given value.


### GetLoanCr

`func (o *Pool) GetLoanCr() string`

GetLoanCr returns the LoanCr field if non-nil, zero value otherwise.

### GetLoanCrOk

`func (o *Pool) GetLoanCrOk() (*string, bool)`

GetLoanCrOk returns a tuple with the LoanCr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoanCr

`func (o *Pool) SetLoanCr(v string)`

SetLoanCr sets LoanCr field to given value.


### GetDerivedDepthBps

`func (o *Pool) GetDerivedDepthBps() string`

GetDerivedDepthBps returns the DerivedDepthBps field if non-nil, zero value otherwise.

### GetDerivedDepthBpsOk

`func (o *Pool) GetDerivedDepthBpsOk() (*string, bool)`

GetDerivedDepthBpsOk returns a tuple with the DerivedDepthBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDerivedDepthBps

`func (o *Pool) SetDerivedDepthBps(v string)`

SetDerivedDepthBps sets DerivedDepthBps field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


