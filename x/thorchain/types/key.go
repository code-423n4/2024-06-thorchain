package types

const (
	// ModuleName name of THORChain module
	ModuleName = "thorchain"
	// DefaultCodespace is the same as ModuleName
	DefaultCodespace = ModuleName
	// ReserveName the module account name to keep reserve
	ReserveName = "reserve"
	// AsgardName the module account name to keep asgard fund
	AsgardName = "asgard"
	// BondName the name of account used to store bond
	BondName = "bond"
	// LendingName
	LendingName = "lending"
	// AffiliateCollectorName the name of the account used to store rune for affiliate fee swaps
	AffiliateCollectorName = "affiliate_collector"
	// TreasuryName the name of the account used for treasury governance
	TreasuryName = "treasury"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
	// RouterKey used in the RPC query
	RouterKey = ModuleName // this was defined in your key.go file
)
