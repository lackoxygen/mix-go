package xsql

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"time"
)

type Fetcher struct {
	r *sql.Rows
}

func (t *Fetcher) First(i interface{}) error {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return errors.New("argument can only be pointer type")
	}
	root := value.Elem()

	rows, err := t.Rows()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return errors.New("rows is empty")
	}
	row := rows[0]

	for n := 0; n < root.NumField(); n++ {
		field := root.Field(n)
		if !field.CanSet() {
			continue
		}
		tag := root.Type().Field(n).Tag.Get("xsql")
		if tag == "-" || tag == "_" {
			continue
		}
		if !row.Exist(tag) {
			continue
		}
		if err = mapped(field, row, tag); err != nil {
			return err
		}
	}

	return nil
}

func (t *Fetcher) Find(i interface{}) error {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return errors.New("argument can only be pointer type")
	}
	root := value.Elem()
	itemType := root.Type().Elem()

	rows, err := t.Rows()
	if err != nil {
		return err
	}

	for r := 0; r < len(rows); r++ {
		newItem := reflect.New(itemType)
		if newItem.Kind() == reflect.Ptr {
			newItem = newItem.Elem()
		}
		for n := 0; n < newItem.NumField(); n++ {
			field := newItem.Field(n)
			if !field.CanSet() {
				continue
			}
			tag := newItem.Type().Field(n).Tag.Get("xsql")
			if tag == "-" || tag == "_" {
				continue
			}
			if !rows[r].Exist(tag) {
				continue
			}
			if err = mapped(field, rows[r], tag); err != nil {
				return err
			}
		}
		root.Set(reflect.Append(root, newItem))
	}

	return nil
}

func (t *Fetcher) Rows() ([]Row, error) {
	// 获取列名
	columns, err := t.r.Columns()
	if err != nil {
		return nil, err
	}

	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	var rows []Row

	for t.r.Next() {
		err = t.r.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, value := range values {
			// Here we can check if the value is nil (NULL value)
			if value != nil {
				rowMap[columns[i]] = value
			}
		}

		rows = append(rows, rowMap)
	}

	return rows, nil
}

type Row map[string]interface{}

func (t Row) Exist(field string) bool {
	_, ok := t[field]
	return ok
}

func (t Row) Get(field string) *Result {
	if v, ok := t[field]; ok {
		return &Result{v: v}
	}
	return &Result{v: ""}
}

func (t Row) Value() map[string]interface{} {
	return t
}

type Result struct {
	v interface{}
}

func (t *Result) Empty() bool {
	if b, ok := t.v.([]uint8); ok {
		return len(b) == 0
	}
	if s, ok := t.v.(string); ok {
		return len(s) == 0
	}
	if t.v == nil {
		return true
	}
	return false
}

func (t *Result) String() string {
	switch reflect.ValueOf(t.v).Kind() {
	case reflect.Int:
		i := t.v.(int)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Int8:
		i := t.v.(int8)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Int16:
		i := t.v.(int16)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Int32:
		i := t.v.(int32)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Int64:
		i := t.v.(int64)
		return strconv.FormatInt(i, 10)
	case reflect.Uint:
		i := t.v.(uint)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Uint8:
		i := t.v.(uint8)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Uint16:
		i := t.v.(uint16)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Uint32:
		i := t.v.(uint32)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Uint64:
		i := t.v.(uint64)
		return strconv.FormatInt(int64(i), 10)
	case reflect.String:
		return t.v.(string)
	default:
		if b, ok := t.v.([]uint8); ok {
			return string(b)
		}
	}
	return ""
}

func (t *Result) Int() int64 {
	switch reflect.ValueOf(t.v).Kind() {
	case reflect.Int:
		i := t.v.(int)
		return int64(i)
	case reflect.Int8:
		i := t.v.(int8)
		return int64(i)
	case reflect.Int16:
		i := t.v.(int16)
		return int64(i)
	case reflect.Int32:
		i := t.v.(int32)
		return int64(i)
	case reflect.Int64:
		i := t.v.(int64)
		return i
	case reflect.Uint:
		i := t.v.(uint)
		return int64(i)
	case reflect.Uint8:
		i := t.v.(uint8)
		return int64(i)
	case reflect.Uint16:
		i := t.v.(uint16)
		return int64(i)
	case reflect.Uint32:
		i := t.v.(uint32)
		return int64(i)
	case reflect.Uint64:
		i := t.v.(uint64)
		return int64(i)
	case reflect.String:
		s := t.v.(string)
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0
		}
		return i
	default:
		if b, ok := t.v.([]uint8); ok {
			s := string(b)
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return 0
			}
			return i
		}
	}
	return 0
}

func (t *Result) Time() time.Time {
	typ := t.Type()
	if typ == "string" || typ == "[]uint8" {
		tt, _ := time.ParseInLocation(TimeParselayout, t.String(), time.Local)
		return tt
	}
	if typ == "time.Time" {
		return t.v.(time.Time)
	}
	return time.Time{}
}

func (t *Result) Value() interface{} {
	return t.v
}

func (t *Result) Type() string {
	return reflect.TypeOf(t.v).String()
}
