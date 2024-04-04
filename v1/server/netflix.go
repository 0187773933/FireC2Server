package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "math"
	// "image/color"
	"strings"
	// url "net/url"
	regexp "regexp"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	// adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// https://stackoverflow.com/questions/35556634/movie-deeplink-for-netflix-android-tv-app-com-netflix-ninja

// Intent netflix = new Intent();
// netflix.setAction(Intent.ACTION_VIEW);
// netflix.setData(Uri.parse("http://www.netflix.com/watch/70202141"));
// netflix.putExtra("source","30"); // careful: String, not int
// netflix.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK|Intent.FLAG_ACTIVITY_CLEAR_TASK);
// getActivity().startActivity(netflix);

// public void OpenNFX() {
//     Intent netflix = new Intent();
//     netflix.setAction(Intent.ACTION_VIEW);
//     netflix.setData(Uri.parse("http://www.netflix.com/watch/70291117"));
//     netflix.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK|Intent.FLAG_ACTIVITY_CLEAR_TASK);
//     getActivity().startActivity(netflix);
// }

// adb shell am start -c android.intent.category.LEANBACK_LAUNCHER -a android.intent.action.VIEW -d https://www.netflix.com/watch/81692458?trackId=255824129 -f 0x10808000 -e source 30 com.netflix.ninja/.MainActivity

// https://www.netflix.com/watch/81692459?trackId=255824129

// 1.) Series Title = https://www.netflix.com/title/80996601
// 2.) Episode = https://www.netflix.com/watch/81692464?trackId=14277283

// ( ()=> {
// 	try {
// 	document.querySelector( 'svg[data-name="ChevronDown"' ).dispatchEvent( new MouseEvent( "click" , {
// 		bubbles: true,
// 		cancelable: true,
// 		view: window
// 	}));
// 	} catch( e ) {}
// 	let episodes = document.querySelectorAll( "div.titleCardList--metadataWrapper" );
// 	let yaml_string = "";
// 	for ( let i = 0; i < episodes.length; ++i ) {
// 		let title_text = episodes[ i ].querySelector( "span.titleCard-title_text" ).innerText;
// 		let info_elem = episodes[ i ].querySelector( "div.ptrack-content" );
// 		// let extra_text = info_elem.innerText;
// 		let meta_data_ue = info_elem.getAttribute( "data-ui-tracking-context" );
// 		let meta_data_str = decodeURIComponent( meta_data_ue );
// 		let meta_data = JSON.parse( meta_data_str );
// 		let video_id = meta_data[ "video_id" ];
// 		let track_id = meta_data[ "track_id" ];
// 		let id = `${video_id}?trackId=${track_id}`;
// 		yaml_string += `          - id: ${id}\n`;
// 		yaml_string += `            name: "${title_text}"\n`;
// 	}
// 	console.log( yaml_string );
// })();

// const NETFLIX_ACTIVITY = "com.netflix.ninja/com.netflix.ninja.MainActivity"
// const NETFLIX_APP_NAME = "com.netflix.ninja"

func ( s *Server ) NetflixReopenApp() {
	log.Debug( "NetflixReopenApp()" )
	s.ADB.StopAllPackages()
	// s.ADB.SetBrightness( 0 )
	s.ADB.ClosePackage( s.Config.ADB.APKS[ "netflix" ][ s.Config.ADB.DeviceType ].Package )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenPackage( s.Config.ADB.APKS[ "netflix" ][ s.Config.ADB.DeviceType ].Package )
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
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeUserProfileIndex ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
		time.Sleep( 1000 * time.Millisecond )
	// } else if s.Status.ADB.Activity == s.Config.APKS[ "netflix" ][ "activity" ] {
	} else if s.Status.ADB.Activity == s.Config.ADB.APKS[ "netflix" ][ s.Config.ADB.DeviceType ].Activities[ "profile_selection" ] {
		log.Debug( "@ netflix profile selection screen" )
		for i := 0; i < s.Config.NetflixTotalUserProfiles; i++ {
			s.ADB.Left()
			time.Sleep( 100 * time.Millisecond )
		}
		for i := 0; i < s.Config.NetflixUserProfileIndex; i++ {
			s.ADB.Right()
			time.Sleep( 100 * time.Millisecond )
		}
		s.ADB.Enter()
	} else if strings.Contains( s.Status.ADB.Activity , "netflix" ) == true {
		log.Debug( "netflix was already open" )
	} else {
		log.Debug( "netflix was NOT already open" )
		log.Debug( s.Config.ADB.APKS[ "netflix" ][ s.Config.ADB.DeviceType ].Activities[ "main" ] )
		s.NetflixReopenApp()
	}
}

