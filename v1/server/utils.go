package server

import (
	"fmt"
	context "context"
	"encoding/json"
	// "strings"
	// tv "github.com/0187773933/FireC2Server/v1/tv"
	// types "github.com/0187773933/FireC2Server/v1/types"
	redis "github.com/redis/go-redis/v9"
	// adb_wrapper "ADBWrapper/v1/wrapper"
	adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
	fiber "github.com/gofiber/fiber/v2"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// weak attempt at sanitizing form input to build a "username"
func SanitizeUsername( first_name string , last_name string ) ( username string ) {
	if first_name == "" { first_name = "Not Provided" }
	if last_name == "" { last_name = "Not Provided" }
	sanitized_first_name := utils.SanitizeInputString( first_name )
	sanitized_last_name := utils.SanitizeInputString( last_name )
	username = fmt.Sprintf( "%s-%s" , sanitized_first_name , sanitized_last_name )
	return
}

func serve_failed_attempt( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>no</h1>" )
	// return context.SendFile( "./v1/server/html/admin_login.html" )
}

func ServeLoginPage( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/admin_login.html" )
}

// func ServeAuthenticatedPage( context *fiber.Ctx ) ( error ) {
// 	if validate_admin_cookie( context ) == false { return serve_failed_attempt( context ) }
// 	x_path := context.Route().Path
// 	url_key := strings.Split( x_path , "/admin" )
// 	if len( url_key ) < 2 { return context.SendFile( "./v1/server/html/admin_login.html" ) }
// 	// fmt.Println( "Sending -->" , url_key[ 1 ] , x_path )
// 	return context.SendFile( ui_html_pages[ url_key[ 1 ] ] )
// }

func ( s *Server ) LogRequest( context *fiber.Ctx ) ( error ) {
	time_string := utils.GetFormattedTimeString()
	ip_address := context.Get( "x-forwarded-for" )
	if ip_address == "" { ip_address = context.IP() }
	log_message := fmt.Sprintf( "%s === %s === %s === %s" , time_string , ip_address , context.Method() , context.Path() )
	fmt.Println( log_message )
	return context.Next()
}

func ( s *Server ) Log( message string ) {
	time_string := utils.GetFormattedTimeString()
	log_message := fmt.Sprintf( "%s === %s" , time_string , message )
	fmt.Println( log_message )
}

func ( s *Server ) Printf( format_string string , args ...interface{} ) {
	time_string := utils.GetFormattedTimeString()
	sent_format := fmt.Sprintf( format_string , args... )
	fmt.Printf( "%s === %s" , time_string , sent_format )
}

func ( s *Server ) ADBConnect() ( connection adb_wrapper.Wrapper ) {
	if s.Config.ADBConnectionType == "tcp" {
		connection = adb_wrapper.ConnectIP( s.Config.ADBPath , s.Config.ADBServerIP , s.Config.ADBServerPort )
	} else if s.Config.ADBConnectionType == "usb" {
		connection = adb_wrapper.ConnectUSB( s.Config.ADBPath , s.Config.ADBSerial )
	}
	s.ADB = connection
	return
}

func ( s *Server ) GetStatus() ( result Status ) {
	log.Debug( "GetStatus()" )

	// 1.) Get Previous State Info from DB
	start_time , start_time_obj := utils.GetFormattedTimeStringOBJ()
	previous_player_name := s.Get( "active_player_name" )
	previous_player_command := s.Get( "active_player_command" )
	previous_start_time := s.Get( "active_player_start_time" )

	active_player_now_playing_id := s.Get( "active_player_now_playing_id" )
	active_player_now_playing_text := s.Get( "active_player_now_playing_text" )

	result.NowPlayingID = active_player_now_playing_id
	result.NowPlayingText = active_player_now_playing_text

	result.PlayerName = previous_player_name
	result.PlayerCommand = previous_player_command

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

	// 2.) Get Current ADB Status
	result.ADB = s.ADB.GetStatus()

	// 3.) Get TV Status
	result.TV = s.TV.Status()

	s.Status = result
	return
}
func ( s *Server ) GetStatusUrl( c *fiber.Ctx ) ( error ) {
	status := s.GetStatus()
	return c.JSON( fiber.Map{
		"url": "/status" ,
		"status": status ,
		"result": true ,
	})
}

func ( s *Server ) StoreLibrary() {

	var ctx = context.Background()

	// Should we delete everyting before hand? .... design question

	// Spotify Songs == TODO

	// Spotify Playlists
	for key , _ := range s.Config.Library.Spotify.Playlists {
		circular_set.Add( s.DB , "LIBRARY.SPOTIFY.PLAYLISTS" , key )
	}

	// Twitch Following - Curated
	s.DB.Del( ctx , "LIBRARY.TWITCH.FOLLOWING.CURRATED" )
	s.DB.Del( ctx , "STATE.TWITCH.FOLLOWING.LIVE" )
	s.DB.Del( ctx , "STATE.TWITCH.FOLLOWING.LIVE.INDEX" )
	for _ , item := range s.Config.Library.Twitch.Following.Currated {
		circular_set.Add( s.DB , "LIBRARY.TWITCH.FOLLOWING.CURRATED" , item )
	}

	// Twitch Following - All
	for _ , item := range s.Config.Library.Twitch.Following.All {
		circular_set.Add( s.DB , "LIBRARY.TWITCH.FOLLOWING.ALL" , item )
	}

	// Set API Keys
	for _ , api_key := range s.Config.YouTubeAPIKeys {
		circular_set.Add( s.DB , "CONFIG.YOUTUBE.API_KEYS" , api_key )
	}

	// Youtube Videos - Live
	// for _ , item := range s.Config.Library.YouTube.Videos.Live {
	// 	circular_set.Add( s.DB , "LIBRARY.YOUTUBE.VIDEOS.LIVE" , item )
	// }

	// // Youtube Videos - Normal
	// for _ , item := range s.Config.Library.YouTube.Videos.Normal {
	// 	circular_set.Add( s.DB , "LIBRARY.YOUTUBE.VIDEOS.NORMAL" , item )
	// }

	// // Youtube Playlists - Normal
	// for _ , item := range s.Config.Library.YouTube.Playlists.Normal {
	// 	circular_set.Add( s.DB , "LIBRARY.YOUTUBE.PLAYLISTS.NORMAL" , item )
	// }

	// // Youtube Playlists - Relaxing
	// for _ , item := range s.Config.Library.YouTube.Playlists.Relaxing {
	// 	circular_set.Add( s.DB , "LIBRARY.YOUTUBE.PLAYLISTS.RELAXING" , item )
	// }

	// // Youtube Following - Live Channels
	// for _ , item := range s.Config.Library.YouTube.Following.Live {
	// 	circular_set.Add( s.DB , "LIBRARY.YOUTUBE.FOLLOWING.LIVE" , item )
	// }

	// // Youtube Following - Channels
	// for _ , item := range s.Config.Library.YouTube.Following.Normal {
	// 	circular_set.Add( s.DB , "LIBRARY.YOUTUBE.FOLLOWING.NORMAL" , item )
	// }

	// Disney Movies - Currated
	disney_movies_currated_shuffled := utils.ShuffleKeys( s.Config.Library.Disney.Movies.Currated )
	// s.DB.Del( context.Background() , "LIBRARY.DISNEY.MOVIES.CURRATED" ) // #design-decision , forces new random
	// s.DB.Del( context.Background() , "LIBRARY.DISNEY.MOVIES.CURRATED.INDEX" ) // #design-decision , forces new random
	for _ , item := range disney_movies_currated_shuffled {
		circular_set.Add( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" , item )
	}

	// VLC == TODO

}

func ( s *Server ) Set( key string , value interface{} ) ( result string ) {
	var ctx = context.Background()
	err := s.DB.Set( ctx , key , value , 0 ).Err()
	if err != nil { panic( err ) }
	result = "success"
	return
}
func ( s *Server ) Get( key string ) ( result string ) {
	var ctx = context.Background()
	val , err := s.DB.Get( ctx , key ).Result()
	if err == redis.Nil { return }
	if err != nil { panic( err ) }
	result = val
	return
}

// have to build
// https://github.com/RedisJSON/RedisJSON/
// and then add loadmodule /Users/morpheous/APPLICATIONS_2/RedisJSON/target/release/librejson.dylib
// to the redis.conf
// brew services restart redis

// https://github.com/RedisJSON/RedisJSON/
// s.SetJSON( "config" , s.Config )
// https://redis.io/commands/json.set/

// https://github.com/RedisJSON/RedisJSON/#community-supported-clients
// fmt.Println( s.DB.JSONSet ) 😭
// https://github.com/redis/go-redis/pull/2704/commits/acf2b714f7b4920f1b910247dc42799b354f62ce
func ( s *Server ) SetJSON( key string , value interface{} ) ( result string ) {
	json_value , _ := json.Marshal( value )
	result = s.Set( key , json_value )
	return
}

// var test = &types.ConfigFile{}
// s.GetJSON( "config" , test )
// fmt.Println( test.BoltDBEncryptionKey )
// https://redis.io/commands/json.get/
func ( s *Server ) GetJSON( key string , target interface{} ) {
	json_value := s.Get( key )
	json.Unmarshal( []byte( json_value ) , target )
	return
}