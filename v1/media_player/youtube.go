package media_player

import (
	"fmt"
	// bolt_api "github.com/boltdb/bolt"
)

type YouTube struct {
	*MediaPlayer
	Name string `yaml:"youtube"`
	// DB *bolt_api.DB `yaml:"-"`
}

func ( mp *YouTube ) PlayPlaylistWithShuffle( playlist_id string ) {
	log.Debug( "YouTube PlayPlaylistWithShuffle()" )
}

func ( mp *YouTube ) PlayPlaylist( playlist_id string ) {
	log.Debug( "YouTube PlayPlaylist()" )
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