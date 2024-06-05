package types

import (
	types "github.com/cosmos/cosmos-sdk/types"
	"gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type MimirTestSuite struct{}

var _ = Suite(&MimirTestSuite{})

func (MimirTestSuite) TestNodeMimir(c *C) {
	m := NodeMimirs{}
	acc1 := GetRandomBech32Addr()
	acc2 := GetRandomBech32Addr()
	acc3 := GetRandomBech32Addr()
	active := []cosmos.AccAddress{acc1, acc2, acc3}
	key1 := "foo"
	key2 := "bar"
	key3 := "baz"

	m.Set(key1, 1, acc1)
	m.Set(key1, 1, acc2)
	m.Set(key1, 1, acc3)
	m.Set(key2, 1, acc1)
	m.Set(key2, 2, acc2)
	m.Set(key1, 3, acc3)
	m.Set(key3, 4, acc1)
	m.Set(key3, 5, acc2)
	m.Set(key3, 5, acc3)

	// test key1
	val, ok := m.HasSuperMajority(key1, active)
	c.Check(val, Equals, int64(1))
	c.Check(ok, Equals, true)
	val, ok = m.HasSimpleMajority(key1, active)
	c.Check(val, Equals, int64(1))
	c.Check(ok, Equals, true)

	// test key2
	_, ok = m.HasSuperMajority(key2, active)
	c.Check(ok, Equals, false)
	_, ok = m.HasSimpleMajority(key2, active)
	c.Check(ok, Equals, false)

	// test key3
	val, ok = m.HasSuperMajority(key3, active)
	c.Check(val, Equals, int64(5))
	c.Check(ok, Equals, true)
	val, ok = m.HasSimpleMajority(key3, active)
	c.Check(val, Equals, int64(5))
	c.Check(ok, Equals, true)
}

func (s *MimirTestSuite) TestValueOfOperational(c *C) {
	var m NodeMimirs
	key := "testKey"
	active := []types.AccAddress{
		types.AccAddress("addr1"),
		types.AccAddress("addr2"),
		types.AccAddress("addr3"),
		types.AccAddress("addr4"),
	}

	// Test 1: Basic test, no duplicate votes, no tie
	m.Mimirs = []NodeMimir{
		{Key: key, Value: 100, Signer: types.AccAddress("addr1")},
		{Key: key, Value: 100, Signer: types.AccAddress("addr2")},
		{Key: key, Value: 200, Signer: types.AccAddress("addr3")},
	}
	c.Assert(m.ValueOfOperational(key, 2, active), Equals, int64(100))

	// Test 2: Duplicate votes from same node, should be ignored
	m.Mimirs = []NodeMimir{
		{Key: key, Value: 100, Signer: types.AccAddress("addr1")},
		{Key: key, Value: 100, Signer: types.AccAddress("addr1")},
		{Key: key, Value: 200, Signer: types.AccAddress("addr2")},
	}
	c.Assert(m.ValueOfOperational(key, 2, active), Equals, int64(-1))

	// Test 3: Tie scenario
	m.Mimirs = []NodeMimir{
		{Key: key, Value: 100, Signer: types.AccAddress("addr1")},
		{Key: key, Value: 100, Signer: types.AccAddress("addr2")},
		{Key: key, Value: 101, Signer: types.AccAddress("addr3")},
		{Key: key, Value: 101, Signer: types.AccAddress("addr4")},
	}
	c.Assert(m.ValueOfOperational(key, 2, active), Equals, int64(-1))

	// Test 4: Not meeting minVotes
	m.Mimirs = []NodeMimir{
		{Key: key, Value: 100, Signer: types.AccAddress("addr1")},
	}
	c.Assert(m.ValueOfOperational(key, 2, active), Equals, int64(-1))

	// Test 5: Non-active node votes, should be ignored
	m.Mimirs = []NodeMimir{
		{Key: key, Value: 100, Signer: types.AccAddress("addr1")},
		{Key: key, Value: 200, Signer: types.AccAddress("addr2")},
		{Key: key, Value: 200, Signer: types.AccAddress("addr3")},
		{Key: key, Value: 100, Signer: types.AccAddress("addr5")},
	}
	c.Assert(m.ValueOfOperational(key, 2, active), Equals, int64(200))
}

func (s *MimirTestSuite) TestValueOfEconomic(c *C) {
	addr1 := types.AccAddress([]byte("addr1"))
	addr2 := types.AccAddress([]byte("addr2"))
	addr3 := types.AccAddress([]byte("addr3"))
	addr4 := types.AccAddress([]byte("addr4"))

	key := "testKey"

	m := NodeMimirs{
		Mimirs: []NodeMimir{
			{Key: key, Value: 10, Signer: addr1},
			{Key: key, Value: 10, Signer: addr2},
			{Key: key, Value: 20, Signer: addr3},
		},
	}

	// Test for no supermajority, should return -1
	c.Assert(m.ValueOfEconomic(key, []types.AccAddress{addr1, addr3}), Equals, int64(-1))

	// Test for supermajority (2/3 vote for 10)
	m.Mimirs = append(m.Mimirs, NodeMimir{Key: key, Value: 10, Signer: addr4})
	c.Assert(m.ValueOfEconomic(key, []types.AccAddress{addr1, addr2, addr3, addr4}), Equals, int64(10))

	// Test for duplicate votes (should consider only one vote from addr1)
	m.Mimirs = append(m.Mimirs, NodeMimir{Key: key, Value: 20, Signer: addr1})
	c.Assert(m.ValueOfEconomic(key, []types.AccAddress{addr1, addr2, addr3, addr4}), Equals, int64(10))

	// Test for key not present
	c.Assert(m.ValueOfEconomic("testKey2", []types.AccAddress{addr1, addr2, addr3, addr4}), Equals, int64(-1))
}
