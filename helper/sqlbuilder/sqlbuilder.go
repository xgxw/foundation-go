package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"
)

/*
	使用须知
	1. 只有 public 字段, 且字段tag中有指定标签才会被扫描
	2. 传入的 interface{} 类型为: []*Struct
*/

const DefaultTag = "assemble"

type SQLBuilder struct{}

func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{}
}

/*
	InsertSQL build insert sql. 使用默认Tag获取字段名称
	只有 public 字段, 且字段tag中有指定标签才会被扫描
*/
func (b *SQLBuilder) InsertSQL(structPtrSlice interface{}, table string) (string, []interface{}, error) {
	return b.insertSQL(structPtrSlice, table, b.defaultGetField)
}

func (b *SQLBuilder) InsertSQLByGORM(structPtrSlice interface{}, table string) (string, []interface{}, error) {
	return b.insertSQL(structPtrSlice, table, b.getFieldByGorm)
}

func (b *SQLBuilder) InsertSQLCustom(structPtrSlice interface{}, table string,
	getField func(reflect.StructField) string) (string, []interface{}, error) {
	return b.insertSQL(structPtrSlice, table, getField)
}

func (b *SQLBuilder) defaultGetField(field reflect.StructField) string {
	return field.Tag.Get(DefaultTag)
}

func (b *SQLBuilder) getFieldByGorm(field reflect.StructField) string {
	tag := field.Tag.Get("gorm")
	tags := strings.Split(tag, ":")
	return tags[1]
}

// -------------- 具体实现 --------------

func (b *SQLBuilder) insertSQL(structPtrSlice interface{}, table string,
	getField func(reflect.StructField) string) (string, []interface{}, error) {

	structPtrSliceValue := reflect.ValueOf(structPtrSlice)
	if structPtrSliceValue.Len() == 0 {
		return "", []interface{}{}, nil
	}

	typ := structPtrSliceValue.Index(0).Type().Elem()
	fields, fieldLocs := b.getFields(typ, getField)

	holders, sqlValues := b.getValues(structPtrSliceValue, fieldLocs)
	sql := b.buildInsertSQL(table, fields, holders)
	return sql, sqlValues, nil
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

func (b *SQLBuilder) getValues(structPtrSliceValue reflect.Value,
	fieldLocs []int) (holders string, sqlValues []interface{}) {

	buf := new(strings.Builder)
	// 占位符单元
	holder := strings.Repeat("?,", len(fieldLocs))
	holder = fmt.Sprintf("(%s)", strings.TrimSuffix(holder, ","))

	// 获取 占位符 和 values
	structPtrSliceValueLen := structPtrSliceValue.Len()
	sqlValues = make([]interface{}, structPtrSliceValueLen*len(fieldLocs))
	cursor := 0
	for i := 0; i < structPtrSliceValueLen; i++ {
		buf.WriteString(holder)
		buf.WriteString(",")

		structValue := structPtrSliceValue.Index(i).Elem()
		for _, fieldIndex := range fieldLocs {
			sqlValues[cursor] = structValue.Field(fieldIndex).Interface()
			cursor++
		}
	}
	holders = buf.String()
	holders = strings.TrimRight(holders, ",")
	return holders, sqlValues
}

func (b *SQLBuilder) buildInsertSQL(table, fields, holders string) string {
	return fmt.Sprintf("insert into %s (%s) values %s", table, fields, holders)
}
