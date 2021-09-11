package main

import (
	"database/sql"
	"errors"
)

// hw
// 我们在数据库操作的时候，
// 比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，
// 是否应该 Wrap 这个 error，
// 抛给上层。为什么，应该怎么做请写出代码？
//
// 从隐藏内部逻辑的角度上，应该重定义这个错误进行上抛，
// 不能直接返回 sql.ErrNoRows
// 当数据库确实没有这个数据时，即获取数据在插入数据之前
// 该错误不应该是一个严重错误，此时 sql.ErrNoRows 应该
// 描述为提示, 这时候的错误不是数据丢失还是数据还未产生

var NoData = errors.New("data: no data")

func hw() error {
    var err error
    // ...
    // 访问数据库

    if errors.Is(err, sql.ErrNoRows) {
        return NoData
    }

    // ...
    return nil
}

func main() {
}
