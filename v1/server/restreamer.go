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
	// url := fmt.Sprintf( "%s/que/url/%s?k=%s" , s.Config.ReStreamServerUrlLocal , x_url , s.Config.ReStreamServerAPIKey )
	url := fmt.Sprintf( "%s/i/%s/que/url/%s" , s.Config.ReStreamServerUrlLocal , s.Config.ReStreamServerAPIKey , x_url )
	log.Debug( url )

	// go utils.GetJSON( url , nil , nil )
	result := utils.GetJSON( url , nil , nil )
	result_map , _ := result.(map[string]interface{})
	stream_url , _ := result_map["stream_url"].(string)

	if stream_url == "" {
		return c.JSON( fiber.Map{
			"url": "/restream/url/:url" ,
			"param_url": x_url ,
			"result": false ,
		})
	}

	log.Debug( stream_url )

	s.ADB.OpenPackage( BROWSER_APP_NAME )
	s.BrowserReopenApp()
	time.Sleep( 10000 * time.Millisecond )
	// s.ADB.Enter()
	// time.Sleep( 10000 * time.Millisecond )
	// s.ADB.Type( x_url )
	s.ADB.OpenURI( stream_url )
	// time.Sleep( 1000 * time.Millisecond )
	// s.ADB.OpenURI( stream_url )
	// s.ADB.Up()


	// 2.) Call VLC Load https://ReStreamURL/hls/stream.m3u8
	// s.VLCContinuousOpen()
	// uri := fmt.Sprintf( "vlc://%s" , stream_url )
	// uri := fmt.Sprintf( "%s" , stream_url )
	// log.Debug( uri )
	// fmt.Println( "waiting 5 seconds for vlc to init" )
	// time.Sleep( 5 * time.Second )
	// s.ADB.OpenURI( uri )
	s.Set( "active_player_name" , "restream" )
	s.Set( "active_player_now_playing_id" , stream_url )
	s.Set( "active_player_now_playing_uri" , stream_url )

	// time.Sleep( 5 * time.Second )
	// s.ADB.Up()
	// 3.) If URL=TikTok Rotate Screen to Landscape ???

	// 4.) Return
	return c.JSON( fiber.Map{
		"url": "/restream/url/:url" ,
		"param_url": x_url ,
		"stream_url": stream_url ,
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