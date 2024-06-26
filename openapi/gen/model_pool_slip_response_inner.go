/*
Thornode API

Thornode REST API.

Contact: devs@thorchain.org
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// PoolSlipResponseInner struct for PoolSlipResponseInner
type PoolSlipResponseInner struct {
	Asset string `json:"asset"`
	// Pool slip for this asset's pool for the current height
	PoolSlip int64 `json:"pool_slip"`
	// Number of stored pool slips contributing to the current stored rollup
	RollupCount int64 `json:"rollup_count"`
	// Median of rollup snapshots over a long period
	LongRollup int64 `json:"long_rollup"`
	// Stored sum of pool slips over a number of previous block heights
	Rollup int64 `json:"rollup"`
	// Summed pool slips over a number of previous block heights, to checksum the stored rollup
	SummedRollup *int64 `json:"summed_rollup,omitempty"`
}

// NewPoolSlipResponseInner instantiates a new PoolSlipResponseInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPoolSlipResponseInner(asset string, poolSlip int64, rollupCount int64, longRollup int64, rollup int64) *PoolSlipResponseInner {
	this := PoolSlipResponseInner{}
	this.Asset = asset
	this.PoolSlip = poolSlip
	this.RollupCount = rollupCount
	this.LongRollup = longRollup
	this.Rollup = rollup
	return &this
}

// NewPoolSlipResponseInnerWithDefaults instantiates a new PoolSlipResponseInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPoolSlipResponseInnerWithDefaults() *PoolSlipResponseInner {
	this := PoolSlipResponseInner{}
	return &this
}

// GetAsset returns the Asset field value
func (o *PoolSlipResponseInner) GetAsset() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Asset
}

// GetAssetOk returns a tuple with the Asset field value
// and a boolean to check if the value has been set.
func (o *PoolSlipResponseInner) GetAssetOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Asset, true
}

// SetAsset sets field value
func (o *PoolSlipResponseInner) SetAsset(v string) {
	o.Asset = v
}

// GetPoolSlip returns the PoolSlip field value
func (o *PoolSlipResponseInner) GetPoolSlip() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.PoolSlip
}

// GetPoolSlipOk returns a tuple with the PoolSlip field value
// and a boolean to check if the value has been set.
func (o *PoolSlipResponseInner) GetPoolSlipOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PoolSlip, true
}

// SetPoolSlip sets field value
func (o *PoolSlipResponseInner) SetPoolSlip(v int64) {
	o.PoolSlip = v
}

// GetRollupCount returns the RollupCount field value
func (o *PoolSlipResponseInner) GetRollupCount() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.RollupCount
}

// GetRollupCountOk returns a tuple with the RollupCount field value
// and a boolean to check if the value has been set.
func (o *PoolSlipResponseInner) GetRollupCountOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RollupCount, true
}

// SetRollupCount sets field value
func (o *PoolSlipResponseInner) SetRollupCount(v int64) {
	o.RollupCount = v
}

// GetLongRollup returns the LongRollup field value
func (o *PoolSlipResponseInner) GetLongRollup() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.LongRollup
}

// GetLongRollupOk returns a tuple with the LongRollup field value
// and a boolean to check if the value has been set.
func (o *PoolSlipResponseInner) GetLongRollupOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LongRollup, true
}

// SetLongRollup sets field value
func (o *PoolSlipResponseInner) SetLongRollup(v int64) {
	o.LongRollup = v
}

// GetRollup returns the Rollup field value
func (o *PoolSlipResponseInner) GetRollup() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Rollup
}

// GetRollupOk returns a tuple with the Rollup field value
// and a boolean to check if the value has been set.
func (o *PoolSlipResponseInner) GetRollupOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Rollup, true
}

// SetRollup sets field value
func (o *PoolSlipResponseInner) SetRollup(v int64) {
	o.Rollup = v
}

// GetSummedRollup returns the SummedRollup field value if set, zero value otherwise.
func (o *PoolSlipResponseInner) GetSummedRollup() int64 {
	if o == nil || o.SummedRollup == nil {
		var ret int64
		return ret
	}
	return *o.SummedRollup
}

// GetSummedRollupOk returns a tuple with the SummedRollup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PoolSlipResponseInner) GetSummedRollupOk() (*int64, bool) {
	if o == nil || o.SummedRollup == nil {
		return nil, false
	}
	return o.SummedRollup, true
}

// HasSummedRollup returns a boolean if a field has been set.
func (o *PoolSlipResponseInner) HasSummedRollup() bool {
	if o != nil && o.SummedRollup != nil {
		return true
	}

	return false
}

// SetSummedRollup gets a reference to the given int64 and assigns it to the SummedRollup field.
func (o *PoolSlipResponseInner) SetSummedRollup(v int64) {
	o.SummedRollup = &v
}

func (o PoolSlipResponseInner) MarshalJSON_deprecated() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["asset"] = o.Asset
	}
	if true {
		toSerialize["pool_slip"] = o.PoolSlip
	}
	if true {
		toSerialize["rollup_count"] = o.RollupCount
	}
	if true {
		toSerialize["long_rollup"] = o.LongRollup
	}
	if true {
		toSerialize["rollup"] = o.Rollup
	}
	if o.SummedRollup != nil {
		toSerialize["summed_rollup"] = o.SummedRollup
	}
	return json.Marshal(toSerialize)
}

type NullablePoolSlipResponseInner struct {
	value *PoolSlipResponseInner
	isSet bool
}

func (v NullablePoolSlipResponseInner) Get() *PoolSlipResponseInner {
	return v.value
}

func (v *NullablePoolSlipResponseInner) Set(val *PoolSlipResponseInner) {
	v.value = val
	v.isSet = true
}

func (v NullablePoolSlipResponseInner) IsSet() bool {
	return v.isSet
}

func (v *NullablePoolSlipResponseInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePoolSlipResponseInner(val *PoolSlipResponseInner) *NullablePoolSlipResponseInner {
	return &NullablePoolSlipResponseInner{value: val, isSet: true}
}

func (v NullablePoolSlipResponseInner) MarshalJSON_deprecated() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePoolSlipResponseInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


