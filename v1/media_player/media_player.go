package media_player

import (
	"fmt"
	"time"
	// "sync"
	logger "github.com/0187773933/FireC2Server/v1/logger"
	bolt_api "github.com/boltdb/bolt"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	types "github.com/0187773933/FireC2Server/v1/types"
	adb_wrapper "ADBWrapper/v1/wrapper"
	tv "github.com/0187773933/FireC2Server/v1/tv"
)

// var mu sync.Mutex
// var active *State

var log = logger.GetLogger()

const ACTIVITY_BLANK = "com.amazon.tv.launcher/com.amazon.tv.launcher.ui.HomeActivity_vNext"

type Functions interface {
	Play()
	Pause()
	Stop()
	Next()
	Previous()
	Teardown()
	Setup()
	Update()
	ShuffleOn()
	ShuffleOff()
	Item( video_id string )
	Playlist( playlist_id string )
	NextItem()
	NextPlaylist()
	PreviousItem()
	PreviousPlaylist()
}

type MediaPlayer struct {
	Config *types.ConfigFile `json:"config"`
	DB *bolt_api.DB `json:"-"`
	MediaPlayers map[string]Functions `json:"-"`
	ADB adb_wrapper.Wrapper `json:"-"`
	TV *tv.TV `json:"-"`
	Status Status `json:"status"`
}

func New( db *bolt_api.DB , config *types.ConfigFile ) ( result *MediaPlayer ) {
	log.Debug( "setting up media-player" )
	result = &MediaPlayer{
		Config: config ,
		DB: db ,
		MediaPlayers: make( map[string]Functions ) ,
	}
	result.ADB = result.ADBConnect()
	result.TV = tv.New( config )
	result.Status = result.GetStatus()
	result.Set( "active_player_name" , "startup" )
	result.Set( "active_player_start_time" , utils.GetFormattedTimeString() )

	twitch := &Twitch{
		Name: "twitch" ,
		MediaPlayer: result ,
	}
	result.MediaPlayers[ twitch.Name ] = twitch

	youtube := &YouTube{
		Name: "youtube" ,
		MediaPlayer: result ,
	}
	result.MediaPlayers[ youtube.Name ] = youtube

	spotify := &Spotify{
		Name: "spotify" ,
		MediaPlayer: result ,
	}
	result.MediaPlayers[ spotify.Name ] = spotify
	return
}


func ( s *MediaPlayer ) Run( player_name string , command string , args ...interface{} ) ( result string ) {
	prepare_commands := map[string]bool {
		"play": true ,
		"pause": true ,
		"stop": true ,
		"next": true ,
		"previous": true ,
	}
	if prepare_commands[ command ] == true {
		go s.Prepare()
		s.Status = s.GetStatus()
	}
	log.Debug( fmt.Sprintf( "%s === %s" , player_name , command ) )
	log.Debug( args )
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

type Status struct {
	StartTime string `json:"start_time"`
	StartTimeOBJ time.Time `json:"-"`
	PreviousPlayerName string `json:"previous_player_name"`
	PreviousPlayerCommand string `json:"previous_player_command"`
	PreviousStartTime string `json:"previous_start_time"`
	PreviousStartTimeOBJ time.Time `json:"-"`
	PreviousStartTimeDuration time.Duration `json:"-"`
	PreviousStartTimeDurationSeconds float64 `json:"previous_start_time_duration_seconds"`
	ADBTopWindow string `json:"adb_top_window"`
	ADBVolume int `json:"adb_volume"`
	TV tv.Status `json:"tv"`
}

func ( s *MediaPlayer ) GetStatus() ( result Status ) {
	log.Debug( "GetStatus()" )

	// 1.) Get Previous State Info from DB
	start_time , start_time_obj := utils.GetFormattedTimeStringOBJ()
	previous_player_name := s.Get( "active_player_name" )
	previous_player_command := s.Get( "active_player_command" )
	previous_start_time := s.Get( "active_player_start_time" )

	result.StartTime = start_time
	result.StartTimeOBJ = start_time_obj
	result.PreviousPlayerName = previous_player_name
	result.PreviousPlayerCommand = previous_player_command
	result.PreviousStartTime = previous_start_time

	if previous_start_time != "" {
		previous_start_time_obj := utils.ParseFormattedTimeString( previous_start_time )
		previous_start_time_duration := start_time_obj.Sub( previous_start_time_obj )
		previous_start_time_duration_seconds := previous_start_time_duration.Seconds()
		result.PreviousStartTimeOBJ = previous_start_time_obj
		result.PreviousStartTimeDuration = previous_start_time_duration
		result.PreviousStartTimeDurationSeconds = previous_start_time_duration_seconds
	}

	// 2.) Get Current ADB Status Info
	adb_windows := s.ADB.GetWindowStack()
	if len( adb_windows ) > 0 {
		result.ADBTopWindow = adb_windows[ 0 ].Activity
	}
	result.ADBVolume = s.ADB.GetVolume()

	// 3.) TV Get Status
	result.TV = s.TV.Status()
	return
}

func ( s *MediaPlayer ) Prepare() {
	s.TV.Prepare()
}

func ( s *MediaPlayer ) ADBConnect() ( connection adb_wrapper.Wrapper ) {
	if s.Config.ADBConnectionType == "tcp" {
		connection = adb_wrapper.ConnectIP( s.Config.ADBPath , s.Config.ADBServerIP , s.Config.ADBServerPort )
	} else if s.Config.ADBConnectionType == "usb" {
		connection = adb_wrapper.ConnectUSB( s.Config.ADBPath , s.Config.ADBSerial )
	}
	s.ADB = connection
	return
}

func ( s *MediaPlayer ) Set( key string , value string ) ( result string ) {
	s.DB.Update( func( tx *bolt_api.Tx ) error {
		bucket := tx.Bucket( []byte( "state" ) )
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