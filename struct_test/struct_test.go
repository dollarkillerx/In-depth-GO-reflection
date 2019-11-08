/**
 * @Author: DollarKiller
 * @Description:
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 15:34 2019-11-08
 */
package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
)

type test struct {
	Id   int
	Name string
	da
}

type da struct {
	Ncr string
}

func TestOne(t *testing.T) {
	one := test{Id: 123, Name: "sadas"}
	//toMap := structToMap(&one)
	//log.Println(toMap)
	result, err := structToMap2(&one)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, e := json.Marshal(result)
	if e == nil {
		log.Println(string(bytes))
	}

	//structToMap1(one)
	//fmt.Println()
	//log.Println("=========")
	//fmt.Println()
	//structToMap1("撒大声地")
}

func structToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	if obj1.Kind() == reflect.Ptr {
		return nil
	} else if obj1.Kind() != reflect.Struct {
		return nil
	}

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

func structToMap1(obj interface{}) map[string]interface{} {
	of := reflect.TypeOf(obj)
	log.Println(of.String())  // 返回类型的字符串描述
	log.Println(of.Name())    // 名称
	log.Println(of.PkgPath()) // 返回package 路径  如果是指针就没有
	log.Println(of.Kind())    // 返回类型的类别  如果是指针 就返回 reflect.Ptr
	log.Println(of.Size())    // 返回值在分配内存时的大小
	return nil
}

func structToMap2(obj interface{}) (result map[string]interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			result = nil
			err = errors.New("Structure definition error")
		}
	}()
	typ := reflect.TypeOf(obj).Elem()
	val := reflect.ValueOf(obj).Elem()

	// 判断是否是结构体
	//if typ.Kind() == reflect.Ptr || typ.Kind() != reflect.Struct {
	//	return nil, errors.New("Not Struct")
	//}
	result = make(map[string]interface{})
	for i := 0; i < typ.NumField(); i++ {
		if val.Field(i).Kind() == reflect.Struct {
			results, err := structToMap2(val.Field(i).Interface())
			if err != nil {
				return nil, errors.New("Not Struct")
			}
			result[typ.Field(i).Name] = results
		} else {
			result[typ.Field(i).Name] = val.Field(i).Interface()
		}
	}

	return result, nil
}

type test1 struct {
	Name string `hock:"ncr,asd" json:"name"`
	Prc  string `hock:"ccr,ger" json:"prc"`
}

func TestTwo(t *testing.T) {
	testData := test1{"this1", "this2"}
	tc1(testData)
}

func tc1(obj interface{}) error {
	typ := reflect.TypeOf(obj)
	//val := reflect.ValueOf(obj)
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag
		log.Println(tag)
		log.Println(tag.Get("hock"))
		log.Println("==========")
	}
	return nil
}

type ic struct {
	Name string `json:"name"`
}

func TestTh(t *testing.T) {
	ts1(&ic{Name: "assdasd"})
}

func ts1(obj interface{}) {
	ac := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(obj).Elem()

	log.Println(ac.Kind())
	log.Println(typ.Kind())

	for i := 0; i < typ.NumField(); i++ {
		fmt.Println(val.Field(i).Interface())
	}

}
