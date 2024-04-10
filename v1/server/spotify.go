package server

import (
	"fmt"
	"time"
	"math"
	"strings"
	"image/color"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// const SPOTIFY_ACTIVITY = "com.spotify.tv.android/com.spotify.tv.android.SpotifyTVActivity"
// const SPOTIFY_APP_NAME = "com.spotify.tv.android"

// func reopen_spotify( adb *adb_wrapper.Wrapper ) {
// 	adb.StopAllPackages()
// 	adb.ClosePackage( "com.spotify.tv.android" )
// 	// time.Sleep( 1 * time.Second )
// 	time.Sleep( 500 * time.Millisecond )
// 	adb.OpenActivity( ACTIVITY_SPOTIFY )
// 	time.Sleep( 500 * time.Millisecond )
// 	// time.Sleep( 1 * time.Second )
// }

func ( s *Server ) SpotifyGetActiveButtonIndex() ( result int ) {
	log.Debug( "SpotifyGetActiveButtonIndex()" )
	result = -1

	// flatten ui-selected pixels
	var ui_selected_pixel_keys []string
	var ui_selected_pixel_colors []color.RGBA
	var ui_selected_pixel_coords []adb_wrapper.Coord
	for key , coord := range s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Pixels[ "ui_selected" ] {
		ui_selected_pixel_keys = append( ui_selected_pixel_keys , key )
		ui_selected_pixel_coords = append( ui_selected_pixel_coords , adb_wrapper.Coord{ X: coord.X , Y: coord.Y } )
		ui_selected_pixel_colors = append( ui_selected_pixel_colors , utils.HexToRGBColor( coord.Color ) )
	}

	screenshot_bytes := s.ADB.ScreenshotToBytes()
	result_list := s.ADB.GetPixelColorsFromImageBytes( &screenshot_bytes , ui_selected_pixel_coords )
	for i , color := range result_list {
		if color == ui_selected_pixel_colors[ i ] {
			fmt.Println( "found ui selected" , i , ui_selected_pixel_keys[ i ] )
			result = i
			return
		}
	}
	return
}

func ( s *Server ) SpotifyGetActiveButtonIndexFromScreenshotBytes( screenshot_bytes *[]byte ) ( result int ) {
	log.Debug( "SpotifyGetActiveButtonIndexFromScreenshotBytes()" )
	result = -1
	ui_selected_map := map[string]int{
		"favorite": 0 ,
		"shuffle": 1 ,
		"previous": 2 ,
		"play_pause": 3 ,
		"next": 4 ,
		"loop": 5 ,
	}
	for key , pixel := range s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Pixels[ "ui_selected" ] {
		color_test := s.ADB.GetPixelColorFromImageBytes( screenshot_bytes , pixel.X , pixel.Y )
		if color_test == utils.HexToRGBColor( pixel.Color ) {
			log.Debug( fmt.Sprintf( "found ui selected === %s === %s" , key , ui_selected_map[ key ] ) )
			result = ui_selected_map[ key ]
			return
		}
	}
	return
}

func ( s *Server ) SpotifyNavigateToUIIndexFromScreenshotBytes( screenshot_bytes *[]byte , ui_button_index int ) {
	log.Debug( "SpotifyNavigateToUIIndexFromScreenshotBytes()" )
	current_index := s.SpotifyGetActiveButtonIndexFromScreenshotBytes( screenshot_bytes )
	distance := int( math.Abs( float64( ui_button_index - current_index ) ) )
	log.Debug( fmt.Sprintf( "Current Index === %d === Distance from Target === %d" , current_index , distance ) )
	if current_index < ui_button_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance; i++ {
			log.Debug( "pressing right" )
			s.ADB.Right()
			time.Sleep( 100 * time.Millisecond )
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			s.ADB.Left()
			time.Sleep( 100 * time.Millisecond )
		}
	}
}