// adb shell am start -c android.intent.category.LEANBACK_LAUNCHER -a android.intent.action.VIEW -d https://www.netflix.com/watch/81692458?trackId=255824129 -f 0x10808000 -e source 30 com.netflix.ninja/.MainActivity

// https://developer.android.com/reference/android/content/Intent

// Bitwise OR
// 0x10808000 =
// 0x10000000 = FLAG_ACTIVITY_NEW_TASK
// 0x00800000 = FLAG_ACTIVITY_EXCLUDE_FROM_RECENTS
// 0x00008000 = FLAG_ACTIVITY_CLEAR_TASK

// 70143664?trackId=14170286

// https://developer.android.com/reference/android/content/Intent#FLAG_ACTIVITY_NEW_TASK

func parse_netflix_sent_id( sent_id string ) ( uri string ) {
	id_part_one := ""
	id_part_two := ""
	is_url , parsed_url := utils.IsURL( sent_id )
	pattern := `^\d+\?trackId=\d+(?:&.+)*$`
	re := regexp.MustCompile( pattern )
	if is_url == true {
		parsed_id := strings.TrimPrefix( parsed_url.Path , "/watch/" )
		id_part_one = strings.TrimSuffix( parsed_id , "?" )
		id_part_two = parsed_url.Query().Get( "trackId" )
	} else {
		matches := re.FindStringSubmatch( sent_id )
		if len( matches ) != 3 {
			log.Debug( "couldn't parse netflix id" , sent_id )
			return
		}
		id_part_one = matches[ 1 ]
		id_part_two = matches[ 2 ]
	}
	if id_part_one == "" || id_part_two == "" {
		log.Debug( fmt.Sprintf( "id_part_one === %s" , id_part_one ) )
		log.Debug( fmt.Sprintf( "id_part_two === %s" , id_part_two ) )
		log.Debug( "couldn't parse netflix id" , sent_id )
		return
	}
	uri = fmt.Sprintf( "https://www.netflix.com/watch/%s?trackId=%s" , id_part_one , id_part_two )
	return
}

func ( s *Server ) NetflixOpenID( sent_id string ) {
	uri := parse_netflix_sent_id( sent_id )
	log.Debug( uri )
	s.ADB.Shell(
		"am" , "start" , "-c" , "android.intent.category.LEANBACK_LAUNCHER" ,
		"-a" , "android.intent.action.VIEW" , "-d" , uri ,
		// "-f" , "0x10808000" ,
		"-f" , "0x10008000" ,
		"-e" , "source" , "30" , s.Config.ADB.APKS[ "netflix" ][ s.Config.ADB.DeviceType ].Activities[ "source" ] ,
	)
	log.Debug( "waiting 20 seconds for netflix player to appear" )
	netflix_players := s.ADB.WaitOnPlayers( "netflix" , 20 )
	if len( netflix_players ) < 1 {
		log.Debug( "never started playing , we might have to try play button" )
	}
	log.Debug( "netflix player should be ready" )
	utils.PrettyPrint( netflix_players )
	log.Debug( "waiting 10 seconds to see if netflix auto starts playing" )
	playing := s.ADB.WaitOnPlayersPlaying( "netflix" , 10 )
	if len( playing ) < 1 {
		log.Debug( "never started playing , we might have to try play button" )
	}
	utils.PrettyPrint( playing )
	log.Debug( fmt.Sprintf( "total now playing === %d" , len( playing ) ) )
	for _ , player := range playing {
		if player.Updated > 0 {
			fmt.Println( "netflix autostarted playing on it's own" )
			return
		}
	}
	log.Debug( "trying to force update adb playback state" )
	for _ , player := range playing {
		playing = s.ADB.WaitOnPlayersUpdatedForce( "netflix" , player.Updated , 60 )
		utils.PrettyPrint( playing )
	}
	log.Debug( "trying to force update adb playback state" )
	for _ , player := range playing {
		playing = s.ADB.WaitOnPlayersUpdatedForce( "netflix" , player.Updated , 60 )
		utils.PrettyPrint( playing )
	}
}

