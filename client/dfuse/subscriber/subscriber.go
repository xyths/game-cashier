package subscriber

import (
	"bytes"
	"fmt"
	"github.com/tidwall/gjson"
	pb "github.com/xyths/game-cashier/client/dfuse/pb"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func getToken(apiKey string) (token string, expiration time.Time, err error) {
	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"api_key":"%s"}`, apiKey)))
	resp, err := http.Post("https://auth.dfuse.io/v1/auth/issue", "application/json", reqBody)
	if err != nil {
		err = fmt.Errorf("unable to obtain token: %s", err)
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("unable to obtain token, status not 200, got %d: %s", resp.StatusCode, reqBody.String())
		return
	}

	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		token = gjson.GetBytes(body, "token").String()
		expiration = time.Unix(gjson.GetBytes(body, "expires_at").Int(), 0)
	}
	return
}

func CreateClient(endpoint string, dfuseAPIKey string) pb.GraphQLClient {
	if dfuseAPIKey == "" {
		panic("you must specify a dfuse API key")
	}

	token, _, err := getToken(dfuseAPIKey)
	if err != nil {
		log.Fatalf("error when getToken: %s", err)
	}

	credential := oauth.NewOauthAccess(&oauth2.Token{AccessToken: token, TokenType: "Bearer"})
	transportCreds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(endpoint,
		grpc.WithPerRPCCredentials(credential),
		grpc.WithTransportCredentials(transportCreds),
	)
	if err != nil {
		log.Fatalf("error when Dial: %s", err)
	}

	return pb.NewGraphQLClient(conn)
}

const (
	Mainnet     = "mainnet.eos.dfuse.io:443"
	CryptoKylin = "kylin.eos.dfuse.io:443"
	Jungle      = "jungle.eos.dfuse.io:443"
)
