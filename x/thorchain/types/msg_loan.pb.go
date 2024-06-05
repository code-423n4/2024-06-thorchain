// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: thorchain/v1/x/thorchain/types/msg_loan.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	common "gitlab.com/thorchain/thornode/common"
	gitlab_com_thorchain_thornode_common "gitlab.com/thorchain/thornode/common"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MsgLoanOpen struct {
	Owner                   gitlab_com_thorchain_thornode_common.Address  `protobuf:"bytes,1,opt,name=owner,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"owner,omitempty"`
	CollateralAsset         common.Asset                                  `protobuf:"bytes,2,opt,name=collateral_asset,json=collateralAsset,proto3" json:"collateral_asset"`
	CollateralAmount        github_com_cosmos_cosmos_sdk_types.Uint       `protobuf:"bytes,3,opt,name=collateral_amount,json=collateralAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"collateral_amount"`
	TargetAddress           gitlab_com_thorchain_thornode_common.Address  `protobuf:"bytes,4,opt,name=target_address,json=targetAddress,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"target_address,omitempty"`
	TargetAsset             common.Asset                                  `protobuf:"bytes,5,opt,name=target_asset,json=targetAsset,proto3" json:"target_asset"`
	MinOut                  github_com_cosmos_cosmos_sdk_types.Uint       `protobuf:"bytes,6,opt,name=min_out,json=minOut,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"min_out"`
	AffiliateAddress        gitlab_com_thorchain_thornode_common.Address  `protobuf:"bytes,7,opt,name=affiliate_address,json=affiliateAddress,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"affiliate_address,omitempty"`
	AffiliateBasisPoints    github_com_cosmos_cosmos_sdk_types.Uint       `protobuf:"bytes,8,opt,name=affiliate_basis_points,json=affiliateBasisPoints,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"affiliate_basis_points"`
	Aggregator              string                                        `protobuf:"bytes,9,opt,name=aggregator,proto3" json:"aggregator,omitempty"`
	AggregatorTargetAddress string                                        `protobuf:"bytes,10,opt,name=aggregator_target_address,json=aggregatorTargetAddress,proto3" json:"aggregator_target_address,omitempty"`
	AggregatorTargetLimit   github_com_cosmos_cosmos_sdk_types.Uint       `protobuf:"bytes,11,opt,name=aggregator_target_limit,json=aggregatorTargetLimit,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"aggregator_target_limit"`
	Signer                  github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,12,opt,name=signer,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"signer,omitempty"`
	TxID                    gitlab_com_thorchain_thornode_common.TxID     `protobuf:"bytes,13,opt,name=tx_id,json=txId,proto3,casttype=gitlab.com/thorchain/thornode/common.TxID" json:"tx_id,omitempty"`
}

func (m *MsgLoanOpen) Reset()         { *m = MsgLoanOpen{} }
func (m *MsgLoanOpen) String() string { return proto.CompactTextString(m) }
func (*MsgLoanOpen) ProtoMessage()    {}
func (*MsgLoanOpen) Descriptor() ([]byte, []int) {
	return fileDescriptor_c67798bf3171e43f, []int{0}
}
func (m *MsgLoanOpen) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgLoanOpen) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgLoanOpen.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgLoanOpen) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgLoanOpen.Merge(m, src)
}
func (m *MsgLoanOpen) XXX_Size() int {
	return m.Size()
}
func (m *MsgLoanOpen) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgLoanOpen.DiscardUnknown(m)
}

var xxx_messageInfo_MsgLoanOpen proto.InternalMessageInfo

func (m *MsgLoanOpen) GetOwner() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *MsgLoanOpen) GetCollateralAsset() common.Asset {
	if m != nil {
		return m.CollateralAsset
	}
	return common.Asset{}
}

func (m *MsgLoanOpen) GetTargetAddress() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.TargetAddress
	}
	return ""
}

func (m *MsgLoanOpen) GetTargetAsset() common.Asset {
	if m != nil {
		return m.TargetAsset
	}
	return common.Asset{}
}

func (m *MsgLoanOpen) GetAffiliateAddress() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.AffiliateAddress
	}
	return ""
}

func (m *MsgLoanOpen) GetAggregator() string {
	if m != nil {
		return m.Aggregator
	}
	return ""
}

func (m *MsgLoanOpen) GetAggregatorTargetAddress() string {
	if m != nil {
		return m.AggregatorTargetAddress
	}
	return ""
}

func (m *MsgLoanOpen) GetSigner() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Signer
	}
	return nil
}

func (m *MsgLoanOpen) GetTxID() gitlab_com_thorchain_thornode_common.TxID {
	if m != nil {
		return m.TxID
	}
	return ""
}

type MsgLoanRepayment struct {
	Owner           gitlab_com_thorchain_thornode_common.Address  `protobuf:"bytes,1,opt,name=owner,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"owner,omitempty"`
	CollateralAsset common.Asset                                  `protobuf:"bytes,2,opt,name=collateral_asset,json=collateralAsset,proto3" json:"collateral_asset"`
	Coin            common.Coin                                   `protobuf:"bytes,3,opt,name=coin,proto3" json:"coin"`
	MinOut          github_com_cosmos_cosmos_sdk_types.Uint       `protobuf:"bytes,4,opt,name=min_out,json=minOut,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"min_out"`
	Signer          github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,5,opt,name=signer,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"signer,omitempty"`
	From            gitlab_com_thorchain_thornode_common.Address  `protobuf:"bytes,6,opt,name=from,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"from,omitempty"`
	TxID            gitlab_com_thorchain_thornode_common.TxID     `protobuf:"bytes,7,opt,name=tx_id,json=txId,proto3,casttype=gitlab.com/thorchain/thornode/common.TxID" json:"tx_id,omitempty"`
}

