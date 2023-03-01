package token

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestTRC20_Balance(t *testing.T) {
	trc20Contract := "TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7" // USDT
	address := "TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9"

	conn := client.NewGrpcClient("grpc.trongrid.io:50051")
	conn.SetAPIKey("244f918d-56b5-4a16-9665-9637598b1223")
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	balance, err := conn.TRC20ContractBalance(address, trc20Contract)
	
	log.Println(balance.String())

	log.Println(conn.TRC20GetDecimals(trc20Contract))
}
