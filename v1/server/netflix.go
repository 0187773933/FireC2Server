package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "math"
	// "image/color"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	// circular_set "github.com/0187773933/RedisCircular/v1/set"
)

const NETFLIX_ACTIVITY = "com.netflix.ninja/com.netflix.ninja.MainActivity"
const NETFLIX_APP_NAME = "com.netflix.ninja"

func ( s *Server ) NetflixReopenApp() {
	log.Debug( "NetflixReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( NETFLIX_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( NETFLIX_APP_NAME )
	log.Debug( "Done" )
}

func ( s *Server ) NetflixContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "NetflixContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "netflix" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	if s.Status.ADB.Activity == ACTIVITY_PROFILE_PICKER {
		// i mean this assumes you only have like 2 profiles , idk man
		// TODO : add config options to specify how many profiles you have
		time.Sleep( 1000 * time.Millisecond )
		s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		time.Sleep( 200 * time.Millisecond )
		s.ADB.PressKeyName( "KEYCODE_ENTER" )
		time.Sleep( 1000 * time.Millisecond )
	} else if s.Status.ADB.Activity == NETFLIX_ACTIVITY {
		log.Debug( "netflix was already open" )
	} else {
		log.Debug( "netflix was NOT already open" )
		s.NetflixReopenApp()
		time.Sleep( 500 * time.Millisecond )
		// TODO , deal with profile selection screen ....
	}
}

func ( s *Server ) NetflixNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "NetflixNext()" )
	s.NetflixContinuousOpen()
	// next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	// uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	// log.Debug( uri )
	// s.ADB.OpenURI( uri )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	// s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	// s.Set( "active_player_now_playing_id" , next_movie )
	// s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	return c.JSON( fiber.Map{
		"url": "/netflix/next" ,
		"result": true ,
	})
}

func ( s *Server ) NetflixPrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "NetflixPrevious()" )
	s.NetflixContinuousOpen()
	return c.JSON( fiber.Map{
		"url": "/netflix/previous" ,
		"result": true ,
	})
}
