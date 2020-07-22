package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"naga/config"
	"naga/logger"
	"naga/orm"
	"naga/routes"
)

func main() {
	// 获取运行的必要参数
	envName := flag.String("env", "", "set the program running env")
	flag.Parse()
	if *envName == "" {
		logrus.Panic("env params error")
	}
	// 读取配置文件
	config.Config.Start("config/" + *envName + ".yml")
	// 开启logger
	logger.Start()
	// 开启orm
	orm.Start()
	// gin api
	router := gin.New()
	routes.Start(router)
	logrus.Infof("asd")
	// 启动端口服务
	err := router.Run(config.Config.ServerAddr)
	if err != nil {
		logrus.Panicf("gin run err:%s", err.Error())
	}
}
