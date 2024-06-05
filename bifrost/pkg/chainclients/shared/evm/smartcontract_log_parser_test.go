package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/thorchain/thornode/bifrost/thorclient/types"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	types2 "gitlab.com/thorchain/thornode/x/thorchain/types"
	. "gopkg.in/check.v1"
)

type SmartContractLogParserTestSuite struct {
	abi *abi.ABI
}

var _ = Suite(&SmartContractLogParserTestSuite{})

func (t *SmartContractLogParserTestSuite) SetUpSuite(c *C) {
	vaultABI, _, err := GetContractABI(routerContractABI, erc20ContractABI)
	c.Assert(err, IsNil)
	t.abi = vaultABI
}

func mockIsValidContractAddr(addr *ecommon.Address, _ bool) bool {
	return addr.String() == "0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44"
}

const (
	errAssetToken = "0x983e2cC84Bb8eA7b75685F285A28Bde2b4D5aCDA" // nolint gosec:G101 not a hardcoded credential
	tknTestToken  = "0X3B7FA4DD21C6F9BA3CA375217EAD7CAB9D6BF483" // nolint gosec:G101 not a hardcoded credential
)

func mockAssetResolver(token string) (common.Asset, error) {
	if strings.EqualFold(token, NativeTokenAddr) {
		return common.ETHAsset, nil
	}
	if strings.EqualFold(token, errAssetToken) {
		return common.EmptyAsset, fmt.Errorf("fail to parse asset")
	}
	if strings.EqualFold(token, tknTestToken) {
		return common.NewAsset("ETH.TKN-0X3B7FA4DD21C6F9BA3CA375217EAD7CAB9D6BF483")
	}
	return common.NewAsset(token)
}

func mockTokenDecimalResolver(_ string) int64 {
	return 8
}

func mockAmountConverter(_ string, amt *big.Int) cosmos.Uint {
	return cosmos.NewUintFromBigInt(amt)
}

func (t *SmartContractLogParserTestSuite) getDepositEvent(smartContractAddr, to, asset string, amount *big.Int, memo string) *etypes.Log {
	evt, err := t.abi.EventByID(ecommon.HexToHash(depositEvent))
	if err != nil {
		return nil
	}
	depositData, err := evt.Inputs.NonIndexed().Pack(amount, memo)
	if err != nil {
		return nil
	}
	return &etypes.Log{
		Address: ecommon.HexToAddress(smartContractAddr),
		Topics: []ecommon.Hash{
			ecommon.HexToHash(depositEvent),
			ecommon.BytesToHash(ecommon.LeftPadBytes(ecommon.HexToAddress(to).Bytes(), ecommon.HashLength)),
			ecommon.BytesToHash(ecommon.LeftPadBytes(ecommon.HexToAddress(asset).Bytes(), ecommon.HashLength)),
		},
		Data: depositData,
	}
}

func (t *SmartContractLogParserTestSuite) getTransferOutEvent(smartContractAddr, vault, to, asset string, amount *big.Int, memo string) *etypes.Log {
	evt, err := t.abi.EventByID(ecommon.HexToHash(transferOutEvent))
	if err != nil {
		return nil
	}
	transferOutData, err := evt.Inputs.NonIndexed().Pack(ecommon.HexToAddress(asset), amount, memo)
	if err != nil {
		return nil
	}
	return &etypes.Log{
		Address: ecommon.HexToAddress(smartContractAddr),
		Topics: []ecommon.Hash{
			ecommon.HexToHash(transferOutEvent),
			ecommon.BytesToHash(ecommon.LeftPadBytes(ecommon.HexToAddress(vault).Bytes(), ecommon.HashLength)),
			ecommon.BytesToHash(ecommon.LeftPadBytes(ecommon.HexToAddress(to).Bytes(), ecommon.HashLength)),
		},
		Data: transferOutData,
	}
}

