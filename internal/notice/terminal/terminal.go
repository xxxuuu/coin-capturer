package terminal

import (
	"coin-capturer/internal/capturer"
	"fmt"
	"time"
)

type Terminal struct {
}

func New() *Terminal {
	return &Terminal{}
}

func (t *Terminal) OnTransfer(event *capturer.TransferEvent) {
	fmt.Printf(
		"[%s] Received transfer event - Tx Hash: %s, From: %s, To: %s, Value: %s USDT\n",
		time.Now().Format("2006-01-02 15:04:05"),
		event.TxHash.Hex(),
		event.From.Hex(),
		event.To.Hex(),
		event.Tokens,
	)
}
