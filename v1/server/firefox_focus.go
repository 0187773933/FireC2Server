package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "math"
	// "image/color"
	// utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	// circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// https://stackoverflow.com/questions/3512198/need-command-line-to-start-web-browser-using-adb

const FIREFOX_FOCUS_ACTIVITY = "org.mozilla.focus/org.mozilla.focus.activity.MainActivity"
const FIREFOX_FOCUS_APP_NAME = "org.mozilla.focus"

func ( s *Server ) FirefoxFocusReopenApp() {
	log.Debug( "FirefoxFocusReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( FIREFOX_FOCUS_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( FIREFOX_FOCUS_APP_NAME )
	log.Debug( "Done" )
}

// TODO HTML UI Wrapper for at least audio
// saves / restores play position
// https://stackoverflow.com/a/48903341
// adb shell am start \
// -n com.android.chrome/com.google.android.apps.chrome.Main \
// -a android.intent.action.VIEW -d 'file:///sdcard/lazer.html'

// adb shell am start \
// -n org.mozilla.focus/org.mozilla.focus.activity.MainActivity \
// -a android.intent.action.VIEW -d 'https://files.34353.org/AudioBooks/CarlosCastaneda/01-The-Teachings-of-Don-Juan-A-Yaqui-Way-of-Knowledge.mp3'

// am start \
// -n org.mozilla.focus/org.mozilla.focus.activity.MainActivity \
// -a android.intent.action.VIEW -d 'https://files.34353.org/AudioBooks/CarlosCastaneda/01-The-Teachings-of-Don-Juan-A-Yaqui-Way-of-Knowledge.mp3'

func ( s *Server ) FireFoxFocusOpenURL( c *fiber.Ctx ) ( error ) {
	x_url := c.Params( "*" )
	log.Debug( fmt.Sprintf( "FireFoxFocusOpenURL( %s )" , x_url ) )
	s.FirefoxFocusReopenApp()
	time.Sleep( 1000 * time.Millisecond )
	s.ADB.Type( x_url )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_ENTER" )
	s.Set( "active_player_now_playing_id" , x_url )
	s.Set( "active_player_now_playing_uri" , "" )
	return c.JSON( fiber.Map{
		"url": "/firefox-focus/url/*" ,
		"param_url": x_url ,
		"result": true ,
	})
}

func ( s *Server ) GetFireFoxFocusAudioPlayer( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/firefox_focus_audio_player.html" )
}

func ( s *Server ) GetFireFoxFocusVideoPlayer( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/firefox_focus_video_player.html" )
}
// func ( s *Server ) FireFoxFocusAudioPlayer( c *fiber.Ctx ) ( error ) {
// 	x_url := c.Params( "*" )
// 	log.Debug( fmt.Sprintf( "FireFoxFocusAudioPlayer( %s )" , x_url ) )
// 	s.FirefoxFocusReopenApp()
// 	time.Sleep( 1000 * time.Millisecond )
// 	s.ADB.Type( x_url )
// 	time.Sleep( 500 * time.Millisecond )
// 	s.ADB.PressKeyName( "KEYCODE_ENTER" )
// 	s.Set( "active_player_now_playing_id" , x_url )
// 	s.Set( "active_player_now_playing_uri" , "" )
// 	return c.JSON( fiber.Map{
// 		"url": "/firefox-focus/url/*" ,
// 		"param_url": x_url ,
// 		"result": true ,
// 	})
// }

