package ascii

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"reflect"
	"strconv"
	"strings"
)

func unmarshal(str string, rType reflect.Type, rValue reflect.Value) error {

	rKind := rValue.Kind()
	//如果传入的不是struct，就退出
	if rKind != reflect.Struct {
		fmt.Println("expect struct")
		return errors.New("expect struct")
	}

	//获取到该结构体有几个字段
	numField := rValue.NumField()

	//变量结构体的所有字段
	var start = 1
	var end int
	for i := 0; i < numField; i++ {
		//fmt.Printf("Field %d: 类型为：%v\n", i, rValue.Field(i).Kind())
		//获取到struct标签, 注意需要通过reflect.Type来获取tag标签的值
		tagValStart := rType.Field(i).Tag.Get("start")
		tagValEnd := rType.Field(i).Tag.Get("end")
		//如果该字段是tag标签就显示，否则就不显示
		//start, _ := strconv.ParseInt(tagValStart[0:], 10, 32)
		//end, _ := strconv.ParseInt(tagValEnd[0:], 10, 32)
		if tagValStart != "" {
			start, _ = strconv.Atoi(tagValStart[0:])
		}

		if tagValEnd == "..." || tagValEnd == "" {
			end = len(str)
		} else {
			end, _ = strconv.Atoi(tagValEnd[0:])
		}

		if end > len(str) {
			return errors.New("message is not enough")
		}

		if start > end || start == 0 || end == 0 {
			return errors.New("the tag is wrong")
		}
		//fmt.Printf("Field %d: tag为start：%v--end：%v\n", i, start, end)

		//	ParseKind(str[start:end],rValue.Field(i))
		/*
			if i > 0 {
				rValue.Field(i).SetString("hhhhhhhh")
			} else {
				rValue.Field(i).Set(reflect.ValueOf(500))
			}
		*/
		switch rValue.Field(i).Kind() {
		case reflect.String:
			rValue.Field(i).SetString(str[start-1 : end])
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			n, _ := strconv.ParseUint(strings.TrimSpace(str[start-1:end]), 10, 64)
			rValue.Field(i).SetUint(n)
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			n, _ := strconv.ParseInt(strings.TrimSpace(str[start-1:end]), 10, 64)
			rValue.Field(i).SetInt(n)
			//rValue.Field(i).Set(reflect.ValueOf(n))
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			n, _ := strconv.ParseFloat(strings.TrimSpace(str[start-1:end]), 64)
			rValue.Field(i).SetFloat(n)
		case reflect.Slice:
			rValue.Field(i).SetBytes([]byte(str[start-1 : end]))
		case reflect.Bool:
			if start == end && str[start-1:end] == "1" {
				rValue.Field(i).Set(reflect.ValueOf(true))
				//rValue.Field(i).SetBool(true)
			} else {
				rValue.Field(i).SetBool(false)
			}

		case reflect.Struct:
			//chidld:=reflect.ValueOf(rValue.Field(i))
			unmarshal(str[start-1:end], rType.Field(i).Type, rValue.Field(i))
		default:
			return errors.New("unexpect type")
		}
		if i == numField-1 && end != len(str) {
			return errors.New("strings input is much longer than struct defined")
		}
	}

	return nil

}

func Unmarshal(str string, h interface{}) error {
	//获取reflect.Type 类型
	rType := reflect.TypeOf(h).Elem()
	//获取reflect.Value 类型
	rValue := reflect.ValueOf(h).Elem()

	err := unmarshal(str, rType, rValue)

	return err
}
