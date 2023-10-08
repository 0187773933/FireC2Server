package server

import (
	"fmt"
	"time"
	"math"
	"image/color"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

const SPOTIFY_ACTIVITY = "com.spotify.tv.android/com.spotify.tv.android.SpotifyTVActivity"
const SPOTIFY_APP_NAME = "com.spotify.tv.android"

// func reopen_spotify( adb *adb_wrapper.Wrapper ) {
// 	adb.StopAllApps()
// 	adb.CloseAppName( "com.spotify.tv.android" )
// 	// time.Sleep( 1 * time.Second )
// 	time.Sleep( 500 * time.Millisecond )
// 	adb.OpenActivity( ACTIVITY_SPOTIFY )
// 	time.Sleep( 500 * time.Millisecond )
// 	// time.Sleep( 1 * time.Second )
// }

func ( s *Server ) SpotifyGetActiveButtonIndex() ( result int ) {
	log.Debug( "SpotifyGetActiveButtonIndex()" )
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
	screenshot := s.ADB.ScreenshotToPNG()
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

func ( s *Server ) SpotifyIsShuffleOn() ( result bool ) {
	log.Debug( "SpotifyIsShuffleOn()" )
	active_color := color.RGBA{ R: 255 , G: 255 , B: 255 , A: 255 }
	coords := []int{ 752 , 964 }
	log.Debug( "pressing left" )
	s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
	time.Sleep( 500 * time.Millisecond )
	result = s.ADB.IsPixelTheSameColor( coords[ 0 ] , coords[ 1 ] , active_color )
	return
}

func ( s *Server ) SpotifyPressPreviousButton() {
	log.Debug( "SpotifyPressPreviousButton()" )
	shuffle_on := s.SpotifyIsShuffleOn()
	button_index := 2
	index := s.SpotifyGetActiveButtonIndex()
	if index == button_index {
		s.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
	distance := int( math.Abs( float64( button_index - index ) ) )
	log.Debug( fmt.Sprintf( "Index === %d === Distance === %d" , index , distance ) )
	if index < button_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing right" )
			s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		}
	}
	log.Debug( "pressing enter" )
	s.ADB.PressKeyName( "KEYCODE_ENTER" )
	if shuffle_on == true {
		log.Debug( "pressing enter" )
		s.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
	return
}

func ( s *Server ) SpotifyReopenApp() {
	log.Debug( "SpotifyReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( SPOTIFY_APP_NAME )
	s.ADB.OpenAppName( SPOTIFY_APP_NAME )
	log.Debug( "Done" )
}

func ( s *Server ) SpotifyContinuousOpen() {
	log.Debug( "SpotifyContinuousOpen()" )
	s.GetStatus()
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "ContinuousOpen()" )
	log.Debug( s.Status )
	s.Set( "active_player_name" , "spotify" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADBTopWindow ) )
	if s.Status.ADBTopWindow != SPOTIFY_ACTIVITY {
		log.Debug( "spotify was NOT already open" )
		s.SpotifyReopenApp()
		time.Sleep( 3 * time.Second )
	} else {
		log.Debug( "spotify was already open" )
	}
}

func ( s *Server ) SpotifyPlaylistWithShuffle( c *fiber.Ctx ) ( error ) {
	playlist_id := c.Params( "playlist_id" )
	log.Debug( fmt.Sprintf( "SpotifyPlaylistWithShuffle( %s )" , playlist_id ) )
	s.SpotifyContinuousOpen()
	// TODO === Need to Add TV Mute and Unmute
	// TODO === If Same Playlist don't open , just press next ? depends
	s.ADB.SetVolume( 0 )
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , playlist_id )
	s.ADB.OpenURI( playlist_uri )
	time.Sleep( 500 * time.Millisecond )
	// was_on := s.ShuffleOn()
	s.SpotifyShuffleOn()
	s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" ) // they sometimes force same song
	s.ADB.SetVolume( s.Status.ADBVolume )

	return c.JSON( fiber.Map{
		"url": "/spotify/playlist-shuffle/:playlist_id" ,
		"playlist_id": playlist_id ,
		"result": true ,
	})
}

// 5Muvh0ooAJkSgBylFyI3su
func ( s *Server ) SpotifySong( c *fiber.Ctx ) ( error ) {
	song_id := c.Params( "song_id" )
	log.Debug( fmt.Sprintf( "SpotifySong( %s )" , song_id ) )

	go s.TV.Prepare()
	s.Status = s.GetStatus()
	s.ADB.SetVolume( 0 )
	s.SpotifyContinuousOpen()
	uri := fmt.Sprintf( "spotify:track:%s:play" , song_id )
	s.ADB.OpenURI( uri )
	s.ADB.SetVolume( s.Status.ADBVolume )

	return c.JSON( fiber.Map{
		"url": "/spotify/song/:song_id" ,
		"song_id": song_id ,
		"result": true ,
	})
}

