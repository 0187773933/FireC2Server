package server

import (
	fmt "fmt"
	time "time"
	"strings"
	// url "net/url"
	// "math"
	// "image/color"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	// circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// const VLC_ACTIVITY = "org.videolan.vlc/org.videolan.television.ui.MainTvActivity"
// const VLC_APP_NAME = "org.videolan.vlc"

func ( s *Server ) VLCReopenApp() {
	log.Debug( "VLCReopenApp()" )
	s.ADB.StopAllPackages()
	// s.ADB.SetBrightness( 0 )
	s.ADB.ClosePackage( s.Config.ADB.APKS[ "vlc" ][ s.Config.ADB.DeviceType ].Package )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenPackage( s.Config.ADB.APKS[ "vlc" ][ s.Config.ADB.DeviceType ].Package )
	log.Debug( "Done" )
}

func ( s *Server ) VLCContinuousOpen() ( was_open bool ) {
	was_open = false
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "VLCContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "vlc" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	s.ADBWakeup()
	windows := s.ADB.GetWindowStack()
	for _ , window := range windows {
		activity_lower := strings.ToLower( window.Activity )
		if strings.Contains( activity_lower , "vlc" ) {
			log.Debug( "vlc was already open" )
			was_open = true
			return
		}
	}
	log.Debug( "vlc was NOT already open" )
	s.VLCReopenApp()
	time.Sleep( 1 * time.Second )
	return
}

func parse_vlc_sent_id( sent_id string ) ( uri string ) {
	if strings.HasPrefix( sent_id , "vlc://" ) {
		return sent_id
	}
	uri = fmt.Sprintf( "vlc://%s" , sent_id )
	return
}

func ( s *Server ) VLCOpenID( sent_id string ) {
	log.Debug( fmt.Sprintf( "VLCOpenID( %s )" , sent_id ) )
	s.VLCContinuousOpen()
	uri := parse_vlc_sent_id( sent_id )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
}

func ( s *Server ) VLCID( c *fiber.Ctx ) ( error ) {
	sent_id := c.Params( "*" )
	sent_query := c.Request().URI().QueryArgs().String()
	if sent_query != "" { sent_id += "?" + sent_query }
	log.Debug( fmt.Sprintf( "VLCID( %s )" , sent_id ) )
	s.VLCOpenID( sent_id )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_uri" , sent_id )
	return c.JSON( fiber.Map{
		"url": "/vlc/:id" ,
		"id": sent_id ,
		"result": true ,
	})
}

func ( s *Server ) VLCNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "VLCNext()" )
	// next_movie := circular_set.Next( s.DB , "LIBRARY.DISNEY.MOVIES.CURRATED" )
	// uri := fmt.Sprintf( "https://www.disneyplus.com/video/%s" , next_movie )
	// log.Debug( uri )
	// s.ADB.OpenURI( uri )
	// s.ADB.Key( "KEYCODE_DPAD_RIGHT" )
	// s.Set( "STATE.DISNEY.NOW_PLAYING" , next_movie )
	// s.Set( "active_player_now_playing_id" , next_movie )
	// s.Set( "active_player_now_playing_text" , s.Config.Library.Disney.Movies.Currated[ next_movie ].Name )
	return c.JSON( fiber.Map{
		"url": "/vlc/next" ,
		"result": true ,
	})
}

func ( s *Server ) VLCPrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "VLCPrevious()" )
	return c.JSON( fiber.Map{
		"url": "/vlc/previous" ,
		"result": true ,
	})
}

func ( s *Server ) VLCPlayURL( c *fiber.Ctx ) ( error ) {
	x_url := c.Params( "*" )
	log.Debug( fmt.Sprintf( "VLCPlayURL( %s )" , x_url ) )
	uri := fmt.Sprintf( "vlc://%s" , x_url )
	s.VLCOpenID( uri )
	s.Set( "active_player_now_playing_id" , x_url )
	s.Set( "active_player_now_playing_uri" , uri )
	return c.JSON( fiber.Map{
		"url": "/vlc/url/:url" ,
		"param_url": x_url ,
		"result": true ,
	})
}

// Custom Playlist Stuff
func ( s *Server ) VLCPlaylistAddURL( c *fiber.Ctx ) ( error ) {
	log.Debug( "VLCPlaylistAddURL()" )
	playlist_name := c.Params( "name" )
	sent_url := c.Params( "*" )
	// key := fmt.Sprintf( "LIBRARY.VLC.PLAYLISTS.%s" , playlist_name )
	// circular_set.Add( s.DB , key , video_id )
	return c.JSON( fiber.Map{
		"url": "/vlc/playlist/:name/add/url/*" ,
		"playlist_name": playlist_name ,
		"sent_url": sent_url ,
		"result": true ,
	})
}