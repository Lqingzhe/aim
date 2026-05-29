package message

import (
	"aim/app/messageservice/model"
	newerror "aim/pkg/error"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func addWhereInfo(info *Message) bson.M {
	newM := bson.M{}
	if info.whereWithUserID {
		newM["user_id"] = info.UserID
	}
	if info.whereWithGroupID {
		newM["group_id"] = info.GroupID
	}
	if info.whereWithMessageID {
		newM["_id"] = info.MessageID
	}
	if info.whereWithMessageTime {
		timeM := bson.M{}
		if info.findStartTimeSecond.Unix() != 0 {
			timeM["$gte"] = info.findStartTimeSecond.Unix()
		}
		if info.findEndTimeSecond.Unix() != 0 {
			timeM["$lte"] = info.findEndTimeSecond.Unix()
		}

		newM["send_time_second"] = timeM
	}
	return newM
}
func setMongo(ctx context.Context, dbContext *model.DBContext, info *Message) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mongoDB:SetMongo")
	collection := dbContext.MongoDB.Client.Database("aim").Collection("message")
	_, err = collection.InsertOne(ctx, &info.MessageInfo)
	if err != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
func getMongo(ctx context.Context, dbContext *model.DBContext, info *Message) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mongoDB:GetMongo")
	collection := dbContext.MongoDB.Client.Database("aim").Collection("message")
	course, err := collection.Find(ctx, addWhereInfo(info), options.Find().SetSort(bson.M{"_id": -1}).SetLimit(50))
	if err != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return false, err2
		}
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	defer course.Close(ctx)
	err = course.All(ctx, &info.Info)
	if err != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return false, err2
		}
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	if len(info.Info) == 0 {
		return false, nil
	}
	return true, nil
}
func deleteMongo(ctx context.Context, dbContext *model.DBContext, info *Message) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mongoDB:DeleteMongo")
	collection := dbContext.MongoDB.Client.Database("aim").Collection("message")
	_, err = collection.DeleteOne(ctx, addWhereInfo(info))
	if err != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
