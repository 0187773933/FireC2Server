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
	"runtime"
	"math/rand"
	color "image/color"
	colorful "github.com/lucasb-eyer/go-colorful"
	filepath "path/filepath"
	"bytes"
	"net/http"
	url "net/url"
	"encoding/json"
	sha256 "crypto/sha256"
	// "strings"
	"io/ioutil"
	uuid "github.com/google/uuid"
	yaml "gopkg.in/yaml.v2"
	// hid "github.com/dh1tw/hid"
	try "github.com/manucorporat/try"
	types "github.com/0187773933/FireC2Server/v1/types"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

func SetupStackTraceReport() {
	if r := recover(); r != nil {
		stacktrace := make( []byte , 1024 )
		runtime.Stack( stacktrace , true )
		fmt.Printf( "%s\n" , stacktrace )
	}
}

func FingerPrint( config *types.ConfigFile ) {
	fmt.Println( GetLocalIPAddresses() )
}

func PrettyPrint( input interface{} ) {
	jd , _ := json.MarshalIndent( input , "" , "  " )
	fmt.Println( string( jd ) )
}

func Sha256( input string ) ( result string ) {
	hasher := sha256.New()
	hasher.Write( []byte( input ) )
	hash_bytes := hasher.Sum( nil )
	result = fmt.Sprintf( "%x" , hash_bytes )
	return
}

func IsUUID( u string ) ( result bool ) {
	_ , err := uuid.Parse( u )
	return err == nil
}

func IsURL( input string ) ( result bool , url *url.URL ) {
	result = false
	try.This( func() {
		if input == "" { return }
		parsed , err := url.Parse( input )
		if err != nil { return }
		if parsed.Scheme == "" { return }
		if parsed.Host == "" { return }
		url = parsed
		result = true
	}).Catch( func( e try.E ) {
		// fmt.Println( e )
		// fmt.Println( input )
	})
	return
}

func HexToRGBColor( hex_color string ) ( result color.RGBA ) {
	c , _ := colorful.Hex( hex_color )
	r , g , b := c.RGB255()
	result = color.RGBA{ R: r , G: g , B: b , A: 255 }
	return
}

func GetLocalIPAddresses() ( ip_addresses []string ) {
	host , _ := os.Hostname()
	addrs , _ := net.LookupIP( host )
	encountered := make( map[ string ]bool )
	for _ , addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip := ipv4.String()
			if !encountered[ ip ] {
				encountered[ ip ] = true
				ip_addresses = append( ip_addresses , ip )
			}
		}
	}
	return
}

// var location , _ = tz.LoadLocation( "America/New_York" )
var location , _ = time.LoadLocation( "America/New_York" )
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

