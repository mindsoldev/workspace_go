package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const inputForRaw = `
{
	"type": "sound",
	"msg": {
		"description": "dynamite",
		"authority": "the Bruce Dickinson"
	}
}
`

type EnvelopeForRaw struct {
	Type string
	Msg  interface{}
}

type SoundforRaw struct {
	Description string
	Authority   string
}

// Test it: go run . jsonWithRawMessage
func jsonWithRawMessage_main() {
	var msg json.RawMessage
	env := Envelope{
		Msg: &msg,
	}
	if err := json.Unmarshal([]byte(inputForRaw), &env); err != nil {
		log.Fatal(err)
	}

	switch env.Type {
	case "sound":
		var s Sound
		if err := json.Unmarshal([]byte(msg), &s); err != nil {
			log.Fatal(err)
		}
		var desc string = s.Description
		fmt.Println(desc)
	default:
		log.Fatalf("unknown message type: %q", env.Type)
	}
}
