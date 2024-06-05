// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: thorchain/v1/x/thorchain/types/msg_add_liquidity.proto

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

type MsgAddLiquidity struct {
	Tx                   common.Tx                                     `protobuf:"bytes,1,opt,name=tx,proto3" json:"tx"`
	Asset                common.Asset                                  `protobuf:"bytes,2,opt,name=asset,proto3" json:"asset"`
	AssetAmount          github_com_cosmos_cosmos_sdk_types.Uint       `protobuf:"bytes,3,opt,name=asset_amount,json=assetAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"asset_amount"`
	RuneAmount           github_com_cosmos_cosmos_sdk_types.Uint       `protobuf:"bytes,4,opt,name=rune_amount,json=runeAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"rune_amount"`
	RuneAddress          gitlab_com_thorchain_thornode_common.Address  `protobuf:"bytes,5,opt,name=rune_address,json=runeAddress,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"rune_address,omitempty"`
	AssetAddress         gitlab_com_thorchain_thornode_common.Address  `protobuf:"bytes,6,opt,name=asset_address,json=assetAddress,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"asset_address,omitempty"`
	AffiliateAddress     gitlab_com_thorchain_thornode_common.Address  `protobuf:"bytes,7,opt,name=affiliate_address,json=affiliateAddress,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"affiliate_address,omitempty"`
	AffiliateBasisPoints github_com_cosmos_cosmos_sdk_types.Uint       `protobuf:"bytes,8,opt,name=affiliate_basis_points,json=affiliateBasisPoints,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"affiliate_basis_points"`
	Signer               github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,9,opt,name=signer,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"signer,omitempty"`
}

func (m *MsgAddLiquidity) Reset()         { *m = MsgAddLiquidity{} }
func (m *MsgAddLiquidity) String() string { return proto.CompactTextString(m) }
func (*MsgAddLiquidity) ProtoMessage()    {}
func (*MsgAddLiquidity) Descriptor() ([]byte, []int) {
	return fileDescriptor_916314bb889296b4, []int{0}
}
func (m *MsgAddLiquidity) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddLiquidity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddLiquidity.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddLiquidity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddLiquidity.Merge(m, src)
}
func (m *MsgAddLiquidity) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddLiquidity) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddLiquidity.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddLiquidity proto.InternalMessageInfo

func (m *MsgAddLiquidity) GetTx() common.Tx {
	if m != nil {
		return m.Tx
	}
	return common.Tx{}
}

func (m *MsgAddLiquidity) GetAsset() common.Asset {
	if m != nil {
		return m.Asset
	}
	return common.Asset{}
}

func (m *MsgAddLiquidity) GetRuneAddress() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.RuneAddress
	}
	return ""
}

func (m *MsgAddLiquidity) GetAssetAddress() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.AssetAddress
	}
	return ""
}

func (m *MsgAddLiquidity) GetAffiliateAddress() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.AffiliateAddress
	}
	return ""
}

func (m *MsgAddLiquidity) GetSigner() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Signer
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgAddLiquidity)(nil), "types.MsgAddLiquidity")
}

func init() {
	proto.RegisterFile("thorchain/v1/x/thorchain/types/msg_add_liquidity.proto", fileDescriptor_916314bb889296b4)
}

var fileDescriptor_916314bb889296b4 = []byte{
	// 424 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0xd3, 0xcf, 0x8a, 0xd3, 0x40,
	0x1c, 0x07, 0xf0, 0xa4, 0x6e, 0xaa, 0x3b, 0xed, 0xa2, 0x86, 0x45, 0xc2, 0x1e, 0x92, 0xe0, 0xc5,
	0x0a, 0x6e, 0xc7, 0x55, 0xf0, 0x9e, 0xdc, 0x16, 0x14, 0x96, 0xe8, 0x5e, 0x04, 0x09, 0xd3, 0x4c,
	0x9a, 0x0c, 0x36, 0x33, 0x35, 0x33, 0x91, 0xf4, 0x2d, 0x7c, 0x25, 0x6f, 0x3d, 0xf6, 0x28, 0x1e,
	0x82, 0xb4, 0x6f, 0xd1, 0x93, 0x64, 0x66, 0xd2, 0x2a, 0x82, 0x2c, 0x39, 0xcd, 0x9f, 0x7c, 0xe7,
	0xc3, 0x6f, 0x7e, 0x61, 0xc0, 0x1b, 0x91, 0xb3, 0x32, 0xc9, 0x11, 0xa1, 0xf0, 0xeb, 0x15, 0xac,
	0xe1, 0x71, 0x29, 0x56, 0xcb, 0x94, 0xc3, 0x82, 0x67, 0x31, 0xc2, 0x38, 0x5e, 0x90, 0x2f, 0x15,
	0xc1, 0x44, 0xac, 0xa6, 0xcb, 0x92, 0x09, 0x66, 0x5b, 0xf2, 0xf3, 0x85, 0xff, 0xd7, 0xf1, 0x84,
	0x15, 0x05, 0xa3, 0x7a, 0x50, 0xc1, 0x8b, 0xf3, 0x8c, 0x65, 0x4c, 0x4e, 0x61, 0x3b, 0x53, 0xbb,
	0x4f, 0xbf, 0x5b, 0xe0, 0xe1, 0x3b, 0x9e, 0x05, 0x18, 0xbf, 0xed, 0x60, 0xdb, 0x07, 0x03, 0x51,
	0x3b, 0xa6, 0x6f, 0x4e, 0x46, 0xaf, 0xc0, 0x54, 0x23, 0x1f, 0xea, 0xf0, 0x64, 0xdd, 0x78, 0x46,
	0x34, 0x10, 0xb5, 0xfd, 0x1c, 0x58, 0x88, 0xf3, 0x54, 0x38, 0x03, 0x19, 0x3a, 0xeb, 0x42, 0x41,
	0xbb, 0xa9, 0x73, 0x2a, 0x61, 0x47, 0x60, 0x2c, 0x27, 0x31, 0x2a, 0x58, 0x45, 0x85, 0x73, 0xcf,
	0x37, 0x27, 0xa7, 0x21, 0x6c, 0x23, 0x3f, 0x1b, 0xef, 0x59, 0x46, 0x44, 0x5e, 0xcd, 0xda, 0xf3,
	0x30, 0x61, 0xbc, 0x60, 0x5c, 0x0f, 0x97, 0x1c, 0x7f, 0x56, 0x37, 0x9f, 0xde, 0x12, 0x2a, 0xa2,
	0x91, 0x44, 0x02, 0x69, 0xd8, 0x37, 0x60, 0x54, 0x56, 0x34, 0xed, 0xc8, 0x93, 0x7e, 0x24, 0x68,
	0x0d, 0x2d, 0xbe, 0x07, 0x63, 0x25, 0x62, 0x5c, 0xa6, 0x9c, 0x3b, 0x96, 0x24, 0x5f, 0xee, 0x1b,
	0xef, 0x45, 0x46, 0xc4, 0x02, 0x29, 0xee, 0x8f, 0x7f, 0x92, 0xb3, 0x92, 0x32, 0x9c, 0x76, 0x2d,
	0x0e, 0xd4, 0xb9, 0x48, 0xd6, 0xa5, 0x17, 0xf6, 0x2d, 0x38, 0xd3, 0x57, 0xd7, 0xea, 0xb0, 0xa7,
	0xaa, 0x3a, 0xd8, 0xb1, 0x9f, 0xc0, 0x63, 0x34, 0x9f, 0x93, 0x05, 0x41, 0xe2, 0x58, 0xf0, 0xfd,
	0x9e, 0xf4, 0xa3, 0x03, 0xd5, 0xf1, 0x29, 0x78, 0x72, 0xe4, 0x67, 0x88, 0x13, 0x1e, 0x2f, 0x19,
	0xa1, 0x82, 0x3b, 0x0f, 0xfa, 0xf5, 0xf9, 0xfc, 0xc0, 0x85, 0xad, 0x76, 0x23, 0x31, 0xfb, 0x1a,
	0x0c, 0x39, 0xc9, 0x68, 0x5a, 0x3a, 0xa7, 0xbe, 0x39, 0x19, 0x87, 0x57, 0xfb, 0xc6, 0xbb, 0xbc,
	0x03, 0x19, 0x24, 0x49, 0x57, 0xbb, 0x06, 0xc2, 0xeb, 0xf5, 0xd6, 0x35, 0x37, 0x5b, 0xd7, 0xfc,
	0xb5, 0x75, 0xcd, 0x6f, 0x3b, 0xd7, 0xd8, 0xec, 0x5c, 0xe3, 0xc7, 0xce, 0x35, 0x3e, 0xc2, 0xff,
	0xf7, 0xe2, 0x9f, 0x57, 0x36, 0x1b, 0xca, 0x57, 0xf1, 0xfa, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xb2, 0xf5, 0x87, 0x68, 0x8e, 0x03, 0x00, 0x00,
}

func (m *MsgAddLiquidity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddLiquidity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddLiquidity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0x4a
	}
	{
		size := m.AffiliateBasisPoints.Size()
		i -= size
		if _, err := m.AffiliateBasisPoints.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	if len(m.AffiliateAddress) > 0 {
		i -= len(m.AffiliateAddress)
		copy(dAtA[i:], m.AffiliateAddress)
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(len(m.AffiliateAddress)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.AssetAddress) > 0 {
		i -= len(m.AssetAddress)
		copy(dAtA[i:], m.AssetAddress)
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(len(m.AssetAddress)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.RuneAddress) > 0 {
		i -= len(m.RuneAddress)
		copy(dAtA[i:], m.RuneAddress)
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(len(m.RuneAddress)))
		i--
		dAtA[i] = 0x2a
	}
	{
		size := m.RuneAmount.Size()
		i -= size
		if _, err := m.RuneAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.AssetAmount.Size()
		i -= size
		if _, err := m.AssetAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.Asset.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.Tx.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMsgAddLiquidity(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintMsgAddLiquidity(dAtA []byte, offset int, v uint64) int {
	offset -= sovMsgAddLiquidity(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgAddLiquidity) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Tx.Size()
	n += 1 + l + sovMsgAddLiquidity(uint64(l))
	l = m.Asset.Size()
	n += 1 + l + sovMsgAddLiquidity(uint64(l))
	l = m.AssetAmount.Size()
	n += 1 + l + sovMsgAddLiquidity(uint64(l))
	l = m.RuneAmount.Size()
	n += 1 + l + sovMsgAddLiquidity(uint64(l))
	l = len(m.RuneAddress)
	if l > 0 {
		n += 1 + l + sovMsgAddLiquidity(uint64(l))
	}
	l = len(m.AssetAddress)
	if l > 0 {
		n += 1 + l + sovMsgAddLiquidity(uint64(l))
	}
	l = len(m.AffiliateAddress)
	if l > 0 {
		n += 1 + l + sovMsgAddLiquidity(uint64(l))
	}
	l = m.AffiliateBasisPoints.Size()
	n += 1 + l + sovMsgAddLiquidity(uint64(l))
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovMsgAddLiquidity(uint64(l))
	}
	return n
}

func sovMsgAddLiquidity(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMsgAddLiquidity(x uint64) (n int) {
	return sovMsgAddLiquidity(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgAddLiquidity) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgAddLiquidity
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
			return fmt.Errorf("proto: MsgAddLiquidity: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddLiquidity: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tx", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Tx.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Asset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AssetAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RuneAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RuneAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RuneAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RuneAddress = gitlab_com_thorchain_thornode_common.Address(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AssetAddress = gitlab_com_thorchain_thornode_common.Address(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AffiliateAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
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
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
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
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgAddLiquidity
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
				return ErrInvalidLengthMsgAddLiquidity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgAddLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = append(m.Signer[:0], dAtA[iNdEx:postIndex]...)
			if m.Signer == nil {
				m.Signer = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgAddLiquidity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMsgAddLiquidity
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMsgAddLiquidity
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
func skipMsgAddLiquidity(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMsgAddLiquidity
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
					return 0, ErrIntOverflowMsgAddLiquidity
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
					return 0, ErrIntOverflowMsgAddLiquidity
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
				return 0, ErrInvalidLengthMsgAddLiquidity
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMsgAddLiquidity
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMsgAddLiquidity
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMsgAddLiquidity        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMsgAddLiquidity          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMsgAddLiquidity = fmt.Errorf("proto: unexpected end of group")
)