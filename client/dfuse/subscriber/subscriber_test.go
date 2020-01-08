package subscriber

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	pb "github.com/xyths/game-cashier/client/dfuse/pb"
	"github.com/xyths/game-cashier/client/dfuse/types"
	"testing"
)

const Query = `subscription {
  searchTransactionsForward(query:"receiver:eosio.token action:transfer account:eosio.token receiver:eosio.token (data.from:eidosonecoin OR data.to:eidosonecoin)") {
    undo cursor isIrreversible irreversibleBlockNum
    block { id num timestamp }
    trace { id status matchingActions { json } }
  }
}`

func TestCreateClient(t *testing.T) {
	ctx := context.Background()
	/* The client can be re-used for all requests, cache it at the appropriate level */
	client := CreateClient(Mainnet, "server_37083f4133041b5fd97a4c7740648d1c")

	executor, err := client.Execute(ctx, &pb.Request{Query: Query})
	require.NoError(t, err)
	defer executor.CloseSend()

	resp, err := executor.Recv()
	require.NoError(t, err)

	if len(resp.Errors) > 0 {
		for _, err := range resp.Errors {
			t.Logf("Request failed: %s\n", err)
		}

		/* We continue here, but you could take another decision here, like exiting the process */
		return
	}

	document := &types.EosioDocument{}
	err = json.Unmarshal([]byte(resp.Data), document)
	require.NoError(t, err)

	result := document.SearchTransactionsForward
	reverted := ""
	if result.Undo {
		reverted = " REVERTED"
	}
	t.Logf("Cursor: %s, Undo: %v\n", result.Cursor, result.Undo)
	for _, action := range result.Trace.MatchingActions {
		data := action.JSON
		t.Logf("Transfer %s -> %s (%s)[%s]%s\n", data["from"], data["to"], data["memo"], data["quantity"], reverted)
	}

}
