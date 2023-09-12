package media_player

import (
	"fmt"
	bolt_api "github.com/boltdb/bolt"
)

type Twitch struct {
	Name string `yaml:"name"`
	DB *bolt_api.DB `yaml:"-"`
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