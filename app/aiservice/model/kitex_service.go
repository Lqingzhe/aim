package model

import (
	"aim/kitex_gen/kitexgroupservice/kitexgroupservice"
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"
)

type ServiceClient struct {
	MessageService kitexmessageservice.Client
	GroupService   kitexgroupservice.Client
}
