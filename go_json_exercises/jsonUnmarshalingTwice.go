/*
Suggested solution
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const inputTorTwice = `
{
	"type": "sound",
	"description": "dynamite",
	"authority": "the Bruce Dickinson"
}
`

type EnvelopeTorTwice struct {
	Type string
}

type SoundTorTwice struct {
	Description string
	Authority   string
}

// Test it: go run . jsonUnmarshalingTwice
func jsonUnmarshalingTwice_main() {
	var env Envelope
	buf := []byte(inputTorTwice)
	if err := json.Unmarshal(buf, &env); err != nil {
		log.Fatal(err)
	}

	switch env.Type {
	case "sound":
		var s struct {
			EnvelopeTorTwice
			SoundTorTwice
		}
		if err := json.Unmarshal(buf, &s); err != nil {
			log.Fatal(err)
		}
		var desc string = s.Description
		fmt.Println(desc)
	default:
		log.Fatalf("unknown message type: %q", env.Type)
	}
}
