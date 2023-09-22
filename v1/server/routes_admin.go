package server

import (
	"fmt"
	"strings"
	fiber "github.com/gofiber/fiber/v2"
)

var media_player_command_routes = []string{ "play" , "pause" , "stop" , "next" , "previous" , "teardown" , "setup" , "update" }
func ( s *Server ) SetupMediaPlayerRoutes( group fiber.Router , player_name string ) {
	for _ , command := range media_player_command_routes {
		x_command := command
		route_suffix := fmt.Sprintf( "/%s" , x_command )
		route_name := fmt.Sprintf( "/%s/%s" , player_name , x_command )
		group.Get( route_suffix , func( c *fiber.Ctx ) error {
			result := s.MediaPlayer.Run( player_name , x_command )
			return c.JSON( fiber.Map{
				"url": route_name ,
				"result": result ,
			})
		})
	}
}


func ( s *Server ) SetupAdminRoutes() {

	var generic_adb_command_routes = []string{ "play" , "pause" , "stop" , "next" , "previous" }
	for _ , command := range generic_adb_command_routes {
		x_command := command
		route := fmt.Sprintf( "/%s" , x_command )
		s.FiberApp.Get( route , func( c *fiber.Ctx ) error {
			s.MediaPlayer.ADB.PressKeyName( fmt.Sprintf( "KEYCODE_MEDIA_%s" , strings.ToUpper( x_command ) ) )
			return c.JSON( fiber.Map{
				"url": route ,
				"result": true ,
			})
		})
	}

	twitch := s.FiberApp.Group( "/twitch" )
	twitch.Use( validate_admin_mw )
	s.SetupMediaPlayerRoutes( twitch , "twitch" )

	youtube := s.FiberApp.Group( "/youtube" )
	youtube.Use( validate_admin_mw )
	s.SetupMediaPlayerRoutes( youtube , "youtube" )

	spotify := s.FiberApp.Group( "/spotify" )
	spotify.Use( validate_admin_mw )
	s.SetupMediaPlayerRoutes( spotify , "spotify" )
	spotify.Get( "/playlist-shuffle/:playlist_id" , func( c *fiber.Ctx ) error {
		playlist_id := c.Params( "playlist_id" )
		s.MediaPlayer.MediaPlayers[ "spotify" ].PlayPlaylistWithShuffle( playlist_id )
		return c.JSON( fiber.Map{
			"url": "/playlist-shuffle/:playlist_id" ,
			"playlist_id": playlist_id ,
			"result": true ,
		})
	})
	spotify.Get( "/playlist/:playlist_id" , func( c *fiber.Ctx ) error {
		playlist_id := c.Params( "playlist_id" )
		s.MediaPlayer.MediaPlayers[ "spotify" ].PlayPlaylist( playlist_id )
		return c.JSON( fiber.Map{
			"url": "/playlist/:playlist_id" ,
			"playlist_id": playlist_id ,
			"result": true ,
		})
	})

}