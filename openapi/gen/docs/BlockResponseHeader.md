# BlockResponseHeader

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | [**BlockResponseHeaderVersion**](BlockResponseHeaderVersion.md) |  | 
**ChainId** | **string** |  | 
**Height** | **int64** |  | 
**Time** | **string** |  | 
**LastBlockId** | [**BlockResponseId**](BlockResponseId.md) |  | 
**LastCommitHash** | **string** |  | 
**DataHash** | **string** |  | 
**ValidatorsHash** | **string** |  | 
**NextValidatorsHash** | **string** |  | 
**ConsensusHash** | **string** |  | 
**AppHash** | **string** |  | 
**LastResultsHash** | **string** |  | 
**EvidenceHash** | **string** |  | 
**ProposerAddress** | **string** |  | 

## Methods

### NewBlockResponseHeader

`func NewBlockResponseHeader(version BlockResponseHeaderVersion, chainId string, height int64, time string, lastBlockId BlockResponseId, lastCommitHash string, dataHash string, validatorsHash string, nextValidatorsHash string, consensusHash string, appHash string, lastResultsHash string, evidenceHash string, proposerAddress string, ) *BlockResponseHeader`

NewBlockResponseHeader instantiates a new BlockResponseHeader object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockResponseHeaderWithDefaults

`func NewBlockResponseHeaderWithDefaults() *BlockResponseHeader`

NewBlockResponseHeaderWithDefaults instantiates a new BlockResponseHeader object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersion

