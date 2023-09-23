package server

import (
	fiber "github.com/gofiber/fiber/v2"
)

func ( s *Server ) ADBPlay() {
	log.Debug( "ADBPlay()" )
	s.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
}
func ( s *Server ) Play( c *fiber.Ctx ) ( error ) {
	log.Debug( "Play()" )
	s.ADBPlay()
	return c.JSON( fiber.Map{
		"url": "/play" ,
		"result": true ,
	})
}

func ( s *Server ) ADBPause() {
	log.Debug( "ADBPause()" )
	s.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
}
func ( s *Server ) Pause( c *fiber.Ctx ) ( error ) {
	log.Debug( "Pause()" )
	s.ADBPause()
	return c.JSON( fiber.Map{
		"url": "/pause" ,
		"result": true ,
	})
}

func ( s *Server ) ADBStop() {
	log.Debug( "ADBStop()" )
	s.ADB.PressKeyName( "KEYCODE_MEDIA_STOP" )
}
func ( s *Server ) Stop( c *fiber.Ctx ) ( error ) {
	log.Debug( "Stop()" )
	s.ADBStop()
	return c.JSON( fiber.Map{
		"url": "/stop" ,
		"result": true ,
	})
}

func ( s *Server ) ADBNext() {
	log.Debug( "ADBNext()" )
	s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
}
func ( s *Server ) Next( c *fiber.Ctx ) ( error ) {
	log.Debug( "Next()" )
	s.ADBNext()
	return c.JSON( fiber.Map{
		"url": "/next" ,
		"result": true ,
	})
}

func ( s *Server ) ADBPrevious() {
	log.Debug( "ADBPrevious()" )
	s.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
}
func ( s *Server ) Previous( c *fiber.Ctx ) ( error ) {
	log.Debug( "Previous()" )
	s.ADBPrevious()
	return c.JSON( fiber.Map{
		"url": "/previous" ,
		"result": true ,
	})
}