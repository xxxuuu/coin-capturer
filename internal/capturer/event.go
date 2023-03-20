package capturer

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type Coin struct {
	Name    string
	Address common.Address
	Balance decimal.Decimal
}

func (c Coin) String() string {
	return fmt.Sprintf("{name: %s, address: %s, balance: %s}", c.Name, c.Address, c.Balance)
}

type TransferEvent struct {
	TxHash              common.Hash
	From                common.Address
	FromBalance         []Coin
	FromTransactionList []Transaction
	To                  common.Address
	Tokens              decimal.Decimal
}
