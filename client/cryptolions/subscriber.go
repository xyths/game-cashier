package cryptolions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xyths/game-cashier/types"
	"log"
	"net/http"
	"net/url"
)

type Cryptolions struct {
	Server  string
	Manager string
	Limit   uint
}

type Response struct {
	ActionCount int64    `json:"action_count"`
	TotalAmount float64  `json:"total_amount"`
	Actions     []Action `json:"actions"`
}

type Action struct {
	Act             Act      `json:"act"`
	Timestamp       string   `json:"@timestamp"`
	BlockNumber     uint64   `json:"block_num"`
	Producer        string   `json:"producer"`
	TxId            string   `json:"trx_id"`
	Global_sequence uint64   `json:"global_sequence"`
	Notified        []string `json:"notified"`
}

type Act struct {
	Authorization []Auth `json:"authorization"`
	Data          Data   `json:"data"`
	Account       string `json:"account"`
	Name          string `json:"name"`
}

type Auth struct {
	Actor      string `json:"actor"`
	Permission string `json:"permission"`
}

type Data struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Amount   float64 `json:"amount"`
	Symbol   string  `json:"symbol"`
	Quantity string  `json:"quantity"`
	Memo     string  `json:"memo"`
}

func (cl *Cryptolions) format(after, before string) string {
	// https://junglehistory.cryptolions.io/v2/history/get_transfers?symbol=EOS&contract=eosio.token&limit=10
	u := cl.Server + "/v2/history/get_transfers"
	//?symbol=EOS&contract=eosio.token
	v := url.Values{}
	v.Set("symbol", "EOS")
	v.Set("contract", "eosio.token")
	if cl.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", cl.Limit))
	}
	if cl.Manager != "" {
		v.Set("to", cl.Manager)
	}
	if after != "" {
		v.Set("after", after)
	}
	if before != "" {
		v.Set("before", before)
	}
	u += "?" + v.Encode()
	return u
}

// after before均为闭区间，搜索结果包含该时间节点，[after, before]
func (cl *Cryptolions) Pull(ctx context.Context, after, before string) (records []types.TransferRecord, err error) {
	url := cl.format(after, before)
	log.Printf("url=%s", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)

	res := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
		extractRecord(&res, &records)
		log.Printf("records size: %d", len(records))
	} else {
		log.Printf("error when decode response: %s", err)
	}
	return
}

func extractRecord(response *Response, records *[]types.TransferRecord) {
	for _, action := range response.Actions {
		r := types.TransferRecord{
			Timestamp:   action.Timestamp,
			Tx:          action.TxId,
			Id:          action.TxId,
			BlockNumber: action.BlockNumber,
			From:        action.Act.Data.From,
			To:          action.Act.Data.To,
			Amount:      action.Act.Data.Amount,
		}
		*records = append(*records, r)
	}
}
