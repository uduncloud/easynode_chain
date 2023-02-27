package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode_chain/config"
	"github.com/uduncloud/easynode_chain/service"
	"log"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(1).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/blockchain", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := service.NewHandler(cfg.Cluster, xLog)
	root.POST("/:chain/send", srv.HandlerReq)

	root.POST("/:chain/account/balance", srv.GetBalance)
	root.POST("/:chain/account/tokenBalance", srv.GetTokenBalance)
	root.POST("/:chain/account/nonce", srv.GetNonce)
	root.POST("/:chain/block/latest", srv.GetLatestBlock)
	root.POST("/:chain/tx/sendRawTransaction", srv.SendRawTx)

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}
