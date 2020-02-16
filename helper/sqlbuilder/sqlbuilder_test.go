package sqlbuilder

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SQLBuilder_InsertSQL(t *testing.T) {
	Convey("Normal", t, func() {
		entities := make([]*TestStruct, 2)
		entities[0] = &TestStruct{0, "w", time.Now()}
		entities[1] = &TestStruct{1, "z", time.Now()}
		rvalues := reflectToValues(entities)
		sql, _, err := NewSQLBuilder().InsertSQL(rvalues, "test")

		So(err, ShouldBeNil)
		So(sql, ShouldEqual, "insert into test (id,name,created_at) values (?,?,?),(?,?,?)")
	})
}
