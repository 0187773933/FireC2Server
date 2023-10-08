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
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

const DISNEY_ACTIVITY = "com.disney.disneyplus/com.bamtechmedia.dominguez.main.MainActivity"
const DISNEY_PLAYING_ACTIVITY = "com.disney.disneyplus/com.bamtechmedia.dominguez.player.ui.experiences.legacy.v1.TvPlaybackActivity"
const DISNEY_APP_NAME = "com.disney.disneyplus"

func ( s *Server ) DisneyReopenApp() {
	log.Debug( "DisneyReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( DISNEY_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( DISNEY_APP_NAME )
	log.Debug( "Done" )
}

func ( s *Server ) DisneyContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "DisneyContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "disney" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADBTopWindow ) )
	if s.Status.ADBTopWindow == DISNEY_PLAYING_ACTIVITY || s.Status.ADBTopWindow == DISNEY_ACTIVITY {
		log.Debug( "disney was already open" )
	} else {
		log.Debug( "disney was NOT already open" )
		s.DisneyReopenApp()
		time.Sleep( 500 * time.Millisecond )
		s.ADB.WaitOnScreen( "./screenshots/disney/profile_selection.png" , ( 20 * time.Second ) )
		s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		time.Sleep( 500 * time.Millisecond )
		s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		time.Sleep( 200 * time.Millisecond )
		s.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
}

func ( s *Server ) DisneyMovieNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "DisneyMovieNext()" )
	s.DisneyContinuousOpen()
	next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	s.ADB.OpenURI( uri )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	return c.JSON( fiber.Map{
		"url": "/disney/movies/next" ,
		"movie": next_movie ,
		"uri": uri ,
		"result": true ,
	})
}
