package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
)

func main() {
	// os.Args = []string{".", "jsonIntoDinamicHell"}
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
	"jsonToEnvelope":        jsonToEnvelope_main,
	"jsonIntoDinamicHell":   jsonIntoDinamicHell_main,
	"jsonWithRawMessage":    jsonWithRawMessage_main,
	"jsonUnmarshalingTwice": jsonUnmarshalingTwice_main,
}

func Call(funcName string, params ...interface{}) (result interface{}, err error) {
	fmt.Printf("funcName: %v\n", funcName)
	fmt.Printf("params: %v\n", params)
	function := reflect.ValueOf(StubStorage[funcName])
	fmt.Printf("function: %v\n", function)
	fmt.Println("Output:")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
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

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

	return result, err
}
