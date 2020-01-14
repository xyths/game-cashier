package cryptolions

import (
	"context"
	"encoding/json"
	"github.com/xyths/game-cashier/types"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestResponse(t *testing.T) {
	cl := Cryptolions{Server: "https://junglehistory.cryptolions.io/v2/history/get_transfers?symbol=EOS&contract=eosio.token&limit=10"}

	resp, err := http.Get(cl.Server)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	t.Log("Response status:", resp.Status)

	res := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
		print(t, res)
	} else {
		t.Log(err)
	}
}

func print(t *testing.T, r Response) {
	t.Logf("action_count: %d", r.ActionCount)
	t.Logf("total_amount: %f", r.TotalAmount)
	t.Logf("action len: %d", len(r.Actions))
	for i, a := range r.Actions {
		t.Logf("[%d] %s: %s -> %s: %s", i,
			a.Timestamp, a.Act.Data.From, a.Act.Data.To, a.Act.Data.Quantity)
	}
}

func TestTime(t *testing.T) {
	now := time.Now().Format("2006-01-02T15:04:05.000-07:00")
	t.Log(now)
	v := url.Values{}
	v.Set("after", now)
	t.Logf("query is: %s", v.Encode())
}

func TestExtract(t *testing.T) {
	cl := Cryptolions{Server: "https://junglehistory.cryptolions.io/v2/history/get_transfers?symbol=EOS&contract=eosio.token&limit=10"}

	resp, err := http.Get(cl.Server)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	t.Log("Response status:", resp.Status)

	res := Response{}
	var records []types.TransferRecord
	if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
		t.Logf("records size should be: %d", len(res.Actions))
		print(t, res)
		extractRecord(&res, &records)
		t.Logf("records size: %d", len(records))
	} else {
		t.Log(err)
	}
}

func TestCryptolion_Pull(t *testing.T) {
	cl := Cryptolions{
		Server:  "https://junglehistory.cryptolions.io",
		Manager: "testtestaaa1",
		//Limit:   10,
	}
	after := "2019-06-04T00:00:00.000+08:00"
	before := "2019-06-05T00:00:00.000+08:00"
	records, _ := cl.Pull(context.TODO(), after, before)
	t.Logf("records size: %d", len(records))
}
