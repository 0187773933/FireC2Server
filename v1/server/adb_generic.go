package server

import (
	fiber "github.com/gofiber/fiber/v2"
)

func ( s *Server ) ADBPlay( c *fiber.Ctx ) ( error ) {
	log.Debug( "ADBPlay()" )
	s.ADB.Key( "KEYCODE_MEDIA_PLAY" )
	return c.JSON( fiber.Map{
		"url": "/adb/play" ,
		"result": true ,
	})
}

func ( s *Server ) ADBPause( c *fiber.Ctx ) ( error ) {
	log.Debug( "ADBPause()" )
	s.ADB.Key( "KEYCODE_MEDIA_PAUSE" )
	return c.JSON( fiber.Map{
		"url": "/adb/pause" ,
		"result": true ,
	})
}

func ( s *Server ) ADBStop( c *fiber.Ctx ) ( error ) {
	log.Debug( "ADBStop()" )
	s.ADB.Key( "KEYCODE_MEDIA_STOP" )
	return c.JSON( fiber.Map{
		"url": "/adb/stop" ,
		"result": true ,
	})
}

func ( s *Server ) ADBNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "ADBNext()" )
	s.ADB.Key( "KEYCODE_MEDIA_NEXT" )
	return c.JSON( fiber.Map{
		"url": "/adb/next" ,
		"result": true ,
	})
}

func ( s *Server ) ADBPrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "ADBPrevious()" )
	s.ADB.Key( "KEYCODE_MEDIA_PREVIOUS" )
	return c.JSON( fiber.Map{
		"url": "/adb/previous" ,
		"result": true ,
	})
}