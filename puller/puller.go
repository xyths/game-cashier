package puller

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/xyths/game-cashier/client/dfuse/pb"
	"github.com/xyths/game-cashier/client/dfuse/subscriber"
	types2 "github.com/xyths/game-cashier/types"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

type Puller struct {
	Network string
	Key     string
	Address string
	query   string
	client  pb.GraphQLClient
	db      *mongo.Database
}

func (p *Puller) Init(network, apiKey, address string, db *mongo.Database) error {
	p.Network = network
	p.Key = apiKey
	p.Address = address
	p.db = db
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
	defer closeExecutor(executor)
	coll := p.db.Collection(types2.TransferCollName)
	if err != nil {
		log.Printf("error when Execute: %s", err)
		return err
	}
	for {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			return nil
		default:
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

			document := &EosioDocument{}

			//log.Printf("resp.Data: %s", resp.Data)

			if err = json.Unmarshal([]byte(resp.Data), document); err != nil {
				log.Printf("error when decode respone: %s", err)
				continue
			}

			result := document.SearchTransactionsForward
			reverted := ""
			if result.Undo {
				reverted = "REVERTED"
			}
			//fmt.Printf("Cursor: %s, Undo: %v\n", result.Cursor, result.Undo)
			//resultByte, _ := json.Marshal(result)
			//log.Printf("result: %s", string(resultByte))
			//if !result.IsIrreversible {
			//	log.Printf("the block's IsIrreversible=false, blockNumber=%d", result.Block.Num)
			//	continue
			//}
			if result.Trace.Status != StatusExecuted {
				log.Printf("tx %s status not %s, is %s", result.Trace.Id, StatusExecuted, result.Trace.Status)
				continue
			}
			for _, action := range result.Trace.MatchingActions {
				// action is an ActionTrace
				// https://docs.dfuse.io/reference/eosio/graphql/#actiontrace
				data := action.JSON
				// map[from:eidosonecoin memo:Refund EOS quantity:0.0020 EOS to:pptqipaelyog]
				// log.Printf("data map: %s", data)
				log.Printf("Transfer %s -> %s (%s)[%s]%s\n", data["from"], data["to"], data["memo"], data["quantity"], reverted)
				r := types2.TransferRecord{
					Id:          fmt.Sprintf("%s_%s", result.Trace.Id, action.Seq), // txid_seq
					Tx:          result.Trace.Id,
					BlockNumber: uint64(result.Block.Num),
					From:        fmt.Sprintf("%s", data["from"]),
					To:          fmt.Sprintf("%s", data["to"]),
					Amount:      toAmount(fmt.Sprintf("%s", data["quantity"])),
					Memo:        fmt.Sprintf("%s", data["memo"]),
					Timestamp:   result.Block.Timestamp.String(),
					TxTime:      result.Block.Timestamp,
					LogTime:     time.Now(),
					// no notifyTime here
				}
				log.Printf("got record: %v", r)
				//_ = coll
				if _, err := coll.InsertOne(ctx, r); err != nil {
					log.Printf("error when insert transfer record to mongo: %s", err)
				}
			}
		}
	}
	return nil
}

func (p *Puller) makeQuery() {
	const format = `subscription {
  searchTransactionsForward(query:"receiver:eosio.token action:transfer account:eosio.token receiver:eosio.token data.to:'%s'") {
    undo cursor isIrreversible irreversibleBlockNum
    block { id num timestamp }
    trace { id status scheduled matchingActions { seq json } }
  }
}`
	//p.query = fmt.Sprintf(format, p.Address)
	p.query = fmt.Sprintf(format, "eidosonecoin")
	log.Printf("query is: \n%s", p.query)

	//	p.query=`subscription {
	//  searchTransactionsForward(query:"receiver:eosio.token action:transfer -data.quantity:'0.0001 EOS'") {
	//    undo cursor
	//    trace { id matchingActions { json } }
	//  }
	//}`
}

func closeExecutor(e pb.GraphQL_ExecuteClient) {
	if e == nil {
		return
	}

	if err := e.CloseSend(); err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("connection closed.")
}

func toAmount(quantity string) float64 {
	tokens := strings.Split(quantity, " ")
	if len(tokens) != 2 || tokens[1] != "EOS" {
		log.Printf("bad format or not EOS transfer, should never happen: %s", quantity)
		return 0
	}
	if amount, ok := big.NewFloat(0).SetString(tokens[0]); ok {
		a, _ := amount.Float64()
		return a
	} else {
		return 0
	}
}
