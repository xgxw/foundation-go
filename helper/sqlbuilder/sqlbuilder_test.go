package sqlbuilder

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SQLBuilder_InsertSQL(t *testing.T) {
	type TestStruct struct {
		ID        int64     `gorm:"column:id" assemble:"id"`
		Name      string    `gorm:"column:name" assemble:"name"`
		CreatedAt time.Time `gorm:"column:created_at" assemble:"created_at"`
	}
	Convey("Normal", t, func() {
		entities := make([]*TestStruct, 2)
		entities[0] = &TestStruct{0, "w", time.Now()}
		entities[1] = &TestStruct{1, "z", time.Now()}
		sql, _, err := NewSQLBuilder().InsertSQL(entities, "test")

		So(err, ShouldBeNil)
		So(sql, ShouldEqual, "insert into test (id,name,created_at) values (?,?,?),(?,?,?)")
	})
}

func Test_SQLBuilder_InsertSQLByGORM(t *testing.T) {
	type TestStruct struct {
		ID        int64     `gorm:"column:id" assemble:"id"`
		Name      string    `gorm:"column:name" assemble:"name"`
		CreatedAt time.Time `gorm:"column:created_at" assemble:"created_at"`
	}
	Convey("Normal", t, func() {
		entities := make([]*TestStruct, 2)
		entities[0] = &TestStruct{0, "w", time.Now()}
		entities[1] = &TestStruct{1, "z", time.Now()}
		sql, _, err := NewSQLBuilder().InsertSQLByGORM(entities, "test")

		So(err, ShouldBeNil)
		So(sql, ShouldEqual, "insert into test (id,name,created_at) values (?,?,?),(?,?,?)")
	})
}
