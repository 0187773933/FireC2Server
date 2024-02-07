package server

import (
	fmt "fmt"
	time "time"
	"encoding/base64"
	// url "net/url"
	// "math"
	// "image/color"
	// utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	// circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// https://stackoverflow.com/questions/3512198/need-command-line-to-start-web-browser-using-adb

// const FIREFOX_FOCUS_ACTIVITY = "org.mozilla.focus/org.mozilla.focus.activity.MainActivity"
// const FIREFOX_FOCUS_APP_NAME = "org.mozilla.focus"
// const FIREFOX_FOCUS_ACTIVITY = "org.mozilla.firefox/org.mozilla.fenix.customtabs.ExternalAppBrowserActivity"
// const FIREFOX_FOCUS_APP_NAME = "org.mozilla.firefox"
const BROWSER_ACTIVITY = "com.amazon.cloud9/com.amazon.slate.fire_tv.FireTvSlateActivity"
const BROWSER_APP_NAME = "com.amazon.cloud9"

func ( s *Server ) BrowserReopenApp() {
	log.Debug( "BrowserReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( BROWSER_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( BROWSER_APP_NAME )
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

func ( s *Server ) BrowserOpenURL( c *fiber.Ctx ) ( error ) {
	x_url := c.Params( "*" )
	log.Debug( fmt.Sprintf( "BrowserOpenURL( %s )" , x_url ) )
	s.BrowserReopenApp()
	time.Sleep( 1000 * time.Millisecond )
	s.ADB.Type( x_url )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_ENTER" )
	s.Set( "active_player_now_playing_id" , x_url )
	s.Set( "active_player_now_playing_uri" , "" )
	return c.JSON( fiber.Map{
		"url": "/browser/url/*" ,
		"param_url": x_url ,
		"result": true ,
	})
}

// http://localhost:5954/firefox-focus/audio/https://files.34353.org/AudioBooks/CarlosCastaneda/01-The-Teachings-of-Don-Juan-A-Yaqui-Way-of-Knowledge.mp3
// http://192.168.4.23:5954/firefox-focus/audio/https://files.34353.org/AudioBooks/CarlosCastaneda/01-The-Teachings-of-Don-Juan-A-Yaqui-Way-of-Knowledge.mp3
func ( s *Server ) GetBrowserAudioPlayer( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/browser_audio_player.html" )
}
func ( s *Server ) BrowserOpenAudioPlayer( c *fiber.Ctx ) ( error ) {
	x_url := c.Params( "*" )
	log.Debug( fmt.Sprintf( "BrowserOpenAudioPlayer( %s )" , x_url ) )
	x_url_b64 := base64.StdEncoding.EncodeToString( []byte( x_url ) )
	target_url := fmt.Sprintf(
		// "%s/firefox-focus/audio?k=%s\\&url=%s" ,
		"\"%s/browser/audio?k=%s&url=%s\"" ,
		s.Config.ServerPublicUrl ,
		s.Config.ServerAPIKey ,
		x_url_b64 ,
	)

	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( BROWSER_APP_NAME )
	s.ADB.Shell( "am" , "start" , "-a" , "android.intent.action.VIEW" , "-d" , target_url )

	fmt.Println( target_url )
	s.ADB.Type( target_url )

	// attempt to trigger play
	time.Sleep( 3000 * time.Millisecond )
	s.ADB.PressKey( 62 )

	s.Set( "active_player_now_playing_id" , x_url )
	s.Set( "active_player_now_playing_uri" , "" )
	return c.SendFile( "./v1/server/html/browser_audio_player.html" )
}

// https://9304d5ed.34353.org/Tracy%20Chapman%20with%20Luke%20Combs%20-%20Fast%20Car.mp4
func ( s *Server ) GetBrowserFocusVideoPlayer( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/browser_video_player.html" )
}
func ( s *Server ) BrowserOpenVideoPlayer( c *fiber.Ctx ) ( error ) {
	x_url := c.Params( "*" )
	log.Debug( fmt.Sprintf( "BrowserOpenVideoPlayer( %s )" , x_url ) )
	s.BrowserReopenApp()
	time.Sleep( 1000 * time.Millisecond )
	s.ADB.Type( x_url )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_ENTER" )
	s.Set( "active_player_now_playing_id" , x_url )
	s.Set( "active_player_now_playing_uri" , "" )
	return c.JSON( fiber.Map{
		"url": "/browser/video/*" ,
		"param_url": x_url ,
		"result": true ,
	})
}

