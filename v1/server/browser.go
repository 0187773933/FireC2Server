package server

import (
	fmt "fmt"
	time "time"
	base64 "encoding/base64"
	json "encoding/json"
	// url "net/url"
	// "math"
	// "image/color"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	websocket "github.com/gofiber/contrib/websocket"
	// redis "github.com/redis/go-redis/v9"
	// circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// https://stackoverflow.com/questions/3512198/need-command-line-to-start-web-browser-using-adb

// const FIREFOX_FOCUS_ACTIVITY = "org.mozilla.focus/org.mozilla.focus.activity.MainActivity"
// const FIREFOX_FOCUS_APP_NAME = "org.mozilla.focus"
// const FIREFOX_ACTIVITY = "org.mozilla.firefox/org.mozilla.fenix.customtabs.ExternalAppBrowserActivity"
// const FIREFOX_APP_NAME = "org.mozilla.firefox"
const BROWSER_ACTIVITY = "com.amazon.cloud9/com.amazon.slate.fire_tv.FireTvSlateActivity"
const BROWSER_APP_NAME = "com.amazon.cloud9"

func ( s *Server ) BrowserReopenApp() {
	log.Debug( "BrowserReopenApp()" )
	if s.Status.ADB.Activity == ACTIVITY_PROFILE_PICKER {
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeUserProfileIndex ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
		time.Sleep( 1000 * time.Millisecond )
	}
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( BROWSER_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( BROWSER_APP_NAME )
	s.Set( "active_player_name" , "browser" )
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
	s.Set( "active_player_now_playing_uri" , "url" )
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

	// 1.) Init
	x_url := c.Params( "*" )
	log.Debug( fmt.Sprintf( "BrowserOpenAudioPlayer( %s )" , x_url ) )
	x_url_hash := utils.Sha256( x_url )
	x_url_b64 := base64.StdEncoding.EncodeToString( []byte( x_url ) )

	// 2.) See if Previously Played
	dbk := fmt.Sprintf( "HISTORY.BROWSER.AUDIO.%s" , x_url_hash )
	dbk_position := ( dbk + ".POSITION" )
	previous_position := s.Get( dbk_position )
	if previous_position == "" {
		previous_position = "0"
		s.Set( dbk_position , previous_position )
		s.Set( dbk , x_url )
	}

	target_url := fmt.Sprintf(
		// "%s/firefox-focus/audio?k=%s\\&url=%s" ,
		"\"%s/browser/audio?k=%s&p=%s&h=%s&url=%s\"" ,
		s.Config.ServerPublicUrl ,
		s.Config.BrowserAPIKey ,
		previous_position ,
		x_url_hash ,
		x_url_b64 ,
	)
	log.Debug( target_url )

	// 2.) Stop And Relaunch with Target URL
	// firefox focus glitches out
	// firefox regular forces portrait display
	// chrome requires google play tools
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( BROWSER_APP_NAME )
	s.ADB.Shell( "am" , "start" , "-a" , "android.intent.action.VIEW" , "-d" , target_url )

	// 3.) Press One of the JS Hooked Event Keys
	// Browsers try and block "autoplay"
	// TODO : browser message passing back here once ready for ADB press
	time.Sleep( 3000 * time.Millisecond )
	s.ADB.PressKey( 62 )

	// 4.) Press "Menu" Key twice to hide browser bar
	time.Sleep( 2000 * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_MENU" )
	time.Sleep( 200 * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_MENU" )

	s.Set( "active_player_now_playing_id" , x_url )
	s.Set( "active_player_now_playing_uri" , "url" )
	return c.JSON( fiber.Map{
		"url": "/browser/audio/*" ,
		"param_url": x_url ,
		"param_url_hash": x_url_hash ,
		"previous_position": previous_position ,
		"result": true ,
	})
}

func ( s *Server ) BrowserAudioPlayerSetPosition( c *fiber.Ctx ) ( error ) {
	x_hash := c.Params( "hash" )
	x_position := c.Params( "position" )
	log.Debug( fmt.Sprintf( "BrowserAudioPlayerSetPosition( %s , %s )" , x_hash , x_position ) )
	dbk := fmt.Sprintf( "HISTORY.BROWSER.AUDIO.%s" , x_hash )
	dbk_position := ( dbk + ".POSITION" )
	old_position := s.Get( dbk_position )
	// if its never had a position stored , because its never been started , then return
	if old_position == "" {
		return c.JSON( fiber.Map{
			"url": "/browser/audio/set/:hash/position/:position" ,
			"hash": x_hash ,
			"position": -1 ,
			"result": false ,
		})
	}
	s.Set( dbk_position , x_position )
	return c.JSON( fiber.Map{
		"url": "/browser/audio/set/:hash/position/:position" ,
		"hash": x_hash ,
		"position": x_position ,
		"result": true ,
	})
}

// https://9304d5ed.34353.org/Tracy%20Chapman%20with%20Luke%20Combs%20-%20Fast%20Car.mp4
func ( s *Server ) GetBrowserVideoPlayer( context *fiber.Ctx ) ( error ) {
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
	s.Set( "active_player_now_playing_uri" , "url" )
	return c.JSON( fiber.Map{
		"url": "/browser/video/*" ,
		"param_url": x_url ,
		"result": true ,
	})
}

// https://docs.gofiber.io/contrib/websocket/
func ( s *Server ) BrowserWebSocketHandler( c *websocket.Conn ) {
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt , msg , err = c.ReadMessage(); err != nil {
			log.Debug( "read:" , err , mt )
			break
		}
		// log.Debug( fmt.Sprintf( "recv: %s", msg ) )
		// if err = c.WriteMessage(mt, msg); err != nil {
		// 	log.Debug( "write:" , err )
		// 	break
		// }
		var decoded_message map[ string ]interface{}
		decode_error := json.Unmarshal( msg , &decoded_message )
		if decode_error != nil { log.Debug( "json decode error" , decode_error ); break; }
		fmt.Println( decoded_message )
	}
}
