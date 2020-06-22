package aop

import (
	"context"

	"github.com/xgxw/foundation-go/database"
	"github.com/xgxw/foundation-go/errors"
)

// 常量定义
const (
	Transaction = "transaction"
)

type Transaction interface {
	Set(context.Context, *database.DB) (context.Context, func(*error), *database.DB, error)
}

var defaultTransHandler = func(err *error) {}

// 当启用事务时, 返回的析构函数. 会对事务进行提交/回滚等操作.
func getTransHandler(tx *database.DB) func(*error) {
	return func(errAddr *error) {
		err := *errAddr
		if r := recover(); r != nil {
			tx.Rollback()
			panic(err)
		}
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
}

// 错误类型常量定义
var (
	TransactionUnformat    = "transaction_unformat"
	TransactionUnformatErr = &errors.Error{Code: TransactionUnformat, Msg: TransactionUnformat}
)

// SetTransactional 设置事务传播. 后续开发参照 Java/Spring @Transactional 注解 的理念
func SetTransactional(ctx context.Context, db *database.DB) (
	newCtx context.Context, transHandler func(*error), tx *database.DB, err error) {
	val := ctx.Value(Transaction)
	if val == nil {
		// 当为入口服务时, 需要执行 commit/Rollback 操作
		tx = db.Begin()
		transHandler = getTransHandler(tx)
	} else {
		// 当不是入口服务时, 不需要 commit/Rollback, 交由入口函数即可.
		transHandler = defaultTransHandler
		var ok bool
		if tx, ok = val.(*database.DB); !ok {
			return ctx, transHandler, nil, TransactionUnformatErr
		}
		// 开启SavePoint
		tx = tx.Begin()
	}
	newCtx = context.WithValue(ctx, Transaction, tx)
	return newCtx, transHandler, tx, nil
}
