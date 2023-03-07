package tron

import (
	"log"
	"testing"
)

func TestEth_GetToken(t *testing.T) {
	log.Println(Eth_GetToken("grpc.trongrid.io:50051", "", "TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7", "TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9"))
}
