package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "io/ioutil"
	// "math"
	color "image/color"
	// "strings"
	filepath "path/filepath"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	image_similarity "github.com/0187773933/ADBWrapper/v1/image-similarity"
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

func ( s *Server ) DisneySelectProfileIfNecessary() ( result bool ) {
	pss := filepath.Join( s.Config.SaveFilesPath , "screenshots" , "disney" , "profile_selection.png" )
	// new complex/raw screenshot stuff
	// screenshot_bytes := s.ADB.ScreenshotToBytes()
	// screenshot_features := s.ADB.ImageBytesToFeatures( &screenshot_bytes )

	distance := s.ADB.CurrentScreenSimilarityToReferenceImage( pss )
	// s.ADB.WaitOnScreen( pss_fp , ( 20 * time.Second ) )
	result = false
	// white_pixel := color.RGBA{ R: 255 , G: 255 , B: 255 , A: 255 }
	cross_add_profile_pixel := color.RGBA{ R: 188 , G: 189 , B: 193 , A: 255 }
	if distance == -1 { log.Debug( "screenshot failed" ); return }
	if distance < 1.5 {
		log.Debug( fmt.Sprintf( "we are on the profile selection screen , %f" , distance ) )
			switch s.Config.ADB.DeviceType {
				case "firecube" , "firestick":
					pixel_color := s.ADB.GetPixelColor( 1420 , 503 )
					if pixel_color == cross_add_profile_pixel {
						log.Debug( "we are on the profile selection screen" )
						s.DisneySelectProfile()
						result = true
					}
					break;
				case "firetablet":
					log.Debug( "TODO = find pixel coords for disney profile add on firetablet. we need it to verify" )
					s.DisneySelectProfile()
					result = true
					break
			}
		return
	}
	return
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
	s.ADBWakeup()
	log.Debug( "forcing relaunch of app" ) // leads to crashes if you don't
	s.DisneyReopenApp()
}

func ( s *Server ) DisneyMovieNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "DisneyMovieNext()" )
	s.DisneyContinuousOpen()
	next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	next_movie_name := s.Config.Library.Disney.Movies.Currated[ next_movie ].Name
	log.Debug( fmt.Sprintf( "%s === %s" , next_movie_name , uri ) )
	s.ADB.OpenURI( uri )
	s.DisneyOpenID( next_movie )
	s.ADB.Right()
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
	s.DisneyOpenID( next_movie )
	s.ADB.Right()
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
	s.ADB.OpenURI( uri )
	s.Set( "STATE.DISNEY.NOW_PLAYING" , sent_id )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_text" , next_movie_name )
	verified_now_playing := false
	verified_now_playing_updated_time := 0
	pss := filepath.Join( s.Config.SaveFilesPath , "screenshots" , "disney" , "profile_selection.png" )
	pss_features := image_similarity.GetFeatureVectorFromFilePath( pss )
	cross_add_profile_pixel := color.RGBA{ R: 188 , G: 189 , B: 193 , A: 255 }
	queries := 20
	stage_one_ready := false
	for i := 0; i < queries; i++ {
		// status := s.ADB.GetStatus()
		if stage_one_ready == true { break; }
		log.Debug( fmt.Sprintf( "checking [%d] of %d for disney to be ready" , ( i + 1 ) , queries ) )
		// activity := s.ADB.GetActivity()
		players := s.ADB.FindPlayers( "disney" )
		if len( players ) > 0 {
			log.Debug( "found disney player" )
			for _ , player := range players {
				if player.Updated > 0 {
					log.Debug( "disney autostarted playing on it's own" )
					verified_now_playing = true
					verified_now_playing_updated_time = player.Updated
					stage_one_ready = true
					break
				}
			}
			break
		}
		// on_profile_screen := s.DisneySelectProfileIfNecessary()
		// new complex/raw screenshot stuff
		screenshot_bytes := s.ADB.ScreenshotToBytes()
		// ioutil.WriteFile( "test2.png" , screenshot_bytes , 0644 )
		screenshot_features := s.ADB.ImageBytesToFeatures( &screenshot_bytes )
		distance := image_similarity.CalculateDistancePoint( &screenshot_features , &pss_features )
		if distance == -1 {
			log.Debug( "screenshot failed" );
		} else if distance < 1.5 {
			log.Debug( fmt.Sprintf( "we think we are on the profile selection screen , %f" , distance ) )
			switch s.Config.ADB.DeviceType {
				case "firecube" , "firestick":
					pixel_color := s.ADB.GetPixelColorFromImageBytes( &screenshot_bytes , 1423 , 492 )
					if pixel_color == cross_add_profile_pixel {
						log.Debug( "we are on the profile selection screen" )
						stage_one_ready = true
						s.DisneySelectProfile()
					}
					break;
				case "firetablet":
					log.Debug( "TODO = find pixel coords for disney profile add on firetablet. we need it to verify" )
					// s.DisneySelectProfile()
					// stage_one_ready = true
					break
			}
		}
		time.Sleep( 500 * time.Millisecond )
	}
	log.Debug( "ready stage 2" , " " , verified_now_playing , " " , verified_now_playing_updated_time )
	playing := s.ADB.WaitOnPlayersPlaying( "disney" , 10 )
	for _ , player := range playing {
		if player.Updated > 0 {
			log.Debug( "disney autostarted playing on it's own" )
			verified_now_playing = true
			verified_now_playing_updated_time = player.Updated
			break
		}
	}
	if verified_now_playing == false {
		log.Debug( "disney didn't auto start playing , we might have to try play button" )
		return
	}
	log.Debug( "trying to force update adb playback state" )
	updated := s.ADB.WaitOnPlayersUpdated( "disney" , verified_now_playing_updated_time , 20 )
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
		"url": "/disney/:id" ,
		"id": sent_id ,
		"result": true ,
	})
}
