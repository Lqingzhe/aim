package model

import (
	"aim/kitex_gen/kitexfileservice/kitexfileservice"
	"aim/kitex_gen/kitexgroupservice/kitexgroupservice"
)

type ServiceClient struct {
	GroupService kitexgroupservice.Client
	FileService  kitexfileservice.Client
	//AiService    kitexaiservice.Client
}
