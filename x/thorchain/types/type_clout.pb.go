// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: thorchain/v1/x/thorchain/types/type_clout.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type SwapperClout struct {
	Address           gitlab_com_thorchain_thornode_common.Address `protobuf:"bytes,1,opt,name=address,proto3,casttype=gitlab.com/thorchain/thornode/common.Address" json:"address,omitempty"`
	Score             github_com_cosmos_cosmos_sdk_types.Uint      `protobuf:"bytes,2,opt,name=score,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"score"`
	Reclaimed         github_com_cosmos_cosmos_sdk_types.Uint      `protobuf:"bytes,3,opt,name=reclaimed,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"reclaimed"`
	Spent             github_com_cosmos_cosmos_sdk_types.Uint      `protobuf:"bytes,4,opt,name=spent,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"spent"`
	LastSpentHeight   int64                                        `protobuf:"varint,5,opt,name=last_spent_height,json=lastSpentHeight,proto3" json:"last_spent_height,omitempty"`
	LastReclaimHeight int64                                        `protobuf:"varint,6,opt,name=last_reclaim_height,json=lastReclaimHeight,proto3" json:"last_reclaim_height,omitempty"`
}

func (m *SwapperClout) Reset()         { *m = SwapperClout{} }
func (m *SwapperClout) String() string { return proto.CompactTextString(m) }
func (*SwapperClout) ProtoMessage()    {}
func (*SwapperClout) Descriptor() ([]byte, []int) {
	return fileDescriptor_84975f584ab6f822, []int{0}
}
func (m *SwapperClout) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SwapperClout) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SwapperClout.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SwapperClout) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwapperClout.Merge(m, src)
}
func (m *SwapperClout) XXX_Size() int {
	return m.Size()
}
func (m *SwapperClout) XXX_DiscardUnknown() {
	xxx_messageInfo_SwapperClout.DiscardUnknown(m)
}

var xxx_messageInfo_SwapperClout proto.InternalMessageInfo

func (m *SwapperClout) GetAddress() gitlab_com_thorchain_thornode_common.Address {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *SwapperClout) GetLastSpentHeight() int64 {
	if m != nil {
		return m.LastSpentHeight
	}
	return 0
}

func (m *SwapperClout) GetLastReclaimHeight() int64 {
	if m != nil {
		return m.LastReclaimHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*SwapperClout)(nil), "types.SwapperClout")
}

func init() {
	proto.RegisterFile("thorchain/v1/x/thorchain/types/type_clout.proto", fileDescriptor_84975f584ab6f822)
}

var fileDescriptor_84975f584ab6f822 = []byte{
	// 327 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2f, 0xc9, 0xc8, 0x2f,
	0x4a, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x2f, 0x33, 0xd4, 0xaf, 0x40, 0xe2, 0x96, 0x54, 0x16, 0xa4,
	0x16, 0x83, 0xc9, 0xf8, 0xe4, 0x9c, 0xfc, 0xd2, 0x12, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21,
	0x56, 0xb0, 0xb8, 0x94, 0x48, 0x7a, 0x7e, 0x7a, 0x3e, 0x58, 0x44, 0x1f, 0xc4, 0x82, 0x48, 0x2a,
	0x4d, 0x66, 0xe6, 0xe2, 0x09, 0x2e, 0x4f, 0x2c, 0x28, 0x48, 0x2d, 0x72, 0x06, 0xe9, 0x11, 0xf2,
	0xe2, 0x62, 0x4f, 0x4c, 0x49, 0x29, 0x4a, 0x2d, 0x2e, 0x96, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x74,
	0x32, 0xf8, 0x75, 0x4f, 0x5e, 0x27, 0x3d, 0xb3, 0x24, 0x27, 0x31, 0x49, 0x2f, 0x39, 0x3f, 0x17,
	0xd9, 0xbe, 0x8c, 0xfc, 0xa2, 0xbc, 0xfc, 0x94, 0x54, 0xfd, 0xe4, 0xfc, 0xdc, 0xdc, 0xfc, 0x3c,
	0x3d, 0x47, 0x88, 0xbe, 0x20, 0x98, 0x01, 0x42, 0xae, 0x5c, 0xac, 0xc5, 0xc9, 0xf9, 0x45, 0xa9,
	0x12, 0x4c, 0x60, 0x93, 0xf4, 0x4f, 0xdc, 0x93, 0x67, 0xb8, 0x75, 0x4f, 0x5e, 0x3d, 0x3d, 0xb3,
	0x24, 0xa3, 0x14, 0x62, 0x5a, 0x72, 0x7e, 0x71, 0x6e, 0x7e, 0x31, 0x94, 0xd2, 0x2d, 0x4e, 0xc9,
	0x86, 0xf8, 0x42, 0x2f, 0x34, 0x33, 0xaf, 0x24, 0x08, 0xa2, 0x5b, 0xc8, 0x97, 0x8b, 0xb3, 0x28,
	0x35, 0x39, 0x27, 0x31, 0x33, 0x37, 0x35, 0x45, 0x82, 0x99, 0x3c, 0xa3, 0x10, 0x26, 0x80, 0x5d,
	0x55, 0x90, 0x9a, 0x57, 0x22, 0xc1, 0x42, 0xae, 0xab, 0x40, 0xba, 0x85, 0xb4, 0xb8, 0x04, 0x73,
	0x12, 0x8b, 0x4b, 0xe2, 0xc1, 0xbc, 0xf8, 0x8c, 0xd4, 0xcc, 0xf4, 0x8c, 0x12, 0x09, 0x56, 0x05,
	0x46, 0x0d, 0xe6, 0x20, 0x7e, 0x90, 0x44, 0x30, 0x48, 0xdc, 0x03, 0x2c, 0x2c, 0xa4, 0xc7, 0x25,
	0x0c, 0x56, 0x0b, 0x75, 0x04, 0x4c, 0x35, 0x1b, 0x58, 0x35, 0xd8, 0x98, 0x20, 0x88, 0x0c, 0x44,
	0xbd, 0x93, 0xe7, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38,
	0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0xe9, 0xe3, 0x8f,
	0x09, 0x8c, 0xe4, 0x90, 0xc4, 0x06, 0x8e, 0x67, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb9,
	0x09, 0x72, 0xf1, 0x37, 0x02, 0x00, 0x00,
}

