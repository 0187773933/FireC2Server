package server

import (
	// fmt "fmt"
	// time "time"
	// context "context"
	// utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	// circular_set "github.com/0187773933/RedisCircular/v1/set"
)

type PlaylistItem struct {
	Type string `json:"type"`
	SubType string `json:"sub_type"`
	URL string `json:"url"`
	URI string `json:"uri"`
	Position int `json:"position"`
	Watched bool `json:"watched"`
	TimesWatched int `json:"times_watched"`
	TimesSkipped int `json:"times_skipped"`
	TimesPlayed int `json:"times_played"`
}
type Playlist []PlaylistItem

func ( s *Server ) PlaylistAdd( c *fiber.Ctx ) ( error ) {
	// log.Debug( fmt.Sprintf( "PlaylistAdd( %s )" , username ) )
	return c.JSON( fiber.Map{
		"url": "/playlist/add" ,
		// "playlist_item": playlist_item ,
		"result": true ,
	})
}