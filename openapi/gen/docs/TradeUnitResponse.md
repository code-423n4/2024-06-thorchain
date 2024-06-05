# TradeUnitResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** | trade account asset with \&quot;~\&quot; separator | 
**Units** | **string** | total units of trade asset | 
**Depth** | **string** | total depth of trade asset | 

## Methods

### NewTradeUnitResponse

`func NewTradeUnitResponse(asset string, units string, depth string, ) *TradeUnitResponse`

NewTradeUnitResponse instantiates a new TradeUnitResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTradeUnitResponseWithDefaults

`func NewTradeUnitResponseWithDefaults() *TradeUnitResponse`

NewTradeUnitResponseWithDefaults instantiates a new TradeUnitResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *TradeUnitResponse) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *TradeUnitResponse) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *TradeUnitResponse) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetUnits

`func (o *TradeUnitResponse) GetUnits() string`

GetUnits returns the Units field if non-nil, zero value otherwise.

### GetUnitsOk

`func (o *TradeUnitResponse) GetUnitsOk() (*string, bool)`

GetUnitsOk returns a tuple with the Units field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnits

`func (o *TradeUnitResponse) SetUnits(v string)`

SetUnits sets Units field to given value.


### GetDepth

`func (o *TradeUnitResponse) GetDepth() string`

GetDepth returns the Depth field if non-nil, zero value otherwise.

### GetDepthOk

`func (o *TradeUnitResponse) GetDepthOk() (*string, bool)`

GetDepthOk returns a tuple with the Depth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDepth

`func (o *TradeUnitResponse) SetDepth(v string)`

SetDepth sets Depth field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


