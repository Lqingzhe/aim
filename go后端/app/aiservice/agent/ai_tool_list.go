package agent

import (
	"aim/app/aiservice/dao"
	"aim/app/aiservice/dao/userprofile"
	"aim/app/aiservice/model"
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexmessageservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
	"time"

	"github.com/bytedance/sonic"
	agenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func InitTools(logger *zap.Logger, TraceWithUserManager *TraceWithUserManager, dbContext *model.DBContext, serviceClient model.ServiceClient) []agenttool.BaseTool {
	return []agenttool.BaseTool{
		NewGetTimeNow(TraceWithUserManager, logger),
		NewGetUserProfile(TraceWithUserManager, logger, dbContext),
		NewSetUserProfile(TraceWithUserManager, logger, dbContext),
		NewGetGroupMessage(TraceWithUserManager, logger, dbContext, serviceClient),
	}
}

type GetTimeNow struct {
	TraceWithUserManager *TraceWithUserManager
	logger               *zap.Logger
}

func NewGetTimeNow(TraceWithUserManager *TraceWithUserManager, logger *zap.Logger) *GetTimeNow {
	return &GetTimeNow{
		TraceWithUserManager: TraceWithUserManager,
		logger:               logger,
	}
}
func (*GetTimeNow) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get-time-now",
		Desc: `获取当前系统时间。当用户询问"现在几点"、"当前时间"、"现在是什么时候"等问题时调用此工具。`,
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

type GetUserProfile struct {
	TraceWithUserManager *TraceWithUserManager
	logger               *zap.Logger
	dbContext            *model.DBContext
}

func NewGetUserProfile(TraceWithUserManager *TraceWithUserManager, logger *zap.Logger, dbContext *model.DBContext) *GetUserProfile {
	return &GetUserProfile{
		TraceWithUserManager: TraceWithUserManager,
		logger:               logger,
		dbContext:            dbContext,
	}
}
func (g *GetUserProfile) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get-user-profile",
		Desc: `获取【单个用户】的个人画像信息。

【适用场景 - 只有这些情况才能使用】
- 用户问："我的画像是什么"
- 用户问："我有什么特点"
- 用户问："我的兴趣爱好是什么"
- 用户问："帮我看看我的个人信息"
- 需要根据用户特征提供个性化推荐时

【绝对不适用场景 - 遇到以下问题禁止使用】
- ❌ 任何包含"群"、"群里"、"群聊"、"他们"、"大家"的问题
- ❌ "群里说了什么" → 这是群聊内容，不是用户画像
- ❌ "他们在讨论什么" → 这是群聊内容，不是用户画像
- ❌ "群里有消息吗" → 这是群聊内容，不是用户画像

【返回内容】
- 有画像：返回用户画像JSON（兴趣爱好、职业、偏好等）
- 无画像：返回"该用户目前没有用户画像"

【重要提醒】
此工具只返回【单个人】的信息，不返回群聊内容！`,
	}, nil
}
func (g *GetUserProfile) InvokableRun(ctx context.Context, _ string, _ ...agenttool.Option) (result string, err error) {
	now := time.Now()
	userID := ctx.Value("user_id").(int64)
	traceID := g.TraceWithUserManager.GetTraceID(userID)
	logger := newlog.AddTraceID(g.logger, traceID)
	userProfileStruct := userprofile.NewStruct(userID, "")
	exist, err := dao.Get(ctx, userProfileStruct, g.dbContext)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddLatencyAndTime(logger, now)
		newlog.Log(logger, err2.LogLevel, "tool:GetUserProfile")
		return "", err
	}
	if !exist {
		logger = newlog.AddLatencyAndTime(logger, now)
		newlog.Log(logger, newerror.LevelInfo, "tool:GetUserProfile")
		return "该用户目前没有用户画像", nil
	}
	logger = newlog.AddLatencyAndTime(logger, now)
	newlog.Log(logger, newerror.LevelInfo, "tool:GetUserProfile")
	return userProfileStruct.Info.Profile, nil
}

type SetUserProfile struct {
	TraceWithUserManager *TraceWithUserManager
	logger               *zap.Logger
	dbContext            *model.DBContext
}

