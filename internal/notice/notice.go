package notice

import (
	"coin-capturer/internal/capturer"
)

type Notice interface {
	OnTransfer(event *capturer.TransferEvent)
}
