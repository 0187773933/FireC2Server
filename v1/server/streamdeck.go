package server

import (
	// "fmt"
	fiber "github.com/gofiber/fiber/v2"
	// utils "github.com/0187773933/FireC2Server/v1/utils"
)

func ( s *Server ) StreamDeckSpotify( c *fiber.Ctx ) ( error ) {
	go s.TV.Prepare()
	s.SpotifyNextPlaylistWithShuffle( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/spotify" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckYouTube( c *fiber.Ctx ) ( error ) {
	go s.TV.Prepare()
	s.YouTubeLiveNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/youtube" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckDisney( c *fiber.Ctx ) ( error ) {
	go s.TV.Prepare()
	s.DisneyMovieNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/disney" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckTwitch( c *fiber.Ctx ) ( error ) {
	go s.TV.Prepare()
	s.TwitchLiveNext( c )
	return c.JSON( fiber.Map{
		"url": "/streamdeck/twitch" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckEscapeRope( c *fiber.Ctx ) ( error ) {
	s.TV.Prepare()
	return c.JSON( fiber.Map{
		"url": "/streamdeck/escape-rope" ,
		"result": true ,
	})
}

func ( s *Server ) StreamDeckHeart( c *fiber.Ctx ) ( error ) {
	s.TV.Prepare()
	return c.JSON( fiber.Map{
		"url": "/streamdeck/heart" ,
		"result": true ,
	})
}