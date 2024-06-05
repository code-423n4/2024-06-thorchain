# MimirV2IDsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int** |  | 
**Name** | **string** |  | 
**VoteKey** | **string** |  | 
**LegacyKey** | **string** |  | 
**Type** | **string** |  | 
**Votes** | **map[string]int64** |  | 

## Methods

### NewMimirV2IDsResponse

`func NewMimirV2IDsResponse(id int, name string, voteKey string, legacyKey string, type_ string, votes map[string]int64, ) *MimirV2IDsResponse`

NewMimirV2IDsResponse instantiates a new MimirV2IDsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMimirV2IDsResponseWithDefaults

`func NewMimirV2IDsResponseWithDefaults() *MimirV2IDsResponse`

NewMimirV2IDsResponseWithDefaults instantiates a new MimirV2IDsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *MimirV2IDsResponse) GetId() int`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *MimirV2IDsResponse) GetIdOk() (*int, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *MimirV2IDsResponse) SetId(v int)`

SetId sets Id field to given value.


### GetName

`func (o *MimirV2IDsResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *MimirV2IDsResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *MimirV2IDsResponse) SetName(v string)`

SetName sets Name field to given value.


### GetVoteKey

`func (o *MimirV2IDsResponse) GetVoteKey() string`

GetVoteKey returns the VoteKey field if non-nil, zero value otherwise.

### GetVoteKeyOk

`func (o *MimirV2IDsResponse) GetVoteKeyOk() (*string, bool)`

GetVoteKeyOk returns a tuple with the VoteKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoteKey

`func (o *MimirV2IDsResponse) SetVoteKey(v string)`

SetVoteKey sets VoteKey field to given value.


### GetLegacyKey

`func (o *MimirV2IDsResponse) GetLegacyKey() string`

GetLegacyKey returns the LegacyKey field if non-nil, zero value otherwise.

### GetLegacyKeyOk

`func (o *MimirV2IDsResponse) GetLegacyKeyOk() (*string, bool)`

GetLegacyKeyOk returns a tuple with the LegacyKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLegacyKey

`func (o *MimirV2IDsResponse) SetLegacyKey(v string)`

SetLegacyKey sets LegacyKey field to given value.


### GetType

`func (o *MimirV2IDsResponse) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *MimirV2IDsResponse) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *MimirV2IDsResponse) SetType(v string)`

SetType sets Type field to given value.


### GetVotes

`func (o *MimirV2IDsResponse) GetVotes() map[string]int64`

GetVotes returns the Votes field if non-nil, zero value otherwise.

### GetVotesOk

`func (o *MimirV2IDsResponse) GetVotesOk() (*map[string]int64, bool)`

GetVotesOk returns a tuple with the Votes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVotes

`func (o *MimirV2IDsResponse) SetVotes(v map[string]int64)`

SetVotes sets Votes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


