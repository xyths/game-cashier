module github.com/xyths/game-cashier

go 1.13

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.26.0

require (
	github.com/golang/protobuf v1.3.2
	github.com/stretchr/testify v1.4.0
	github.com/tidwall/gjson v1.3.5
	github.com/urfave/cli/v2 v2.1.1
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	google.golang.org/grpc v1.23.0
	gopkg.in/urfave/cli.v2 v2.0.0-20190806201727-b62605953717
)
