package sqlbuilder

import (
	"fmt"
	"reflect"
	"time"
)

type TestStruct struct {
	ID        int64     `assemble:"id"`
	Name      string    `assemble:"name"`
	CreatedAt time.Time `assemble:"created_at"`
}

func buildInsertSQL() {
	entities := make([]*TestStruct, 2)
	entities[0] = &TestStruct{0, "w", time.Now()}
	entities[1] = &TestStruct{1, "z", time.Now()}
	values := reflectToValues(entities)
	fmt.Println(NewSQLBuilder().InsertSQL(values, "gensql"))
}

// 将struct转换为 reflect.Value
func reflectToValues(entities []*TestStruct) []reflect.Value {
	values := make([]reflect.Value, len(entities))
	for i, _ := range entities {
		values[i] = reflect.ValueOf(entities[i])
	}
	return values
}
