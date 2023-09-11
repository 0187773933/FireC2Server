package state

import (
	"fmt"
	// "sync"
)

// var mu sync.Mutex
// var active *State

type MediaPlayer interface {
	Play()
	Pause()
	Stop()
	Next()
	Previous()
	Teardown()
	Setup()
	Update()
}

type State struct {
	ActivePlayerName string
	MediaPlayers map[string]MediaPlayer
}

func New() ( result *State ) {
	result = &State{
		ActivePlayerName: "" ,
		MediaPlayers: make( map[string]MediaPlayer ) ,
	}
	twitch := &Twitch{ Name: "twitch" }
	result.MediaPlayers[ twitch.Name ] = twitch
	youtube := &YouTube{ Name: "youtube" }
	result.MediaPlayers[ youtube.Name ] = youtube
	return
}

type CommandFunc func(s *State, player_name string) string
var commands = map[string]CommandFunc{
	"play": func( s *State , player_name string ) string { return s.Play( player_name ) } ,
	"pause": func( s *State , player_name string ) string { return s.Pause( player_name ) } ,
	"stop": func( s *State , player_name string ) string { return s.Stop( player_name ) } ,
	"next": func( s *State , player_name string ) string { return s.Next( player_name ) } ,
	"previous": func( s *State , player_name string ) string { return s.Previous( player_name ) } ,
	"teardown": func( s *State , player_name string ) string { return s.Teardown( player_name ) } ,
	"setup": func( s *State , player_name string ) string { return s.Setup( player_name ) } ,
	"update": func( s *State , player_name string ) string { return s.Update( player_name ) } ,
}
func ( s *State ) ExecuteCommand( player_name string , command_name string ) ( result string ) {
	fmt.Printf("Executing command: %s for player: %s\n", command_name, player_name)
	command , exists := commands[ command_name ]
	if !exists { panic( "command does not exist" ) }
	result = command( s , player_name )
	return
}

func ( s *State ) Play( player_name string ) ( result string ) {
	// TODO: Implement Teardown() and Setup() Decision Logic
	// TODO: TV Setup / Decision Logic
	// TODO: ADB Setup / Decision Logic
	s.ActivePlayerName = player_name
	fmt.Printf( "State --> Play( %s )\n" , s.ActivePlayerName )
	s.MediaPlayers[ s.ActivePlayerName ].Play()
	result = "ok ???"
	return
}

func ( s *State ) Pause( player_name string ) ( result string ) {
	s.ActivePlayerName = player_name
	fmt.Printf( "State --> Pause( %s )\n" , s.ActivePlayerName )
	s.MediaPlayers[ s.ActivePlayerName ].Pause()
	return
}

func ( s *State ) Stop( player_name string ) ( result string ) {
	s.ActivePlayerName = player_name
	fmt.Printf( "State --> Stop( %s )\n" , s.ActivePlayerName )
	s.MediaPlayers[ s.ActivePlayerName ].Stop()
	return
}

func ( s *State ) Next( player_name string ) ( result string ) {
	s.ActivePlayerName = player_name
	fmt.Printf( "State --> Next( %s )\n" , s.ActivePlayerName )
	s.MediaPlayers[ s.ActivePlayerName ].Next()
	return
}

func ( s *State ) Previous( player_name string ) ( result string ) {
	s.ActivePlayerName = player_name
	fmt.Printf( "State --> Previous( %s )\n" , s.ActivePlayerName )
	s.MediaPlayers[ s.ActivePlayerName ].Previous()
	return
}

func ( s *State ) Teardown( player_name string ) ( result string ) {
	s.ActivePlayerName = player_name
	fmt.Printf( "State --> Teardown( %s )\n" , s.ActivePlayerName )
	s.MediaPlayers[ s.ActivePlayerName ].Teardown()
	return
}

func ( s *State ) Setup( player_name string ) ( result string ) {
	s.ActivePlayerName = player_name
	fmt.Printf( "State --> Setup( %s )\n" , s.ActivePlayerName )
	s.MediaPlayers[ s.ActivePlayerName ].Setup()
	return
}

func ( s *State ) Update( player_name string ) ( result string ) {
	s.ActivePlayerName = player_name
	fmt.Printf( "State --> Update( %s )\n" , s.ActivePlayerName )
	s.MediaPlayers[ s.ActivePlayerName ].Update()
	return
}