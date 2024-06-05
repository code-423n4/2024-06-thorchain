# SwapperCloutResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address** | **string** | address associated with this clout account | 
**Score** | Pointer to **string** | clout score, which is the amount of rune spent on swap fees | [optional] 
**Reclaimed** | Pointer to **string** | amount of clout that has been reclaimed in total over time (observed clout spent) | [optional] 
**Spent** | Pointer to **string** | amount of clout that has been spent in total over time | [optional] 
**LastSpentHeight** | Pointer to **int64** | last block height that clout was spent | [optional] 
**LastReclaimHeight** | Pointer to **int64** | last block height that clout was reclaimed | [optional] 

## Methods

### NewSwapperCloutResponse

`func NewSwapperCloutResponse(address string, ) *SwapperCloutResponse`

NewSwapperCloutResponse instantiates a new SwapperCloutResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwapperCloutResponseWithDefaults

`func NewSwapperCloutResponseWithDefaults() *SwapperCloutResponse`

NewSwapperCloutResponseWithDefaults instantiates a new SwapperCloutResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAddress

`func (o *SwapperCloutResponse) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *SwapperCloutResponse) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *SwapperCloutResponse) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetScore

`func (o *SwapperCloutResponse) GetScore() string`

GetScore returns the Score field if non-nil, zero value otherwise.

### GetScoreOk

`func (o *SwapperCloutResponse) GetScoreOk() (*string, bool)`

GetScoreOk returns a tuple with the Score field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScore

`func (o *SwapperCloutResponse) SetScore(v string)`

SetScore sets Score field to given value.

### HasScore

`func (o *SwapperCloutResponse) HasScore() bool`

HasScore returns a boolean if a field has been set.

### GetReclaimed

`func (o *SwapperCloutResponse) GetReclaimed() string`

GetReclaimed returns the Reclaimed field if non-nil, zero value otherwise.

### GetReclaimedOk

`func (o *SwapperCloutResponse) GetReclaimedOk() (*string, bool)`

GetReclaimedOk returns a tuple with the Reclaimed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReclaimed

`func (o *SwapperCloutResponse) SetReclaimed(v string)`

SetReclaimed sets Reclaimed field to given value.

### HasReclaimed

`func (o *SwapperCloutResponse) HasReclaimed() bool`

HasReclaimed returns a boolean if a field has been set.

### GetSpent

`func (o *SwapperCloutResponse) GetSpent() string`

GetSpent returns the Spent field if non-nil, zero value otherwise.

### GetSpentOk

`func (o *SwapperCloutResponse) GetSpentOk() (*string, bool)`

GetSpentOk returns a tuple with the Spent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpent

`func (o *SwapperCloutResponse) SetSpent(v string)`

SetSpent sets Spent field to given value.

### HasSpent

`func (o *SwapperCloutResponse) HasSpent() bool`

HasSpent returns a boolean if a field has been set.

### GetLastSpentHeight

`func (o *SwapperCloutResponse) GetLastSpentHeight() int64`

GetLastSpentHeight returns the LastSpentHeight field if non-nil, zero value otherwise.

### GetLastSpentHeightOk

`func (o *SwapperCloutResponse) GetLastSpentHeightOk() (*int64, bool)`

GetLastSpentHeightOk returns a tuple with the LastSpentHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastSpentHeight

`func (o *SwapperCloutResponse) SetLastSpentHeight(v int64)`

SetLastSpentHeight sets LastSpentHeight field to given value.

### HasLastSpentHeight

`func (o *SwapperCloutResponse) HasLastSpentHeight() bool`

HasLastSpentHeight returns a boolean if a field has been set.

### GetLastReclaimHeight

`func (o *SwapperCloutResponse) GetLastReclaimHeight() int64`

GetLastReclaimHeight returns the LastReclaimHeight field if non-nil, zero value otherwise.

### GetLastReclaimHeightOk

`func (o *SwapperCloutResponse) GetLastReclaimHeightOk() (*int64, bool)`

GetLastReclaimHeightOk returns a tuple with the LastReclaimHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastReclaimHeight

`func (o *SwapperCloutResponse) SetLastReclaimHeight(v int64)`

SetLastReclaimHeight sets LastReclaimHeight field to given value.

### HasLastReclaimHeight

`func (o *SwapperCloutResponse) HasLastReclaimHeight() bool`

HasLastReclaimHeight returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


