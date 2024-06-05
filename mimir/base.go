package mimir

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

//go:generate stringer -type=MimirType
type MimirType uint8

const (
	UnknownMimir MimirType = iota
	EconomicMimir
	OperationalMimir
)

type Mimir interface {
	LegacyKey(_ string) string
	Tags() []string
	Description() string
	Name() string
	Reference() string
	DefaultValue() int64
	Type() MimirType
	FetchValue(_ cosmos.Context, _ keeper.Keeper) int64
	IsOn(_ cosmos.Context, _ keeper.Keeper) bool
	IsOff(_ cosmos.Context, _ keeper.Keeper) bool
}

type mimir struct {
	name           string
	defaultValue   int64
	reference      string
	id             Id
	mimirType      MimirType
	tags           []string
	description    string
	legacyMimirKey func(ref string) string // mimir v1 key/constant
}

func (m *mimir) LegacyKey(ref string) string {
	return strings.ToUpper(m.legacyMimirKey(ref))
}

func (m *mimir) Tags() []string {
	return m.tags
}

func (m *mimir) Description() string {
	return m.description
}

func (m *mimir) Name() string {
	return strings.ToUpper(fmt.Sprintf("%s%s", m.name, m.Reference()))
}

func (m *mimir) DefaultValue() int64 {
	return m.defaultValue
}

func (m *mimir) Type() MimirType {
	return m.mimirType
}

func (m *mimir) Reference() string {
	if m.reference == "" {
		return "GLOBAL"
	}
	return strings.ToUpper(m.reference)
}

func (m *mimir) key() string {
	return fmt.Sprintf("%d-%s", m.id, strings.ToUpper(m.Reference()))
}

func (m *mimir) FetchValue(ctx cosmos.Context, keeper keeper.Keeper) (value int64) {
	var err error
	version := keeper.GetVersion()
	// fetch mimir v2
	value = int64(-1)
	switch {
	case version.GTE(semver.MustParse("1.125.0")):
		value = m.fetchValueV125(ctx, keeper)
	case version.GTE(semver.MustParse("1.124.0")):
		value = m.fetchValueV124(ctx, keeper)
	}
	if value >= 0 {
		return value
	}
	legacyKey := m.LegacyKey(m.Reference())
	if version.GTE(semver.MustParse("1.125.0")) {
		// return if legacy key does not exist (case of v2 only mimir)
		if len(legacyKey) == 0 {
			return m.DefaultValue()
		}
	}
	// fetch legacy mimir (v1)
	value, err = keeper.GetMimir(ctx, legacyKey)
	if err != nil {
		ctx.Logger().Error("failed to get mimir V1", "error", err)
		return -1
	}
	if value >= 0 {
		return value
	}

	// use default
	return m.DefaultValue()
}

func (m *mimir) fetchValueV125(ctx cosmos.Context, keeper keeper.Keeper) (value int64) {
	var (
		err    error
		active types.NodeAccounts
		key    string
	)
	active, err = keeper.ListActiveValidators(ctx)
	if err != nil {
		ctx.Logger().Error("failed to get active validator set", "error", err)
	}

	key = m.key()
	var mimirs types.NodeMimirs
	mimirs, err = keeper.GetNodeMimirsV2(ctx, key)
	if err != nil {
		ctx.Logger().Error("failed to get node mimir v2", "error", err)
	}
	value = int64(-1)
	switch m.Type() {
	case EconomicMimir:
		value = mimirs.ValueOfEconomic(key, active.GetNodeAddresses())
		if value < 0 {
			// no value, fallback to last economic value (if present)
			value, err = keeper.GetMimirV2(ctx, key)
			if err != nil {
				ctx.Logger().Error("failed to get mimir v2", "error", err)
			}
		}
	case OperationalMimir:
		value = mimirs.ValueOfOperational(key, constants.MinMimirV2Vote, active.GetNodeAddresses())
	}
	return
}

func (m *mimir) IsOn(ctx cosmos.Context, keeper keeper.Keeper) bool {
	value := m.FetchValue(ctx, keeper)
	return value > 0
}

func (m *mimir) IsOff(ctx cosmos.Context, keeper keeper.Keeper) bool {
	value := m.FetchValue(ctx, keeper)
	return value <= 0
}
