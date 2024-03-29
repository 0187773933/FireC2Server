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
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeUserProfileIndex ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
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

// adb shell am start -c android.intent.category.LEANBACK_LAUNCHER -a android.intent.action.VIEW -d https://www.netflix.com/watch/81692458?trackId=255824129 -f 0x10808000 -e source 30 com.netflix.ninja/.MainActivity

// https://developer.android.com/reference/android/content/Intent

// Bitwise OR
// 0x10808000 =
// 0x10000000 = FLAG_ACTIVITY_NEW_TASK
// 0x00800000 = FLAG_ACTIVITY_EXCLUDE_FROM_RECENTS
// 0x00008000 = FLAG_ACTIVITY_CLEAR_TASK

// https://developer.android.com/reference/android/content/Intent#FLAG_ACTIVITY_NEW_TASK
func ( s *Server ) NetflixOpenID( id string ) {
	uri := fmt.Sprintf( "https://www.netflix.com/watch/%s" , id )
	log.Debug( uri )
	s.ADB.Shell(
		"am" , "start" , "-c" , "android.intent.category.LEANBACK_LAUNCHER" ,
		"-a" , "android.intent.action.VIEW" , "-d" , uri ,
		// "-f" , "0x10808000" ,
		"-f" , "0x10008000" ,
		"-e" , "source" , "30" , "com.netflix.ninja/.MainActivity" ,
	)
	positions := s.ADB.GetPlaybackPositions()
	position , ok := positions[ "netflix" ]
	timeout_init_seconds := 20
	if ok == false {
		log.Debug( "1 , not ready , waiting up to 20 seconds" )
		for i := 0; i < timeout_init_seconds; i++ {
			time.Sleep( 1 * time.Second )
			positions = s.ADB.GetPlaybackPositions()
			position , ok = positions[ "netflix" ]
			if ok == true {
				log.Debug( "2 , ready" )
				break
			}
			if i & 3 == 0 {
				log.Debug( "3 , pressing play button , just in case its stalled" )
				s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			}
		}
	}
	if ok == false {
		log.Debug( "4 , timed out , app crashed ?" )
		return
	}
	if position.State == "playing" {
		log.Debug( "5 , supposedly already playing" )
	} else {
		log.Debug( "6 , not playing , pressing play button" )
		s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
		for i := 0; i < timeout_init_seconds; i++ {
			time.Sleep( 1 * time.Second )
			positions = s.ADB.GetPlaybackPositions()
			position , ok = positions[ "netflix" ]
			if ok == true {
				log.Debug( "7 , ready" )
				break
			}
			if i & 3 == 0 {
				log.Debug( "8 , pressing play button , just in case its stalled" )
				s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			}
		}
	}
	fmt.Println( position )
	log.Debug( "9 , double checking we are playing , quick check" )
	for i := 0; i < 10; i++ {
		updated := s.ADB.GetUpdatedPlaybackPosition( position )
		fmt.Println( updated )
		if updated.Position != position.Position {
			log.Debug( "10 , playing" )
			return
		}
		time.Sleep( 500 * time.Millisecond )
	}

	log.Debug( "11 , still not playing , pressing play button again" )
	s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
	for i := 0; i < timeout_init_seconds; i++ {
		time.Sleep( 1 * time.Second )
		positions = s.ADB.GetPlaybackPositions()
		position , ok = positions[ "netflix" ]
		if ok == true {
			log.Debug( "12 , ready" )
			break
		}
	}

	time.Sleep( 1 * time.Second )

	// longer check
	log.Debug( "13 , double checking we are playing , longer check" )
	max_retries := 10
	for r := 0; r < max_retries; r++ {
		for i := 0; i < 10; i++ {
			updated := s.ADB.GetUpdatedPlaybackPosition( position )
			fmt.Println( updated )
			if updated.Position != position.Position {
				log.Debug( "14 , playing" )
				return
			}
			time.Sleep( 500 * time.Millisecond )
		}
		time.Sleep( 1 * time.Second )
	}
	log.Debug( "15 , Timeout reached, exiting." )
}

func ( s *Server ) NetflixMovieNext( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	log.Debug( "NetflixMovieNext()" )
	r_movie := "LIBRARY.NETFLIX.MOVIES"
	next_movie := circular_set.Next( s.DB , r_movie )
	next_movie_name := s.Get( fmt.Sprintf( "LIBRARY.NETFLIX.MOVIES.%s" , next_movie ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( next_movie )
	s.Set( "active_player_now_playing_id" , next_movie )
	s.Set( "active_player_now_playing_uri" , next_movie )
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/netflix/next" ,
		"next_movie_id": next_movie ,
		"next_movie_name": next_movie_name ,
		"result": true ,
	})
}

func ( s *Server ) NetflixMoviePrevious( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	log.Debug( "NetflixMoviePrevious()" )
	s.NetflixContinuousOpen()
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/netflix/previous" ,
		"result": true ,
	})
}

func ( s *Server ) NetflixTVID( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	series_id := c.Params( "series_id" )
	if series_id == "" {
		s.StateMutex.Unlock()
		return c.JSON( fiber.Map{
			"url": "/netflix/tv/:series_id" ,
			"series_id": series_id ,
			"result": false ,
		})
	}
	_ , series_exists := s.Config.Library.Netflix.TV[ series_id ]
	if series_exists == false {
		s.StateMutex.Unlock()
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
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/netflix/tv/:series_id" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) NetflixTVNext( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	log.Debug( "NetflixTVNext()" )
	series_id := s.Get( "STATE.NETFLIX.NOW_PLAYING.TV.SERIES_ID" )
	next_episode := circular_set.Next( s.DB , fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "NetflixTVNext( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( next_episode )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , next_episode )
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/netflix/tv/:id/next" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) NetflixTVPrevious( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	log.Debug( "NetflixTVPrevious()" )
	series_id := s.Get( "STATE.NETFLIX.NOW_PLAYING.TV.SERIES_ID" )
	previous_episode := circular_set.Previous( s.DB , fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s" , series_id ) )
	previous_episode_name := s.Get( fmt.Sprintf( "LIBRARY.NETFLIX.TV.%s.%s" , series_id , previous_episode ) )
	log.Debug( fmt.Sprintf( "NetflixTVPrevious( %s ) --> %s === " , series_id , previous_episode , previous_episode_name ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( previous_episode )
	s.Set( "active_player_now_playing_id" , previous_episode )
	s.Set( "active_player_now_playing_uri" , previous_episode )
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/netflix/tv/:id/previous" ,
		"series_id": series_id ,
		"previous_episode_id": previous_episode ,
		"previous_episode_name": previous_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) NetflixID( c *fiber.Ctx ) ( error ) {
	s.StateMutex.Lock()
	sent_id := c.Params( "*" )
	log.Debug( fmt.Sprintf( "NetflixID( %s )" , sent_id ) )
	s.NetflixContinuousOpen()
	s.NetflixOpenID( sent_id )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_uri" , sent_id )
	s.StateMutex.Unlock()
	return c.JSON( fiber.Map{
		"url": "/netflix/:id" ,
		"id": sent_id ,
		"result": true ,
	})
}