func (m *MsgLoanRepayment) Reset()         { *m = MsgLoanRepayment{} }
func (m *MsgLoanRepayment) String() string { return proto.CompactTextString(m) }
func (*MsgLoanRepayment) ProtoMessage()    {}
func (*MsgLoanRepayment) Descriptor() ([]byte, []int) {
	return fileDescriptor_c67798bf3171e43f, []int{1}
}
func (m *MsgLoanRepayment) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgLoanRepayment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgLoanRepayment.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgLoanRepayment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgLoanRepayment.Merge(m, src)
}
func (m *MsgLoanRepayment) XXX_Size() int {
	return m.Size()
}
func (m *MsgLoanRepayment) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgLoanRepayment.DiscardUnknown(m)
}

var xxx_messageInfo_MsgLoanRepayment proto.InternalMessageInfo

func (m *MsgLoanRepayment) GetOwner() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *MsgLoanRepayment) GetCollateralAsset() common.Asset {
	if m != nil {
		return m.CollateralAsset
	}
	return common.Asset{}
}

func (m *MsgLoanRepayment) GetCoin() common.Coin {
	if m != nil {
		return m.Coin
	}
	return common.Coin{}
}

func (m *MsgLoanRepayment) GetSigner() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Signer
	}
	return nil
}

func (m *MsgLoanRepayment) GetFrom() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *MsgLoanRepayment) GetTxID() gitlab_com_thorchain_thornode_common.TxID {
	if m != nil {
		return m.TxID
	}
	return ""
}

func init() {
	proto.RegisterType((*MsgLoanOpen)(nil), "types.MsgLoanOpen")
	proto.RegisterType((*MsgLoanRepayment)(nil), "types.MsgLoanRepayment")
}

func init() {
	proto.RegisterFile("thorchain/v1/x/thorchain/types/msg_loan.proto", fileDescriptor_c67798bf3171e43f)
}

