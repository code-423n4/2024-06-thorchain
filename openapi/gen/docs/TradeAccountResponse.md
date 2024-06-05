# TradeAccountResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** | trade account asset with \&quot;~\&quot; separator | 
**Units** | **string** | units of trade asset belonging to this owner | 
**Owner** | **string** | thor address of trade account owner | 
**LastAddHeight** | Pointer to **int64** | last thorchain height trade assets were added to trade account | [optional] 
**LastWithdrawHeight** | Pointer to **int64** | last thorchain height trade assets were withdrawn from trade account | [optional] 

## Methods

### NewTradeAccountResponse

`func NewTradeAccountResponse(asset string, units string, owner string, ) *TradeAccountResponse`

NewTradeAccountResponse instantiates a new TradeAccountResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTradeAccountResponseWithDefaults

`func NewTradeAccountResponseWithDefaults() *TradeAccountResponse`

NewTradeAccountResponseWithDefaults instantiates a new TradeAccountResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *TradeAccountResponse) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *TradeAccountResponse) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *TradeAccountResponse) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetUnits

`func (o *TradeAccountResponse) GetUnits() string`

GetUnits returns the Units field if non-nil, zero value otherwise.

### GetUnitsOk

`func (o *TradeAccountResponse) GetUnitsOk() (*string, bool)`

GetUnitsOk returns a tuple with the Units field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnits

`func (o *TradeAccountResponse) SetUnits(v string)`

SetUnits sets Units field to given value.


### GetOwner

`func (o *TradeAccountResponse) GetOwner() string`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *TradeAccountResponse) GetOwnerOk() (*string, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *TradeAccountResponse) SetOwner(v string)`

SetOwner sets Owner field to given value.


### GetLastAddHeight

`func (o *TradeAccountResponse) GetLastAddHeight() int64`

GetLastAddHeight returns the LastAddHeight field if non-nil, zero value otherwise.

### GetLastAddHeightOk

`func (o *TradeAccountResponse) GetLastAddHeightOk() (*int64, bool)`

GetLastAddHeightOk returns a tuple with the LastAddHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastAddHeight

`func (o *TradeAccountResponse) SetLastAddHeight(v int64)`

SetLastAddHeight sets LastAddHeight field to given value.

### HasLastAddHeight

`func (o *TradeAccountResponse) HasLastAddHeight() bool`

HasLastAddHeight returns a boolean if a field has been set.

### GetLastWithdrawHeight

`func (o *TradeAccountResponse) GetLastWithdrawHeight() int64`

GetLastWithdrawHeight returns the LastWithdrawHeight field if non-nil, zero value otherwise.

### GetLastWithdrawHeightOk

`func (o *TradeAccountResponse) GetLastWithdrawHeightOk() (*int64, bool)`

GetLastWithdrawHeightOk returns a tuple with the LastWithdrawHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastWithdrawHeight

`func (o *TradeAccountResponse) SetLastWithdrawHeight(v int64)`

SetLastWithdrawHeight sets LastWithdrawHeight field to given value.

### HasLastWithdrawHeight

`func (o *TradeAccountResponse) HasLastWithdrawHeight() bool`

HasLastWithdrawHeight returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


