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
	// circular_set "github.com/0187773933/RedisCircular/v1/set"
)


func ( s *Server ) ReStreamURL( c *fiber.Ctx ) ( error ) {
	x_url := c.Params( "*" )
	log.Debug( fmt.Sprintf( "ReStreamURL( %s )" , x_url ) )

	// 1.) Call https://ReStreamURL/que/url
	url := fmt.Sprintf( "%s/que/url/%s?k=%s" , s.Config.ReStreamServerUrl , s.Config.ReStreamServerAPIKey )
	log.Debug( url )
	utils.GetJSON( url , nil , nil )
	time.Sleep( 5 * time.Second )

	// 2.) Call VLC Load https://ReStreamURL/hls/stream.m3u8
	s.VLCContinuousOpen()
	uri := fmt.Sprintf( "vlc://%s/hls/stream.m3u8" , s.Config.ReStreamServerUrl )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.Set( "active_player_now_playing_id" , x_url )
	s.Set( "active_player_now_playing_uri" , uri )

	// 3.) If URL=TikTok Rotate Screen to Landscape ???

	// 4.) Return
	return c.JSON( fiber.Map{
		"url": "/restream/url/:url" ,
		"param_url": x_url ,
		"result": true ,
	})
}