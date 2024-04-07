package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "math"
	"strings"
	color "image/color"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// // // https://www.hulu.com/series/502bbc34-fa19-48fb-89c6-074da28335d3
// ( async ()=> {
// 	let images = document.querySelectorAll('img[class^="StandardEmphasisHorizontalTileThumbnail"]');
// 	let yaml_string = "";
// 	for ( let i = 0; i < images.length; ++i ) {
// 		let uuid = images[ i ].src.split( "artwork/" )[ 1 ].split( "?" )[ 0 ]
// 		let alt = images[ i ].alt.split( "Cover art for " )[ 1 ];
// 		if ( alt === undefined ) { continue; }
// 		yaml_string += `          - id: ${uuid}\n`;
// 		yaml_string += `            name: "${alt}"\n`;
// 	}
// 	console.log( yaml_string );
// })();

// https://community.home-assistant.io/t/androidtv-autoplay-hulu-series-where-you-left-off/531105
// adb shell am start -n com.hulu.livingroomplus/.WKFactivity
// -a hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3

// am start -n hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3
// start -n com.hulu.plus/.WKFactivity -a hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3

// const HULU_ACTIVITY = "com.hulu.plus/com.hulu.plus.MainActivity"
// const HULU_APP_NAME = "com.hulu.plus"

func ( s *Server ) HuluReopenApp() {
	log.Debug( "HuluReopenApp()" )
	s.ADB.StopAllPackages()
	// s.ADB.SetBrightness( 0 )
	s.ADB.ClosePackage( s.Config.ADB.APKS[ "hulu" ][ s.Config.ADB.DeviceType ].Package )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenPackage( s.Config.ADB.APKS[ "hulu" ][ s.Config.ADB.DeviceType ].Package )
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
	s.ADBWakeup()
	if strings.Contains( s.Status.ADB.Activity , "hulu" ) {
		log.Debug( fmt.Sprintf( "hulu was already open with activity %s" , s.Status.ADB.Activity ) )
		return
	}
	log.Debug( "hulu was NOT already open" )
	// log.Debug( "we have to force app to reopen to get correct deep uri linking" )
	s.HuluReopenApp()
	switch s.Config.ADB.DeviceType {
		case "firecube":
			log.Debug( "waiting on green pixel in top right" )
			result := s.ADB.WaitOnPixelColor( 1694 , 96 , color.RGBA{ R: 28 , G: 231 , B: 131 , A: 255 } , 10 * time.Second )
			if result == false {
				log.Debug( "never found green pixel" )
				break;
			}
			fmt.Println( "assuming this means we are on the profile selection screen" )
			log.Debug( "selecting hulu profile" , " " , s.Config.HuluUserProfileIndex )
			for i := 0; i < s.Config.HuluTotalUserProfiles; i++ {
				s.ADB.Down()
				time.Sleep( 200 * time.Millisecond )
			}
			s.ADB.Up()
			time.Sleep( 200 * time.Millisecond )
			s.ADB.Enter()
			break
		case "firestick":
			log.Debug( "sleeping 10 seconds to wait on app init" )
			time.Sleep( 10 * time.Second )
			break;
		case "firetablet":
			log.Debug( "sleeping 10 seconds to wait on app init" )
			time.Sleep( 10 * time.Second )
			break;
	}
}

func parse_hulu_sent_id( sent_id string ) ( uri string ) {
	if utils.IsUUID( sent_id ) {
		uri = fmt.Sprintf( "https://www.hulu.com/watch/%s" , sent_id )
		return
	}
	is_url , _ := utils.IsURL( sent_id )
	if is_url {
		uri = sent_id
		return
	}
	return
}

