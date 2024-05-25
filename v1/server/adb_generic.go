package server

import (
	base64 "encoding/base64"
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

func ( s *Server ) ADBGetScreenshot( c *fiber.Ctx ) error {
	log.Debug( "ADBGetScreenshot()" )
	screenshot_bytes := s.ADB.ScreenshotToBytes()
	base_64_image := base64.StdEncoding.EncodeToString( screenshot_bytes )
	html_response := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Screenshot</title>
			<style>
				body {
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
					background-color: #f0f0f0;
				}
				img {
					max-width: 100%;
					max-height: 100%;
				}
			</style>
		</head>
		<body>
			<img src="data:image/png;base64,` + base_64_image + `" alt="Screenshot"/>
		</body>
		</html>
	`
	return c.Type( "html" ).SendString( html_response )
}