func (m *SwapperClout) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SwapperClout) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SwapperClout) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastReclaimHeight != 0 {
		i = encodeVarintTypeClout(dAtA, i, uint64(m.LastReclaimHeight))
		i--
		dAtA[i] = 0x30
	}
	if m.LastSpentHeight != 0 {
		i = encodeVarintTypeClout(dAtA, i, uint64(m.LastSpentHeight))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.Spent.Size()
		i -= size
		if _, err := m.Spent.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypeClout(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.Reclaimed.Size()
		i -= size
		if _, err := m.Reclaimed.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypeClout(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Score.Size()
		i -= size
		if _, err := m.Score.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypeClout(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTypeClout(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypeClout(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypeClout(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SwapperClout) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTypeClout(uint64(l))
	}
	l = m.Score.Size()
	n += 1 + l + sovTypeClout(uint64(l))
	l = m.Reclaimed.Size()
	n += 1 + l + sovTypeClout(uint64(l))
	l = m.Spent.Size()
	n += 1 + l + sovTypeClout(uint64(l))
	if m.LastSpentHeight != 0 {
		n += 1 + sovTypeClout(uint64(m.LastSpentHeight))
	}
	if m.LastReclaimHeight != 0 {
		n += 1 + sovTypeClout(uint64(m.LastReclaimHeight))
	}
	return n
}

func sovTypeClout(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypeClout(x uint64) (n int) {
	return sovTypeClout(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SwapperClout) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypeClout
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
			return fmt.Errorf("proto: SwapperClout: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SwapperClout: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypeClout
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
				return ErrInvalidLengthTypeClout
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypeClout
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = gitlab_com_thorchain_thornode_common.Address(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Score", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypeClout
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
				return ErrInvalidLengthTypeClout
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypeClout
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Score.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reclaimed", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypeClout
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
				return ErrInvalidLengthTypeClout
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypeClout
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Reclaimed.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Spent", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypeClout
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
				return ErrInvalidLengthTypeClout
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypeClout
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Spent.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastSpentHeight", wireType)
			}
			m.LastSpentHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypeClout
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastSpentHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastReclaimHeight", wireType)
			}
			m.LastReclaimHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypeClout
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastReclaimHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypeClout(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypeClout
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypeClout
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
func skipTypeClout(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypeClout
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
					return 0, ErrIntOverflowTypeClout
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
					return 0, ErrIntOverflowTypeClout
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
				return 0, ErrInvalidLengthTypeClout
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypeClout
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypeClout
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypeClout        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypeClout          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypeClout = fmt.Errorf("proto: unexpected end of group")
)