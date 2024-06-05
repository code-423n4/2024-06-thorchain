# QuoteFees

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** | the target asset used for all fees | 
**Affiliate** | Pointer to **string** | affiliate fee in the target asset | [optional] 
**Outbound** | Pointer to **string** | outbound fee in the target asset | [optional] 
**Liquidity** | **string** | liquidity fees paid to pools in the target asset | 
**Total** | **string** | total fees in the target asset | 
**SlippageBps** | **int64** | the swap slippage in basis points | 
**TotalBps** | **int64** | total basis points in fees relative to amount out | 

## Methods

### NewQuoteFees

`func NewQuoteFees(asset string, liquidity string, total string, slippageBps int64, totalBps int64, ) *QuoteFees`

NewQuoteFees instantiates a new QuoteFees object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQuoteFeesWithDefaults

`func NewQuoteFeesWithDefaults() *QuoteFees`

NewQuoteFeesWithDefaults instantiates a new QuoteFees object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *QuoteFees) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *QuoteFees) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *QuoteFees) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetAffiliate

`func (o *QuoteFees) GetAffiliate() string`

GetAffiliate returns the Affiliate field if non-nil, zero value otherwise.

### GetAffiliateOk

`func (o *QuoteFees) GetAffiliateOk() (*string, bool)`

GetAffiliateOk returns a tuple with the Affiliate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAffiliate

`func (o *QuoteFees) SetAffiliate(v string)`

SetAffiliate sets Affiliate field to given value.

### HasAffiliate

`func (o *QuoteFees) HasAffiliate() bool`

HasAffiliate returns a boolean if a field has been set.

### GetOutbound

`func (o *QuoteFees) GetOutbound() string`

GetOutbound returns the Outbound field if non-nil, zero value otherwise.

### GetOutboundOk

`func (o *QuoteFees) GetOutboundOk() (*string, bool)`

GetOutboundOk returns a tuple with the Outbound field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutbound

`func (o *QuoteFees) SetOutbound(v string)`

SetOutbound sets Outbound field to given value.

### HasOutbound

`func (o *QuoteFees) HasOutbound() bool`

HasOutbound returns a boolean if a field has been set.

### GetLiquidity

`func (o *QuoteFees) GetLiquidity() string`

GetLiquidity returns the Liquidity field if non-nil, zero value otherwise.

### GetLiquidityOk

`func (o *QuoteFees) GetLiquidityOk() (*string, bool)`

GetLiquidityOk returns a tuple with the Liquidity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLiquidity

`func (o *QuoteFees) SetLiquidity(v string)`

SetLiquidity sets Liquidity field to given value.


### GetTotal

`func (o *QuoteFees) GetTotal() string`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *QuoteFees) GetTotalOk() (*string, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *QuoteFees) SetTotal(v string)`

SetTotal sets Total field to given value.


### GetSlippageBps

`func (o *QuoteFees) GetSlippageBps() int64`

GetSlippageBps returns the SlippageBps field if non-nil, zero value otherwise.

### GetSlippageBpsOk

`func (o *QuoteFees) GetSlippageBpsOk() (*int64, bool)`

GetSlippageBpsOk returns a tuple with the SlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlippageBps

`func (o *QuoteFees) SetSlippageBps(v int64)`

SetSlippageBps sets SlippageBps field to given value.


### GetTotalBps

`func (o *QuoteFees) GetTotalBps() int64`

GetTotalBps returns the TotalBps field if non-nil, zero value otherwise.

### GetTotalBpsOk

`func (o *QuoteFees) GetTotalBpsOk() (*int64, bool)`

GetTotalBpsOk returns a tuple with the TotalBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalBps

`func (o *QuoteFees) SetTotalBps(v int64)`

SetTotalBps sets TotalBps field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


