package model

import (
	"aim/kitex_gen/kitexaiservice/kitexaiservice"
	"aim/kitex_gen/kitexfileservice/kitexfileservice"
	"aim/kitex_gen/kitexgroupservice/kitexgroupservice"
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"
	"aim/kitex_gen/kitexuserservice/kitexuserservice"
)

type ServiceClient struct {
	UserClient    kitexuserservice.Client
	GroupClient   kitexgroupservice.Client
	MessageClient kitexmessageservice.Client
	FileClient    kitexfileservice.Client
	AiClient      kitexaiservice.Client
}