func NewSetUserProfile(TraceWithUserManager *TraceWithUserManager, logger *zap.Logger, dbContext *model.DBContext) *SetUserProfile {
	return &SetUserProfile{
		TraceWithUserManager: TraceWithUserManager,
		logger:               logger,
		dbContext:            dbContext,
	}
}
func (g *SetUserProfile) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "set-user-profile",
		Desc: `设置或更新【当前用户】的个人画像信息。

【适用场景】
- 用户主动提供个人信息时："我喜欢看电影"、"我是程序员"
- 用户表达偏好时："我喜欢简洁的回复"
- 用户要求保存信息时："帮我记住我不喜欢吃辣"

【输入参数格式】
JSON字符串，例如：{"interests": ["篮球", "电影"], "occupation": "工程师"}

【返回值】
成功返回"success"，失败返回错误信息。

【注意】此工具只操作用户画像，与群聊无关。`,
	}, nil
}
func (g *SetUserProfile) InvokableRun(ctx context.Context, input string, _ ...agenttool.Option) (result string, err error) {
	now := time.Now()
	userID := ctx.Value("user_id").(int64)
	traceID := g.TraceWithUserManager.GetTraceID(userID)
	logger := newlog.AddTraceID(g.logger, traceID)
	userProfileStruct := userprofile.NewStruct(userID, input)
	exist, err := dao.Update(ctx, userProfileStruct, g.dbContext)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddLatencyAndTime(logger, now)
		newlog.Log(logger, err2.LogLevel, "tool:SetUserProfile")
		return "", err
	}
	if !exist {
		err = dao.Add(ctx, userProfileStruct, g.dbContext)
		if err != nil {
			err2 := newerror.TranslateError(err)
			logger = newlog.AddError(logger, err, err2.StatusCode)
			logger = newlog.AddLatencyAndTime(logger, now)
			newlog.Log(logger, err2.LogLevel, "tool:SetUserProfile")
			return "", err
		}
	}
	logger = newlog.AddLatencyAndTime(logger, now)
	newlog.Log(logger, newerror.LevelInfo, "tool:SetUserProfile")
	return "success", nil
}

type GetGroupMessage struct {
	TraceWithUserManager *TraceWithUserManager
	logger               *zap.Logger
	dbContext            *model.DBContext
	serviceClient        model.ServiceClient
}

func NewGetGroupMessage(TraceWithUserManager *TraceWithUserManager, logger *zap.Logger, dbContext *model.DBContext, serviceClient model.ServiceClient) *GetGroupMessage {
	return &GetGroupMessage{
		TraceWithUserManager: TraceWithUserManager,
		logger:               logger,
		dbContext:            dbContext,
		serviceClient:        serviceClient,
	}
}
func (g *GetGroupMessage) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get-group-message",
		Desc: `【最重要】这是【唯一】能查看群聊消息的工具！

【核心功能】获取群组内的历史聊天记录，查看群里的人说了什么。

【强制使用场景 - 遇到以下任何问题必须使用此工具】
- "群里说了什么" / "他们在说什么"
- "群里有消息吗" / "群里讨论了什么"
- "帮我看看群聊" / "大家说了什么"
- "回顾一下群里的聊天" / "群里有提到XX吗"
- "他们聊了什么" / "群内消息"

【关键词触发】
只要用户问题中包含以下任意关键词，【必须】使用此工具：
- 群、群里、群聊、群内
- 他们、大家、各位（在群聊上下文中）
- 讨论、说、聊、发消息（配合群聊语境）

【绝对禁止】
- ❌ 不要用 get-user-profile 回答群聊问题
- ❌ 不要用 set-user-profile 回答群聊问题
- ❌ 不要用 get-time-now 回答群聊问题

【返回内容】
群聊消息列表，每条消息包含：
- message_content: 消息内容
- user_id: 发送者
- is_ai: 是否AI消息
- send_time_second: 发送时间

【示例】
用户问："他们在群里面说什么呢" → 必须调用此工具
用户问："群里刚才聊了什么" → 必须调用此工具
用户问："大家有没有讨论项目" → 必须调用此工具

【失败处理】
如果无法获取消息，会返回提示信息。`,
	}, nil
}
func (g *GetGroupMessage) InvokableRun(ctx context.Context, input string, _ ...agenttool.Option) (result string, err error) {
	now := time.Now()
	userID := ctx.Value("user_id").(int64)
	groupID := ctx.Value("group_id").(int64)
	traceID := g.TraceWithUserManager.GetTraceID(userID)
	logger := newlog.AddTraceID(g.logger, traceID)
	kitexReq := kitexmessageservice.GetMessageListReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: traceID,
		},
		GroupId:       groupID,
		UserId:        userID,
		EndTimeSecond: time.Now().Unix(),
	}
	kitexResp, err := g.serviceClient.MessageService.GetMessageList(ctx, &kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(newerror.UnMarshalError(err))
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddLatencyAndTime(logger, now)
		newlog.Log(logger, err2.LogLevel, "tool:GetGroupMessage")
		return "无法获取到群聊内容，委婉地告诉用户自己暂时无法查看群聊", nil
	}
	logger = newlog.AddLatencyAndTime(logger, now)
	newlog.Log(logger, newerror.LevelInfo, "tool:GetGroupMessage")
	output, _ := sonic.Marshal(kitexResp)
	return string(output), nil
}
