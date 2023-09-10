package utils

import (
	"os"
	"fmt"
	"sync"
	"time"
	"strings"
	"unicode"
	"bufio"
	tz "4d63.com/tz"
	"encoding/json"
	// "strings"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
	// hid "github.com/dh1tw/hid"
	types "github.com/0187773933/FireC2Server/v1/types"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	encryption "github.com/0187773933/FireC2Server/v1/encryption"
)

func GetFormattedTimeString() ( result string ) {
	location , _ := tz.LoadLocation( "America/New_York" )
	time_object := time.Now().In( location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

func RemoveNonASCII( input string ) ( result string ) {
	for _ , i := range input {
		if i > unicode.MaxASCII { continue }
		result += string( i )
	}
	return
}

const SanitizedStringSizeLimit = 100
func SanitizeInputString( input string ) ( result string ) {
	trimmed := strings.TrimSpace( input )
    if len( trimmed ) > SanitizedStringSizeLimit { trimmed = strings.TrimSpace( trimmed[ 0 : SanitizedStringSizeLimit ] ) }
	result = RemoveNonASCII( trimmed )
	return
}

func WriteJSON( filePath string , data interface{} ) {
	file, _ := json.MarshalIndent( data , "" , " " )
	_ = ioutil.WriteFile( filePath , file , 0644 )
}

func ParseConfig( file_path string ) ( result types.ConfigFile ) {
	file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( file , &result )
	if error != nil { panic( error ) }
	return
}

func GenerateNewKeys() {
	fiber_cookie_key := fiber_cookie.GenerateKey()
	bolt_db_key := encryption.GenerateRandomString( 32 )
	server_api_key := encryption.GenerateRandomString( 16 )
	admin_username := encryption.GenerateRandomString( 16 )
	admin_password := encryption.GenerateRandomString( 16 )
	fmt.Println( "Generated New Keys :" )
	fmt.Printf( "\tFiber Cookie Key === %s\n" , fiber_cookie_key )
	fmt.Printf( "\tBolt DB Key === %s\n" , bolt_db_key )
	fmt.Printf( "\tServer API Key === %s\n" , server_api_key )
	fmt.Printf( "\tAdmin Username === %s\n" , admin_username )
	fmt.Printf( "\tAdmin Password === %s\n\n" , admin_password )
}

func WriteLoginURLPrefixWG( wg *sync.WaitGroup , server_login_url_prefix string ) {
	file, _ := os.OpenFile( "./v1/server/html/login.html" , os.O_RDWR , 0 )
	reader := bufio.NewReader( file )
	line_number := 1
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if line_number == 17 {
			x_line := fmt.Sprintf( "\t\t\t\t\t<form id=\"form-login\" action=\"/%s\" onSubmit=\"return on_submit();\" method=\"post\">\n" , server_login_url_prefix )
			lines = append(lines, x_line)
		} else {
			lines = append(lines, line)
		}

		if err != nil { break; }
		line_number++
	}
	file.Seek( 0 , 0 )
	file.Truncate( 0 )
	for _ , line := range lines {
		file.WriteString( line )
	}
	file.Close()
	wg.Done()
}
func WriteLoginURLPrefix( server_login_url_prefix string ) {
	file, _ := os.OpenFile( "./v1/server/html/login.html" , os.O_RDWR , 0 )
	defer file.Close()
	reader := bufio.NewReader( file )
	line_number := 1
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if line_number == 17 {
			x_line := fmt.Sprintf( "\t\t\t\t\t<form id=\"form-login\" action=\"/%s\" onSubmit=\"return on_submit();\" method=\"post\">\n" , server_login_url_prefix )
			lines = append(lines, x_line)
		} else {
			lines = append(lines, line)
		}

		if err != nil { break; }
		line_number++
	}
	file.Seek( 0 , 0 )
	file.Truncate( 0 )
	for _ , line := range lines {
		file.WriteString( line )
	}
}