func ( s *Server ) SpotifyIsShuffleOn() ( result bool ) {
	result = false
	log.Debug( "SpotifyIsShuffleOn()" )
	log.Debug( "pressing left" )
	s.ADB.Left()
	shuffle_on := s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Pixels[ "shuffle" ][ "on" ]
	shuffle_smart := s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Pixels[ "shuffle" ][ "smart_shuffle" ]
	shuffle_position_color := s.ADB.GetPixelColor( shuffle_on.X , shuffle_on.Y )
	if shuffle_position_color == utils.HexToRGBColor( shuffle_on.Color ) {
		result = true
		return
	}
	if shuffle_position_color == utils.HexToRGBColor( shuffle_smart.Color ) {
		result = true
		return
	}
	return
}

func ( s *Server ) SpotifyPressPreviousButton() {
	log.Debug( "SpotifyPressPreviousButton()" )
	shuffle_on := s.SpotifyIsShuffleOn()
	button_index := 2
	index := s.SpotifyGetActiveButtonIndex()
	if index == button_index {
		s.ADB.Enter()
	}
	distance := int( math.Abs( float64( button_index - index ) ) )
	log.Debug( fmt.Sprintf( "Index === %d === Distance === %d" , index , distance ) )
	if index < button_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing right" )
			s.ADB.Right()
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			s.ADB.Left()
		}
	}
	log.Debug( "pressing enter" )
	s.ADB.Enter()
	if shuffle_on == true {
		log.Debug( "pressing enter" )
		s.ADB.Enter()
	}
	return
}

func ( s *Server ) SpotifyReopenApp() {
	log.Debug( "SpotifyReopenApp()" )
	s.ADB.StopAllPackages()
	// s.ADB.SetBrightness( 0 )
	s.ADB.ClosePackage( s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Package )
	s.ADB.OpenPackage( s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Package )
	// time.Sleep( 1500 * time.Millisecond )
	time.Sleep( 3000 * time.Millisecond )
	log.Debug( "Done" )
}

func ( s *Server ) SpotifyContinuousOpen() ( was_open bool ) {
	was_open = false
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "SpotifyContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "spotify" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	s.ADBWakeup()
	windows := s.ADB.GetWindowStack()
	for _ , window := range windows {
		activity_lower := strings.ToLower( window.Activity )
		if strings.Contains( activity_lower , "spotify" ) {
			log.Debug( "spotify was already open" )
			was_open = true
			return
		}
	}
	log.Debug( "spotify was NOT already open" )
	s.SpotifyReopenApp()
	time.Sleep( 6 * time.Second )
	return
}

func parse_spotify_sent_id( sent_id string ) ( uri string ) {
	if strings.HasPrefix( sent_id , "spotify:" ) {
		uri = sent_id
		return
	}
	is_url , _ := utils.IsURL( sent_id )
	if is_url {
		fmt.Println( "is url" )
		uri = sent_id
		return
	}
	return
}

func ( s *Server ) SpotifyOpenID( sent_id string ) {
	log.Debug( fmt.Sprintf( "SpotifyOpenID( %s )" , sent_id ) )
	was_open := s.SpotifyContinuousOpen()
	uri := parse_spotify_sent_id( sent_id )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	if was_open == true {
		s.ADB.Right()
	}
	s.SpotifyWaitOnNowPlaying()
}

func ( s *Server ) SpotifyWaitOnNowPlaying() {
	log.Debug( "waiting on now playing pixel" )
	now_playing_pixel := s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Pixels[ "now_playing" ][ "play_pause" ]
	s.ADB.WaitOnPixelColor(
		now_playing_pixel.X , now_playing_pixel.Y ,
		utils.HexToRGBColor( now_playing_pixel.Color ) ,
		( 10 * time.Second ) ,
	)
	log.Debug( "should be now playing" )
}

