package aop

import (
	"context"

	"github.com/everywan/foundation-go/database/mysql"
	"github.com/everywan/foundation-go/errors"
)

const (
	Transaction = "transaction"
)

var (
	// 当启用事务时, 返回的析构函数. 会对事务进行提交/回滚等操作.
	teardownTrans = func(tx *mysql.DB, err *error) {
		if r := recover(); r != nil {
			*err = errors.InvalidSourceErr
			tx.Rollback()
			return
		}
		if *err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}

	teardownDefault = func(tx *mysql.DB, err *error) {}

	TransactionUnformat    = "transaction_unformat"
	TransactionUnformatErr = &errors.Error{Code: TransactionUnformat, Msg: TransactionUnformat}
)

// SetTransactional 设置事务传播. 后续开发参照 Java/Spring @Transactional 注解 的理念
func SetTransactional(ctx context.Context, db *mysql.DB) (newCtx context.Context, teardown func(*mysql.DB, *error), tx *mysql.DB, err error) {
	val := ctx.Value(Transaction)
	if val == nil {
		// 当为入口服务时, 需要执行 commit/Rollback 操作
		tx = db.Begin()
		teardown = teardownTrans
	} else {
		// 当不是入口服务时, 不需要 commit/Rollback, 交由入口函数即可.
		teardown = teardownDefault
		var ok bool
		if tx, ok = val.(*mysql.DB); !ok {
			return ctx, teardown, nil, TransactionUnformatErr
		}
		tx = tx.Begin()
	}
	newCtx = context.WithValue(ctx, Transaction, tx)
	return newCtx, teardown, tx, nil
}
