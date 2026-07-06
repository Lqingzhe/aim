package main

import (
	"aim/app/userservice/config"
	"aim/app/userservice/dao"
	"aim/app/userservice/handler"
	"aim/app/userservice/model"
	"aim/kitex_gen/kitexuserservice/kitexuserservice"
	commonconfig "aim/pkg/config"
	"aim/pkg/id"
	newlog "aim/pkg/log"
	"net"
	"strconv"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {

	Config := config.InitConfig()

	logger := newlog.InitLog(Config.Service, Config.EquipID)
	defer logger.Sync()

	snowNode := id.InitSnowNode(Config.EquipID, logger)

	dbContext := dao.InitDB(&Config.DBConfig.Mysql, logger)
	defer dao.CloseDB(dbContext)
	commonconfig.AutoMysql(dbContext.Mysql, &model.UserInfo{}, &model.UserLoginInfo{}, &model.RemarkInfo{})

	//listener, err := net.Listen("tcp", Config.ServiceConfig.Host+":"+strconv.FormatInt(Config.ServiceConfig.Port, 10))
	addr, err := net.ResolveTCPAddr("tcp", Config.ServiceConfig.ServiceAddr.Host+":"+strconv.FormatInt(Config.ServiceAddr.Port, 10))
	if err != nil {
		newlog.LogInitFatal(logger, err, "Make Addr Failed")
	}

	svr := kitexuserservice.NewServer(
		handler.NewUserServiceImpl(
			snowNode,
			dbContext,
			logger,
			Config.UserConfig,
		),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "user_service",
			},
		),
		server.WithServiceAddr(addr),
		commonconfig.RegisterService(
			Config.NacosConfig,
			logger,
		),
	)
	err = svr.Run()
	if err != nil {
		newlog.LogInitFatal(logger, err, "Grcp Begin Error")
	}
}
