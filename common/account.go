package common

// Account define a struct to hold account information across all chain
type Account struct {
	Sequence      int64
	AccountNumber int64
	Coins         Coins
	HasMemoFlag   bool
}

// NewAccount create a new instance of Account
func NewAccount(sequence, accountNumber int64, coins Coins, hasMemoFlag bool) Account {
	return Account{
		Sequence:      sequence,
		AccountNumber: accountNumber,
		Coins:         coins,
		HasMemoFlag:   hasMemoFlag,
	}
}
