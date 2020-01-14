package puller

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/xyths/game-cashier/client/dfuse/pb"
	"github.com/xyths/game-cashier/client/dfuse/subscriber"
	"github.com/xyths/game-cashier/client/dfuse/types"
	"log"
	"os"
)

type Puller struct {
	Network string
	Key     string
	Address string
	query   string
	client  pb.GraphQLClient
}

func (p *Puller) Init(network, apiKey, address string) error {
	p.Network = network
	p.Key = apiKey
	p.Address = address
	p.makeQuery()

	switch network {
	case "jungle":
		p.client = subscriber.CreateClient(subscriber.Jungle, "server_37083f4133041b5fd97a4c7740648d1c")
	case "mainnet":
		p.client = subscriber.CreateClient(subscriber.Mainnet, "server_37083f4133041b5fd97a4c7740648d1c")
	case "cryptokylin":
		p.client = subscriber.CreateClient(subscriber.CryptoKylin, "server_37083f4133041b5fd97a4c7740648d1c")
	default:
		log.Fatalf("unknown network type: %s", network)
	}

	return nil
}

func (p *Puller) Pull(ctx context.Context) error {
	executor, err := p.client.Execute(ctx, &pb.Request{Query: p.query})
	defer closeExcutor(executor)
	if err != nil {
		log.Printf("error when Execute: %s", err)
		return err
	}
	for {
		resp, err := executor.Recv()
		if err != nil {
			log.Printf("error when rev: %s", err)
			break
		}

		if len(resp.Errors) > 0 {
			for _, err := range resp.Errors {
				log.Printf("Request failed: %s\n", err)
			}

			/* We continue here, but you could take another decision here, like exiting the process */
			continue
		}

		document := &types.EosioDocument{}

		if err = json.Unmarshal([]byte(resp.Data), document); err != nil {
			log.Printf("error when decode respone: %s", err)
			continue
		}

		result := document.SearchTransactionsForward
		reverted := ""
		if result.Undo {
			reverted = " REVERTED"
		}
		fmt.Printf("Cursor: %s, Undo: %v\n", result.Cursor, result.Undo)
		for _, action := range result.Trace.MatchingActions {
			data := action.JSON
			fmt.Printf("Transfer %s -> %s (%s)[%s]%s\n", data["from"], data["to"], data["memo"], data["quantity"], reverted)
		}
	}
	return nil
}

func (p *Puller) makeQuery() {
	const format = `subscription {
  searchTransactionsForward(query:"receiver:eosio.token action:transfer account:eosio.token receiver:eosio.token (data.from:%s OR data.to:%s)") {
    undo cursor isIrreversible irreversibleBlockNum
    block { id num timestamp }
    trace { id status matchingActions { json } }
  }
}`
	p.query = fmt.Sprintf(format, p.Address, p.Address)
	log.Printf("query is: \n%s", p.query)
}

func closeExcutor(e pb.GraphQL_ExecuteClient) {
	if e == nil {
		return
	}

	if err := e.CloseSend(); err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("connection closed.")
}
