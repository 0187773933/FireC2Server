package server

import (
	"fmt"
	"time"
	"strings"
	fiber "github.com/gofiber/fiber/v2"
	bcrypt "golang.org/x/crypto/bcrypt"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

func validate_login_credentials( context *fiber.Ctx ) ( result bool ) {
	result = false
	uploaded_username := context.FormValue( "username" )
	if uploaded_username == "" { fmt.Println( "username empty" ); return }
	if uploaded_username != GlobalServer.Config.AdminUsername { fmt.Println( "username not correct" ); return }
	uploaded_password := context.FormValue( "password" )
	if uploaded_password == "" { fmt.Println( "password empty" ); return }
	fmt.Println( "uploaded_username ===" , uploaded_username )
	fmt.Println( "uploaded_password ===" , uploaded_password )
	password_matches := bcrypt.CompareHashAndPassword( []byte( uploaded_password ) , []byte( GlobalServer.Config.AdminPassword ) )
	if password_matches != nil { fmt.Println( "bcrypted password doesn't match" ); return }
	fmt.Println( "password matched" )
	result = true
	return
}

func HandleLogout( context *fiber.Ctx ) ( error ) {
	context.Cookie( &fiber.Cookie{
		Name: GlobalServer.Config.ServerCookieName ,
		Value: "" ,
		Expires: time.Now().Add( -time.Hour ) , // set the expiration to the past
		HTTPOnly: true ,
		Secure: true ,
	})
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>Logged Out</h1>" )
}

// POST http://localhost:5950/admin/login
func HandleLogin( context *fiber.Ctx ) ( error ) {
	valid_login := validate_login_credentials( context )
	if valid_login == false { return serve_failed_attempt( context ) }
	host := context.Hostname()
	domain := strings.Split( host , ":" )[ 0 ]
	context.Cookie(
		&fiber.Cookie{
			Name: GlobalServer.Config.ServerCookieName ,
			Value: encryption.SecretBoxEncrypt( GlobalServer.Config.EncryptionKey , GlobalServer.Config.ServerCookieAdminSecretMessage ) ,
			Secure: true ,
			Path: "/" ,
			// Domain: "blah.ngrok.io" , // probably should set this for webkit
			Domain: domain ,
			HTTPOnly: true ,
			SameSite: "Lax" ,
			Expires: time.Now().AddDate( 10 , 0 , 0 ) , // aka 10 years from now
		} ,
	)
	return context.Redirect( "/" )
}

func validate_admin_cookie( context *fiber.Ctx ) ( result bool ) {
	result = false
	admin_cookie := context.Cookies( GlobalServer.Config.ServerCookieName )
	if admin_cookie == "" { fmt.Println( "admin cookie was blank" ); return }
	admin_cookie_value := encryption.SecretBoxDecrypt( GlobalServer.Config.EncryptionKey , admin_cookie )
	if admin_cookie_value != GlobalServer.Config.ServerCookieAdminSecretMessage { fmt.Println( "admin cookie secret message was not equal" ); return }
	result = true
	return
}

func validate_admin( context *fiber.Ctx ) ( result bool ) {
	result = false
	admin_cookie := context.Cookies( GlobalServer.Config.ServerCookieName )
	if admin_cookie != "" {
		admin_cookie_value := encryption.SecretBoxDecrypt( GlobalServer.Config.EncryptionKey , admin_cookie )
		if admin_cookie_value == GlobalServer.Config.ServerCookieAdminSecretMessage {
			result = true
			return
		}
	}
	admin_api_key_header := context.Get( "key" )
	if admin_api_key_header != "" {
		if admin_api_key_header == GlobalServer.Config.ServerAPIKey {
			result = true
			return
		}
	}
	admin_api_key_query := context.Query( "k" )
	if admin_api_key_query != "" {
		if admin_api_key_query == GlobalServer.Config.ServerAPIKey {
			result = true
			return
		}
	}
	return
}

func validate_admin_mw( context *fiber.Ctx ) ( error ) {
	admin_cookie := context.Cookies( GlobalServer.Config.ServerCookieName )
	if admin_cookie != "" {
		admin_cookie_value := encryption.SecretBoxDecrypt( GlobalServer.Config.EncryptionKey , admin_cookie )
		if admin_cookie_value == GlobalServer.Config.ServerCookieAdminSecretMessage {
			return context.Next()
		}
	}
	admin_api_key_header := context.Get( "key" )
	if admin_api_key_header != "" {
		if admin_api_key_header == GlobalServer.Config.ServerAPIKey {
			return context.Next()
		}
	}
	admin_api_key_query := context.Query( "k" )
	if admin_api_key_query != "" {
		if admin_api_key_query == GlobalServer.Config.ServerAPIKey {
			return context.Next()
		}
	}
	return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
}

func validate_browser_mw(context *fiber.Ctx) error {
	browser_api_key := context.Get("key")
	if browser_api_key == "" {
		browser_api_key = context.Get("k")
	}
	if browser_api_key == "" {
		browser_api_key = context.Query("k")
	}
	if browser_api_key == GlobalServer.Config.BrowserAPIKey {
		return context.Next()
	}
	return context.Status(fiber.StatusUnauthorized).SendString("why")
}
