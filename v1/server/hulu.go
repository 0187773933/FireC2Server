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

// // // https://www.hulu.com/series/502bbc34-fa19-48fb-89c6-074da28335d3
// ( async ()=> {
// 	let images = document.querySelectorAll('img[class^="StandardEmphasisHorizontalTileThumbnail"]');
// 	let yaml_string = "";
// 	for ( let i = 0; i < images.length; ++i ) {
// 		let uuid = images[ i ].src.split( "artwork/" )[ 1 ].split( "?" )[ 0 ]
// 		let alt = images[ i ].alt.split( "Cover art for " )[ 1 ];
// 		if ( alt === undefined ) { continue; }
// 		yaml_string += `${uuid}:\n`;
// 		yaml_string += `  name: "${alt}"\n`;
// 	}
// 	console.log( yaml_string );
// })();

// https://community.home-assistant.io/t/androidtv-autoplay-hulu-series-where-you-left-off/531105
// adb shell am start -n com.hulu.livingroomplus/.WKFactivity
// -a hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3

// am start -n hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3
// start -n com.hulu.plus/.WKFactivity -a hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3

const HULU_ACTIVITY = "com.hulu.plus/com.hulu.plus.MainActivity"
const HULU_APP_NAME = "com.hulu.plus"

func ( s *Server ) HuluReopenApp() {
	log.Debug( "HuluReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( HULU_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( HULU_APP_NAME )
	log.Debug( "Done" )
}

func ( s *Server ) HuluContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "HuluContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "hulu" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	if s.Status.ADB.Activity == ACTIVITY_PROFILE_PICKER {
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeUserProfileIndex ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
		time.Sleep( 1000 * time.Millisecond )
	} else if s.Status.ADB.Activity == HULU_ACTIVITY {
		log.Debug( "hulu was already open" )
	} else {
		log.Debug( "hulu was NOT already open" )
		s.HuluReopenApp()
		time.Sleep( 500 * time.Millisecond )
		for i := 0; i < s.Config.HuluTotalUserProfiles; i++ {
			s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
			time.Sleep( 100 * time.Millisecond )
		}
		s.ADB.PressKeyName( "KEYCODE_DPAD_UP" )
		time.Sleep( 100 * time.Millisecond )
		s.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
}

func ( s *Server ) HuluMovieNext( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	log.Debug( "HuluMovieNext()" )
	s.HuluContinuousOpen()
	// next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	// uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	// log.Debug( uri )
	// s.ADB.OpenURI( uri )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	// s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	// s.Set( "active_player_now_playing_id" , next_movie )
	// s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/hulu/next" ,
		"result": true ,
	})
}

func ( s *Server ) HuluMoviePrevious( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	log.Debug( "HuluMoviePrevious()" )
	s.HuluContinuousOpen()
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/hulu/previous" ,
		"result": true ,
	})
}

func ( s *Server ) HuluTVNext( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	log.Debug( "HuluTVNext()" )
	s.HuluContinuousOpen()
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/hulu/tv/:id/next" ,
		"result": true ,
	})
}

func ( s *Server ) HuluTVPrevious( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	log.Debug( "HuluTVPrevious()" )
	s.HuluContinuousOpen()
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/hulu/tv/:id/previous" ,
		"result": true ,
	})
}

func ( s *Server ) HuluID( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	sent_id := c.Params( "*" )
	if utils.IsUUID( sent_id ) {
		sent_id = fmt.Sprintf( "https://www.hulu.com/watch/%s" , sent_id )
	}
	log.Debug( fmt.Sprintf( "HuluPlaylistAddURL( %s )" , sent_id ) )
	s.HuluContinuousOpen()
	s.ADB.OpenURI( sent_id )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_uri" , sent_id )
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/hulu/:id" ,
		"id": sent_id ,
		"result": true ,
	})
}