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
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	if s.Status.ADB.Activity == ACTIVITY_PROFILE_PICKER {
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeTotalUserProfiles ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
		time.Sleep( 1000 * time.Millisecond )
	} else if s.Status.ADB.Activity == DISNEY_PLAYING_ACTIVITY || s.Status.ADB.Activity == DISNEY_ACTIVITY {
		log.Debug( "disney was already open" )
	} else {
		log.Debug( "disney was NOT already open" )
		s.DisneyReopenApp()
		time.Sleep( 500 * time.Millisecond )
		s.ADB.WaitOnScreen( "./screenshots/disney/profile_selection.png" , ( 20 * time.Second ) )
		time.Sleep( 500 * time.Millisecond )
		s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
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
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	s.Set( "active_player_now_playing_id" , next_movie )
	s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	return c.JSON( fiber.Map{
		"url": "/disney/next" ,
		"uuid": next_movie ,
		"name": s.Config.Library.Disney.Movies.Currated[ next_movie ].Name ,
		"result": true ,
	})
}

func ( s *Server ) DisneyMoviePrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "DisneyMoviePrevious()" )
	s.DisneyContinuousOpen()
	next_movie := circular_set.Previous( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	s.Set( "active_player_now_playing_id" , next_movie )
	s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	return c.JSON( fiber.Map{
		"url": "/disney/previous" ,
		"uuid": next_movie ,
		"name": s.Config.Library.Disney.Movies.Currated[ next_movie ].Name ,
		"result": true ,
	})
}

func ( s *Server ) DisneyMovie( c *fiber.Ctx ) ( error ) {
	movie_id := c.Params( "movie_id" )
	log.Debug( fmt.Sprintf( "DisneyMovie( %s )" , movie_id ) )
	s.DisneyContinuousOpen()
	uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , movie_id )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	name := "unknown"
	if movie , ok := s.Config.Library.Disney.Movies.Currated[ movie_id ]; ok {
		name = movie.Name
	}
	s.Set( "STATE.DISNEY.NOW_PLAYING" , movie_id )
	s.Set( "active_player_now_playing_id" , movie_id )
	s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ movie_id ].Name )
	return c.JSON( fiber.Map{
		"url": "/disney/previous" ,
		"uuid": movie_id ,
		"name": name ,
		"result": true ,
	})
}