func ( s *Server ) SpotifyPlaylist( c *fiber.Ctx ) ( error ) {
	playlist_id := c.Params( "playlist_id" )
	log.Debug( fmt.Sprintf( "SpotifyPlaylist( %s )" , playlist_id ) )

	go s.TV.Prepare()
	s.Status = s.GetStatus()
	s.ADB.SetVolume( 0 )
	s.SpotifyContinuousOpen()
	uri := fmt.Sprintf( "spotify:playlist:%s:play" , playlist_id )
	s.ADB.OpenURI( uri )
	s.ADB.SetVolume( s.Status.ADBVolume )

	return c.JSON( fiber.Map{
		"url": "/spotify/playlist/:playlist_id" ,
		"playlist_id": playlist_id ,
		"result": true ,
	})
}

// circular_set.Add( s.DB , "LIBRARY.SPOTIFY.PLAYLISTS" , key )

func ( s *Server ) SpotifyNextSong( song_id string ) {
	log.Debug( "SpotifyNextSong()" )
}

func ( s *Server ) SpotifyNextPlaylist() {
	log.Debug( "SpotifyNextPlaylist()" )
}

func ( s *Server ) SpotifyNextPlaylistWithShuffle( c *fiber.Ctx ) ( error ) {
	playlist_id := circular_set.Next( s.DB , "LIBRARY.SPOTIFY.PLAYLISTS" )
	log.Debug( fmt.Sprintf( "SpotifyPlaylistWithShuffle( %s )" , playlist_id ) )
	s.SpotifyContinuousOpen()
	// TODO === Need to Add TV Mute and Unmute
	// TODO === If Same Playlist don't open , just press next ? depends
	s.ADB.SetVolume( 0 )
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , playlist_id )
	s.ADB.OpenURI( playlist_uri )
	// was_on := s.ShuffleOn()
	s.SpotifyShuffleOn()
	s.ADB.PressKeyName( "KEYCODE_MEDIA_NEXT" ) // they sometimes force same song
	s.ADB.SetVolume( s.Status.ADBVolume )

	return c.JSON( fiber.Map{
		"url": "/spotify/next/playlist-shuffle" ,
		"playlist_id": playlist_id ,
		"result": true ,
	})
}

func ( s *Server ) SpotifyPreviousSong( song_id string ) {
	log.Debug( "SpotifyPreviousSong()" )
}

func ( s *Server ) SpotifyPreviousPlaylist( playlist_id string ) {
	log.Debug( "SpotifyPreviousPlaylist()" )
}

func ( s *Server ) SpotifyShuffleOn() ( was_on bool ) {
	log.Debug( "SpotifyShuffleOn()" )
	was_on = s.SpotifyIsShuffleOn()
	if was_on == true {
		log.Debug( "Shuffle === ON" )
		return
	}
	shuffle_index := 1
	index := s.SpotifyGetActiveButtonIndex()
	if index == shuffle_index {
		s.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
	distance := int( math.Abs( float64( shuffle_index - index ) ) )
	log.Debug( "Shuffle === OFF" )
	log.Debug( fmt.Sprintf( "Index === %d === Distance === %d" , index , distance ) )
	if index < shuffle_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing right" )
			s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		}
	}
	log.Debug( "pressing enter" )
	s.ADB.PressKeyName( "KEYCODE_ENTER" )
	log.Debug( "Shuffle === ON" )
	return
}

func ( s *Server ) SpotifyShuffleOff() ( was_on bool ) {
	log.Debug( "SpotifyShuffleOff()" )
	was_on = s.SpotifyIsShuffleOn()
	if was_on == false {
		log.Debug( "Shuffle === OFF" )
		return
	}
	shuffle_index := 1
	index := s.SpotifyGetActiveButtonIndex()
	if index == shuffle_index {
		s.ADB.PressKeyName( "KEYCODE_ENTER" )
	}
	distance := int( math.Abs( float64( shuffle_index - index ) ) )
	log.Debug( fmt.Sprintf( "Index === %d === Distance === %d" , index , distance ) )
	if index < shuffle_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing right" )
			s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
		}
	}
	log.Debug( "pressing enter" )
	s.ADB.PressKeyName( "KEYCODE_ENTER" )
	log.Debug( "Shuffle === OFF" )
	return
}
func ( s *Server ) SpotifySetShuffle( c *fiber.Ctx ) error {
	log.Debug( "SpotifySetShuffle()" )
	state := c.Params( "state" )
	if state == "on" {
		s.SpotifyShuffleOn()
	} else {
		s.SpotifyShuffleOff()
	}
	return c.JSON( fiber.Map{
		"url": "/spotify/shuffle/:state" ,
		"state": state ,
		"result": true ,
	})

}

func ( s *Server ) SpotifyUpdate() {
	log.Debug( "SpotifyUpdate()" )
}