package model

import "aim/kitex_gen/kitexmessageservice/kitexmessageservice"

type ServiceClient struct {
	MessageClient kitexmessageservice.Client
}
