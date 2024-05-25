package server

import (
	"fmt"
	"time"
	fiber "github.com/gofiber/fiber/v2"
	// utils "github.com/0187773933/FireC2Server/v1/utils"
)

func ( s *Server ) StreamDeckSpotify( c *fiber.Ctx ) ( error ) {
	go s.TV.QuickResetVideo()
	s.SpotifyNextPlaylistWithShuffle( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/spotify" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckYouTube( c *fiber.Ctx ) ( error ) {
	go s.TV.QuickResetVideo()
	s.YouTubeLiveNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/youtube" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckDisney( c *fiber.Ctx ) ( error ) {
	go s.TV.QuickResetVideo()
	s.DisneyMovieNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/disney" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckTwitch( c *fiber.Ctx ) ( error ) {
	go s.TV.QuickResetVideo()
	time_since_last_start := s.TimeSinceLastStart()
	fmt.Println( "time since last start ===" , time_since_last_start )
	if time_since_last_start > 30 * time.Minute {
		log.Debug( "Refreshing Twitch Environment" )
		s.TwitchLiveRefresh()
	}
	s.TwitchLiveNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/twitch" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckTwitchBackground( c *fiber.Ctx ) ( error ) {
	return c.JSON( fiber.Map{
		"url": "/streamdeck/twitch-background" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckTwitchUser( c *fiber.Ctx ) ( error ) {
	go s.TV.QuickResetVideo()
	user_id := c.Params( "user" )
	s.TwitchOpenID( user_id )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/twitch/:user" ,
		"user": user_id ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckNetflix( c *fiber.Ctx ) ( error ) {
	go s.TV.QuickResetVideo()
	s.NetflixMovieNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/netflix" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckHulu( c *fiber.Ctx ) ( error ) {
	go s.TV.QuickResetVideo()
	s.HuluMovieNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/hulu" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckAudioBook( c *fiber.Ctx ) ( error ) {
	return c.JSON( fiber.Map{
		"url": "/streamdeck/audio-book" ,
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


func ( s *Server ) StreamDeckTest( c *fiber.Ctx ) ( error ) {
	go s.TV.QuickResetVideo()
	s.SpotifyNextPlaylistWithShuffle( c )
	time.Sleep( 5 * time.Second )
	s.YouTubeLiveNext( c )
	time.Sleep( 5 * time.Second )
	s.DisneyMovieNext( c )
	time.Sleep( 5 * time.Second )
	s.TwitchLiveNext( c )
	time.Sleep( 5 * time.Second )
	s.NetflixMovieNext( c )
	time.Sleep( 5 * time.Second )
	s.HuluMovieNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/test" ,
		"result": true ,
	})
}