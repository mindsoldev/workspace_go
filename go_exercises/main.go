package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
)

func main() {
	fuctionName := os.Args[1]
	arguments := []interface{}{}
	fmt.Println(len(os.Args))
	for i, arg := range os.Args[2:] {
		fmt.Printf("%v. arg: %v\n", i, arg)
		arguments = append(arguments, arg)
		fmt.Println(arguments...)

	}
	Call(fuctionName, arguments...)
}

// source: https://medium.com/@vicky.kurniawan/go-call-a-function-from-string-name-30b41dcb9e12
type stubMapping map[string]interface{}

var StubStorage = stubMapping{
	"crawl_main":      crawl_main,
	"treeComper_main": treeComper_main,
	"fetch15":         fetch15,
	"fetch17":         fetch17,
	"fetch18":         fetch18,
	"fetch19":         fetch19,
	"fetchAll16":      fetchAll16,
	"fetchAll10":      fetchAll10,
	"fetchAll11":      fetchAll11,
	"server1":         server1_main,
	"server2":         server2_main,
	"server3":         server3_main,
	"lissajous":       lissajous_main,
	"movie":           movie_main,
	"issues":          issues_main,
	"issuesByAges":    issuesByAges_main,
}

func Call(funcName string, params ...interface{}) (result interface{}, err error) {
	fmt.Printf("funcName: %v\n", funcName)
	fmt.Printf("params: %v\n", params)
	function := reflect.ValueOf(StubStorage[funcName])
	fmt.Printf("function: %v\n", function)
	if !function.IsValid() {
		log.Fatalf("%v function not found!", funcName)
		return
	}

	if !function.Type().IsVariadic() && len(params) != function.Type().NumIn() {
		err = errors.New("the number of params is out of index")
		log.Fatalln(err)
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	var res []reflect.Value = function.Call(in)
	if len(res) > 0 {
		result = res[0].Interface()
	} else {
		result = res
	}

	return result, err
}
