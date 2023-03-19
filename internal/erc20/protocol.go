package erc20

import (
	_ "embed"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

//go:embed erc20.abi.json
var ERC20_ABI string

const TRANSFER_EVENT = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

type Transfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
