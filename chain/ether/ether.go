package ether

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode_chain/chain/token"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func Eth_WriteMsgToChain(host string, token string, query string) (string, error) {

	start := time.Now()
	defer func() {
		log.Printf("Eth_GetBlockByHash | duration=%v", time.Now().Sub(start))
	}()

	if len(token) > 1 {
		host = fmt.Sprintf("%v/%v", host, token)
	}
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	//req.Header.Add("Postman-Token", "181e4572-a9db-453a-b7d4-17974f785de0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if gjson.ParseBytes(body).Get("error").Exists() {
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func Eth_GetToken(host string, key string, contractAddress string, userAddress string) (map[string]interface{}, error) {
	if len(key) > 1 {
		host = fmt.Sprintf("%v/%v", host, key)
	}
	client, err := ethclient.Dial(host)
	if err != nil {
		return nil, err
	}

	// Golem (GNT) Address
	tokenAddress := common.HexToAddress(contractAddress)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(userAddress)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{}, 2)
	mp["balance"] = bal.String()

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["name"] = name
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["symbol"] = symbol
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["decimals"] = decimals
	}
	return mp, nil
}
