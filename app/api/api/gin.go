package api

//2060391348787224576
//2060656008488820736
import (
	"aim/app/api/handler"
	"aim/app/api/middleware"
	"aim/app/api/model"
	"aim/commonmodel"
	newlog "aim/pkg/log"
	"time"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ApiConfig struct {
	logger            *zap.Logger
	snowNode          *snowflake.Node
	dbContext         *model.DBContext
	limiterConfig     commonmodel.LimiterConfig
	tokenConfig       commonmodel.TokenConfig
	serviceClient     model.ServiceClient
	equipID           int64
	RoutTimeOut       map[string]time.Duration
	websocketUpgrader websocket.Upgrader
	consumer          sarama.Consumer
}

func NewConfig(logger *zap.Logger, snowNode *snowflake.Node, dbContext *model.DBContext, limiterConfig commonmodel.LimiterConfig, tokenConfig commonmodel.TokenConfig, equipID int64, RoutTimeout map[string]time.Duration, serviceClient model.ServiceClient, websocketUpgrader websocket.Upgrader, consumer sarama.Consumer) *ApiConfig {
	return &ApiConfig{
		logger:            logger,
		snowNode:          snowNode,
		dbContext:         dbContext,
		limiterConfig:     limiterConfig,
		tokenConfig:       tokenConfig,
		serviceClient:     serviceClient,
		equipID:           equipID,
		RoutTimeOut:       RoutTimeout,
		websocketUpgrader: websocketUpgrader,
		consumer:          consumer,
	}
}
func (A *ApiConfig) Begin(port string) {
	handlerConfig := handler.NewHandlerConfig(A.snowNode, A.dbContext, A.tokenConfig, A.serviceClient, A.websocketUpgrader, A.consumer)
	handlerConfig.Consumer(A.logger, 75) //启动消息消费者

	g := gin.New()
	gin.SetMode(gin.ReleaseMode)
	g.MaxMultipartMemory = 10 << 20
	g.Use(gin.Recovery())

	g.Use(middleware.Log(A.logger), middleware.Cors())

	g.LoadHTMLGlob("templates/*.html")
	g.Static("/static", "./static")
	g.GET("/", handlerConfig.IndexPage)
	g.GET("/login", handlerConfig.LoginPage)

	g.GET("/ping", handlerConfig.Ping) //测试

	g.GET("/ws", middleware.SetTimeOut(A.RoutTimeOut["/ws"]), handlerConfig.ConnectWebsocket)
	needLogin := g.Group("")
	needLogin.Use(
		middleware.AnalyseToken(A.tokenConfig, A.dbContext),
		middleware.Limiter(A.dbContext, A.limiterConfig),
	)
	g.GET("/chat", handlerConfig.ChatPage)
	{
		user := needLogin.Group("/user")
		{
			g.POST("/user/register", middleware.SetTimeOut(A.RoutTimeOut["/user/register"]), handlerConfig.Register)
			g.POST("/user/login", middleware.SetTimeOut(A.RoutTimeOut["/user/login"]), handlerConfig.Login)
			g.POST("/user/refresh-token", middleware.SetTimeOut(A.RoutTimeOut["/user/refresh-token"]), handlerConfig.RefreshToken)
			user.POST("/logout-all-device", middleware.SetTimeOut(A.RoutTimeOut["/user/logout-all-device"]), handlerConfig.LogoutAll, handlerConfig.DisconnectWebsocketAllDevice())
			user.POST("/logout-a-device", middleware.SetTimeOut(A.RoutTimeOut["/user/logout-a-device"]), handlerConfig.LogoutOne, handlerConfig.DisconnectWebsocketOneDevice())

			user.POST("/get-user-info", middleware.SetTimeOut(A.RoutTimeOut["/user/get-user-info"]), handlerConfig.GetUserInfo)
			user.POST("/get-other-user-info", middleware.SetTimeOut(A.RoutTimeOut["/user/get-other-user-info"]), handlerConfig.GetOtherUserInfo)
			user.POST("/update-user-info", middleware.SetTimeOut(A.RoutTimeOut["/user/update-user-info"]), handlerConfig.UpdateUserInfo)
			user.POST("/remark", middleware.SetTimeOut(A.RoutTimeOut["/user/remark"]), handlerConfig.Remark)
		}
		group := needLogin.Group("/group")
		{
			{
				group.POST("/create-group", middleware.SetTimeOut(A.RoutTimeOut["/group/create-group"]), handlerConfig.CreateGroup)

				group.POST("/delete-group", middleware.SetTimeOut(A.RoutTimeOut["/group/delete-group"]), handlerConfig.DeleteGroup)
			}
			{
				group.POST("/search-group", middleware.SetTimeOut(A.RoutTimeOut["/group/search-group"]), handlerConfig.SearchGroup)
				group.POST("/set-group-apply", middleware.SetTimeOut(A.RoutTimeOut["/group/set-group-apply"]), handlerConfig.SetGroupApply)
				group.POST("/get-group-apply-list", middleware.SetTimeOut(A.RoutTimeOut["/group/get-group-apply-list"]), handlerConfig.GetGroupApplyList)
				group.POST("/agree-group-apply", middleware.SetTimeOut(A.RoutTimeOut["/group/agree-group-apply"]), handlerConfig.AgreeGroupApply)
				group.POST("/refuse-group-apply", middleware.SetTimeOut(A.RoutTimeOut["/group/refuse-group-apply"]), handlerConfig.RefuseGroupApply)
				group.POST("/leave-group", middleware.SetTimeOut(A.RoutTimeOut["/group/leave-group"]), handlerConfig.LeaveGroup)
			}
			{
				group.POST("/get-group-info", middleware.SetTimeOut(A.RoutTimeOut["/group/get-group-info"]), handlerConfig.GetGroupInfo)
				group.POST("/change-group-info", middleware.SetTimeOut(A.RoutTimeOut["/group/change-group-info"]), handlerConfig.ChangeGroupInfo)
				group.POST("/get-group-info-with-user", middleware.SetTimeOut(A.RoutTimeOut["/group/get-group-info-with-user"]), handlerConfig.GetGroupInfoWithUser)
				group.POST("/update-group-info-with-user", middleware.SetTimeOut(A.RoutTimeOut["/group/update-group-info-with-user"]), handlerConfig.UpdateGroupInfoWithUser)
				group.POST("/get-group-and-session-id", middleware.SetTimeOut(A.RoutTimeOut["/group/get-group-and-session-id"]), handlerConfig.GetGroupAndSessionID)
			}
			{
				group.POST("/transform-group-owner", middleware.SetTimeOut(A.RoutTimeOut["/group/transform-group-owner"]), handlerConfig.TransformGroupOwner)
				group.POST("/set-manager", middleware.SetTimeOut(A.RoutTimeOut["/group/set-manager"]), handlerConfig.SetManager)
				group.POST("/revoke-manager", middleware.SetTimeOut(A.RoutTimeOut["/group/revoke-manager"]), handlerConfig.RevokeManager)
				group.POST("/get-last-visit-time", middleware.SetTimeOut(A.RoutTimeOut["/group/get-last-visit-time"]), handlerConfig.GetLastVisitTime)
				group.POST("/set-last-visit-time", middleware.SetTimeOut(A.RoutTimeOut["/group/set-last-visit-time"]), handlerConfig.SetLastVisitTime)
				group.POST("/kick-out-group", middleware.SetTimeOut(A.RoutTimeOut["/group/kick-out-group"]), handlerConfig.KickOutGroup)
			}
			{
				group.POST("/apply-for-friend", middleware.SetTimeOut(A.RoutTimeOut["/group/apply-for-friend"]), handlerConfig.ApplyForFriend)
				group.POST("/refuse-friend-apply", middleware.SetTimeOut(A.RoutTimeOut["/group/refuse-friend-apply"]), handlerConfig.RefuseFriendApply)
				group.POST("/get-friend-apply-list", middleware.SetTimeOut(A.RoutTimeOut["/group/get-friend-apply-list"]), handlerConfig.GetFriendApplyList)
				group.POST("/creat-session", middleware.SetTimeOut(A.RoutTimeOut["/group/creat-session"]), handlerConfig.CreatSession)

				//添加功能：删除聊天记录
				group.POST("/delete-session", middleware.SetTimeOut(A.RoutTimeOut["/group/delete-session"]), handlerConfig.DeleteSession)
				group.POST("/get-friend-last-visit-time", middleware.SetTimeOut(A.RoutTimeOut["/group/get-friend-last-visit-time"]), handlerConfig.GetFriendLastVisitTime)
			}
			{
				group.POST("/set-mute", middleware.SetTimeOut(A.RoutTimeOut["/group/set-mute"]), handlerConfig.SetMute)
				group.POST("/release-mute", middleware.SetTimeOut(A.RoutTimeOut["/group/release-mute"]), handlerConfig.ReleaseMute)
			}
		}
		message := needLogin.Group("/message")
		{
			message.POST("/send-message", middleware.SetTimeOut(A.RoutTimeOut["/message/send-message"]), handlerConfig.SendMessage)
			message.POST("/send-file", middleware.SetTimeOut(A.RoutTimeOut["/message/send-file"]), handlerConfig.SendFile)
			message.POST("/send-voice", middleware.SetTimeOut(A.RoutTimeOut["/message/send-voice"]), handlerConfig.SendVoice)
			message.POST("/send-picture", middleware.SetTimeOut(A.RoutTimeOut["/message/send-picture"]), handlerConfig.SendPicture)
			message.POST("/withdraw-message", middleware.SetTimeOut(A.RoutTimeOut["/message/withdraw-message"]), handlerConfig.WithdrawMessage)
			message.POST("/get-message-list", middleware.SetTimeOut(A.RoutTimeOut["/message/get-message-list"]), handlerConfig.GetMessageList)
			message.POST("/get-new-message", middleware.SetTimeOut(A.RoutTimeOut["/message/get-new-message"]), handlerConfig.GetNewMessage)
			message.POST("/get-file-content", middleware.SetTimeOut(A.RoutTimeOut["/message/get-file-content"]), handlerConfig.GetFileContent)
			message.POST("/send-group-notice", middleware.SetTimeOut(A.RoutTimeOut["/message/send-group-notice"]), handlerConfig.SendGroupNotice)
		}
		ai := needLogin.Group("/ai")
		{
			ai.POST("/delete-chat-context", middleware.SetTimeOut(A.RoutTimeOut["/delete-chat-context"]), handlerConfig.DeleteChatContext)
			ai.POST("/get-ai-config", middleware.SetTimeOut(A.RoutTimeOut["/get-ai-config"]), handlerConfig.GetAiConfig)
			ai.POST("/update-ai-config", middleware.SetTimeOut(A.RoutTimeOut["/update-ai-config"]), handlerConfig.UpdateAiConfig)
			ai.POST("/delete-ai-config", middleware.SetTimeOut(A.RoutTimeOut["/delete-ai-config"]), handlerConfig.DeleteAiConfig)
		}
	}
	//路由注册

	err := g.Run(":" + port)
	if err != nil {
		newlog.LogInitFatal(A.logger, err, "http begin error")
		return
	}
}
