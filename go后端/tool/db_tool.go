package tool

import "aim/commonmodel"

func BeginMysqlTransaction(dbContext *commonmodel.MysqlContext) *commonmodel.MysqlContext {
	return &commonmodel.MysqlContext{
		Client: dbContext.Client.Begin(),
	}
}