var fileDescriptor_c67798bf3171e43f = []byte{
	// 592 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x55, 0x41, 0x6f, 0xd3, 0x30,
	0x18, 0x6d, 0x58, 0xda, 0x32, 0xb7, 0x83, 0x2e, 0x1a, 0x2c, 0xec, 0x90, 0x54, 0x3b, 0x40, 0x91,
	0x68, 0xc3, 0x86, 0xe0, 0xc0, 0x01, 0xa9, 0x65, 0x42, 0x14, 0x0d, 0x0d, 0x45, 0x43, 0x48, 0x08,
	0x14, 0xb9, 0x89, 0xeb, 0x5a, 0x24, 0x76, 0x15, 0xbb, 0xd0, 0xfd, 0x0b, 0x7e, 0xd6, 0x8e, 0x3b,
	0x22, 0x0e, 0x11, 0x6a, 0x4f, 0xdc, 0x38, 0xf7, 0x84, 0xe2, 0x24, 0x4d, 0xb7, 0x09, 0x34, 0xda,
	0x0b, 0xa7, 0xd8, 0x5f, 0xbe, 0xf7, 0xfc, 0x3d, 0x3f, 0x3d, 0x19, 0x34, 0xc5, 0x80, 0x85, 0xee,
	0x00, 0x12, 0x6a, 0x7d, 0xde, 0xb3, 0xc6, 0x56, 0xbe, 0x15, 0x27, 0x43, 0xc4, 0xad, 0x80, 0x63,
	0xc7, 0x67, 0x90, 0xb6, 0x86, 0x21, 0x13, 0x4c, 0x2b, 0xca, 0xea, 0x4e, 0xfd, 0x1c, 0xca, 0x65,
	0x41, 0xc0, 0x68, 0xfa, 0x49, 0x1a, 0x77, 0xb6, 0x30, 0xc3, 0x4c, 0x2e, 0xad, 0x78, 0x95, 0x54,
	0x77, 0x7f, 0x95, 0x41, 0xe5, 0x35, 0xc7, 0x87, 0x0c, 0xd2, 0xa3, 0x21, 0xa2, 0xda, 0x0b, 0x50,
	0x64, 0x5f, 0x28, 0x0a, 0x75, 0xa5, 0xae, 0x34, 0xd6, 0x3b, 0x0f, 0x67, 0x91, 0xf9, 0x00, 0x13,
	0xe1, 0xc3, 0x5e, 0xcb, 0x65, 0xc1, 0xe2, 0x30, 0x03, 0x16, 0x52, 0xe6, 0xa1, 0xec, 0x90, 0xb6,
	0xe7, 0x85, 0x88, 0x73, 0x3b, 0x81, 0x6b, 0xcf, 0x40, 0xcd, 0x65, 0xbe, 0x0f, 0x05, 0x0a, 0xa1,
	0xef, 0x40, 0xce, 0x91, 0xd0, 0xaf, 0xd5, 0x95, 0x46, 0x65, 0x7f, 0xa3, 0x95, 0x21, 0xe2, 0x62,
	0x47, 0x3d, 0x8d, 0xcc, 0x82, 0x7d, 0x33, 0x6f, 0x96, 0x65, 0xed, 0x03, 0xd8, 0x5c, 0xc4, 0x07,
	0x6c, 0x44, 0x85, 0xbe, 0x26, 0x67, 0xb2, 0x62, 0xc4, 0xf7, 0xc8, 0xbc, 0x87, 0x89, 0x18, 0x8c,
	0x92, 0xb9, 0x5c, 0xc6, 0x03, 0xc6, 0xd3, 0x4f, 0x93, 0x7b, 0x9f, 0x92, 0xcb, 0x6a, 0xbd, 0x25,
	0x54, 0xd8, 0x0b, 0x93, 0xb4, 0x25, 0x91, 0xf6, 0x0e, 0xdc, 0x10, 0x30, 0xc4, 0x48, 0x38, 0x30,
	0x19, 0x5b, 0x57, 0x97, 0x94, 0xbb, 0x91, 0xf0, 0xa4, 0x5b, 0xed, 0x09, 0xa8, 0x66, 0xc4, 0x52,
	0x72, 0xf1, 0xcf, 0x92, 0x2b, 0x29, 0x52, 0xca, 0x7d, 0x09, 0xca, 0x01, 0xa1, 0x0e, 0x1b, 0x09,
	0xbd, 0xb4, 0x9c, 0xc8, 0x52, 0x40, 0xe8, 0xd1, 0x48, 0x68, 0x1f, 0xc1, 0x26, 0xec, 0xf7, 0x89,
	0x4f, 0xa0, 0x40, 0x73, 0x75, 0xe5, 0x25, 0xd5, 0xd5, 0xe6, 0x54, 0x99, 0x40, 0x04, 0x6e, 0xe7,
	0xf4, 0x3d, 0xc8, 0x09, 0x77, 0x86, 0x8c, 0x50, 0xc1, 0xf5, 0xeb, 0xcb, 0xcd, 0xbd, 0x35, 0xa7,
	0xeb, 0xc4, 0x6c, 0x6f, 0x24, 0x99, 0x66, 0x00, 0x00, 0x31, 0x0e, 0x11, 0x86, 0x82, 0x85, 0xfa,
	0x7a, 0x4c, 0x6d, 0x2f, 0x54, 0xb4, 0xa7, 0xe0, 0x4e, 0xbe, 0x73, 0x2e, 0x78, 0x09, 0x64, 0xfb,
	0x76, 0xde, 0x70, 0x7c, 0xce, 0x23, 0x0c, 0xb6, 0x2f, 0x63, 0x7d, 0x12, 0x10, 0xa1, 0x57, 0x96,
	0xd3, 0x70, 0xeb, 0xe2, 0x51, 0x87, 0x31, 0x9b, 0xd6, 0x05, 0x25, 0x4e, 0x70, 0x1c, 0xa6, 0x6a,
	0x5d, 0x69, 0x54, 0x3b, 0x7b, 0xb3, 0xc8, 0x6c, 0x5e, 0x81, 0xb3, 0xed, 0xba, 0x99, 0x01, 0x29,
	0x81, 0xf6, 0x0a, 0x14, 0xc5, 0xd8, 0x21, 0x9e, 0xbe, 0x21, 0x27, 0x7c, 0x3c, 0x89, 0x4c, 0xf5,
	0x78, 0xdc, 0x3d, 0x98, 0x45, 0xe6, 0xfd, 0x2b, 0x39, 0x1a, 0x37, 0xdb, 0xaa, 0x18, 0x77, 0xbd,
	0xdd, 0x9f, 0x6b, 0xa0, 0x96, 0x46, 0xde, 0x46, 0x43, 0x78, 0x12, 0x20, 0x2a, 0xfe, 0x9b, 0xdc,
	0xdf, 0x05, 0xaa, 0xcb, 0x08, 0x95, 0x51, 0xaf, 0xec, 0x57, 0x33, 0xcc, 0x73, 0x46, 0x68, 0x0a,
	0x91, 0xff, 0x17, 0x03, 0xa3, 0xae, 0x16, 0x98, 0xdc, 0xa5, 0xe2, 0xaa, 0x2e, 0x1d, 0x00, 0xb5,
	0x1f, 0xb2, 0x20, 0x8d, 0xf0, 0xbf, 0xdf, 0xa1, 0x44, 0xe7, 0x5e, 0x97, 0x57, 0xf6, 0xba, 0xd3,
	0x3d, 0x9d, 0x18, 0xca, 0xd9, 0xc4, 0x50, 0x7e, 0x4c, 0x0c, 0xe5, 0xeb, 0xd4, 0x28, 0x9c, 0x4d,
	0x8d, 0xc2, 0xb7, 0xa9, 0x51, 0x78, 0x6f, 0xfd, 0x9d, 0xea, 0xd2, 0xbb, 0xd3, 0x2b, 0xc9, 0x07,
	0xe3, 0xd1, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x56, 0x0c, 0x86, 0xa0, 0x06, 0x00, 0x00,
}

