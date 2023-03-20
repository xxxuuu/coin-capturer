package capturer

import (
	"coin-capturer/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"net/http"
)

const TRANSACTION_LIST_API = "https://www.oklink.com/api/v5/explorer/address/transaction-list?chainShortName=bsc&address=%s&protocolType=token_20"

type Transaction struct {
	TxID                 string `json:"txId"`
	BlockHash            string `json:"blockHash"`
	Height               string `json:"height"`
	TransactionTime      string `json:"transactionTime"`
	From                 string `json:"from"`
	To                   string `json:"to"`
	Amount               string `json:"amount"`
	TransactionSymbol    string `json:"transactionSymbol"`
	TokenContractAddress string `json:"tokenContractAddress"`
	TxFee                string `json:"txFee"`
	State                string `json:"state"`
}

type TransactionListResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Page             string        `json:"page"`
		Limit            string        `json:"limit"`
		TotalPage        string        `json:"totalPage"`
		ChainFullName    string        `json:"chainFullName"`
		ChainShortName   string        `json:"chainShortName"`
		TransactionLists []Transaction `json:"transactionLists"`
	} `json:"data"`
}

func GetTransactionList(config *config.Config, addr common.Address) ([]Transaction, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(TRANSACTION_LIST_API, addr.String()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ok-Access-Key", config.OklinkToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	transactionResp := &TransactionListResponse{}
	if err := json.Unmarshal(body, transactionResp); err != nil {
		return nil, err
	}

	if transactionResp.Code != "0" {
		return nil, errors.New(fmt.Sprintf("oklink transaction list api response failed: %+v", transactionResp))
	}
	if len(transactionResp.Data) == 0 {
		return nil, errors.New(fmt.Sprintf("oklink transaction list api response data is empty: %+v", transactionResp))
	}

	return transactionResp.Data[0].TransactionLists, nil
}
