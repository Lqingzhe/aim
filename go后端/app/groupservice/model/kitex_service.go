package model

import (
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"
	"aim/kitex_gen/kitexuserservice/kitexuserservice"
)

type ServiceClient struct {
	MessageClient kitexmessageservice.Client
	UserClient    kitexuserservice.Client
}