func (t *SmartContractLogParserTestSuite) getTransferAllowanceEvent(smartContractAddr, vault, to, asset string, amount *big.Int, memo string) *etypes.Log {
	evt, err := t.abi.EventByID(ecommon.HexToHash(transferAllowanceEvent))
	if err != nil {
		return nil
	}
	transferAllowanceData, err := evt.Inputs.NonIndexed().Pack(ecommon.HexToAddress(asset), amount, memo)
	if err != nil {
		return nil
	}
	return &etypes.Log{
		Address: ecommon.HexToAddress(smartContractAddr),
		Topics: []ecommon.Hash{
			ecommon.HexToHash(transferAllowanceEvent),
			ecommon.BytesToHash(ecommon.LeftPadBytes(ecommon.HexToAddress(vault).Bytes(), ecommon.HashLength)),
			ecommon.BytesToHash(ecommon.LeftPadBytes(ecommon.HexToAddress(to).Bytes(), ecommon.HashLength)),
		},
		Data: transferAllowanceData,
	}
}

func (t *SmartContractLogParserTestSuite) TestGetTxInItem_DepositEvents(c *C) {
	vaultABI, _, err := GetContractABI(routerContractABI, erc20ContractABI)
	c.Assert(err, IsNil)
	parser := NewSmartContractLogParser(mockIsValidContractAddr, mockAssetResolver, mockTokenDecimalResolver, mockAmountConverter, vaultABI, common.ETHAsset, 2)

	// when logs are empty
	isVaultTransfer, err := parser.GetTxInItem(nil, &types.TxInItem{})
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(err, IsNil)

	// when log is not emit by router contract
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		{
			Address: ecommon.HexToAddress("0xe17d9cf3620ea447eed9089a22096a63bbd63eb4"),
			Topics:  nil,
		},
	}, &types.TxInItem{})
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(err, IsNil)

	// Deposit with zero amount should be ignored
	txInItem := &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", NativeTokenAddr, big.NewInt(0), "ADD:ETH.ETH:tthor16xxn0cadruuw6a2qwpv35av0mehryvdzzjz3af"),
	}, txInItem)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(err, IsNil)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")

	// Deposit, invalid asset should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", errAssetToken, big.NewInt(1024000), "ADD:ETH.ETH:tthor16xxn0cadruuw6a2qwpv35av0mehryvdzzjz3af"),
	}, txInItem)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(err, IsNil)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")

	// Normal Deposit
	txInItem = &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", NativeTokenAddr, big.NewInt(1024000), "ADD:ETH.ETH:tthor16xxn0cadruuw6a2qwpv35av0mehryvdzzjz3af"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "0x6C4a2eEB8531E3C18BcA51104Df7eb2377708263")
	c.Assert(txInItem.Memo, Equals, "ADD:ETH.ETH:tthor16xxn0cadruuw6a2qwpv35av0mehryvdzzjz3af")
	c.Assert(txInItem.Coins.EqualsEx(common.NewCoins(common.NewCoin(common.ETHAsset, cosmos.NewUint(1024000)))), Equals, true)

	// multiple Deposit events , which has different to address should result in an error
	txInItem = &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			NativeTokenAddr,
			big.NewInt(1024000),
			"ADD:ETH.ETH:tthor16xxn0cadruuw6a2qwpv35av0mehryvdzzjz3af"),
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x6c4a2eeb8531e3c18bca51104df7eb2377708263",
			NativeTokenAddr,
			big.NewInt(1024000),
			"ADD:ETH.ETH:tthor16xxn0cadruuw6a2qwpv35av0mehryvdzzjz3af"),
	}, txInItem)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(err, NotNil)

	// multiple Deposit events , which has different memo should result in an error
	txInItem = &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			NativeTokenAddr,
			big.NewInt(1024000),
			"ADD:ETH.ETH:tthor16xxn0cadruuw6a2qwpv35av0mehryvdzzjz3af"),
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			NativeTokenAddr,
			big.NewInt(1024000),
			"whatever"),
	}, txInItem)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(err, NotNil)

	// multiple Deposit events , if one of the deposit event has an asset that can't be resolved , it should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			NativeTokenAddr,
			big.NewInt(1024000),
			"MIGRATE:1024"),
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			errAssetToken,
			big.NewInt(1024000),
			"MIGRATE:1024"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473")
	c.Assert(txInItem.Memo, Equals, "MIGRATE:1024")
	c.Assert(txInItem.Coins.EqualsEx(common.NewCoins(
		common.NewCoin(common.ETHAsset, cosmos.NewUint(1024000)),
	)), Equals, true)

	// multiple Deposit events , same
	txInItem = &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			NativeTokenAddr,
			big.NewInt(1024000),
			"MIGRATE:1024"),
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			"0X3B7FA4DD21C6F9BA3CA375217EAD7CAB9D6BF483",
			big.NewInt(2048000),
			"MIGRATE:1024"),
	}, txInItem)
	c.Assert(err, IsNil)
	tknAsset, err := common.NewAsset("ETH.TKN-0X3B7FA4DD21C6F9BA3CA375217EAD7CAB9D6BF483")
	c.Assert(err, IsNil)
	c.Assert(tknAsset.IsEmpty(), Equals, false)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473")
	c.Assert(txInItem.Memo, Equals, "MIGRATE:1024")
	c.Assert(txInItem.Coins.EqualsEx(common.NewCoins(
		common.NewCoin(common.ETHAsset, cosmos.NewUint(1024000)),
		common.NewCoin(tknAsset, cosmos.NewUint(2048000)),
	)), Equals, true)
}