func ShuffleKeys[ K comparable , V any ]( m map[ K ]V ) ( result []K ) {
	result = make( []K , 0 , len( m ) )
	for key := range m {
		result = append( result , key )
	}
	rand.Seed( time.Now().UnixNano() )
	rand.Shuffle( len( result ) , func( i , j int ) {
		result[ i ] , result[ j ] = result[ j ] , result[ i ]
	})
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

	library_base_path := filepath.Join( result.SaveFilesPath , "library" )
	fmt.Println( library_base_path )

	var spotify_library types.SpotifyLibrary
	spotify_library_file , _ := ioutil.ReadFile( filepath.Join( library_base_path , "spotify.yaml" ) )
	error = yaml.Unmarshal( spotify_library_file , &spotify_library )
	if error != nil { panic( error ) }
	result.Library.Spotify = spotify_library

	var twitch_library types.TwitchLibrary
	twitch_library_file , _ := ioutil.ReadFile( filepath.Join( library_base_path , "twitch.yaml" ) )
	error = yaml.Unmarshal( twitch_library_file , &twitch_library )
	if error != nil { panic( error ) }
	result.Library.Twitch = twitch_library

	var disney_library types.DisneyLibrary
	disney_library_file , _ := ioutil.ReadFile( filepath.Join( library_base_path , "disney.yaml" ) )
	error = yaml.Unmarshal( disney_library_file , &disney_library )
	if error != nil { panic( error ) }
	result.Library.Disney = disney_library

	var youtube_library types.YouTubeLibrary
	youtube_library_file , _ := ioutil.ReadFile( filepath.Join( library_base_path , "youtube.yaml" ) )
	error = yaml.Unmarshal( youtube_library_file , &youtube_library )
	if error != nil { panic( error ) }
	result.Library.YouTube = youtube_library

	var vlc_library types.VLCLibrary
	vlc_library_file , _ := ioutil.ReadFile( filepath.Join( library_base_path , "vlc.yaml" ) )
	error = yaml.Unmarshal( vlc_library_file , &vlc_library )
	if error != nil { panic( error ) }
	result.Library.VLC = vlc_library

	var hulu_library types.HuluLibrary
	hulu_library_file , hulu_library_file_read_error := ioutil.ReadFile( filepath.Join( library_base_path , "hulu.yaml" ) )
	if hulu_library_file_read_error != nil { panic( hulu_library_file_read_error ) }
	error = yaml.Unmarshal( hulu_library_file , &hulu_library )
	if error != nil { fmt.Println( error ); panic( error ) }
	result.Library.Hulu = hulu_library

	var netflix_library types.NetflixLibrary
	netflix_library_file , netflix_library_file_read_error := ioutil.ReadFile( filepath.Join( library_base_path , "netflix.yaml" ) )
	if netflix_library_file_read_error != nil { panic( netflix_library_file_read_error ) }
	error = yaml.Unmarshal( netflix_library_file , &netflix_library )
	if error != nil { fmt.Println( error ); panic( error ) }
	result.Library.Netflix = netflix_library

	return
}

func GenerateNewKeys() {
	fiber_cookie_key := fiber_cookie.GenerateKey()
	encryption_key := encryption.GenerateRandomString( 32 )
	server_api_key := encryption.GenerateRandomString( 16 )
	admin_username := encryption.GenerateRandomString( 16 )
	admin_password := encryption.GenerateRandomString( 16 )
	login_url := encryption.GenerateRandomString( 16 )
	browser_api_key := encryption.GenerateRandomString( 16 )
	fmt.Println( "Generated New Keys :" )
	fmt.Printf( "\tFiber Cookie Key === %s\n" , fiber_cookie_key )
	fmt.Printf( "\tEncryption Key === %s\n" , encryption_key )
	fmt.Printf( "\tServer API Key === %s\n" , server_api_key )
	fmt.Printf( "\tAdmin Username === %s\n" , admin_username )
	fmt.Printf( "\tAdmin Password === %s\n" , admin_password )
	fmt.Printf( "\tLogin URL === %s\n" , login_url )
	fmt.Printf( "\tBrowser API Key === %s\n" , browser_api_key )
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


func GetJSON( baseURL string , headers map[string]string , params map[string]string ) ( target interface{} ) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	u , err := url.Parse( baseURL )
	if err != nil { fmt.Println( err ); return }
	q := u.Query()
	for key, value := range params {
		q.Add( key , value )
	}
	u.RawQuery = q.Encode()
	req , err := http.NewRequest("GET", u.String(), nil)
	if err != nil { fmt.Println( err ); return }
	for key , value := range headers {
		req.Header.Set( key , value )
	}
	resp , err := client.Do( req )
	if err != nil { fmt.Println( err ); return }
	defer resp.Body.Close()
	body , err := ioutil.ReadAll( resp.Body )
	if err != nil { fmt.Println( err ); return }
	json.Unmarshal( body , &target )
	return
}

func PostJSON( url string , headers map[string]string , payload interface{} ) ( result interface{} ) {
	client := &http.Client{}
	payload_bytes , err := json.Marshal( payload )
	if err != nil { fmt.Println( err ); return }
	req , err := http.NewRequest( "POST" , url , bytes.NewBuffer( payload_bytes ) )
	if err != nil { fmt.Println( err ); return }
	req.Header.Set( "Content-Type" , "application/json" )
	for key, value := range headers {
		req.Header.Set( key , value )
	}
	resp , err := client.Do( req )
	if err != nil { fmt.Println( err ); return }
	defer resp.Body.Close()
	body , err := ioutil.ReadAll( resp.Body )
	if err != nil { fmt.Println( err ); return }
	json.Unmarshal( body , &result )
	return
}