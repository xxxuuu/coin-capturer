package dingtalk

import (
	"bytes"
	"coin-capturer/internal/capturer"
	_ "embed"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
)

const URL = "https://oapi.dingtalk.com/robot/send?access_token=%s"

//go:embed post.json
var POST_CONTENT string

//go:embed msg.md
var MSG_CONTENT string

//go:embed coins-msg.md
var COINS_CONTENT string

type DingTalk struct {
	token string
}

func New(token string) *DingTalk {
	return &DingTalk{token: token}
}

func (d *DingTalk) OnTransfer(event *capturer.TransferEvent) {
	url := fmt.Sprintf(URL, d.token)

	// 筛选出最近有交易的币种
	balance := map[common.Address]capturer.Coin{}
	for _, c := range event.FromBalance {
		balance[c.Address] = c
	}
	recentTransferCoin := map[common.Address]capturer.Transaction{}
	for _, t := range event.FromTransactionList {
		recentTransferCoin[common.HexToAddress(t.TokenContractAddress)] = t
	}
	var coinsMsg string
	for addr := range recentTransferCoin {
		if _, exists := balance[addr]; exists {
			coinsMsg += fmt.Sprintf(COINS_CONTENT, balance[addr].Name, balance[addr].Balance, balance[addr].Address)
		} else {
			coinsMsg += fmt.Sprintf(COINS_CONTENT,
				recentTransferCoin[addr].TransactionSymbol, decimal.Zero, recentTransferCoin[addr].TokenContractAddress)
		}
	}

	post := fmt.Sprintf(POST_CONTENT,
		event.Tokens,
		fmt.Sprintf(MSG_CONTENT,
			event.TxHash,
			event.TxHash,
			event.TxHash,
			event.From,
			event.From,
			event.To,
			event.To,
			event.Tokens,
			coinsMsg,
		),
	)

	var jsonStr = []byte(post)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
}
