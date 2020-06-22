# 切面方法

## SetTransactional|事务传播
SetTransactional 中会根据传入的 context 中是否已经开启事务来决定当前函数返回的 teardown 中是否提交事务, 无需调用者关注.

当程序panic时, 事务回滚.

因为Go一般使用err作为错误返回, 所以增加: 根据 handler 传入err判断,
当 err!=nil 时回滚.

调用者只需在需要开启事务传播的位置插入以下代码即可, 需要传入 数据库连接db 和 上下文 ctx
```Go
ctx, transHandler, tx, err := SetTransactional(ctx, db)
if err != nil {
	return err
}
defer transHandler(&err)
```
