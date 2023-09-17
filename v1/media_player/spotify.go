package media_player

import (
	"fmt"
	"time"
	"strings"
	// bolt_api "github.com/boltdb/bolt"
	adb_wrapper "ADBWrapper/v1/wrapper"
	// utils "github.com/0187773933/FireC2Server/v1/utils"
)

func enable_shuffle( adb *adb_wrapper.Wrapper ) {

}

func playlist_with_shuffle( adb *adb_wrapper.Wrapper ) {
	adb.StopAllApps()
	adb.Brightness( 0 )
	adb.CloseAppName( "com.spotify.tv.android" )
	time.Sleep( 1 * time.Second )
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , "6ZFJWltDYCI0OVyXNleN9e" )
	adb.OpenURI( playlist_uri )

	// Enable Shuffle
	adb.WaitOnScreen( "./screenshots/spotify/playing.png" , ( 10 * time.Second ) , 945 , 925 , 30 , 30 )
	fmt.Println( "Ready" )
	time.Sleep( 1 * time.Second )
	shuffle_test := adb.ClosestScreenInList( []string{
			"./screenshots/spotify/shuffle_off.png" ,
			"./screenshots/spotify/shuffle_on.png" ,
		} ,
		735 , 925 , 35 , 45 ,
	)
	if strings.Contains( shuffle_test , "off" ) {
		adb.PressKeyName( "KEYCODE_DPAD_LEFT" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_DPAD_LEFT" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_ENTER" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_MEDIA_NEXT" )
		// adb.SetVolumePercent( 100 )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	}

}

type Spotify struct {
	*MediaPlayer
	Name string `yaml:"name"`
	// DB *bolt_api.DB `yaml:"-"`
}

func ( mp *Spotify ) Play() {
	fmt.Println( "Spotify -->Play()" )
	fmt.Println( mp.Status )
	// mp.Set( "active_player_name" , "spotify" )
	// mp.Set( "active_player_command" , "play" )
	// mp.Set( "active_player_start_time" , start_time )
	// playlist_with_shuffle( &mp.ADB )
}

func ( mp *Spotify ) Pause() {
	fmt.Println( "Spotify paused" )
}

func ( mp *Spotify ) Stop() {
	fmt.Println( "Spotify stopped" )
}

func ( mp *Spotify ) Next() {
	fmt.Println( "Spotify next" )
}

func ( mp *Spotify ) Previous() {
	fmt.Println( "Spotify previous" )
}

func ( mp *Spotify ) Teardown() {
	fmt.Println( "Spotify previous" )
}

func ( mp *Spotify ) Setup() {
	fmt.Println( "Spotify previous" )
}

func ( mp *Spotify ) Update() {
	fmt.Println( "Spotify previous" )
}