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
	url := fmt.Sprintf( "%s/que/url/%s?k=%s" , s.Config.ReStreamServerUrlLocal , x_url , s.Config.ReStreamServerAPIKey )
	log.Debug( url )
	go utils.GetJSON( url , nil , nil )
	log.Debug( "Sleeping 15 Seconds for ReStream to Hopefully Be Ready" )
	time.Sleep( 15 * time.Second )

	// 2.) Call VLC Load https://ReStreamURL/hls/stream.m3u8
	s.VLCContinuousOpen()
	uri := fmt.Sprintf( "vlc://%s/%s/stream.m3u8" , s.Config.ReStreamServerUrl , s.Config.ReStreamServerHLSURLPrefix )
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

func ( s *Server ) ReStreamRestart( c *fiber.Ctx ) ( error ) {
	log.Debug( "ReStreamRestart()" )

	s.VLCContinuousOpen()
	uri := fmt.Sprintf( "vlc://%s/%s/stream.m3u8" , s.Config.ReStreamServerUrl , s.Config.ReStreamServerHLSURLPrefix )
	log.Debug( uri )
	s.ADB.OpenURI( uri )

	// it works , but you don't want it. PepeWideMode
	// 3.) If URL=TikTok Rotate Screen to Landscape ???
	// s.ADB.Shell( "settings" , "put" , "system" , "accelerometer_rotation" , "0" )
	// s.ADB.Landscape()

	// 4.) Return
	return c.JSON( fiber.Map{
		"url": "/restream/restart" ,
		"result": true ,
	})
}

func ( s *Server ) ReStreamStop( c *fiber.Ctx ) ( error ) {
	log.Debug( "ReStreamStop()" )
	url := fmt.Sprintf( "%s/stop?k=%s" , s.Config.ReStreamServerUrlLocal , s.Config.ReStreamServerAPIKey )
	log.Debug( url )
	utils.GetJSON( url , nil , nil )
	return c.JSON( fiber.Map{
		"url": "/restream/stop" ,
		"result": true ,
	})
}