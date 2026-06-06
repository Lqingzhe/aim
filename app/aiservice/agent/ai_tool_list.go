package agent

import (
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
	"time"

	agenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func InitTools(logger *zap.Logger, TraceWithUserManager *TraceWithUserManager) []agenttool.BaseTool {
	return []agenttool.BaseTool{
		NewGetTimeNow(logger, TraceWithUserManager),
	}
}

type GetTimeNow struct {
	TraceWithUserManager *TraceWithUserManager
	logger               *zap.Logger
}

func NewGetTimeNow(logger *zap.Logger, TraceWithUserManager *TraceWithUserManager) *GetTimeNow {
	return &GetTimeNow{
		TraceWithUserManager: TraceWithUserManager,
		logger:               logger,
	}
}
func (*GetTimeNow) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get-time-now",
		Desc: "这是一个用来输出当前时间的工具，当需要获取当前系统时间的时候调用这个工会",
	}, nil
}

func (g *GetTimeNow) InvokableRun(ctx context.Context, _ string, _ ...agenttool.Option) (result string, err error) {
	now := time.Now()
	userID := ctx.Value("user_id").(int64)
	traceID := g.TraceWithUserManager.GetTraceID(userID)
	logger := newlog.AddTraceID(g.logger, traceID)
	logger = newlog.AddLatencyAndTime(logger, now)
	newlog.Log(logger, newerror.LevelInfo, "tool:GetTimeNow")
	return time.Now().String(), nil
}

//type GetUserProfile struct {
//	dbContext *model.DBContext
//	logger    *zap.Logger
//}
//
//func (g *GetUserProfile) Info(_ context.Context) (*schema.ToolInfo, error) {
//	return &schema.ToolInfo{
//		Name: "get-user-profile",
//		Desc: "",
//	}, nil
//}
//func (g *GetUserProfile) InvokableRun(_ context.Context, _ string, _ ...agenttool.Option) (result string, err error) {
//
//}
