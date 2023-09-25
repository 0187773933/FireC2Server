package utils

import (
	"os"
	"fmt"
	"sync"
	"time"
	"strings"
	"strconv"
	"unicode"
	"bufio"
	"net"
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

var location , _ = tz.LoadLocation( "America/New_York" )
var month_map = map[string]time.Month{
	"JAN": time.January , "FEB": time.February , "MAR": time.March ,
	"APR": time.April , "MAY": time.May , "JUN": time.June ,
	"JUL": time.July , "AUG": time.August , "SEP": time.September ,
	"OCT": time.October , "NOV": time.November , "DEC": time.December ,
}

func GetFormattedTimeString() ( result string ) {
	time_object := time.Now().In( location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}
func GetFormattedTimeStringOBJ() ( result_string string , result_time time.Time ) {
	result_time = time.Now().In( location )
	month_name := strings.ToUpper( result_time.Format( "Jan" ) )
	milliseconds := result_time.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , result_time.Day() , month_name , result_time.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , result_time.Hour() , result_time.Minute() , result_time.Second() , milliseconds )
	result_string = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

func FormatTime( input_time *time.Time ) ( result string ) {
	time_object := input_time.In( location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

// Year , Month , Day
func FormatDBLogPrefix( input_time *time.Time ) ( result string ) {
	time_object := input_time.In( location )
	result = fmt.Sprintf( "%d.%02d.%02d" , time_object.Year() , time_object.Month() , time_object.Day()  )
	return
}

func ParseFormattedTimeString( time_str string ) ( result time.Time ) {
	parts := strings.Split( time_str , " === " )
	date_part := parts[ 0 ]
	day , _ := strconv.Atoi( date_part[ 0 : 2 ] )
	month_abbr := date_part[ 2 : 5 ]
	month := month_map[ month_abbr ]
	year , _ := strconv.Atoi( date_part[ 5 : ] )
	time_part := parts[ 1 ]
	hour , _ := strconv.Atoi( time_part[ 0 : 2 ] )
	minute , _ := strconv.Atoi( time_part[ 3 : 5 ] )
	second , _ := strconv.Atoi( time_part[ 6 : 8 ] )
	millisecond , _ := strconv.Atoi( time_part[ 9 : ] )
	result = time.Date( year , month , day , hour , minute , second , ( millisecond * 1e6 ) , location )
	return
}

func StringToInt( input string ) ( result int ) {
	result , _ = strconv.Atoi( input )
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
	config_file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( config_file , &result )
	if error != nil { panic( error ) }

	var spotify_library types.SpotifyLibrary
	spotify_library_file , _ := ioutil.ReadFile( "./library/spotify.yaml" )
	error = yaml.Unmarshal( spotify_library_file , &spotify_library )
	if error != nil { panic( error ) }
	result.Library.Spotify = spotify_library

	var twitch_library types.TwitchLibrary
	twitch_library_file , _ := ioutil.ReadFile( "./library/twitch.yaml" )
	error = yaml.Unmarshal( twitch_library_file , &twitch_library )
	if error != nil { panic( error ) }
	result.Library.Twitch = twitch_library

	var disney_library types.DisneyLibrary
	disney_library_file , _ := ioutil.ReadFile( "./library/disney.yaml" )
	error = yaml.Unmarshal( disney_library_file , &disney_library )
	if error != nil { panic( error ) }
	result.Library.Disney = disney_library

	var youtube_library types.YouTubeLibrary
	youtube_library_file , _ := ioutil.ReadFile( "./library/disney.yaml" )
	error = yaml.Unmarshal( youtube_library_file , &youtube_library )
	if error != nil { panic( error ) }
	result.Library.YouTube = youtube_library

	var vlc_library types.VLCLibrary
	vlc_library_file , _ := ioutil.ReadFile( "./library/vlc.yaml" )
	error = yaml.Unmarshal( vlc_library_file , &vlc_library )
	if error != nil { panic( error ) }
	result.Library.VLC = vlc_library

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

func WakeOnLan( mac_address string ) {
	mac_bytes , _ := net.ParseMAC( mac_address )
	magic_packet := []byte{}
	for i := 0; i < 6; i++ {
		magic_packet = append( magic_packet , 0xFF )
	}
	for i := 0; i < 16; i++ {
		magic_packet = append( magic_packet , mac_bytes... )
	}
	addr := &net.UDPAddr{
		IP: net.IPv4bcast ,
		Port: 9 ,
	}
	conn , _ := net.DialUDP( "udp" , nil , addr )
	defer conn.Close()
	conn.Write( magic_packet )
}