func (m *MsgLoanOpen) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgLoanOpen) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgLoanOpen) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TxID) > 0 {
		i -= len(m.TxID)
		copy(dAtA[i:], m.TxID)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.TxID)))
		i--
		dAtA[i] = 0x6a
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0x62
	}
	{
		size := m.AggregatorTargetLimit.Size()
		i -= size
		if _, err := m.AggregatorTargetLimit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x5a
	if len(m.AggregatorTargetAddress) > 0 {
		i -= len(m.AggregatorTargetAddress)
		copy(dAtA[i:], m.AggregatorTargetAddress)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.AggregatorTargetAddress)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.Aggregator) > 0 {
		i -= len(m.Aggregator)
		copy(dAtA[i:], m.Aggregator)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.Aggregator)))
		i--
		dAtA[i] = 0x4a
	}
	{
		size := m.AffiliateBasisPoints.Size()
		i -= size
		if _, err := m.AffiliateBasisPoints.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	if len(m.AffiliateAddress) > 0 {
		i -= len(m.AffiliateAddress)
		copy(dAtA[i:], m.AffiliateAddress)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.AffiliateAddress)))
		i--
		dAtA[i] = 0x3a
	}
	{
		size := m.MinOut.Size()
		i -= size
		if _, err := m.MinOut.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size, err := m.TargetAsset.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.TargetAddress) > 0 {
		i -= len(m.TargetAddress)
		copy(dAtA[i:], m.TargetAddress)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.TargetAddress)))
		i--
		dAtA[i] = 0x22
	}
	{
		size := m.CollateralAmount.Size()
		i -= size
		if _, err := m.CollateralAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.CollateralAsset.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgLoanRepayment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgLoanRepayment) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgLoanRepayment) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TxID) > 0 {
		i -= len(m.TxID)
		copy(dAtA[i:], m.TxID)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.TxID)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.From) > 0 {
		i -= len(m.From)
		copy(dAtA[i:], m.From)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.From)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0x2a
	}
	{
		size := m.MinOut.Size()
		i -= size
		if _, err := m.MinOut.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.Coin.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.CollateralAsset.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMsgLoan(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintMsgLoan(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMsgLoan(dAtA []byte, offset int, v uint64) int {
	offset -= sovMsgLoan(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgLoanOpen) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = m.CollateralAsset.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = m.CollateralAmount.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = len(m.TargetAddress)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = m.TargetAsset.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = m.MinOut.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = len(m.AffiliateAddress)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = m.AffiliateBasisPoints.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = len(m.Aggregator)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = len(m.AggregatorTargetAddress)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = m.AggregatorTargetLimit.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = len(m.TxID)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	return n
}

func (m *MsgLoanRepayment) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = m.CollateralAsset.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = m.Coin.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = m.MinOut.Size()
	n += 1 + l + sovMsgLoan(uint64(l))
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	l = len(m.TxID)
	if l > 0 {
		n += 1 + l + sovMsgLoan(uint64(l))
	}
	return n
}

func sovMsgLoan(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMsgLoan(x uint64) (n int) {
	return sovMsgLoan(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgLoanOpen) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgLoan
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgLoanOpen: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgLoanOpen: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = gitlab_com_thorchain_thornode_common.Address(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralAsset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralAsset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TargetAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TargetAddress = gitlab_com_thorchain_thornode_common.Address(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TargetAsset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TargetAsset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinOut", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AffiliateAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AffiliateAddress = gitlab_com_thorchain_thornode_common.Address(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AffiliateBasisPoints", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AffiliateBasisPoints.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Aggregator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Aggregator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AggregatorTargetAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AggregatorTargetAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AggregatorTargetLimit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AggregatorTargetLimit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = append(m.Signer[:0], dAtA[iNdEx:postIndex]...)
			if m.Signer == nil {
				m.Signer = []byte{}
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxID = gitlab_com_thorchain_thornode_common.TxID(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgLoan(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgLoanRepayment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgLoan
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgLoanRepayment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgLoanRepayment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = gitlab_com_thorchain_thornode_common.Address(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralAsset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralAsset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coin", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Coin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinOut", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = append(m.Signer[:0], dAtA[iNdEx:postIndex]...)
			if m.Signer == nil {
				m.Signer = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = gitlab_com_thorchain_thornode_common.Address(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgLoan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxID = gitlab_com_thorchain_thornode_common.TxID(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgLoan(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMsgLoan
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMsgLoan(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMsgLoan
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMsgLoan
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMsgLoan
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMsgLoan
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMsgLoan
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMsgLoan        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMsgLoan          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMsgLoan = fmt.Errorf("proto: unexpected end of group")
)
