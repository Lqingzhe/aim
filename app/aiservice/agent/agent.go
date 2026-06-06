package agent

import (
	"aim/app/aiservice/model"
	newerror "aim/pkg/error"
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func CreateAiAgent(ctx context.Context, ModelName string, BaseURL string, APIKey string, toolInfo []tool.BaseTool, maxStep int) (agent *react.Agent, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("agent:CreateAiAgent")
	aiModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   ModelName,
		BaseURL: BaseURL,
		APIKey:  APIKey,
	})
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "no such host"):
			return nil, newerror.MakeError(http.StatusServiceUnavailable, newerror.CodeServiceUnavailable, "The BaseURL Analyses Error, Please Try Again Later", err, newerror.LevelInfo)
		case strings.Contains(errMsg, "connection refused"):
			return nil, newerror.MakeError(http.StatusServiceUnavailable, newerror.CodeServiceUnavailable, "Connect Error, Please Try Again Later", err, newerror.LevelInfo)
		case strings.Contains(errMsg, "timeout"):
			return nil, newerror.MakeError(http.StatusGatewayTimeout, newerror.CodeUpstreamTimeout, "AI Service Request Timeout, Please Try Again Later", err, newerror.LevelError)
		case strings.Contains(errMsg, "API key"):
			return nil, newerror.MakeError(http.StatusUnauthorized, newerror.CodeUnauthorized, "API Key Error", err, newerror.LevelInfo)
		case strings.Contains(errMsg, "model"):
			return nil, newerror.MakeError(http.StatusBadRequest, newerror.CodeResourceNotFound, "AI Model Is Not Exist", err, newerror.LevelInfo)
		case strings.Contains(errMsg, "certificate"):
			return nil, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDependencyError, "AI Service SSL Certificate Error", err, newerror.LevelInfo)
		case strings.Contains(errMsg, "rate limit"):
			return nil, newerror.MakeError(http.StatusTooManyRequests, newerror.CodeRateLimitExceeded, "AI Service Rate Limit Exceeded, Please Try Again Later", err, newerror.LevelInfo)
		default:
			return nil, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDependencyError, "AI Service Error", err, newerror.LevelWarn)
		}
	}
	agent, err = react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: aiModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: toolInfo,
		},
		MaxStep: maxStep,
	})
	if err != nil {
		return nil, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "AI Service Error", err, newerror.LevelWarn)
	}
	return agent, err
}
func CreateHistory(maxChatTurns int) []*schema.Message {
	return make([]*schema.Message, 0, (maxChatTurns+1)*2)
}
func CreateFormate(role string, agentPrompt string) *prompt.DefaultChatTemplate {
	systemMessage := "你是一个" + role + "," + agentPrompt
	return prompt.FromMessages(schema.FString, schema.SystemMessage(systemMessage), schema.UserMessage("{user_message}"), schema.MessagesPlaceholder("chat_history", true))
}
func TranslateHistoryToContext(messageInfo []*model.Message) (context []*schema.Message) {
	context = make([]*schema.Message, 0, len(messageInfo))
	for _, msg := range messageInfo {
		if msg.IsAi {
			context = append(context, schema.AssistantMessage(msg.Msg, msg.ToolCalls))
		} else {
			context = append(context, schema.UserMessage(msg.Msg))
		}
	}
	return context
}
func TranslateContextToHistory(context []*schema.Message) (history []*model.Message) {
	history = make([]*model.Message, 0, len(context))
	for _, msg := range context {
		if msg.Role == "user" {
			history = append(history, &model.Message{IsAi: false, ToolCalls: nil, Msg: msg.Content})
		} else {
			history = append(history, &model.Message{IsAi: true, ToolCalls: msg.ToolCalls, Msg: msg.Content})
		}
	}
	return history
}
func CreateMessage(ctx context.Context, formate *prompt.DefaultChatTemplate, userMessage string, history []*schema.Message) (message []*schema.Message, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("agent:CreateMessage")
	message, err = formate.Format(ctx, map[string]any{"user_message": userMessage, "chat_history": history})
	if err != nil {
		return nil, newerror.MakeError(-1, newerror.CodeInternalError, "Ai Service Error, Please Repeat Later", err, newerror.LevelWarn)
	}
	return message, nil
}
func AddHistory(history []*schema.Message, userMessage string, assistantMessage *schema.Message) (afterAddHistory []*schema.Message, addByteLength int64) {
	return append(history, schema.UserMessage(userMessage), assistantMessage), int64(len(userMessage) + len(assistantMessage.Content))
}
func CleanHistory(keepNumber int, history []*schema.Message, byteLength int64, maxByteLength int64) (ramInHistory []*schema.Message, remainByteLength int64) {
	length := len(history)
	if length <= 2*keepNumber && byteLength <= maxByteLength {
		return history, byteLength
	}
	var subByteLength int
	remainLength := length
	for remainLength <= remainLength && byteLength <= maxByteLength {
		remainLength -= 2
		subByteLength = len(history[length-subByteLength-1].Content) + len(history[length-subByteLength-2].Content)
		byteLength -= int64(subByteLength)
	}
	history = history[length-remainLength:]
	return history, byteLength
}
func AgentGenerateChat(ctx context.Context, message []*schema.Message, Agent *react.Agent, userID int64) (aiMessage *schema.Message, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("agent:AgentGenerateChat")
	ctx = context.WithValue(ctx, "user_id", userID)
	aiMessage, err = Agent.Generate(ctx, message)
	if err != nil {
		return nil, newerror.MakeError(-1, newerror.CodeThirdPartyError, "Ai Service Error, Please Repeat Later", err, newerror.LevelWarn)
	}
	return aiMessage, nil
}
