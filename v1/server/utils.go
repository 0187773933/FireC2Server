package server

import (
	"fmt"
	// "time"
	// "strings"
	bolt_api "github.com/boltdb/bolt"
	// tv "github.com/0187773933/FireC2Server/v1/tv"
	adb_wrapper "ADBWrapper/v1/wrapper"
	fiber "github.com/gofiber/fiber/v2"
	utils "github.com/0187773933/FireC2Server/v1/utils"
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
	s.Status = result
	return
}

func ( s *Server ) Set( key string , value string ) ( result string ) {
	s.DB.Update( func( tx *bolt_api.Tx ) error {
		bucket := tx.Bucket( []byte( "state" ) )
		bucket.Put( []byte( key ) , []byte( value ) )
		return nil
	})
	return "success"
}
func ( s *Server ) Get( key string ) ( result string ) {
	s.DB.View( func( tx *bolt_api.Tx ) error {
		bucket := tx.Bucket( []byte( "state" ) )
		value := bucket.Get( []byte( key ) )
		if value != nil { result = string( value ) }
		return nil
	})
	return result
}