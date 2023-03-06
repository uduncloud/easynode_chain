package tron

import (
	"errors"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Eth_WriteMsgToChain(host string, token string, query string) (string, error) {

	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("TRON_PRO_API_KEY", token)

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

func Eth_SendRawTx(host string, token string, query string) (string, error) {
	//host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", host, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("TRON_PRO_API_KEY", token)

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

	//todo 待设定
	//host = "grpc.trongrid.io:50051"
	conn := client.NewGrpcClient(host)
	_ = conn.SetAPIKey(key) // todo 没有发现设置意义
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{}, 2)
	balance, err := conn.TRC20ContractBalance(userAddress, contractAddress)

	if err != nil {
		log.Println("err=", err)
	} else {
		mp["balance"] = balance.String()
	}

	name, err := conn.TRC20GetName(contractAddress)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["name"] = name
	}

	symbol, err := conn.TRC20GetSymbol(contractAddress)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["symbol"] = symbol
	}

	decimals, err := conn.TRC20GetDecimals(contractAddress)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["decimals"] = decimals
	}
	return mp, nil
}
