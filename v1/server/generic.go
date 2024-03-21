package server

import (
	"fmt"
	"strings"
	fiber "github.com/gofiber/fiber/v2"
	utils "github.com/0187773933/FireC2Server/v1/utils"
)

var APP_NAMES = []string{ "twitch" , "disney" , "youtube" , "spotify" , "vlc" , "hulu" , "netflix" }

type GenericInfo struct {
	ADBStatus interface{} `json:"adb_status"`
	PlayerName string `json:"player_name"`
	Activity string `json:"activity"`
	Package string `json:"package"`
	PlayerState string `json:"player_state"`
	CachedPlayerName string `json:"cached_player_name"`
	CachedPlayerCommand string `json:"cached_player_name"`
}
func get_app_name( activity string ) ( result string ) {
	for _ , app_name := range APP_NAMES {
		if strings.Contains( activity , app_name ) { return app_name }
	}
	return
}

func ( s *Server ) generic_get_current_info() ( result GenericInfo ) {
	result.CachedPlayerName = s.Get( "active_player_name" )
	result.CachedPlayerCommand = s.Get( "active_player_command" )
	adb_status := s.ADB.GetStatus()
	result.ADBStatus = adb_status
	result.Activity = adb_status.Activity
	result.PlayerName = get_app_name( adb_status.Activity )
	result.Package = adb_status.MediaSession.Package
	result.PlayerState = adb_status.MediaSession.State
	fmt.Println( result )
	return
}

func ( s *Server ) Play( c *fiber.Ctx ) ( error ) {
	log.Debug( "Play()" )
	go s.TV.Prepare()
	info := s.generic_get_current_info()
	switch info.PlayerName {
		case "youtube":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			break;
		case "twitch":
			switch info.PlayerState {
				case "playing":
					log.Debug( "already playing" )
					break;
				case "stopped":
					last_played := s.Get( "STATE.TWITCH.LIVE.NOW_PLAYING" )
					fmt.Println( "last opened stream ===" , last_played )
					uri := fmt.Sprintf( "twitch://stream/%s" , last_played )
					s.ADB.OpenURI( uri )
					s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
					break;
			}
			break;
		case "disney":
			switch info.PlayerState {
				case "playing":
					log.Debug( "already playing" )
					break;
				case "paused":
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
				default:
					last_played := s.Get( "STATE.DISNEY.NOW_PLAYING" )
					uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , last_played )
					log.Debug( uri )
					s.ADB.OpenURI( uri )
					s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
					break;
			}
			break;
		case "spotify":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			break;
		case "vlc":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			break;
		case "hulu":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			break;
		case "netflix":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			break;
		case "browser":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			break;
		default:
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			break;
	}
	return c.JSON( fiber.Map{
		"url": "/play" ,
		"result": true ,
	})
}

// functional pause-play-resume
// func ( s *Server ) Pause( c *fiber.Ctx ) ( error ) {
// 	log.Debug( "Pause()" )
// 	info := s.generic_get_current_info()
// 	switch info.PlayerName {
// 		case "youtube":
// 			switch info.PlayerState {
// 				case "playing":
// 					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
// 					break;
// 				case "stopped":
// 					log.Debug( "youtube already stopped , not pausing" )
// 					break;
// 				case "paused":
// 					log.Debug( "youtube already paused , resuming" )
// 					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
// 					break;
// 			}
// 			break;
// 		case "twitch":
// 			switch info.PlayerState {
// 				case "playing":
// 					s.ADB.PressKeyName( "KEYCODE_BACK" )
// 					break;
// 				case "stopped":
// 					log.Debug( "twitch already paused , resuming" )
// 					last_played := s.Get( "STATE.TWITCH.LIVE.NOW_PLAYING" )
// 					fmt.Println( "last opened stream ===" , last_played )
// 					uri := fmt.Sprintf( "twitch://stream/%s" , last_played )
// 					s.ADB.OpenURI( uri )
// 					s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
// 					break;
// 			}
// 			break;
// 		case "disney":
// 			switch info.PlayerState {
// 				case "playing":
// 					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
// 					break;
// 				case "stopped":
// 					log.Debug( "disney already stopped , not pausing" )
// 					break;
// 				case "paused":
// 					log.Debug( "disney already paused , resuming" )
// 					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
// 					break;
// 			}
// 			break;
// 		case "spotify":
// 			switch info.PlayerState {
// 				case "playing":
// 					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
// 					break;
// 				case "stopped":
// 					log.Debug( "spotify already stopped , not pausing" )
// 					break;
// 				case "paused":
// 					log.Debug( "spotify already paused , resuming" )
// 					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
// 					break;
// 			}
// 			break;
// 		case "vlc":
// 			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
// 			break;
// 	}
// 	return c.JSON( fiber.Map{
// 		"url": "/pause" ,
// 		"result": true ,
// 	})
// }