func ( s *Server ) NetflixMovieNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "NetflixMovieNext()" )
	r_movie := "LIBRARY.NETFLIX.MOVIES"
	next_movie := circular_set.Next( s.DB , r_movie )
	next_movie_name := s.Get( fmt.Sprintf( "LIBRARY.NETFLIX.MOVIES.%s" , next_movie ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( next_movie )
	s.Set( "active_player_now_playing_id" , next_movie )
	s.Set( "active_player_now_playing_uri" , next_movie )
	return c.JSON( fiber.Map{
		"url": "/netflix/next" ,
		"next_movie_id": next_movie ,
		"next_movie_name": next_movie_name ,
		"result": true ,
	})
}

func ( s *Server ) NetflixMoviePrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "NetflixMoviePrevious()" )
	s.NetflixContinuousOpen()
	return c.JSON( fiber.Map{
		"url": "/netflix/previous" ,
		"result": true ,
	})
}

func ( s *Server ) NetflixTVID( c *fiber.Ctx ) ( error ) {
	series_id := c.Params( "series_id" )
	if series_id == "" {
		return c.JSON( fiber.Map{
			"url": "/netflix/tv/:series_id" ,
			"series_id": series_id ,
			"result": false ,
		})
	}
	_ , series_exists := s.Config.Library.Netflix.TV[ series_id ]
	if series_exists == false {
		return c.JSON( fiber.Map{
			"url": "/netflix/tv/:series_id" ,
			"series_id": series_id ,
			"error": "series doesn't exist in library" ,
			"result": false ,
		})
	}
	// s.Set( "STATE.NETFLIX.NOW_PLAYING.MODE" , "TV" )
	s.Set( "STATE.NETFLIX.NOW_PLAYING.TV.SERIES_ID" , series_id )
	next_episode := circular_set.Next( s.DB , fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "NetflixTVID( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( next_episode )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , next_episode )
	return c.JSON( fiber.Map{
		"url": "/netflix/tv/:series_id" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) NetflixTVNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "NetflixTVNext()" )
	series_id := s.Get( "STATE.NETFLIX.NOW_PLAYING.TV.SERIES_ID" )
	next_episode := circular_set.Next( s.DB , fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "NetflixTVNext( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( next_episode )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , next_episode )
	return c.JSON( fiber.Map{
		"url": "/netflix/tv/:id/next" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) NetflixTVPrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "NetflixTVPrevious()" )
	series_id := s.Get( "STATE.NETFLIX.NOW_PLAYING.TV.SERIES_ID" )
	previous_episode := circular_set.Previous( s.DB , fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s" , series_id ) )
	previous_episode_name := s.Get( fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s.%s" , series_id , previous_episode ) )
	log.Debug( fmt.Sprintf( "NetflixTVPrevious( %s ) --> %s === " , series_id , previous_episode , previous_episode_name ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( previous_episode )
	s.Set( "active_player_now_playing_id" , previous_episode )
	s.Set( "active_player_now_playing_uri" , previous_episode )
	return c.JSON( fiber.Map{
		"url": "/netflix/tv/:id/previous" ,
		"series_id": series_id ,
		"previous_episode_id": previous_episode ,
		"previous_episode_name": previous_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) NetflixID( c *fiber.Ctx ) ( error ) {
	sent_id := c.Params( "*" )
	sent_query := c.Request().URI().QueryArgs().String()
	if sent_query != "" { sent_id += "?" + sent_query }
	log.Debug( fmt.Sprintf( "NetflixID( %s )" , sent_id ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( sent_id )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_uri" , sent_id )
	return c.JSON( fiber.Map{
		"url": "/netflix/:id" ,
		"id": sent_id ,
		"result": true ,
	})
}