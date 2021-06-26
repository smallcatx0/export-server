package helper

import (
	"os"
	"path/filepath"
	"reflect"
)

// Empty 判断是否为空
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// GetDefInt 获取默认值
func GetDefInt(val, def int) int {
	if val == 0 {
		return def
	}
	return val
}

// GetDefStr 获取默认值
func GetDefStr(val, def string) string {
	if len(val) == 0 {
		return def
	}
	return val
}

// FileExists 检查文件是否存在
func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true //文件或者文件夹存在
	}
	if os.IsNotExist(err) {
		return false //不存在
	}
	return false //不存在，这里的err可以查到具体的错误信息
}

// TouchDir 创建文件夹
func TouchDir(path string) error {
	dir, _ := filepath.Split(path)
	if FileExists(dir) {
		return nil
	}
	err := os.MkdirAll(dir, 0666)
	return err
}

// AppendFile 追加文件
func AppendFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
