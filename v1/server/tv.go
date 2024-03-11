package server

import (
	// "fmt"
	fiber "github.com/gofiber/fiber/v2"
	utils "github.com/0187773933/FireC2Server/v1/utils"
)

func ( s *Server ) TVPrepare( c *fiber.Ctx ) ( error ) {
	s.TV.Prepare()
	return c.JSON( fiber.Map{
		"url": "/tv/prepare" ,
		"result": true ,
	})
}

func ( s *Server ) TVPowerOn( c *fiber.Ctx ) ( error ) {
	s.TV.PowerOn()
	return c.JSON( fiber.Map{
		"url": "/tv/power/on" ,
		"result": true ,
	})
}

func ( s *Server ) TVPowerOff( c *fiber.Ctx ) ( error ) {
	s.TV.PowerOff()
	return c.JSON( fiber.Map{
		"url": "/tv/power/off" ,
		"result": true ,
	})
}

func ( s *Server ) TVPowerStatus( c *fiber.Ctx ) ( error ) {
	result := s.TV.GetPowerStatus()
	return c.JSON( fiber.Map{
		"url": "/tv/power/status" ,
		"result": result ,
	})
}

func ( s *Server ) TVGetInput( c *fiber.Ctx ) ( error ) {
	result := s.TV.GetInput()
	return c.JSON( fiber.Map{
		"url": "/tv/input" ,
		"result": result ,
	})
}

func ( s *Server ) TVSetInput( c *fiber.Ctx ) ( error ) {
	input := c.Params( "input" )
	s.TV.SetInput( utils.StringToInt( input ) )
	return c.JSON( fiber.Map{
		"url": "/tv/input/:input" ,
		"result": true ,
	})
}

func ( s *Server ) TVMuteOn( c *fiber.Ctx ) ( error ) {
	s.TV.MuteOn()
	return c.JSON( fiber.Map{
		"url": "/tv/mute/on" ,
		"result": true ,
	})
}

func ( s *Server ) TVMuteOff( c *fiber.Ctx ) ( error ) {
	s.TV.MuteOff()
	return c.JSON( fiber.Map{
		"url": "/tv/mute/off" ,
		"result": true ,
	})
}

func ( s *Server ) TVGetVolume( c *fiber.Ctx ) ( error ) {
	result := s.TV.GetVolume()
	return c.JSON( fiber.Map{
		"url": "/tv/volume" ,
		"result": result ,
	})
}

func ( s *Server ) TVSetVolume( c *fiber.Ctx ) ( error ) {
	volume := c.Params( "volume" )
	s.TV.SetVolume( utils.StringToInt( volume ) )
	return c.JSON( fiber.Map{
		"url": "/tv/volume/:volume" ,
		"result": true ,
	})
}

func ( s *Server ) TVIRSendCode( c *fiber.Ctx ) ( error ) {
	code := c.Params( "code" )
	ir_code := s.Config.TV.IRConfig.Remotes[ s.Config.TV.IRConfig.DefaultRemote ].Keys[ code ].Code
	s.TV.IR.Transmit( ir_code )
	return c.JSON( fiber.Map{
		"url": "/tv/ir/:code" ,
		"code": code ,
		"result": true ,
	})
}