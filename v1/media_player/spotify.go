package media_player

import (
	"fmt"
	"time"
	"strings"
	bolt_api "github.com/boltdb/bolt"
	adb_wrapper "ADBWrapper/v1/wrapper"
)

func example_spotify( adb *adb_wrapper.Wrapper ) {

	// TODO : TV Volume Off
	// adb.SetVolumePercent( 0 )
	adb.StopAllApps()
	adb.Brightness( 0 )
	adb.CloseAppName( "com.spotify.tv.android" )
	time.Sleep( 1 * time.Second )
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , "6ZFJWltDYCI0OVyXNleN9e" )
	adb.OpenURI( playlist_uri )

	// Enable Shuffle
	// time.Sleep( 10 * time.Second )
	// ( 0 , 0 ) = Top-Left
	// adb.Screenshot( "./screenshots/spotify/shuffle_off.png" , 735 , 925 , 35 , 45 )

	// adb.PressKeyName( "KEYCODE_ENTER" )
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
	} else {
		// TODO : Turn TV Volume On
		// adb.SetVolumePercent( 100 )
	}

	time.Sleep( 10 * time.Second )
	adb.OpenURI( fmt.Sprintf( "spotify:playlist:%s:play" , "3UMDmO2YJb8DgUjpSBu8y9" ) )
	time.Sleep( 500 * time.Millisecond )
	adb.PressKeyName( "KEYCODE_MEDIA_NEXT" )

}

type Spotify struct {
	Name string `yaml:"name"`
	DB *bolt_api.DB `yaml:"-"`
}

func ( t *Spotify ) Play() {
	fmt.Println( "Spotify playing" )
	adb := adb_wrapper.ConnectIP(
		"/usr/local/bin/adb" ,
		"192.168.4.174" ,
		"5555" ,
	)
	example_spotify( &adb )
}

func ( t *Spotify ) Pause() {
	fmt.Println( "Spotify paused" )
}

func ( t *Spotify ) Stop() {
	fmt.Println( "Spotify stopped" )
}

func ( t *Spotify ) Next() {
	fmt.Println( "Spotify next" )
}

func ( t *Spotify ) Previous() {
	fmt.Println( "Spotify previous" )
}

func ( t *Spotify ) Teardown() {
	fmt.Println( "Spotify previous" )
}

func ( t *Spotify ) Setup() {
	fmt.Println( "Spotify previous" )
}

func ( t *Spotify ) Update() {
	fmt.Println( "Spotify previous" )
}