package server

import (
	"fmt"
	// "strings"
	fiber "github.com/gofiber/fiber/v2"
	utils "github.com/0187773933/FireC2Server/v1/utils"
)

// weak attempt at sanitizing form input to build a "username"
func SanitizeUsername( first_name string , last_name string ) ( username string ) {
	if first_name == "" { first_name = "Not Provided" }
	if last_name == "" { last_name = "Not Provided" }
	sanitized_first_name := utils.SanitizeInputString( first_name )
	sanitized_last_name := utils.SanitizeInputString( last_name )
	username = fmt.Sprintf( "%s-%s" , sanitized_first_name , sanitized_last_name )
	return
}

func serve_failed_attempt( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>no</h1>" )
	// return context.SendFile( "./v1/server/html/admin_login.html" )
}

func ServeLoginPage( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/admin_login.html" )
}

// func ServeAuthenticatedPage( context *fiber.Ctx ) ( error ) {
// 	if validate_admin_cookie( context ) == false { return serve_failed_attempt( context ) }
// 	x_path := context.Route().Path
// 	url_key := strings.Split( x_path , "/admin" )
// 	if len( url_key ) < 2 { return context.SendFile( "./v1/server/html/admin_login.html" ) }
// 	// fmt.Println( "Sending -->" , url_key[ 1 ] , x_path )
// 	return context.SendFile( ui_html_pages[ url_key[ 1 ] ] )
// }

func ( s *Server ) LogRequest( context *fiber.Ctx ) ( error ) {
	time_string := utils.GetFormattedTimeString()
	ip_address := context.Get( "x-forwarded-for" )
	if ip_address == "" { ip_address = context.IP() }
	log_message := fmt.Sprintf( "%s === %s === %s === %s" , time_string , ip_address , context.Method() , context.Path() )
	fmt.Println( log_message )
	return context.Next()
}

func ( s *Server ) Log( message string ) {
	time_string := utils.GetFormattedTimeString()
	log_message := fmt.Sprintf( "%s === %s" , time_string , message )
	fmt.Println( log_message )
}

func ( s *Server ) Printf( format_string string , args ...interface{} ) {
	time_string := utils.GetFormattedTimeString()
	sent_format := fmt.Sprintf( format_string , args... )
	fmt.Printf( "%s === %s" , time_string , sent_format )
}