func ( s *Server ) Pause( c *fiber.Ctx ) ( error ) {
	log.Debug( "Pause()" )
	info := s.generic_get_current_info()
	utils.PrettyPrint( info )
	switch info.PlayerName {
		case "youtube":
			switch info.PlayerState {
				case "playing":
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
					break;
				case "stopped":
					log.Debug( "youtube already stopped , not pausing" )
					break;
				case "paused":
					log.Debug( "youtube already paused" )
					break;
			}
			break;
		case "twitch":
			switch info.PlayerState {
				case "playing":
					s.ADB.PressKeyName( "KEYCODE_BACK" )
					break;
				case "stopped":
					log.Debug( "twitch already paused" )
					break;
			}
			break;
		case "disney":
			switch info.PlayerState {
				case "playing":
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
					break;
				case "stopped":
					log.Debug( "disney already stopped , not pausing" )
					break;
				case "paused":
					log.Debug( "disney already paused" )
					break;
			}
			break;
		case "spotify":
			switch info.PlayerState {
				case "playing":
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
					break;
				case "stopped":
					log.Debug( "spotify already stopped , not pausing" )
					break;
				case "paused":
					log.Debug( "spotify already paused" )
					break;
			}
			break;
		case "vlc":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
		case "hulu":
			switch info.PlayerState {
				case "playing":
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
					break;
				case "stopped":
					// log.Debug( "hulu already stopped , not pausing" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
					break;
				case "paused":
					log.Debug( "hulu already paused" )
					break;
			}
			break;
		case "netflix":
			switch info.PlayerState {
				case "playing":
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
					break;
				case "stopped":
					// log.Debug( "netflix already stopped , not pausing" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
					break;
				case "paused":
					log.Debug( "netflix already paused" )
					break;
			}
			break;
		case "browser":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
		default:
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
	}
	return c.JSON( fiber.Map{
		"url": "/pause" ,
		"result": true ,
	})
}

func ( s *Server ) Resume( c *fiber.Ctx ) ( error ) {
	log.Debug( "Resume()" )
	info := s.generic_get_current_info()
	switch info.PlayerName {
		case "youtube":
			switch info.PlayerState {
				case "playing":
					log.Debug( "youtube already playing" )
					break;
				case "stopped":
					log.Debug( "youtube stopped , not resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
				case "paused":
					log.Debug( "youtube paused , resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
			}
			break;
		case "twitch":
			switch info.PlayerState {
				case "playing":
					log.Debug( "twitch already playing" )
					break;
				case "stopped":
					log.Debug( "twitch was stopped , resuming" )
					last_played := s.Get( "STATE.TWITCH.LIVE.NOW_PLAYING" )
					fmt.Println( "last opened stream ===" , last_played )
					uri := fmt.Sprintf( "twitch://stream/%s" , last_played )
					s.ADB.OpenURI( uri )
					s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
					break;
				case "paused":
					log.Debug( "twitch paused , resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
			}
			break;
		case "disney":
			switch info.PlayerState {
				case "playing":
					log.Debug( "disney already playing" )
					break;
				case "stopped":
					log.Debug( "disney stopped , not resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
				case "paused":
					log.Debug( "disney paused , resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
			}
			break;
		case "spotify":
			switch info.PlayerState {
				case "playing":
					log.Debug( "spotify already playing" )
					break;
				case "stopped":
					log.Debug( "spotify stopped , not resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
				case "paused":
					log.Debug( "spotify paused , resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
			}
			break;
		case "vlc":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
		case "hulu":
			switch info.PlayerState {
				case "playing":
					log.Debug( "hulu already playing" )
					break;
				case "stopped":
					log.Debug( "hulu stopped , not resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
				case "paused":
					log.Debug( "hulu paused , resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
			}
			break;
		case "netflix":
			switch info.PlayerState {
				case "playing":
					log.Debug( "netflix already playing" )
					break;
				case "stopped":
					log.Debug( "netflix stopped , not resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
				case "paused":
					log.Debug( "netflix paused , resuming" )
					s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
					break;
			}
			break;
		default:
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
			break;
	}
	return c.JSON( fiber.Map{
		"url": "/pause" ,
		"result": true ,
	})
}

func ( s *Server ) Stop( c *fiber.Ctx ) ( error ) {
	log.Debug( "Stop()" )
	info := s.generic_get_current_info()
	utils.PrettyPrint( info )
	switch info.PlayerName {
		case "youtube":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_STOP" )
			break;
		case "twitch":
			switch info.PlayerState {
				case "playing":
					s.ADB.PressKeyName( "KEYCODE_BACK" )
					break;
				case "stopped":
					log.Debug( "twitch already stopped" )
			}
			break;
		case "disney":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
		case "spotify":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
		case "vlc":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
		case "hulu":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
		case "netflix":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
			break;
		case "browser":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_STOP" )
			break;
		default:
			s.ADB.PressKeyName( "KEYCODE_MEDIA_STOP" )
			break;
	}
	return c.JSON( fiber.Map{
		"url": "/stop" ,
		"result": true ,
	})
}

func ( s *Server ) Next( c *fiber.Ctx ) ( error ) {
	log.Debug( "Next()" )
	go s.TV.Prepare()
	info := s.generic_get_current_info()
	switch info.PlayerName {
		case "youtube":
			return s.YouTubeLiveNext( c )
			break;
		case "twitch":
			return s.TwitchLiveNext( c )
			break;
		case "disney":
			return s.DisneyMovieNext( c )
			break;
		case "spotify":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
			break;
		case "vlc":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
			break;
		case "hulu":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
			break;
		case "netflix":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
			break;
		case "browser":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
			break;
		default:
			s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
			break;
	}
	return c.JSON( fiber.Map{
		"url": "/next" ,
		"result": true ,
	})
}

func ( s *Server ) Previous( c *fiber.Ctx ) ( error ) {
	log.Debug( "Previous()" )
	go s.TV.Prepare()
	info := s.generic_get_current_info()
	switch info.PlayerName {
		case "youtube":
			return s.YouTubeLivePrevious( c )
			break;
		case "twitch":
			return s.TwitchLivePrevious( c )
			break;
		case "disney":
			return s.DisneyMoviePrevious( c )
			break;
		case "spotify":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
			break;
		case "vlc":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
			break;
		case "hulu":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
			break;
		case "netflix":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
			break;
		case "browser":
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
			break;
		default:
			s.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
			break;
	}
	return c.JSON( fiber.Map{
		"url": "/previous" ,
		"result": true ,
	})
}