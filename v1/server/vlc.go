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

const VLC_ACTIVITY = "org.videolan.vlc/org.videolan.television.ui.MainTvActivity"
const VLC_APP_NAME = "org.videolan.vlc"

func ( s *Server ) VLCReopenApp() {
	log.Debug( "VLCReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( VLC_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( VLC_APP_NAME )
	log.Debug( "Done" )
}

func ( s *Server ) VLCContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "VLCContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "vlc" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	if s.Status.ADB.Activity == VLC_ACTIVITY {
		log.Debug( "vlc was already open" )
	} else {
		log.Debug( "vlc was NOT already open" )
		s.DisneyReopenApp()
		time.Sleep( 500 * time.Millisecond )
	}
}

func ( s *Server ) VLCNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "VLCNext()" )
	s.VLCContinuousOpen()
	// next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	// uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	// log.Debug( uri )
	// s.ADB.OpenURI( uri )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	// s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	// s.Set( "active_player_now_playing_id" , next_movie )
	// s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	return c.JSON( fiber.Map{
		"url": "/vlc/next" ,
		"result": true ,
	})
}

func ( s *Server ) VLCPrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "VLCPrevious()" )
	s.VLCContinuousOpen()
	return c.JSON( fiber.Map{
		"url": "/vlc/previous" ,
		"result": true ,
	})
}