func ( s *Server ) HuluOpenURI( uri string ) {
	log.Debug( fmt.Sprintf( "HuluOpenURI( %s )" , uri ) )
	s.HuluContinuousOpen()
	s.ADB.OpenURI( uri )
	// verified_now_playing := false
	// verified_now_playing_updated_time := 0
	log.Debug( "waiting 20 seconds for hulu player to appear" )
	players := s.ADB.WaitOnPlayers( "hulu" , 20 )
	if len( players ) < 1 {
		log.Debug( "never started playing , we might have to try play button" )
	}
	log.Debug( "hulu player should be ready" )
	utils.PrettyPrint( players )
	switch s.Config.ADB.DeviceType {
		case "firecube":
			log.Debug( "hulu app never reports playback positions" )
			log.Debug( "we need hdmi passthrough frames and sound" )
			// ... so we have to do complex screenshot to verify if we are already playing ,
			// or if we are staged on the info page
			// press up arrow to trigger potential ui overlay , then check pixels
			time.Sleep( 3 * time.Second )
			s.ADB.Up()
			time.Sleep( 100 * time.Millisecond )
			s.ADB.Down()
			time.Sleep( 100 * time.Millisecond )
			s.ADB.Enter()
			break;
		case "firestick":
			break;
		case "firetablet":
			log.Debug( "entering fullscreen" )
			s.ADB.Tap( 302 , 191 )
			time.Sleep( 500 * time.Millisecond )
			s.ADB.Tap( 573 , 58 )
			break;
	}
	// media_session := s.ADB.GetMediaSessionInfo()
	// playback_positions := s.ADB.GetPlaybackPositions()
	// utils.PrettyPrint( media_session )
	// utils.PrettyPrint( playback_positions )
	// log.Debug( "waiting 20 seconds to see if hulu auto starts playing" )
	// playing := s.ADB.WaitOnPlayersPlaying( "hulu" , 20 )
	// if len( playing ) < 1 {
	// 	log.Debug( "never started playing , we might have to try play button" )
	// }
	// utils.PrettyPrint( playing )
	// log.Debug( fmt.Sprintf( "total now playing === %d" , len( playing ) ) )
	// for _ , player := range playing {
	// 	if player.Updated > 0 {
	// 		log.Debug( "hulu autostarted playing on it's own" )
	// 		verified_now_playing = true
	// 		verified_now_playing_updated_time = player.Updated
	// 		break
	// 	}
	// }
	// if verified_now_playing == false {
	// 	log.Debug( "hulu didn't auto start playing , we might have to try play button" )
	// 	return
	// }
	// log.Debug( "waiting now for player progress" )
	// s.ADB.WaitOnPlayersUpdated( "hulu" , verified_now_playing_updated_time , 60 )
	// log.Debug( "player progress should be ready" )
	// time.Sleep( 3 * time.Second )
}

func ( s *Server ) HuluID( c *fiber.Ctx ) ( error ) {
	sent_id := c.Params( "*" )
	sent_query := c.Request().URI().QueryArgs().String()
	if sent_query != "" { sent_id += "?" + sent_query }
	uri := parse_hulu_sent_id( sent_id )
	log.Debug( fmt.Sprintf( "HuluID( %s )" , uri ) )
	s.HuluOpenURI( uri )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_uri" , sent_id )
	return c.JSON( fiber.Map{
		"url": "/hulu/:id" ,
		"id": sent_id ,
		"result": true ,
	})
}

func ( s *Server ) HuluMovieNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "HuluMovieNext()" )
	// next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	// uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	// log.Debug( uri )
	// s.ADB.OpenURI( uri )
	// s.ADB.Key( "KEYCODE_DPAD_RIGHT" )
	// s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	// s.Set( "active_player_now_playing_id" , next_movie )
	// s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	return c.JSON( fiber.Map{
		"url": "/hulu/next" ,
		"result": true ,
	})
}

func ( s *Server ) HuluMoviePrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "HuluMoviePrevious()" )
	return c.JSON( fiber.Map{
		"url": "/hulu/previous" ,
		"result": true ,
	})
}

func ( s *Server ) HuluTVID( c *fiber.Ctx ) ( error ) {
	series_id := c.Params( "series_id" )
	// log.Debug( fmt.Sprintf( "HuluTVID( %s )" , series_id ) )
	if series_id == "" {
		return c.JSON( fiber.Map{
			"url": "/hulu/tv/:series_id" ,
			"series_id": series_id ,
			"result": false ,
		})
	}
	_ , series_exists := s.Config.Library.Hulu.TV[ series_id ]
	if series_exists == false {
		return c.JSON( fiber.Map{
			"url": "/hulu/tv/:series_id" ,
			"series_id": series_id ,
			"error": "series doesn't exist in library" ,
			"result": false ,
		})
	}
	// s.Set( "STATE.HULU.NOW_PLAYING.MODE" , "TV" )
	s.Set( "STATE.HULU.NOW_PLAYING.TV.SERIES_ID" , series_id )
	next_episode := circular_set.Next( s.DB , fmt.Sprintf( "LIBRARY.HULU.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.HULU.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "HuluTVID( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	uri := fmt.Sprintf( "https://www.hulu.com/watch/%s" , next_episode )
	log.Debug( uri )
	s.HuluOpenURI( uri )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , uri )
	return c.JSON( fiber.Map{
		"url": "/hulu/tv/:series_id" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) HuluTVNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "HuluTVNext()" )
	series_id := s.Get( "STATE.HULU.NOW_PLAYING.TV.SERIES_ID" )
	next_episode := circular_set.Next( s.DB , fmt.Sprintf( "LIBRARY.HULU.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.HULU.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "HuluTVNext( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	uri := fmt.Sprintf( "https://www.hulu.com/watch/%s" , next_episode )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , uri )
	return c.JSON( fiber.Map{
		"url": "/hulu/tv/:id/next" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) HuluTVPrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "HuluTVPrevious()" )
	series_id := s.Get( "STATE.HULU.NOW_PLAYING.TV.SERIES_ID" )
	next_episode := circular_set.Previous( s.DB , fmt.Sprintf( "LIBRARY.HULU.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.HULU.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "HuluTVPrevious( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	uri := fmt.Sprintf( "https://www.hulu.com/watch/%s" , next_episode )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , uri )
	return c.JSON( fiber.Map{
		"url": "/hulu/tv/:id/previous" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}