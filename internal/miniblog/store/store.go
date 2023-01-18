package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once sync.Once
	//全局变量，方便其它包直接调用已初始化好的 S 实例
	S *datastore
)

type IStore interface {
	Users() UserStore
}

type datastore struct {
	db *gorm.DB
}

// 确保 datastore 实现IStore 接口
var _ IStore = (*datastore)(nil)

// NewStore 创建一个IStore 类型的实例
func NewStore(db *gorm.DB) *datastore {
	//确保 S 只被初始化一次
	once.Do(func() {
		S = &datastore{
			db: db,
		}
	})
	return S
}

func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}
