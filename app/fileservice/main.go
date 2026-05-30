package main

import (
	"aim/app/fileservice/config"
	"aim/app/fileservice/dao"
	"aim/app/fileservice/handler"
	"aim/app/fileservice/model"
	"aim/kitex_gen/kitexfileservice/kitexfileservice"
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
	commonconfig.AutoMysql(
		dbContext.Mysql,
		&model.FileModel{},
	)

	addr, err := net.ResolveTCPAddr("tcp", Config.ServiceConfig.ServiceAddr.Host+":"+strconv.FormatInt(Config.ServiceAddr.Port, 10))
	if err != nil {
		newlog.LogInitFatal(logger, err, "Make Addr Failed")
	}

	svr := kitexfileservice.NewServer(
		handler.NewFileServiceImpl(
			Config.FileConfig,
			snowNode,
			dbContext,
		),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "file_service",
			},
		),
		server.WithServiceAddr(addr),
		//server.WithServiceAddr(&net.TCPAddr{
		//	IP:   net.ParseIP(Config.ServiceConfig.Host),
		//	Port: int(Config.ServiceConfig.Port),
		//}),
		//commonconfig.RegisterService(
		//	Config.NacosConfig,
		//	logger,
		//),
	)
	err = svr.Run()
	if err != nil {
		newlog.LogInitFatal(logger, err, "Grcp Begin Error")
	}
}
