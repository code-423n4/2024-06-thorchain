# PoolSlipResponseInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** |  | 
**PoolSlip** | **int64** | Pool slip for this asset&#39;s pool for the current height | 
**RollupCount** | **int64** | Number of stored pool slips contributing to the current stored rollup | 
**LongRollup** | **int64** | Median of rollup snapshots over a long period | 
**Rollup** | **int64** | Stored sum of pool slips over a number of previous block heights | 
**SummedRollup** | Pointer to **int64** | Summed pool slips over a number of previous block heights, to checksum the stored rollup | [optional] 

## Methods

### NewPoolSlipResponseInner

`func NewPoolSlipResponseInner(asset string, poolSlip int64, rollupCount int64, longRollup int64, rollup int64, ) *PoolSlipResponseInner`

NewPoolSlipResponseInner instantiates a new PoolSlipResponseInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPoolSlipResponseInnerWithDefaults

`func NewPoolSlipResponseInnerWithDefaults() *PoolSlipResponseInner`

NewPoolSlipResponseInnerWithDefaults instantiates a new PoolSlipResponseInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *PoolSlipResponseInner) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *PoolSlipResponseInner) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *PoolSlipResponseInner) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetPoolSlip

`func (o *PoolSlipResponseInner) GetPoolSlip() int64`

GetPoolSlip returns the PoolSlip field if non-nil, zero value otherwise.

### GetPoolSlipOk

`func (o *PoolSlipResponseInner) GetPoolSlipOk() (*int64, bool)`

GetPoolSlipOk returns a tuple with the PoolSlip field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolSlip

`func (o *PoolSlipResponseInner) SetPoolSlip(v int64)`

SetPoolSlip sets PoolSlip field to given value.


### GetRollupCount

`func (o *PoolSlipResponseInner) GetRollupCount() int64`

GetRollupCount returns the RollupCount field if non-nil, zero value otherwise.

### GetRollupCountOk

`func (o *PoolSlipResponseInner) GetRollupCountOk() (*int64, bool)`

GetRollupCountOk returns a tuple with the RollupCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRollupCount

`func (o *PoolSlipResponseInner) SetRollupCount(v int64)`

SetRollupCount sets RollupCount field to given value.


### GetLongRollup

`func (o *PoolSlipResponseInner) GetLongRollup() int64`

GetLongRollup returns the LongRollup field if non-nil, zero value otherwise.

### GetLongRollupOk

`func (o *PoolSlipResponseInner) GetLongRollupOk() (*int64, bool)`

GetLongRollupOk returns a tuple with the LongRollup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLongRollup

`func (o *PoolSlipResponseInner) SetLongRollup(v int64)`

SetLongRollup sets LongRollup field to given value.


### GetRollup

`func (o *PoolSlipResponseInner) GetRollup() int64`

GetRollup returns the Rollup field if non-nil, zero value otherwise.

### GetRollupOk

`func (o *PoolSlipResponseInner) GetRollupOk() (*int64, bool)`

GetRollupOk returns a tuple with the Rollup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRollup

`func (o *PoolSlipResponseInner) SetRollup(v int64)`

SetRollup sets Rollup field to given value.


### GetSummedRollup

`func (o *PoolSlipResponseInner) GetSummedRollup() int64`

GetSummedRollup returns the SummedRollup field if non-nil, zero value otherwise.

### GetSummedRollupOk

`func (o *PoolSlipResponseInner) GetSummedRollupOk() (*int64, bool)`

GetSummedRollupOk returns a tuple with the SummedRollup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSummedRollup

`func (o *PoolSlipResponseInner) SetSummedRollup(v int64)`

SetSummedRollup sets SummedRollup field to given value.

### HasSummedRollup

`func (o *PoolSlipResponseInner) HasSummedRollup() bool`

HasSummedRollup returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


