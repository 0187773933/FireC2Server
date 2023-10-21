package media_player

import (
	"fmt"
	"time"
	"math"
	// "strings"
	// bolt_api "github.com/boltdb/bolt"
	"image/color"
	// adb_wrapper "ADBWrapper/v1/wrapper"
	adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
	utils "github.com/0187773933/FireC2Server/v1/utils"
)

// func enable_shuffle( adb *adb_wrapper.Wrapper ) {

// }

const ACTIVITY_SPOTIFY = "com.spotify.tv.android/com.spotify.tv.android.SpotifyTVActivity"

func reopen_spotify( adb *adb_wrapper.Wrapper ) {
	adb.StopAllApps()
	adb.CloseAppName( "com.spotify.tv.android" )
	// time.Sleep( 1 * time.Second )
	time.Sleep( 500 * time.Millisecond )
	adb.OpenActivity( ACTIVITY_SPOTIFY )
	time.Sleep( 500 * time.Millisecond )
	// time.Sleep( 1 * time.Second )
}

type Spotify struct {
	*MediaPlayer
	Name string `yaml:"name"`
	// DB *bolt_api.DB `yaml:"-"`
}

func ( mp *Spotify ) GetActiveButtonIndex() ( result int ) {
	result = -1
	active_color := color.RGBA{ R: 255 , G: 255 , B: 255 , A: 255 }
	indexes := [][]int{
		{ 201 , 940 } ,
		{ 789 , 940 } ,
		{ 893 , 940 } ,
		{ 997 , 940 } ,
		{ 1101 , 940 } ,
		{ 1205 , 940 } ,
		{ 1793 , 940 } ,
	}
	screenshot := mp.ADB.ScreenshotToPNG()
	for index , coords := range indexes {
		pixel := screenshot.At( coords[ 0 ] , coords[ 1 ] )
		r , g , b , a := pixel.RGBA()
		pixel_rgba := color.RGBA{ R: uint8( r ) , G: uint8( g ) , B: uint8( b ) , A: uint8( a ) }
		if pixel_rgba == active_color {
			result = index
			return
		}
	}
	return
}

func ( mp *Spotify ) IsShuffleOn() ( result bool ) {
	active_color := color.RGBA{ R: 255 , G: 255 , B: 255 , A: 255 }
	coords := []int{ 752 , 964 }
	log.Debug( "pressing left" )
	mp.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
	time.Sleep( 500 * time.Millisecond )
	result = mp.ADB.IsPixelTheSameColor( coords[ 0 ] , coords[ 1 ] , active_color )
	return
}


func ( mp *Spotify ) PressPreviousButton() {
	shuffle_on := mp.IsShuffleOn()
	button_index := 2
	index := mp.GetActiveButtonIndex()
	if index == button_index {
		mp.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
	distance := int( math.Abs( float64( button_index - index ) ) )
	log.Debug( fmt.Sprintf( "Index === %d === Distance === %d" , index , distance ) )
	if index < button_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing right" )
			mp.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			mp.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		}
	}
	log.Debug( "pressing enter" )
	mp.ADB.PressKeyName( "KEYCODE_ENTER" )
	if shuffle_on == true {
		log.Debug( "pressing enter" )
		mp.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
	return
}

func ( mp *Spotify ) ReopenSpotifyApp() {
	mp.ADB.StopAllApps()
	mp.ADB.Brightness( 0 )
	mp.ADB.CloseAppName( "com.spotify.tv.android" )
	log.Debug( "Done" )
}

func ( mp *Spotify ) ContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "ContinuousOpen()" )
	log.Debug( mp.Status )
	mp.Set( "active_player_name" , "spotify" )
	mp.Set( "active_player_command" , "play" )
	mp.Set( "active_player_start_time" , start_time_string )
	if mp.Status.ADBTopWindow != ACTIVITY_SPOTIFY {
		log.Debug( "spotify was NOT already open" )
		mp.ReopenSpotifyApp()
	} else {
		log.Debug( "spotify was already open" )
	}
}

func ( mp *Spotify ) PlayPlaylistWithShuffle( playlist_id string ) {
	mp.ContinuousOpen()
	log.Debug( mp.Status.ADBVolume )
	// TODO === Need to Add TV Mute and Unmute
	// TODO === If Same Playlist don't open , just press next ? depends
	mp.ADB.SetVolume( 0 )
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , playlist_id )
	mp.ADB.OpenURI( playlist_uri )
	log.Debug( "Opened Playlist === " , playlist_id )
	// was_on := mp.ShuffleOn()
	mp.ShuffleOn()
	mp.Next() // they sometimes force same song
	mp.ADB.SetVolume( mp.Status.ADBVolume )
}

