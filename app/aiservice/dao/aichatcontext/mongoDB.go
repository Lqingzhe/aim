package aichatcontext

import (
	"aim/app/aiservice/model"
	newerror "aim/pkg/error"
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setMongo(ctx context.Context, dbContext *model.DBContext, info *AiChatContext) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mongoDB:SetMongo")
	collection := dbContext.MongoDB.Client.Database("aim").Collection("ai_chat_context")
	_, err = collection.InsertOne(ctx, &info.MessageContext)
	if err != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
func getMongo(ctx context.Context, dbContext *model.DBContext, info *AiChatContext) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mongoDB:GetMongo")
	collection := dbContext.MongoDB.Client.Database("aim").Collection("ai_chat_context")
	err = collection.FindOne(ctx, bson.M{"_id": info.UserID}).Decode(info.Info)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil // 文档不存在，不是错误
		}
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return false, err2
		}
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return true, nil
}
func deleteMongo(ctx context.Context, dbContext *model.DBContext, info *AiChatContext) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mongoDB:DeleteMongo")
	collection := dbContext.MongoDB.Client.Database("aim").Collection("ai_chat_context")
	_, err = collection.DeleteOne(ctx, bson.M{"_id": info.UserID})
	if err != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
