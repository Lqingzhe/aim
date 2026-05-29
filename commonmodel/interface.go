package commonmodel

import "context"

type DataModel interface {
	Data()
} //mysql自动迁移的约束接口.需要迁移的struct均需满足这个方法

type DBOperater interface {
	AddInfo(context.Context, any) error
	UpdateInfo(context.Context, any) (bool, error)
	DeleteInfo(context.Context, any) error
	GetInfo(context.Context, any) (bool, error)
}
