package media_player

import (
	// "fmt"
	// "sync"
	bolt_api "github.com/boltdb/bolt"
	// utils "github.com/0187773933/FireC2Server/v1/utils"
)

// var mu sync.Mutex
// var active *State

type Functions interface {
	Play()
	Pause()
	Stop()
	Next()
	Previous()
	Teardown()
	Setup()
	Update()
}

type MediaPlayer struct {
	ActivePlayerName string `json:"active_player_name"`
	MediaPlayers map[string]Functions `json:"-"`
	DB *bolt_api.DB `json:"-"`
}

func New( db *bolt_api.DB ) ( result *MediaPlayer ) {
	result = &MediaPlayer{
		MediaPlayers: make( map[string]Functions ) ,
		DB: db ,
	}
	twitch := &Twitch{
		Name: "twitch" ,
		DB: db ,
	}
	result.MediaPlayers[ twitch.Name ] = twitch
	youtube := &YouTube{
		Name: "youtube" ,
		DB: db ,
	}
	result.MediaPlayers[ youtube.Name ] = youtube
	spotify := &Spotify{
		Name: "spotify" ,
		DB: db ,
	}
	result.MediaPlayers[ spotify.Name ] = spotify
	return
}

func ( s *MediaPlayer ) Run( player_name string , command string ) (result string) {
	if mp , exists := s.MediaPlayers[ player_name ]; exists {
		switch command {
			case "play":
				mp.Play()
			case "pause":
				mp.Pause()
			case "stop":
				mp.Stop()
			case "next":
				mp.Next()
			case "previous":
				mp.Previous()
			case "teardown":
				mp.Teardown()
			case "setup":
				mp.Setup()
			case "update":
				mp.Update()
			default:
				return "Invalid Command"
		}
		return "Command Executed"
	}
	return "Player Not Found"
}

func ( s *MediaPlayer ) Set( key string , value string ) ( result string ) {
	s.DB.Update( func( tx *bolt_api.Tx ) error {
		bucket , err := tx.CreateBucketIfNotExists( []byte( "state" ) )
		if err != nil { return err }
		bucket.Put( []byte( key ) , []byte( value ) )
		return nil
	})
	return "success"
}
func ( s *MediaPlayer ) Get( key string ) ( result string ) {
	s.DB.View( func( tx *bolt_api.Tx ) error {
		bucket := tx.Bucket( []byte( "state" ) )
		value := bucket.Get( []byte( key ) )
		if value != nil { result = string( value ) }
		return nil
	})
	return result
}