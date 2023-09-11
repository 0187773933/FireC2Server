package state

import (
	"fmt"
)

type YouTube struct {
	Name string `yaml:"youtube"`
}

func ( t *YouTube ) Play() {
	fmt.Println( "YouTube playing" )
}

func ( t *YouTube ) Pause() {
	fmt.Println( "YouTube paused" )
}

func ( t *YouTube ) Stop() {
	fmt.Println( "YouTube stopped" )
}

func ( t *YouTube ) Next() {
	fmt.Println( "YouTube next" )
}

func ( t *YouTube ) Previous() {
	fmt.Println( "YouTube previous" )
}

func ( t *YouTube ) Teardown() {
	fmt.Println( "YouTube previous" )
}

func ( t *YouTube ) Setup() {
	fmt.Println( "YouTube previous" )
}

func ( t *YouTube ) Update() {
	fmt.Println( "YouTube previous" )
}