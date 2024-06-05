# TradeAccountsUnitResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** | trade account asset with \&quot;~\&quot; separator | 
**Unit** | Pointer to **string** | total units of trade asset | [optional] 
**Depth** | **string** | total depth of trade asset | 

## Methods

### NewTradeAccountsUnitResponse

`func NewTradeAccountsUnitResponse(asset string, depth string, ) *TradeAccountsUnitResponse`

NewTradeAccountsUnitResponse instantiates a new TradeAccountsUnitResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTradeAccountsUnitResponseWithDefaults

`func NewTradeAccountsUnitResponseWithDefaults() *TradeAccountsUnitResponse`

NewTradeAccountsUnitResponseWithDefaults instantiates a new TradeAccountsUnitResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *TradeAccountsUnitResponse) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *TradeAccountsUnitResponse) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *TradeAccountsUnitResponse) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetUnit

`func (o *TradeAccountsUnitResponse) GetUnit() string`

GetUnit returns the Unit field if non-nil, zero value otherwise.

### GetUnitOk

`func (o *TradeAccountsUnitResponse) GetUnitOk() (*string, bool)`

GetUnitOk returns a tuple with the Unit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnit

`func (o *TradeAccountsUnitResponse) SetUnit(v string)`

SetUnit sets Unit field to given value.

### HasUnit

`func (o *TradeAccountsUnitResponse) HasUnit() bool`

HasUnit returns a boolean if a field has been set.

### GetDepth

`func (o *TradeAccountsUnitResponse) GetDepth() string`

GetDepth returns the Depth field if non-nil, zero value otherwise.

### GetDepthOk

`func (o *TradeAccountsUnitResponse) GetDepthOk() (*string, bool)`

GetDepthOk returns a tuple with the Depth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDepth

`func (o *TradeAccountsUnitResponse) SetDepth(v string)`

SetDepth sets Depth field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