func (t *SmartContractLogParserTestSuite) TestGetTxInItem_TransferOutEvents(c *C) {
	vaultABI, _, err := GetContractABI(routerContractABI, erc20ContractABI)
	c.Assert(err, IsNil)
	parser := NewSmartContractLogParser(mockIsValidContractAddr, mockAssetResolver, mockTokenDecimalResolver, mockAmountConverter, vaultABI, common.ETHAsset, 2)
	// corrupted transferOutEvent should be ignored
	txInItem := &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}
	logItem := t.getTransferOutEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", tknTestToken, big.NewInt(1024000), "OUT:"+types2.GetRandomTxHash().String())
	logItem.Data = []byte("whatever")
	isVaultTransfer, err := parser.GetTxInItem([]*etypes.Log{
		logItem,
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")

	// invalid asset in transfer out should result in an error
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferOutEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", errAssetToken, big.NewInt(1024000), "OUT:"+types2.GetRandomTxHash().String()),
	}, txInItem)
	c.Assert(err, NotNil)
	c.Assert(isVaultTransfer, Equals, false)

	// invalid memo in transfer out should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferOutEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", NativeTokenAddr, big.NewInt(1024000), "whatever"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")
	c.Assert(txInItem.Coins.IsEmpty(), Equals, true)
	// incorrect memo in transfer out should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferOutEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", NativeTokenAddr, big.NewInt(1024000), "ADD:ETH.ETH"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")
	c.Assert(txInItem.Coins.IsEmpty(), Equals, true)
}

func (t *SmartContractLogParserTestSuite) TestGetTxInItem_TransferAllowance(c *C) {
	vaultABI, _, err := GetContractABI(routerContractABI, erc20ContractABI)
	c.Assert(err, IsNil)
	parser := NewSmartContractLogParser(mockIsValidContractAddr, mockAssetResolver, mockTokenDecimalResolver, mockAmountConverter, vaultABI, common.ETHAsset, 2)
	// corrupted transferAllowance should be ignored
	txInItem := &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}
	logItem := t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", tknTestToken, big.NewInt(1024000), "MIGRATE:1024")
	logItem.Data = []byte("whatever")
	isVaultTransfer, err := parser.GetTxInItem([]*etypes.Log{
		logItem,
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")

	// invalid asset in transferAllowance should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", errAssetToken, big.NewInt(1024000), "MIGRATE:1024"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")

	// TransferAllowance with zero amount should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", NativeTokenAddr, big.NewInt(0), "MIGRATE:1024"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")

	// TransferAllowance with different sender, should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", NativeTokenAddr, big.NewInt(102400), "MIGRATE:1024"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")

	// TransferAllowance with different to address, should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
		To:     "0x9EcA25ee04FDCc9d9CDFF377aa8da019Dba38437",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", NativeTokenAddr, big.NewInt(102400), "MIGRATE:1024"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "0x9EcA25ee04FDCc9d9CDFF377aa8da019Dba38437")
	c.Assert(txInItem.Memo, Equals, "")

	// TransferAllowance with different sender, should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", NativeTokenAddr, big.NewInt(102400), "MIGRATE:1024"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "")
	c.Assert(txInItem.Memo, Equals, "")

	// TransferAllowance with different memo, should be ignored
	txInItem = &types.TxInItem{
		Sender: "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473",
	}

	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", NativeTokenAddr, big.NewInt(102400), "MIGRATE:1024"),
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x3fd2D4cE97B082d4BcE3f9fee2A3D60668D2f473", "0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", NativeTokenAddr, big.NewInt(102400), "whatever"),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "0x9EcA25ee04FDCc9d9CDFF377aa8da019Dba38437")
	c.Assert(txInItem.Memo, Equals, "MIGRATE:1024")
	c.Assert(txInItem.Coins.EqualsEx(common.NewCoins(
		common.NewCoin(common.ETHAsset, cosmos.NewUint(102400)),
	)), Equals, true)
}

func (t *SmartContractLogParserTestSuite) TestFakeDeposit(c *C) {
	vaultABI, _, err := GetContractABI(routerContractABI, erc20ContractABI)
	c.Assert(err, IsNil)
	parser := NewSmartContractLogParser(mockIsValidContractAddr, mockAssetResolver, mockTokenDecimalResolver, mockAmountConverter, vaultABI, common.ETHAsset, 2)
	// When user deposit , if user use a malicious contract to trigger transferAllowance before Deposit
	// it should fail , because memo is not the same
	txInItem := &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	isVaultTransfer, err := parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263",
			tknTestToken, big.NewInt(102400), "MIGRATE:1024"),
		t.getDepositEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x6c4a2eeb8531e3c18bca51104df7eb2377708263", tknTestToken, big.NewInt(204800), "SWAP:BTC.BTC:tb1q2mwvhprtsgtamrte23mptehnsjwe3a7j4yvdn4"),
	}, txInItem)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(err, NotNil)
}

