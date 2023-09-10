package server

import (
	"fmt"
	"time"
	fiber "github.com/gofiber/fiber/v2"
	rate_limiter "github.com/gofiber/fiber/v2/middleware/limiter"
	// try "github.com/manucorporat/try"
)

var public_limiter = rate_limiter.New( rate_limiter.Config{
	Max: 30 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: func( c *fiber.Ctx ) error {
		ip_address := c.IP()
		log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
		GlobalServer.Log( log_message )
		c.Set( "Content-Type" , "text/html" )
		return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6);</script></html>" )
	} ,
})

func ( s *Server ) SetupPublicRoutes() {
	s.FiberApp.Static( "/cdn" , "./v1/server/cdn" )
	s.FiberApp.Get( "/" , public_limiter , RenderHomePage )
}

func RenderHomePage( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	// admin_logged_in := check_if_admin_cookie_exists( context )
	// if admin_logged_in == true {
	// 	// fmt.Println( "RenderHomePage() --> Admin" )
	// 	return context.SendFile( "./v1/server/html/admin.html" )
	// }
	return context.SendFile( "./v1/server/html/home.html" )
}