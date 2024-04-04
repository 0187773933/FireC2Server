package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "math"
	// "image/color"
	filepath "path/filepath"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

func ( s *Server ) DisneyReopenApp() {
	log.Debug( "DisneyReopenApp()" )
	s.ADB.StopAllPackages()
	s.ADB.ClosePackage( s.Config.ADB.APKS[ "disney" ][ s.Config.ADB.DeviceType ].Package )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenPackage( s.Config.ADB.APKS[ "disney" ][ s.Config.ADB.DeviceType ].Package )
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
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeUserProfileIndex ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
		time.Sleep( 1000 * time.Millisecond )
	}
	for _ , v := range s.Config.ADB.APKS[ "disney" ][ s.Config.ADB.DeviceType ].Activities {
		if s.Status.ADB.Activity == v {
			log.Debug( fmt.Sprintf( "disney was already open with activity %s" , v ) )
			return
		}
	}
	log.Debug( "disney was NOT already open" )
	s.DisneyReopenApp()
	time.Sleep( 500 * time.Millisecond )
	pss_fp := filepath.Join( s.Config.SaveFilesPath , "screenshots" , "disney" , "profile_selection.png" )
	log.Debug( "waiting on profile selection screen" )
	s.ADB.WaitOnScreen( pss_fp , ( 20 * time.Second ) )
	log.Debug( "selecting hardcoded profile" ) // todo , disney profile index numbers
	time.Sleep( 500 * time.Millisecond )
	s.ADB.Right()
	s.ADB.Right()
	s.ADB.Right()
	s.ADB.Right()
	s.ADB.Right()
	time.Sleep( 500 * time.Millisecond )
	s.ADB.Left()
	time.Sleep( 200 * time.Millisecond )
	s.ADB.Enter()
}

func ( s *Server ) DisneyMovieNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "DisneyMovieNext()" )
	s.DisneyContinuousOpen()
	next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.ADB.Key( "KEYCODE_DPAD_RIGHT" )
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
	s.ADB.Key( "KEYCODE_DPAD_RIGHT" )
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
	movie_id := c.Params( "id" )
	log.Debug( fmt.Sprintf( "DisneyMovie( %s )" , movie_id ) )
	s.DisneyContinuousOpen()
	uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , movie_id )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.ADB.Key( "KEYCODE_DPAD_RIGHT" )
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
