package state

import (
	"fmt"
)

type Twitch struct {
	Name string `yaml:"name"`
}

func ( t *Twitch ) Play() {
	fmt.Println( "Twitch playing" )
}

func ( t *Twitch ) Pause() {
	fmt.Println( "Twitch paused" )
}

func ( t *Twitch ) Stop() {
	fmt.Println( "Twitch stopped" )
}

func ( t *Twitch ) Next() {
	fmt.Println( "Twitch next" )
}

func ( t *Twitch ) Previous() {
	fmt.Println( "Twitch previous" )
}

func ( t *Twitch ) Teardown() {
	fmt.Println( "Twitch previous" )
}

func ( t *Twitch ) Setup() {
	fmt.Println( "Twitch previous" )
}

func ( t *Twitch ) Update() {
	fmt.Println( "Twitch previous" )
}