// 5Muvh0ooAJkSgBylFyI3su
func ( mp *Spotify ) Item( item_id string ) {
	mp.ADB.SetVolume( 0 )
	mp.ContinuousOpen()
	playlist_uri := fmt.Sprintf( "spotify:track:%s:play" , item_id )
	mp.ADB.OpenURI( playlist_uri )
	log.Debug( "Opened Playlist === " , playlist_id )
	mp.Next() // they sometimes force same song
	mp.ADB.SetVolume( mp.Status.ADBVolume )
}

func ( mp *Spotify ) Playlist( playlist_id string ) {
	go mp.TV.Prepare()
	mp.Status = s.GetStatus()
	mp.ADB.SetVolume( 0 )
	mp.ContinuousOpen()
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , playlist_id )
	mp.ADB.OpenURI( playlist_uri )
	log.Debug( "Opened Playlist === " , playlist_id )
	mp.ADB.SetVolume( mp.Status.ADBVolume )
}

func ( mp *Spotify ) NextItem( song_id string ) {
	log.Debug( "NextItem()" )
}

func ( mp *Spotify ) NextPlaylist( playlist_id string ) {
	log.Debug( "NextPlaylist()" )
}

func ( mp *Spotify ) PreviousItem( song_id string ) {
	log.Debug( "PreviousItem()" )
}

func ( mp *Spotify ) PreviousPlaylist( playlist_id string ) {
	log.Debug( "PreviousPlaylist()" )
}

func ( mp *Spotify ) ShuffleOn() ( was_on bool ) {
	log.Debug( "ShuffleOn()" )
	was_on = mp.IsShuffleOn()
	if was_on == true {
		log.Debug( "Shuffle === ON" )
		return
	}
	shuffle_index := 1
	index := mp.GetActiveButtonIndex()
	if index == shuffle_index {
		mp.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
	distance := int( math.Abs( float64( shuffle_index - index ) ) )
	log.Debug( "Shuffle === OFF" )
	log.Debug( fmt.Sprintf( "Index === %d === Distance === %d" , index , distance ) )
	if index < shuffle_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing right" )
			mp.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			mp.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		}
	}
	log.Debug( "pressing enter" )
	mp.ADB.PressKeyName( "KEYCODE_ENTER" )
	log.Debug( "Shuffle === ON" )
	return
}

func ( mp *Spotify ) ShuffleOff() ( was_on bool ) {
	log.Debug( "ShuffleOff()" )
	was_on = mp.IsShuffleOn()
	if was_on == false {
		log.Debug( "Shuffle === OFF" )
		return
	}
	shuffle_index := 1
	index := mp.GetActiveButtonIndex()
	if index == shuffle_index {
		mp.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
	distance := int( math.Abs( float64( shuffle_index - index ) ) )
	log.Debug( fmt.Sprintf( "Index === %d === Distance === %d" , index , distance ) )
	if index < shuffle_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing right" )
			mp.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			mp.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		}
	}
	log.Debug( "pressing enter" )
	mp.ADB.PressKeyName( "KEYCODE_ENTER" )
	log.Debug( "Shuffle === OFF" )
	return
}

// =============================================

func ( mp *Spotify ) Play() {
	log.Debug( "Play()" )
	mp.ADB.PressKeyName( "KEYCODE_MEDIA_PLAY" )
}

func ( mp *Spotify ) Pause() {
	log.Debug( "Pause()" )
	mp.ADB.PressKeyName( "KEYCODE_MEDIA_PAUSE" )
}

func ( mp *Spotify ) Stop() {
	log.Debug( "Stop()" )
	mp.ADB.PressKeyName( "KEYCODE_MEDIA_STOP" )
}

func ( mp *Spotify ) Next() {
	log.Debug( "Next()" )
	mp.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" )
}

func ( mp *Spotify ) Previous() {
	log.Debug( "Previous()" )
	// mp.ADB.PressKeyName( "KEYCODE_MEDIA_PREVIOUS" )
	mp.PressPreviousButton()
}

func ( mp *Spotify ) Teardown() {
	log.Debug( "Teardown()" )
}

func ( mp *Spotify ) Setup() {
	log.Debug( "Setup()" )
}

func ( mp *Spotify ) Update() {
	log.Debug( "Update()" )
}