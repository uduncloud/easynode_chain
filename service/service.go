package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode_chain/chain/ether"
	"github.com/uduncloud/easynode_chain/chain/tron"
	"github.com/uduncloud/easynode_chain/config"
	"io/ioutil"
	"math/rand"
	"strconv"
)

type Handler struct {
	nodeCluster map[int64][]*config.NodeCluster
	log         *xlog.XLog
}

func NewHandler(cluster map[int64][]*config.NodeCluster, xlog *xlog.XLog) *Handler {
	return &Handler{
		log:         xlog,
		nodeCluster: cluster,
	}
}

func (h *Handler) BalanceCluster(blockChain int64) *config.NodeCluster {
	cluster, ok := h.nodeCluster[blockChain]
	if !ok {
		//不存在节点
		return nil
	}
	//todo 后期重构节点筛选算法
	//根据 采集节点、任务节点的节点使用数据，综合判断出最佳节点
	//目前暂使用随机算法 找到节点
	if len(cluster) > 1 {
		l := len(cluster)
		return cluster[rand.Intn(l)]
	} else if len(cluster) == 1 {
		return cluster[0]
	} else {
		return nil
	}
}
func (h *Handler) SendReq(blockChain int64, reqBody string) (string, error) {
	cluster := h.BalanceCluster(blockChain)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	if blockChain == 200 {
		return ether.Eth_WriteMsgToChain(cluster.NodeUrl, cluster.NodeToken, reqBody)
	} else if blockChain == 205 {
		return tron.Eth_WriteMsgToChain(cluster.NodeUrl, cluster.NodeToken, reqBody)
	}

	return "", errors.New("blockChainCode is error")
}

func (h *Handler) GetBalance(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	addr := gjson.ParseBytes(b).Get("address").String()
	tag := gjson.ParseBytes(b).Get("tag").String()

	req := `
 {
     "id": 1,
     "jsonrpc": "2.0",
     "params": [
          "%v",
          "%v"
     ],
     "method": "eth_getBalance"
}
`
	req = fmt.Sprintf(req, addr, tag)

	res, err := h.SendReq(blockChainCode, req)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, res, ctx.Request.RequestURI)

}

// GetTokenBalance todo ERC20协议代币余额，后期补充
func (h *Handler) GetTokenBalance(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	addr := gjson.ParseBytes(b).Get("address").String()
	codeHash := gjson.ParseBytes(b).Get("codeHash").String()
	tag := gjson.ParseBytes(b).Get("tag").String()

	req := `
 {
     "id": 1,
     "jsonrpc": "2.0",
     "params": [
          "%v",
          "%v"
     ],
     "method": "eth_getBalance"
}
`
	req = fmt.Sprintf(req, addr, tag, codeHash)

	res, err := h.SendReq(blockChainCode, req)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, res, ctx.Request.RequestURI)

}

// GetNonce todo 仅适用于 ether,tron 暂不支持
func (h *Handler) GetNonce(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	addr := gjson.ParseBytes(b).Get("address").String()
	tag := gjson.ParseBytes(b).Get("tag").String() //pending,latest

	req := `
 {
     "id": 1,
     "jsonrpc": "2.0",
     "params": [
          "%v",
          "%v"
     ],
     "method": "eth_getTransactionCount"
}
`
	req = fmt.Sprintf(req, addr, tag)

	res, err := h.SendReq(blockChainCode, req)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, res, ctx.Request.RequestURI)

}

func (h *Handler) GetLatestBlock(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	req := `
 {
     "id": 1,
     "jsonrpc": "2.0",
     "method": "eth_blockNumber"
}
`
	res, err := h.SendReq(blockChainCode, req)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, res, ctx.Request.RequestURI)

}

func (h *Handler) SendRawTx(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	signedTx := gjson.ParseBytes(b).Get("signed_tx").String()

	req := `
{
     "id": 1,
     "jsonrpc": "2.0",
     "params": [
          "%v"
     ],
     "method": "eth_sendRawTransaction"
}
`
	req = fmt.Sprintf(req, signedTx)

	res, err := h.SendReq(blockChainCode, req)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, res, ctx.Request.RequestURI)
}

func (h *Handler) HandlerReq(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	res, err := h.SendReq(blockChainCode, string(b))
	if err != nil {
		h.Error(ctx, ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, res, ctx.Request.RequestURI)
}

const (
	SUCCESS = 0
	FAIL    = 1
)

func (h *Handler) Success(c *gin.Context, resp interface{}, path string) {
	h.log.Printf("path=%v,body=%v\n", path, resp)
	mp := make(map[string]interface{})
	mp["code"] = SUCCESS
	mp["data"] = resp
	c.JSON(200, mp)
}

func (h *Handler) Error(c *gin.Context, path string, err string) {
	h.log.Printf("path=%v,err=%v\n", path, err)
	mp := make(map[string]interface{})
	mp["code"] = FAIL
	mp["data"] = err
	c.JSON(200, mp)
}
