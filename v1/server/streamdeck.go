package server

import (
	"fmt"
	"time"
	fiber "github.com/gofiber/fiber/v2"
	// utils "github.com/0187773933/FireC2Server/v1/utils"
)

func ( s *Server ) StreamDeckPrepare() {
	time_since_last_start := s.TimeSinceLastStart()
	fmt.Println( "time since last start ===" , time_since_last_start )
	if time_since_last_start > 30 * time.Minute {
		go s.TV.QuickResetVideo()
		s.ADB.PressKeyName( "KEYCODE_WAKEUP" )
	}
}

func ( s *Server ) StreamDeckSpotify( c *fiber.Ctx ) ( error ) {
	s.StreamDeckPrepare()
	s.SpotifyNextPlaylistWithShuffle( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/spotify" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckYouTube( c *fiber.Ctx ) ( error ) {
	s.StreamDeckPrepare()
	s.YouTubeLiveNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/youtube" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckDisney( c *fiber.Ctx ) ( error ) {
	s.StreamDeckPrepare()
	s.DisneyMovieNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/disney" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckTwitch( c *fiber.Ctx ) ( error ) {
	time_since_last_start := s.TimeSinceLastStart()
	fmt.Println( "time since last start ===" , time_since_last_start )
	if time_since_last_start > 30 * time.Minute {
		go s.TV.QuickResetVideo()
		s.TwitchLiveRefresh()
		s.ADB.PressKeyName( "KEYCODE_WAKEUP" )
	}
	s.TwitchLiveNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/twitch" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckEscapeRope( c *fiber.Ctx ) ( error ) {
	return c.JSON( fiber.Map{
		"url": "/streamdeck/escape-rope" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckHeart( c *fiber.Ctx ) ( error ) {
	return c.JSON( fiber.Map{
		"url": "/streamdeck/heart" ,
		"result": true ,
	})
}