`func (o *BlockResponseHeader) GetVersion() BlockResponseHeaderVersion`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *BlockResponseHeader) GetVersionOk() (*BlockResponseHeaderVersion, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *BlockResponseHeader) SetVersion(v BlockResponseHeaderVersion)`

SetVersion sets Version field to given value.


### GetChainId

`func (o *BlockResponseHeader) GetChainId() string`

GetChainId returns the ChainId field if non-nil, zero value otherwise.

### GetChainIdOk

`func (o *BlockResponseHeader) GetChainIdOk() (*string, bool)`

GetChainIdOk returns a tuple with the ChainId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChainId

`func (o *BlockResponseHeader) SetChainId(v string)`

SetChainId sets ChainId field to given value.


### GetHeight

`func (o *BlockResponseHeader) GetHeight() int64`

GetHeight returns the Height field if non-nil, zero value otherwise.

### GetHeightOk

`func (o *BlockResponseHeader) GetHeightOk() (*int64, bool)`

GetHeightOk returns a tuple with the Height field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeight

`func (o *BlockResponseHeader) SetHeight(v int64)`

SetHeight sets Height field to given value.


### GetTime

`func (o *BlockResponseHeader) GetTime() string`

GetTime returns the Time field if non-nil, zero value otherwise.

### GetTimeOk

`func (o *BlockResponseHeader) GetTimeOk() (*string, bool)`

GetTimeOk returns a tuple with the Time field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTime

`func (o *BlockResponseHeader) SetTime(v string)`

SetTime sets Time field to given value.


### GetLastBlockId

`func (o *BlockResponseHeader) GetLastBlockId() BlockResponseId`

GetLastBlockId returns the LastBlockId field if non-nil, zero value otherwise.

### GetLastBlockIdOk

`func (o *BlockResponseHeader) GetLastBlockIdOk() (*BlockResponseId, bool)`

GetLastBlockIdOk returns a tuple with the LastBlockId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastBlockId

`func (o *BlockResponseHeader) SetLastBlockId(v BlockResponseId)`

SetLastBlockId sets LastBlockId field to given value.


### GetLastCommitHash

`func (o *BlockResponseHeader) GetLastCommitHash() string`

GetLastCommitHash returns the LastCommitHash field if non-nil, zero value otherwise.

### GetLastCommitHashOk

`func (o *BlockResponseHeader) GetLastCommitHashOk() (*string, bool)`

GetLastCommitHashOk returns a tuple with the LastCommitHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastCommitHash

`func (o *BlockResponseHeader) SetLastCommitHash(v string)`

SetLastCommitHash sets LastCommitHash field to given value.


### GetDataHash

`func (o *BlockResponseHeader) GetDataHash() string`

GetDataHash returns the DataHash field if non-nil, zero value otherwise.

### GetDataHashOk

`func (o *BlockResponseHeader) GetDataHashOk() (*string, bool)`

GetDataHashOk returns a tuple with the DataHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDataHash

`func (o *BlockResponseHeader) SetDataHash(v string)`

SetDataHash sets DataHash field to given value.


### GetValidatorsHash

`func (o *BlockResponseHeader) GetValidatorsHash() string`

GetValidatorsHash returns the ValidatorsHash field if non-nil, zero value otherwise.

### GetValidatorsHashOk

`func (o *BlockResponseHeader) GetValidatorsHashOk() (*string, bool)`

GetValidatorsHashOk returns a tuple with the ValidatorsHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidatorsHash

`func (o *BlockResponseHeader) SetValidatorsHash(v string)`

SetValidatorsHash sets ValidatorsHash field to given value.


### GetNextValidatorsHash

`func (o *BlockResponseHeader) GetNextValidatorsHash() string`

GetNextValidatorsHash returns the NextValidatorsHash field if non-nil, zero value otherwise.

### GetNextValidatorsHashOk

`func (o *BlockResponseHeader) GetNextValidatorsHashOk() (*string, bool)`

GetNextValidatorsHashOk returns a tuple with the NextValidatorsHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextValidatorsHash

`func (o *BlockResponseHeader) SetNextValidatorsHash(v string)`

SetNextValidatorsHash sets NextValidatorsHash field to given value.


### GetConsensusHash

`func (o *BlockResponseHeader) GetConsensusHash() string`

GetConsensusHash returns the ConsensusHash field if non-nil, zero value otherwise.

### GetConsensusHashOk

`func (o *BlockResponseHeader) GetConsensusHashOk() (*string, bool)`

GetConsensusHashOk returns a tuple with the ConsensusHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConsensusHash

`func (o *BlockResponseHeader) SetConsensusHash(v string)`

SetConsensusHash sets ConsensusHash field to given value.


### GetAppHash

`func (o *BlockResponseHeader) GetAppHash() string`

GetAppHash returns the AppHash field if non-nil, zero value otherwise.

### GetAppHashOk

`func (o *BlockResponseHeader) GetAppHashOk() (*string, bool)`

GetAppHashOk returns a tuple with the AppHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAppHash

`func (o *BlockResponseHeader) SetAppHash(v string)`

SetAppHash sets AppHash field to given value.


### GetLastResultsHash

`func (o *BlockResponseHeader) GetLastResultsHash() string`

GetLastResultsHash returns the LastResultsHash field if non-nil, zero value otherwise.

### GetLastResultsHashOk

`func (o *BlockResponseHeader) GetLastResultsHashOk() (*string, bool)`

GetLastResultsHashOk returns a tuple with the LastResultsHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastResultsHash

`func (o *BlockResponseHeader) SetLastResultsHash(v string)`

SetLastResultsHash sets LastResultsHash field to given value.


### GetEvidenceHash

`func (o *BlockResponseHeader) GetEvidenceHash() string`

GetEvidenceHash returns the EvidenceHash field if non-nil, zero value otherwise.

### GetEvidenceHashOk

`func (o *BlockResponseHeader) GetEvidenceHashOk() (*string, bool)`

GetEvidenceHashOk returns a tuple with the EvidenceHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvidenceHash

`func (o *BlockResponseHeader) SetEvidenceHash(v string)`

SetEvidenceHash sets EvidenceHash field to given value.


### GetProposerAddress

`func (o *BlockResponseHeader) GetProposerAddress() string`

GetProposerAddress returns the ProposerAddress field if non-nil, zero value otherwise.

### GetProposerAddressOk

`func (o *BlockResponseHeader) GetProposerAddressOk() (*string, bool)`

GetProposerAddressOk returns a tuple with the ProposerAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProposerAddress

`func (o *BlockResponseHeader) SetProposerAddress(v string)`

SetProposerAddress sets ProposerAddress field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


