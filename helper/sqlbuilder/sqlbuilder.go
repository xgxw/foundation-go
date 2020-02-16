package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"
)

/*
	使用须知
	1. 只有 public 字段, 且字段tag中有 assemble 标签才会被扫描
	2. 为了使代码更简洁, 限制传入的 reflect.Value 必须是指针类型的反射值.
		reflect.ValueOf(Ptr).Elem() == reflect.ValueOf(Struct) (效果上)
*/

const DefaultTag = "assemble"

type SQLBuilder struct{}

func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{}
}

// InsertSQL build insert sql. 使用默认Tag获取字段名称
func (b *SQLBuilder) InsertSQL(rvalues []reflect.Value, table string) (string, []interface{}, error) {
	return b.insertSQL(rvalues, table, b.defaultGetField)
}

// InsertSQLCustom build insert sql. 使用自定义方法获取字段名称
func (b *SQLBuilder) InsertSQLCustom(rvalues []reflect.Value, table string,
	getField func(reflect.StructField) string) (string, []interface{}, error) {

	return b.insertSQL(rvalues, table, getField)
}

func (b *SQLBuilder) defaultGetField(field reflect.StructField) string {
	// Lookup 可以判断是否找到tag
	return field.Tag.Get(DefaultTag)
}

func (b *SQLBuilder) insertSQL(rvalues []reflect.Value, table string,
	getField func(reflect.StructField) string) (string, []interface{}, error) {

	var typ reflect.Type
	typ = rvalues[0].Elem().Type()

	fields, fieldLocs := b.getFields(typ, getField)
	holders, values := b.getValues(rvalues, fieldLocs)
	sql := b.buildInsertSQL(table, fields, holders)
	return sql, values, nil
}

func (b *SQLBuilder) getFields(typ reflect.Type,
	getField func(reflect.StructField) string) (fields string, fieldLocs []int) {

	buf := new(strings.Builder)
	fieldLocs = []int{}
	for i := 0; i < typ.NumField(); i++ {
		field := getField(typ.Field(i))
		buf.WriteString(field)
		buf.WriteString(",")
		fieldLocs = append(fieldLocs, i)
	}
	fields = buf.String()
	fields = strings.TrimRight(fields, ",")
	return fields, fieldLocs
}

func (b *SQLBuilder) getValues(rvalues []reflect.Value,
	fieldLocs []int) (holders string, values []interface{}) {

	buf := new(strings.Builder)
	// 占位符单元
	holder := strings.Repeat("?,", len(fieldLocs))
	holder = fmt.Sprintf("(%s)", strings.TrimSuffix(holder, ","))

	// 获取 占位符 和 values
	values = make([]interface{}, len(rvalues)*len(fieldLocs))
	cursor := 0
	for _, v := range rvalues {
		buf.WriteString(holder)
		buf.WriteString(",")
		for _, fieldIndex := range fieldLocs {
			values[cursor] = v.Elem().Field(fieldIndex).Interface()
			cursor++
		}
	}
	holders = buf.String()
	holders = strings.TrimRight(holders, ",")
	return holders, values
}

func (b *SQLBuilder) buildInsertSQL(table, fields, holders string) string {
	return fmt.Sprintf("insert into %s (%s) values %s", table, fields, holders)
}
