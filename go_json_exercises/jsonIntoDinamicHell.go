/*
This is a bad example !!!
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const input = `
{
	"type": "sound",
	"msg": {
		"description": "dynamite",
		"authority": "the Bruce Dickinson"
	}
}
`

type EnvelopeForHell struct {
	Type string
	Msg  interface{}
}

// Test it: go run . jsonIntoDinamicHell
func jsonIntoDinamicHell_main() {
	var env EnvelopeForHell
	if err := json.Unmarshal([]byte(input), &env); err != nil {
		log.Fatal(err)
	}

	// for the love of Gopher DO NOT DO THIS
	var desc string = env.Msg.(map[string]interface{})["description"].(string)
	fmt.Println(desc)
	desc = env.Msg.(map[string]interface{})["authority"].(string)
	fmt.Println(desc)
}
