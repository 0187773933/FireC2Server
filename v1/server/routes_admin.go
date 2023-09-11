package server

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
)

var media_player_command_routes = []string{ "play" , "pause" , "stop" , "next" , "previous" , "teardown" , "setup" , "update" }
func ( s *Server ) SetupMediaPlayerRoutes( group fiber.Router , player_name string ) {
	for _ , command := range media_player_command_routes {
		x_command := command
		route_suffix := fmt.Sprintf( "/%s" , x_command )
		route_name := fmt.Sprintf( "/%s/%s" , player_name , x_command )
		group.Get( route_suffix , func( c *fiber.Ctx ) error {
			result := s.State.ExecuteCommand( player_name , x_command )
			return c.JSON( fiber.Map{
				"url": route_name ,
				"result": result ,
			})
		})
	}
}

func ( s *Server ) SetupAdminRoutes() {

	twitch := s.FiberApp.Group( "/twitch" )
	twitch.Use( validate_admin_mw )
	s.SetupMediaPlayerRoutes( twitch , "twitch" )

	youtube := s.FiberApp.Group( "/youtube" )
	youtube.Use( validate_admin_mw )
	s.SetupMediaPlayerRoutes( youtube , "youtube" )
}