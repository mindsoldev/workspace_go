package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Envelope struct {
	Type string
	Msg  interface{}
}

type Sound struct {
	Description string
	Authority   string
}

type Cowbell struct {
	More bool
}

// Test it: go run . jsonToEnvelope
func jsonToEnvelope_main() {
	s := EnvelopeForHell{
		Type: "sound",
		Msg: Sound{
			Description: "dynamite",
			Authority:   "the Bruce Dickinson",
		},
	}
	buf, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)

	c := EnvelopeForHell{
		Type: "cowbell",
		Msg: Cowbell{
			More: true,
		},
	}
	buf, err = json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)
}
