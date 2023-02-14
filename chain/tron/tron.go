package tron

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
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

	//url := "https://eth-mainnet.g.alchemy.com/v2/demo"

	host = fmt.Sprintf("%v/%v", host, "jsonrpc")
	//	query := `
	//{
	//    "id": 1,
	//    "jsonrpc": "2.0",
	//    "method": "eth_getBlockByHash",
	//    "params": [
	//        "%v",
	//        true
	//    ]
	//}
	//`
	//
	//	query = fmt.Sprintf(query, blockHash)
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
