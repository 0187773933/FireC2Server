package server

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
)

func ( s *Server ) Play( c *fiber.Ctx ) ( error ) {
	log.Debug( "Play()" )
	player_name := s.Get( "active_player_name" )
	player_command := s.Get( "active_player_command" )
	adb_status := s.ADB.GetStatus()
	fmt.Println( player_name , player_command , adb_status )
	switch player_name {
		case "youtube":
			break;
		case "twitch":
			now_playing := s.Get( "STATE.TWITCH.LIVE.NOW_PLAYING" )
			fmt.Println( "last opened stream ===" , now_playing )
			break;
		case "disney":
			break;
		case "spotify":
			break;
		case "vlc":
			break;
	}
	// s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
	return c.JSON( fiber.Map{
		"url": "/play" ,
		"result": true ,
	})
}

// functional pause-play-resume
func ( s *Server ) Pause( c *fiber.Ctx ) ( error ) {
	log.Debug( "Pause()" )
	player_name := s.Get( "active_player_name" )
	player_command := s.Get( "active_player_command" )
	adb_status := s.ADB.GetStatus()
	player_state := adb_status.MediaSession.State
	fmt.Println( player_name , player_command , player_state )
	switch player_name {
		case "youtube":
			break;
		case "twitch":
			switch player_state {
				case "playing":
					s.ADB.PressKeyName( "KEYCODE_BACK" )
					break;
				case "stopped":
					// assume resuming
					last_played := s.Get( "STATE.TWITCH.LIVE.NOW_PLAYING" )
					fmt.Println( "last opened stream ===" , last_played )
					uri := fmt.Sprintf( "twitch://stream/%s" , last_played )
					s.ADB.OpenURI( uri )
					s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
					break;
			}
			break;
		case "disney":
			break;
		case "spotify":
			break;
		case "vlc":
			break;
	}
	// s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
	return c.JSON( fiber.Map{
		"url": "/pause" ,
		"result": true ,
	})
}

func ( s *Server ) Stop( c *fiber.Ctx ) ( error ) {
	log.Debug( "Stop()" )
	// s.ADB.PressKeyName( "KEYCODE_MEDIA_STOP" )
	return c.JSON( fiber.Map{
		"url": "/stop" ,
		"result": true ,
	})
}

func ( s *Server ) Next( c *fiber.Ctx ) ( error ) {
	log.Debug( "Next()" )
	// s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
	return c.JSON( fiber.Map{
		"url": "/next" ,
		"result": true ,
	})
}

func ( s *Server ) Previous( c *fiber.Ctx ) ( error ) {
	log.Debug( "Previous()" )
	// s.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
	return c.JSON( fiber.Map{
		"url": "/previous" ,
		"result": true ,
	})
}