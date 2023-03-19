package erc20

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
)

func ParseTransfer(log types.Log) (*Transfer, error) {
	contractAbi, err := abi.JSON(strings.NewReader(ERC20_ABI))
	if err != nil {
		return nil, err
	}
	var transferEvent Transfer

	if err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", log.Data); err != nil {
		return nil, err
	}

	transferEvent.From = common.HexToAddress(log.Topics[1].Hex())
	transferEvent.To = common.HexToAddress(log.Topics[2].Hex())

	return &transferEvent, nil
}