func (t *SmartContractLogParserTestSuite) TestFakeTransferOut(c *C) {
	vaultABI, _, err := GetContractABI(routerContractABI, erc20ContractABI)
	c.Assert(err, IsNil)
	parser := NewSmartContractLogParser(mockIsValidContractAddr, mockAssetResolver, mockTokenDecimalResolver, mockAmountConverter, vaultABI, common.ETHAsset, 2)
	tknAsset, err := common.NewAsset("ETH.TKN-0X3B7FA4DD21C6F9BA3CA375217EAD7CAB9D6BF483")
	c.Assert(err, IsNil)

	// When transfer out , user put in a smart contract as destination address,
	// when the smart contract get called , it call in to Router's transferAllowance function , trigger a transferAllowance event
	// trying to trigger bifrost to ignore TransferOutEvent
	// the transferAllowance Event should be override by TransferOut, 1) Sender will not match, 2)
	txInItem := &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	outboundHash := types2.GetRandomTxHash().String()
	isVaultTransfer, err := parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			"0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", tknTestToken, big.NewInt(2048000), "MIGRATE:1024"),
		t.getTransferOutEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x6c4a2eeb8531e3c18bca51104df7eb2377708263",
			"0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", tknTestToken,
			big.NewInt(1024000), "OUT:"+outboundHash),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "0x9EcA25ee04FDCc9d9CDFF377aa8da019Dba38437")
	c.Assert(txInItem.Memo, Equals, "OUT:"+outboundHash)
	c.Assert(txInItem.Coins.EqualsEx(common.NewCoins(
		common.NewCoin(tknAsset, cosmos.NewUint(1024000)),
	)), Equals, true)

	// This test a scenario that user triggered both transfer allowance, and a transfer out event
	// transfer out event should be observed
	txInItem = &types.TxInItem{
		Sender: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
	}
	outboundHash = types2.GetRandomTxHash().String()
	isVaultTransfer, err = parser.GetTxInItem([]*etypes.Log{
		t.getTransferAllowanceEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44", "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			"0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", tknTestToken, big.NewInt(2048000), "MIGRATE:1024"),
		t.getTransferOutEvent("0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
			"0x6c4a2eeb8531e3c18bca51104df7eb2377708263",
			"0x9eca25ee04fdcc9d9cdff377aa8da019dba38437", tknTestToken,
			big.NewInt(1024000), "OUT:"+outboundHash),
	}, txInItem)
	c.Assert(err, IsNil)
	c.Assert(isVaultTransfer, Equals, false)
	c.Assert(txInItem.To, Equals, "0x9EcA25ee04FDCc9d9CDFF377aa8da019Dba38437")
	c.Assert(txInItem.Memo, Equals, "OUT:"+outboundHash)
	c.Assert(txInItem.Coins.EqualsEx(common.NewCoins(
		common.NewCoin(tknAsset, cosmos.NewUint(1024000)),
	)), Equals, true)
}