func ( s *Server ) SpotifyPlaylistWithShuffle( c *fiber.Ctx ) ( error ) {
	playlist_id := c.Params( "playlist_id" )
	log.Debug( fmt.Sprintf( "SpotifyPlaylistWithShuffle( %s )" , playlist_id ) )
	s.SpotifyContinuousOpen()
	// TODO === Need to Add TV Mute and Unmute
	// TODO === If Same Playlist don't open , just press next ? depends
	// s.ADB.SetVolume( 0 )
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , playlist_id )
	s.ADB.OpenURI( playlist_uri )
	s.SpotifyWaitOnNowPlaying()
	s.SpotifyShuffleOn()
	log.Debug( "pressing next" )
	s.ADB.Key( "KEYCODE_MEDIA_NEXT" ) // they sometimes force same song
	s.Set( "active_player_now_playing_id" , playlist_id )
	s.Set( "active_player_now_playing_text" , fmt.Sprintf( "playlist === %s" , s.Config.Library.Spotify.Playlists[ playlist_id ].Name ) )
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
	uri := fmt.Sprintf( "spotify:track:%s:play" , song_id )
	s.SpotifyOpenID( uri )
	s.Set( "active_player_now_playing_id" , song_id )
	s.Set( "active_player_now_playing_text" , fmt.Sprintf( "song id === %s" , song_id ) )
	return c.JSON( fiber.Map{
		"url": "/spotify/song/:song_id" ,
		"song_id": song_id ,
		"result": true ,
	})
}

func ( s *Server ) SpotifyPlaylist( c *fiber.Ctx ) ( error ) {
	playlist_id := c.Params( "playlist_id" )
	log.Debug( fmt.Sprintf( "SpotifyPlaylist( %s )" , playlist_id ) )
	uri := fmt.Sprintf( "spotify:playlist:%s:play" , playlist_id )
	s.SpotifyOpenID( uri )
	s.Set( "active_player_now_playing_text" , fmt.Sprintf( "playlist === %s" , s.Config.Library.Spotify.Playlists[ playlist_id ].Name ) )
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
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , playlist_id )
	s.SpotifyOpenID( playlist_uri )
	s.SpotifyShuffleOn()
	log.Debug( "pressing next" )
	s.ADB.Key( "KEYCODE_MEDIA_NEXT" ) // they sometimes force same song
	s.Set( "active_player_now_playing_id" , playlist_id )
	s.Set( "active_player_now_playing_text" , fmt.Sprintf( "playlist === %s" , s.Config.Library.Spotify.Playlists[ playlist_id ].Name ) )
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
	s.ADB.Left() // just activates ui
	was_on = false
	// shuffle_on := false
	// shuffle_smart := false
	shuffle_on_pixel := s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Pixels[ "shuffle" ][ "on" ]
	shuffle_smart_pixel := s.Config.ADB.APKS[ "spotify" ][ s.Config.ADB.DeviceType ].Pixels[ "shuffle" ][ "smart_shuffle" ]
	screenshot_bytes := s.ADB.ScreenshotToBytes()
	shuffle_position_color := s.ADB.GetPixelColorFromImageBytes( &screenshot_bytes , shuffle_on_pixel.X , shuffle_on_pixel.Y )
	shuffle_ui_index := 1
	if shuffle_position_color == utils.HexToRGBColor( shuffle_on_pixel.Color ) {
		log.Debug( "shuffle was already on" )
		// shuffle_on = true
		was_on = true
		return
	}
	s.SpotifyNavigateToUIIndexFromScreenshotBytes( &screenshot_bytes , shuffle_ui_index )
	if shuffle_position_color == utils.HexToRGBColor( shuffle_smart_pixel.Color ) {
		log.Debug( "smart shuffle was on" )
		// shuffle_smart = true
		was_on = true
		s.ADB.Enter()
		time.Sleep( 1 * time.Second )
		s.ADB.Enter()
		return
	}
	log.Debug( "shuffle was off , turning on" )
	s.ADB.Enter()
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
		s.ADB.Key( "KEYCODE_ENTER" )
	}
	distance := int( math.Abs( float64( shuffle_index - index ) ) )
	log.Debug( fmt.Sprintf( "Index === %d === Distance === %d" , index , distance ) )
	if index < shuffle_index {
		log.Debug( fmt.Sprintf( "pressing right %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing right" )
			s.ADB.Key( "KEYCODE_DPAD_RIGHT" )
		}
	} else {
		log.Debug( fmt.Sprintf( "pressing left %d times" , distance ) )
		for i := 0 ; i < distance ; i++ {
			log.Debug( "pressing left" )
			s.ADB.Key( "KEYCODE_DPAD_LEFT" )
		}
	}
	log.Debug( "pressing enter" )
	s.ADB.Key( "KEYCODE_ENTER" )
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