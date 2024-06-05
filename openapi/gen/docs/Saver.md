# Saver

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** |  | 
**AssetAddress** | **string** |  | 
**LastAddHeight** | Pointer to **int64** |  | [optional] 
**LastWithdrawHeight** | Pointer to **int64** |  | [optional] 
**Units** | **string** |  | 
**AssetDepositValue** | **string** |  | 
**AssetRedeemValue** | **string** |  | 
**GrowthPct** | **string** |  | 

## Methods

### NewSaver

`func NewSaver(asset string, assetAddress string, units string, assetDepositValue string, assetRedeemValue string, growthPct string, ) *Saver`

NewSaver instantiates a new Saver object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSaverWithDefaults

`func NewSaverWithDefaults() *Saver`

NewSaverWithDefaults instantiates a new Saver object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *Saver) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *Saver) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *Saver) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetAssetAddress

`func (o *Saver) GetAssetAddress() string`

GetAssetAddress returns the AssetAddress field if non-nil, zero value otherwise.

### GetAssetAddressOk

`func (o *Saver) GetAssetAddressOk() (*string, bool)`

GetAssetAddressOk returns a tuple with the AssetAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssetAddress

`func (o *Saver) SetAssetAddress(v string)`

SetAssetAddress sets AssetAddress field to given value.


### GetLastAddHeight

`func (o *Saver) GetLastAddHeight() int64`

GetLastAddHeight returns the LastAddHeight field if non-nil, zero value otherwise.

### GetLastAddHeightOk

`func (o *Saver) GetLastAddHeightOk() (*int64, bool)`

GetLastAddHeightOk returns a tuple with the LastAddHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastAddHeight

`func (o *Saver) SetLastAddHeight(v int64)`

SetLastAddHeight sets LastAddHeight field to given value.

### HasLastAddHeight

`func (o *Saver) HasLastAddHeight() bool`

HasLastAddHeight returns a boolean if a field has been set.

### GetLastWithdrawHeight

`func (o *Saver) GetLastWithdrawHeight() int64`

GetLastWithdrawHeight returns the LastWithdrawHeight field if non-nil, zero value otherwise.

### GetLastWithdrawHeightOk

`func (o *Saver) GetLastWithdrawHeightOk() (*int64, bool)`

GetLastWithdrawHeightOk returns a tuple with the LastWithdrawHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastWithdrawHeight

`func (o *Saver) SetLastWithdrawHeight(v int64)`

SetLastWithdrawHeight sets LastWithdrawHeight field to given value.

### HasLastWithdrawHeight

`func (o *Saver) HasLastWithdrawHeight() bool`

HasLastWithdrawHeight returns a boolean if a field has been set.

### GetUnits

`func (o *Saver) GetUnits() string`

GetUnits returns the Units field if non-nil, zero value otherwise.

### GetUnitsOk

`func (o *Saver) GetUnitsOk() (*string, bool)`

GetUnitsOk returns a tuple with the Units field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnits

`func (o *Saver) SetUnits(v string)`

SetUnits sets Units field to given value.


### GetAssetDepositValue

`func (o *Saver) GetAssetDepositValue() string`

GetAssetDepositValue returns the AssetDepositValue field if non-nil, zero value otherwise.

### GetAssetDepositValueOk

`func (o *Saver) GetAssetDepositValueOk() (*string, bool)`

GetAssetDepositValueOk returns a tuple with the AssetDepositValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssetDepositValue

`func (o *Saver) SetAssetDepositValue(v string)`

SetAssetDepositValue sets AssetDepositValue field to given value.


### GetAssetRedeemValue

`func (o *Saver) GetAssetRedeemValue() string`

GetAssetRedeemValue returns the AssetRedeemValue field if non-nil, zero value otherwise.

### GetAssetRedeemValueOk

`func (o *Saver) GetAssetRedeemValueOk() (*string, bool)`

GetAssetRedeemValueOk returns a tuple with the AssetRedeemValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssetRedeemValue

`func (o *Saver) SetAssetRedeemValue(v string)`

SetAssetRedeemValue sets AssetRedeemValue field to given value.


### GetGrowthPct

`func (o *Saver) GetGrowthPct() string`

GetGrowthPct returns the GrowthPct field if non-nil, zero value otherwise.

### GetGrowthPctOk

`func (o *Saver) GetGrowthPctOk() (*string, bool)`

GetGrowthPctOk returns a tuple with the GrowthPct field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrowthPct

`func (o *Saver) SetGrowthPct(v string)`

SetGrowthPct sets GrowthPct field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


