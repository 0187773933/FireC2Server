package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "math"
	// "image/color"
	// "strings"
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

func ( s *Server ) DisneySelectProfile() {
	log.Debug( "DisneySelectProfile()" )
	for i := 0; i < s.Config.DisneyTotalUserProfiles; i++ {
		s.ADB.Right()
		time.Sleep( 200 * time.Millisecond )
	}
	s.ADB.Left()
	time.Sleep( 200 * time.Millisecond )
	s.ADB.Enter()
}

func ( s *Server ) DisneySelectProfileIfNecessary() {
	pss := filepath.Join( s.Config.SaveFilesPath , "screenshots" , "disney" , "profile_selection.png" )
	distance := s.ADB.CurrentScreenSimilarityToReferenceImage( pss )
	if distance == -1 { log.Debug( "screenshot failed" ); return }
	if distance < 1.5 {
		log.Debug( fmt.Sprintf( "we are on the profile selection screen , %f" , distance ) )
		s.DisneySelectProfile()
	}
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


	// "switch" case on what the current activity is
	// if the power/screen is off?
	if s.Status.ADB.DisplayOn == false {
		log.Debug( "display was off , turning on" )
		s.ADB.Wakeup()
		s.ADB.ForceScreenOn()
		time.Sleep( 500 * time.Millisecond )
		switch s.Config.ADB.DeviceType {
			case "firecube":
			case "firestick":
				s.ADB.Home()
				break;
			case "firetablet":
				s.ADB.Swipe( 513 , 564 , 553 , 171 )
		}
	}

	// if its the profile picker , select the profile
	if s.Status.ADB.Activity == ACTIVITY_PROFILE_PICKER {
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeUserProfileIndex ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
		time.Sleep( 1000 * time.Millisecond )
	}
	// if we are already on the disney app , just return
	for _ , v := range s.Config.ADB.APKS[ "disney" ][ s.Config.ADB.DeviceType ].Activities {
		if s.Status.ADB.Activity == v {
			log.Debug( fmt.Sprintf( "disney was already open with activity %s" , v ) )
			s.DisneySelectProfileIfNecessary()
			return
		}
	}

	// if we are on th ehome screen , open the app
	log.Debug( "disney was NOT already open" )
	s.DisneyReopenApp()
	// time.Sleep( 500 * time.Millisecond )
	// pss_fp := filepath.Join( s.Config.SaveFilesPath , "screenshots" , "disney" , "profile_selection.png" )
	// log.Debug( "waiting 20 seconds on profile selection screen" )
	// s.ADB.WaitOnScreen( pss_fp , ( 20 * time.Second ) )
	// time.Sleep( 500 * time.Millisecond )
	// s.DisneySelectProfile()
}

func ( s *Server ) DisneyMovieNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "DisneyMovieNext()" )
	s.DisneyContinuousOpen()
	next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	next_movie_name := s.Config.Library.Disney.Movies.Currated[ next_movie ].Name
	log.Debug( fmt.Sprintf( "%s === %s" , next_movie_name , uri ) )
	s.ADB.OpenURI( uri )
	s.ADB.Right()
	s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	s.Set( "active_player_now_playing_id" , next_movie )
	s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	return c.JSON( fiber.Map{
		"url": "/disney/next" ,
		"uuid": next_movie ,
		"name": next_movie_name ,
		"result": true ,
	})
}

func ( s *Server ) DisneyMoviePrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "DisneyMoviePrevious()" )
	s.DisneyContinuousOpen()
	next_movie := circular_set.Previous( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	next_movie_name := s.Config.Library.Disney.Movies.Currated[ next_movie ].Name
	log.Debug( fmt.Sprintf( "%s === %s" , next_movie_name , uri ) )
	s.ADB.OpenURI( uri )
	s.ADB.Key( "KEYCODE_DPAD_RIGHT" )
	s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	s.Set( "active_player_now_playing_id" , next_movie )
	s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	return c.JSON( fiber.Map{
		"url": "/disney/previous" ,
		"uuid": next_movie ,
		"name": next_movie_name ,
		"result": true ,
	})
}

// because it behaves differently on /video vs /play
// its impossible to know which ids work with which prefix
// movie_id := parse_disney_sent_id( sent_id )
func parse_disney_sent_id( sent_id string ) ( result string ) {
	if utils.IsUUID( sent_id ) {
		result = fmt.Sprintf( "https://www.disneyplus.com/video/%s" , sent_id )
		return sent_id
	}
	is_url , _ := utils.IsURL( sent_id )
	if is_url {
		return sent_id
	}
	return
}

func ( s *Server ) DisneyOpenID( sent_id string ) {
	uri := parse_disney_sent_id( sent_id )
	next_movie_name := ""
	if _ , _ok := s.Config.Library.Disney.Movies.Currated[ sent_id ]; _ok {
		next_movie_name = s.Config.Library.Disney.Movies.Currated[ sent_id ].Name
	}
	log.Debug( fmt.Sprintf( "%s === %s" , next_movie_name , uri ) )
	s.Set( "STATE.DISNEY.NOW_PLAYING" , sent_id )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_text" , next_movie_name )
	verified_now_playing := false
	verified_now_playing_updated_time := 0
	// s.ADB.Shell(
	// 	"am" , "start" , "-c" , "android.intent.category.LEANBACK_LAUNCHER" ,
	// 	"-a" , "android.intent.action.VIEW" , "-d" , uri ,
	// 	// "-f" , "0x10808000" ,
	// 	// "-f" , "0x10008000" ,
	// 	"-e" , "source" , "30" , s.Config.ADB.APKS[ "disney" ][ s.Config.ADB.DeviceType ].Activities[ "source" ] ,
	// )
	s.ADB.OpenURI( uri )
	log.Debug( "waiting 20 seconds for disney player to appear" )
	players := s.ADB.WaitOnPlayers( "disney" , 20 )
	if len( players ) < 1 {
		log.Debug( "never started playing , we might have to try play button" )
	}
	log.Debug( "disney player should be ready" )
	utils.PrettyPrint( players )
	log.Debug( "waiting 10 seconds to see if disney auto starts playing" )
	playing := s.ADB.WaitOnPlayersPlaying( "disney" , 10 )
	if len( playing ) < 1 {
		log.Debug( "never started playing , we might have to try play button" )
	}
	utils.PrettyPrint( playing )
	log.Debug( fmt.Sprintf( "total now playing === %d" , len( playing ) ) )
	for _ , player := range playing {
		if player.Updated > 0 {
			fmt.Println( "disney autostarted playing on it's own" )
			verified_now_playing = true
			verified_now_playing_updated_time = player.Updated
			return
		}
	}
	if verified_now_playing == false {
		log.Debug( "disney didn't auto start playing , we might have to try play button" )
		return
	}
	log.Debug( "trying to force update adb playback state" )
	updated := s.ADB.WaitOnPlayersUpdated( "disney" , verified_now_playing_updated_time , 60 )
	utils.PrettyPrint( updated )
}

func ( s *Server ) DisneyID( c *fiber.Ctx ) ( error ) {
	sent_id := c.Params( "*" )
	sent_query := c.Request().URI().QueryArgs().String()
	if sent_query != "" { sent_id += "?" + sent_query }
	log.Debug( fmt.Sprintf( "DisneyID( %s )" , sent_id ) )
	s.DisneyContinuousOpen()
	s.DisneyOpenID( sent_id )
	return c.JSON( fiber.Map{
		"url": "/disney/movie/:id" ,
		// "uuid": movie_id ,
		// "name": name ,
		"result": true ,
	})
}
