# 切面方法

## SetTransactional|事务传播
SetTransactional 中会根据传入的 context 中是否已经开启事务来决定当前函数返回的 teardown 中是否提交事务, 无需调用者关注.

调用者只需在需要开启事务传播的位置插入以下代码即可, 需要传入 数据库连接db 和 上下文 ctx
```Go
ctx, teardown, tx, err := SetTransactional(ctx, db)
if err != nil {
	return err
}
defer teardown(tx, err)
```

### 使用示例
```Go
// SKU 测试结构体
type SKU struct {
	ID        int
	ProductID int
}

func create(ctx context.Context, sku *SKU) (err error) {
	// 首先执行断言
	if sku == nil || sku.ProductID == 0 {
		return errors.InvalidSourceErr
	}
	db, _ := database.NewDatabase(database.Options{})

	// 加入这一段代码即可, 程序自动负责开启事务, Commit/Rollback
	// 必须传入 ctx 和 db.
	ctx, teardown, tx, err := SetTransactional(ctx, db)
	if err != nil {
		return err
	}
	defer teardown(tx, err)

	// 要执行的代码
	// ....

	return nil
}
```
