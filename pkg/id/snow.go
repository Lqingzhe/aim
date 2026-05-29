package id

import (
	newlog "aim/pkg/log"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

func InitSnowNode(equipID int, logger *zap.Logger) *snowflake.Node {
	newNode, err := snowflake.NewNode(int64(equipID))
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init SnowNode Failed")
	}
	return newNode
}
func MakeID(Node *snowflake.Node) int64 {
	return Node.Generate().Int64()
}
