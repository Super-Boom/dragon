package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
)

//过滤字段，map中只保留需要的键值对
func OnlyCols(cols []string, data map[string]string) {
	for k := range data {
		//判断k 是否在需要的cols中， 如果不在，则对应的键值对
		have := false //不在
		for _, col := range cols {
			if k == col {
				have = true
				break
			}
		}
		if have == false {
			delete(data, k)
		}
	}
}

//过滤字段，返回只包含了cols的map
func OnlyColumns(cols []string, data map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for _, col := range cols {
		// 遍历字段切片，如果data中包含了需要的字段，则放入新的res中
		if v, ok := data[col]; ok {
			res[col] = v
		}
	}
	return res
}

// hmac-sha1 input with key
func HmacSha1(input, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(input))
	return hex.EncodeToString(mac.Sum(nil))
}

// hmac-md5 input with key
func HmacMD5(input, key string) string {
	mac := hmac.New(md5.New, []byte(key))
	mac.Write([]byte(input))
	return hex.EncodeToString(mac.Sum(nil))
}

// return json string
func ToJsonString(data interface{}) string {
	j, _ := json.Marshal(data)
	return string(j)
}

// slice and trim white space
func SliceAndTrim(str string) []string {
	res := make([]string, 0)
	strs := strings.Split(str, ",")
	for _, s := range strs {
		res = append(res, strings.Trim(s, " "))
	}
	return res
}

// obj 不能为指针
func StructJsonTagToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Tag.Get("json")] = v.Field(i).Interface()
	}
	return data
}

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func GetStructFields(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	fields := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	return fields
}

func GetStructJsonTags(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	fields := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Tag.Get("json"))
	}
	return fields
}