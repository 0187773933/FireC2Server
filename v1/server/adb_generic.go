package server

import (
	"fmt"
	"time"
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

func ( s * Server ) ADBWakeup() {
	if s.Status.ADB.DisplayOn == false {
		log.Debug( "display was off , turning on" )
		s.ADB.Wakeup()
		s.ADB.ForceScreenOn()
		time.Sleep( 500 * time.Millisecond )
		switch s.Config.ADB.DeviceType {
			case "firecube" , "firestick":
				s.ADB.Home()
				break;
			case "firetablet":
				s.ADB.Swipe( 513 , 564 , 553 , 171 )
		}
		time.Sleep( 1 * time.Second )
		s.GetStatus()
	}

	// if its the profile picker , select the profile
	if s.Status.ADB.Activity == ACTIVITY_PROFILE_PICKER {
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeUserProfileIndex ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
		time.Sleep( 1000 * time.Millisecond )
		s.GetStatus()
	}
}