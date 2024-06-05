# Borrower

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Owner** | **string** |  | 
**Asset** | **string** |  | 
**DebtIssued** | **string** |  | 
**DebtRepaid** | **string** |  | 
**DebtCurrent** | **string** |  | 
**CollateralDeposited** | **string** |  | 
**CollateralWithdrawn** | **string** |  | 
**CollateralCurrent** | **string** |  | 
**LastOpenHeight** | **int64** |  | 
**LastRepayHeight** | **int64** |  | 

## Methods

### NewBorrower

`func NewBorrower(owner string, asset string, debtIssued string, debtRepaid string, debtCurrent string, collateralDeposited string, collateralWithdrawn string, collateralCurrent string, lastOpenHeight int64, lastRepayHeight int64, ) *Borrower`

NewBorrower instantiates a new Borrower object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBorrowerWithDefaults

`func NewBorrowerWithDefaults() *Borrower`

NewBorrowerWithDefaults instantiates a new Borrower object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOwner

`func (o *Borrower) GetOwner() string`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *Borrower) GetOwnerOk() (*string, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *Borrower) SetOwner(v string)`

SetOwner sets Owner field to given value.


### GetAsset

`func (o *Borrower) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *Borrower) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *Borrower) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetDebtIssued

`func (o *Borrower) GetDebtIssued() string`

GetDebtIssued returns the DebtIssued field if non-nil, zero value otherwise.

### GetDebtIssuedOk

`func (o *Borrower) GetDebtIssuedOk() (*string, bool)`

GetDebtIssuedOk returns a tuple with the DebtIssued field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDebtIssued

`func (o *Borrower) SetDebtIssued(v string)`

SetDebtIssued sets DebtIssued field to given value.


### GetDebtRepaid

`func (o *Borrower) GetDebtRepaid() string`

GetDebtRepaid returns the DebtRepaid field if non-nil, zero value otherwise.

### GetDebtRepaidOk

`func (o *Borrower) GetDebtRepaidOk() (*string, bool)`

GetDebtRepaidOk returns a tuple with the DebtRepaid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDebtRepaid

`func (o *Borrower) SetDebtRepaid(v string)`

SetDebtRepaid sets DebtRepaid field to given value.


### GetDebtCurrent

`func (o *Borrower) GetDebtCurrent() string`

GetDebtCurrent returns the DebtCurrent field if non-nil, zero value otherwise.

### GetDebtCurrentOk

`func (o *Borrower) GetDebtCurrentOk() (*string, bool)`

GetDebtCurrentOk returns a tuple with the DebtCurrent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDebtCurrent

`func (o *Borrower) SetDebtCurrent(v string)`

SetDebtCurrent sets DebtCurrent field to given value.


### GetCollateralDeposited

`func (o *Borrower) GetCollateralDeposited() string`

GetCollateralDeposited returns the CollateralDeposited field if non-nil, zero value otherwise.

### GetCollateralDepositedOk

`func (o *Borrower) GetCollateralDepositedOk() (*string, bool)`

GetCollateralDepositedOk returns a tuple with the CollateralDeposited field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCollateralDeposited

`func (o *Borrower) SetCollateralDeposited(v string)`

SetCollateralDeposited sets CollateralDeposited field to given value.


### GetCollateralWithdrawn

`func (o *Borrower) GetCollateralWithdrawn() string`

GetCollateralWithdrawn returns the CollateralWithdrawn field if non-nil, zero value otherwise.

### GetCollateralWithdrawnOk

`func (o *Borrower) GetCollateralWithdrawnOk() (*string, bool)`

GetCollateralWithdrawnOk returns a tuple with the CollateralWithdrawn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCollateralWithdrawn

`func (o *Borrower) SetCollateralWithdrawn(v string)`

SetCollateralWithdrawn sets CollateralWithdrawn field to given value.


### GetCollateralCurrent

`func (o *Borrower) GetCollateralCurrent() string`

GetCollateralCurrent returns the CollateralCurrent field if non-nil, zero value otherwise.

### GetCollateralCurrentOk

`func (o *Borrower) GetCollateralCurrentOk() (*string, bool)`

GetCollateralCurrentOk returns a tuple with the CollateralCurrent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCollateralCurrent

`func (o *Borrower) SetCollateralCurrent(v string)`

SetCollateralCurrent sets CollateralCurrent field to given value.


### GetLastOpenHeight

`func (o *Borrower) GetLastOpenHeight() int64`

GetLastOpenHeight returns the LastOpenHeight field if non-nil, zero value otherwise.

### GetLastOpenHeightOk

`func (o *Borrower) GetLastOpenHeightOk() (*int64, bool)`

GetLastOpenHeightOk returns a tuple with the LastOpenHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastOpenHeight

`func (o *Borrower) SetLastOpenHeight(v int64)`

SetLastOpenHeight sets LastOpenHeight field to given value.


### GetLastRepayHeight

`func (o *Borrower) GetLastRepayHeight() int64`

GetLastRepayHeight returns the LastRepayHeight field if non-nil, zero value otherwise.

### GetLastRepayHeightOk

`func (o *Borrower) GetLastRepayHeightOk() (*int64, bool)`

GetLastRepayHeightOk returns a tuple with the LastRepayHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastRepayHeight

`func (o *Borrower) SetLastRepayHeight(v int64)`

SetLastRepayHeight sets LastRepayHeight field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


