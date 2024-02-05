package server

import (
	fmt "fmt"
	time "time"
	context "context"
	sort "sort"
	rand "math/rand"
	// "reflect"
	strconv "strconv"
	// json "encoding/json"
	// url "net/url"
	// "math"
	// "image/color"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

const TWITCH_ACTIVITY = "tv.twitch.android.viewer/tv.twitch.starshot64.app.StarshotActivity"
const TWITCH_APP_NAME = "tv.twitch.android.viewer"
const R_KEY_STATE_TWITCH_FOLLOWING_LIVE = "STATE.TWITCH.FOLLOWING.LIVE"

func ( s *Server ) TwitchReopenApp() {
	log.Debug( "TwitchReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( TWITCH_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( TWITCH_APP_NAME )
	log.Debug( "Done" )
}

func ( s *Server ) TwitchContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "TwitchContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "twitch" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	if s.Status.ADB.Activity == ACTIVITY_PROFILE_PICKER {
		log.Debug( fmt.Sprintf( "Choosing Profile Index === %d" , s.Config.FireCubeUserProfileIndex ) )
		time.Sleep( 1000 * time.Millisecond )
		s.SelectFireCubeProfile()
		time.Sleep( 1000 * time.Millisecond )
	} else if s.Status.ADB.Activity != TWITCH_ACTIVITY {
		log.Debug( "twitch was NOT already open" )
		s.TwitchReopenApp()
		time.Sleep( 500 * time.Millisecond )
	} else {
		log.Debug( "twitch was already open" )
		// s.TwitchLiveUpdate()
		// we don't want a destructive update , we just need a cleansing of stale offline usernames

	}
}

func ( s *Server ) TwitchLiveNext( c *fiber.Ctx ) ( error ) {

	log.Debug( "TwitchLiveNext()" )
	s.TwitchContinuousOpen()

	next_stream := circular_set.Next( s.DB , R_KEY_STATE_TWITCH_FOLLOWING_LIVE )
	log.Debug( "Next === " , next_stream )
	if next_stream == "" {
		log.Debug( "Empty , Refreshing" )
		s.TwitchLiveUpdate()
		next_stream = circular_set.Current( s.DB , R_KEY_STATE_TWITCH_FOLLOWING_LIVE )
		if next_stream == "" {
			log.Debug( "nobody is live ...." )
			return c.JSON( fiber.Map{
				"url": "/twitch/live/next" ,
				"stream": "nobody is live ...." ,
				"result": false ,
			})
		}
	}
	// force refresh on last
	first_in_set_z , _ := s.DB.ZRangeWithScores( context.Background() , R_KEY_STATE_TWITCH_FOLLOWING_LIVE , 0 , 0 ).Result()
	first_in_set := first_in_set_z[0].Member.(string)
	if first_in_set == next_stream {
		log.Debug( "recycled list" )
		s.TwitchLiveUpdate()
	}
	log.Debug( fmt.Sprintf( "TwitchLiveNext( %s )" , next_stream ) )
	uri := fmt.Sprintf( "twitch://stream/%s" , next_stream )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	s.Set( "STATE.TWITCH.LIVE.NOW_PLAYING" , next_stream )
	s.Set( "active_player_now_playing_id" , next_stream )
	s.Set( "active_player_now_playing_text" , "" )
	// Force Highest Quality
	// The Problem is buffering could delay when this menu appears
	// you have to wait on button screen
	// look for white heart
	// time.Sleep( 3 * time.Second )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	// s.ADB.PressKeyName( "KEYCODE_ENTER" )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
	// s.ADB.PressKeyName( "KEYCODE_ENTER" )
	return c.JSON( fiber.Map{
		"url": "/twitch/live/next" ,
		"stream": next_stream ,
		"result": true ,
	})
}

func ( s *Server ) TwitchLivePrevious( c *fiber.Ctx ) ( error ) {

	log.Debug( "TwitchLivePrevious()" )
	s.TwitchContinuousOpen()

	next_stream := circular_set.Previous( s.DB , R_KEY_STATE_TWITCH_FOLLOWING_LIVE )
	log.Debug( "Next === " , next_stream )
	if next_stream == "" {
		log.Debug( "Empty , Refreshing" )
		s.TwitchLiveUpdate()
		next_stream = circular_set.Previous( s.DB , R_KEY_STATE_TWITCH_FOLLOWING_LIVE )
		if next_stream == "" {
			log.Debug( "nobody is live ...." )
			return c.JSON( fiber.Map{
				"url": "/twitch/live/next" ,
				"stream": "nobody is live ...." ,
				"result": false ,
			})
		}
	}
	// force refresh on last
	last_in_set_z , _ := s.DB.ZRangeWithScores( context.Background() , R_KEY_STATE_TWITCH_FOLLOWING_LIVE , -1 , -1 ).Result()
	last_in_set := last_in_set_z[0].Member.(string)
	if last_in_set == next_stream {
		log.Debug( "recycled list" )
		s.TwitchLiveUpdate()
	}
	log.Debug( fmt.Sprintf( "TwitchLivePrevious( %s )" , next_stream ) )
	uri := fmt.Sprintf( "twitch://stream/%s" , next_stream )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	s.Set( "STATE.TWITCH.LIVE.NOW_PLAYING" , next_stream )
	s.Set( "active_player_now_playing_id" , next_stream )
	s.Set( "active_player_now_playing_text" , "" )
	// Force Highest Quality
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	// s.ADB.PressKeyName( "KEYCODE_ENTER" )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
	// s.ADB.PressKeyName( "KEYCODE_ENTER" )
	// Right , Down , Right , Enter , Down , Enter
	return c.JSON( fiber.Map{
		"url": "/twitch/live/next" ,
		"stream": next_stream ,
		"result": true ,
	})
}


// Auth Option - 1 = Client Credentials
// curl -X POST 'https://id.twitch.tv/oauth2/token' \
// -H 'Content-Type: application/x-www-form-urlencoded' \
// -d 'client_id=asdf&client_secret=asdf&grant_type=client_credentials'

// Auth Option - 2 - Step - 1 = Create a Toekn that Can Automatically be Refreshed
// 2 steps , first step is creating a "code"
// https://id.twitch.tv/oauth2/authorize
//     ?response_type=code
//     &client_id=asdf
//     &redirect_uri=http://localhost:9371
//     &scope=user:read:follows

// Auth Option - 2 - Step - 2 = Second step is exchaning the "code" for a "token"
// This HAS to be a post request
// curl -X POST 'https://id.twitch.tv/oauth2/token' \
// -H 'Content-Type: application/x-www-form-urlencoded' \
// -d 'client_id=asdf' \
// -d 'client_secret=asdf' \
// -d 'code=asdf' \
// -d 'grant_type=authorization_code' \
// -d 'redirect_uri=http://localhost:9371'
	// returns :
	// {"access_token":"asdf","expires_in":13832,"refresh_token":"asdf","scope":["user:read:follows"],"token_type":"bearer"}

// Auth Option - 2 - Repeated Step - Refresh = Now we can just call refresh every time we need to make a batch of api calls
// curl -X POST https://id.twitch.tv/oauth2/token \
// -H 'Content-Type: application/x-www-form-urlencoded' \
// -d 'grant_type=refresh_token' \
// -d 'refresh_token=asdf' \
// -d 'client_id=asdf' \
// -d 'client_secret=asdf'
	// So you HAVE to store this refresh_token
	// update it every time

// Auth Option - 3 = Create Specific Auth Tokens
// These are temporary , and can't
// https://dev.twitch.tv/docs/authentication/scopes/
// https://id.twitch.tv/oauth2/authorize
//     ?response_type=token
//     &client_id=asdf
//     &redirect_uri=http://localhost:9371
//     &scope=user:read:follows

// Token Validation
// curl -X GET 'https://id.twitch.tv/oauth2/validate' \
// -H 'Authorization: OAuth asdf'

// Get User-ID of Username
// curl -X GET 'https://api.twitch.tv/helix/users?login=asdf' \
// -H 'Client-Id: asdf' \
// -H 'Authorization: Bearer Client-Credentials-Auth-Token'

// Get the channels you follow which are currently live
// curl -X GET 'https://api.twitch.tv/helix/streams/followed?user_id=asdf' \
// -H 'Authorization: Bearer asdf' \
// -H 'Client-Id: asdf'

func ( s *Server ) TwitchRefreshAuthToken() ( access_token string ) {
	time_now := time.Now()
	access_token = s.Get( "STATE.TWITCH.ACCESS_TOKEN" )
	refresh_token := s.Get( "STATE.TWITCH.REFRESH_TOKEN" )
	if access_token == "" || refresh_token == "" {
		access_token = s.Config.TwitchAccessToken
		refresh_token = s.Config.TwitchRefreshToken
	}
	refresh_token_expires_at := s.Get( "STATE.TWITCH.REFRESH_TOKEN.EXPIRES_AT" )
	if refresh_token_expires_at != "" {
		refresh_token_expires_at_int64 , _ := strconv.ParseInt( refresh_token_expires_at , 10 , 64 )
		refresh_token_expire_time := time.Unix( refresh_token_expires_at_int64 , 0 )
		remaining_time := refresh_token_expire_time.Sub( time_now )
		buffer_window := ( 30 * time.Second )
		if remaining_time >= buffer_window {
			// fmt.Println( "the token is still valid , and doesn't expire in the next 30 seconds , reusing" )
			// fmt.Println( ( remaining_time - 30 ) , "=== seconds remaining until refresh needed" )
			return
		}
	}
	// fmt.Println( "access token expired , refreshing" )
	headers := map[string]string{}
	data := map[string]string{
		"grant_type": "refresh_token" ,
		"refresh_token": refresh_token ,
		"client_id": s.Config.TwitchClientID ,
		"client_secret": s.Config.TwitchClientSecret ,
	}
	new := utils.PostJSON( "https://id.twitch.tv/oauth2/token" , headers , data );

	// adding - 13JAN2024
	// new_json_print , _ := json.MarshalIndent( new , "" , "    " )
	// fmt.Println( string( new_json_print ) )

	expires_in_seconds := new.(map[string]interface{})[ "expires_in" ].(float64)
	time_expires := time_now.Add( time.Second * time.Duration( expires_in_seconds ) )
	time_expires_unix := time_expires.Unix()
	// fmt.Println( "expires in ===" , expires_in_seconds )
	// fmt.Println( "expires @ ===" , time_expires , time_expires_unix )
	s.Set( "STATE.TWITCH.REFRESH_TOKEN.EXPIRES_AT" , time_expires_unix )

	new_access_token := new.(map[string]interface{})[ "access_token" ]
	new_refresh_token := new.(map[string]interface{})[ "refresh_token" ]
	if new_access_token != nil && new_refresh_token != nil {
		access_token = new_access_token.(string)
		refresh_token = new_refresh_token.(string)
	}
	s.Set( "STATE.TWITCH.ACCESS_TOKEN" , access_token )
	s.Set( "STATE.TWITCH.REFRESH_TOKEN" , refresh_token )
	return
}

func ( s *Server ) TwitchGetLiveFollowers() ( result []string ) {
	access_token := s.TwitchRefreshAuthToken()
	headers := map[string]string{
		"Authorization": fmt.Sprintf( "Bearer %s" , access_token ) ,
		"Client-Id": s.Config.TwitchClientID ,
	}
	params := map[string]string{
		"user_id": s.Config.TwitchUserID ,
	}
	live_followers_json := utils.GetJSON( "https://api.twitch.tv/helix/streams/followed" , headers , params )
	live_followers := live_followers_json.(map[string]interface{})[ "data" ].( []interface{} )
	for _ , user := range live_followers {
		user_name := user.(map[string]interface{})["user_login"].(string)
		result = append( result , user_name )
	}
	return
}

// Update DB With List of Currated Live Followers
func ( s *Server ) TwitchLiveUpdate() ( result []string ) {
	// 1.) Get Currated List
	var context = context.Background()
	currated_followers , _ := s.DB.ZRangeByScore(
		context ,
		"LIBRARY.TWITCH.FOLLOWING.CURRATED" ,
		&redis.ZRangeBy{
			Min: "-inf" ,
			Max: "+inf" ,
		} ,
	).Result()

	// 2.) Get Live Followers
	access_token := s.TwitchRefreshAuthToken()
	headers := map[string]string{
		"Authorization": fmt.Sprintf( "Bearer %s" , access_token ) ,
		"Client-Id": s.Config.TwitchClientID ,
	}
	params := map[string]string{
		"user_id": s.Config.TwitchUserID ,
	}
	live_followers_json := utils.GetJSON( "https://api.twitch.tv/helix/streams/followed" , headers , params )
	live_followers := live_followers_json.(map[string]interface{})[ "data" ].( []interface{} )

	// 3.) Get Live Currated List
	live_index_map := make(map[string]int)
	for i , user := range live_followers {
		user_name := user.(map[string]interface{})["user_login"].(string)
		live_index_map[user_name] = i;

	}
	// lookup table
	for _ , user := range currated_followers {
		if _ , exists := live_index_map[ user ]; exists {
			result = append( result , user )
			// s.DB.RPush( context , "STATE.TWITCH.FOLLOWING.LIVE" , user )
		}
	}
	currated_map := make(map[string]int)
	for i, user := range currated_followers {
		currated_map[user] = i
	}
	sort.Slice(result, func(i, j int) bool {
		return currated_map[result[i]] < currated_map[result[j]]
	})
	s.DB.Del( context , R_KEY_STATE_TWITCH_FOLLOWING_LIVE )
	for _ , user := range result {
		circular_set.Add( s.DB , R_KEY_STATE_TWITCH_FOLLOWING_LIVE , user )
	}
	log.Debug( result )
	return
}

func ( s *Server ) TwitchFilterCurratedFollers( input_list []string ) ( result []string ) {
	input_list_map := make( map[ string ]int )
	for i , user := range input_list {
		input_list_map[ user ] = i;
	}
	var context = context.Background()
	currated_followers , _ := s.DB.ZRangeByScore(
		context ,
		"LIBRARY.TWITCH.FOLLOWING.CURRATED" ,
		&redis.ZRangeBy{
			Min: "-inf" ,
			Max: "+inf" ,
		} ,
	).Result()
	for _ , user := range currated_followers {
		if _ , exists := input_list_map[ user ]; exists {
			result = append( result , user )
		}
	}
	currated_followers_map := make( map[ string ]int )
	for i , user := range currated_followers {
		currated_followers_map[ user ] = i
	}
	sort.Slice( result , func( i , j int ) bool {
		return currated_followers_map[ result[ i ] ] < currated_followers_map[ result[ j ] ]
	})
	return
}

func ( s *Server ) TwitchLiveRefresh() ( result []string ) {

	// 1.A) Get Live Followers
	live_followers := s.TwitchGetLiveFollowers()
	fmt.Println( "live followers ===" , live_followers )
	// 1.B.) Filter Live Followers to Only those in Currated Library List
	live_followers_currated := s.TwitchFilterCurratedFollers( live_followers )
	live_followers_currated_map := make( map[ string ]int )
	for i , user := range live_followers_currated {
		live_followers_currated_map[ user ] = i
	}
	fmt.Println( "live followers currated ===" , live_followers_currated )

	// 2.) Get Cached Followers
	var context = context.Background()
	cached_followers , _ := s.DB.ZRangeWithScores( context , R_KEY_STATE_TWITCH_FOLLOWING_LIVE , 0 , -1 ).Result()
	fmt.Println( "cached ===" , cached_followers )
	cached_index := circular_set.Index( s.DB , ( R_KEY_STATE_TWITCH_FOLLOWING_LIVE + ".INDEX" ) )
	fmt.Println( "cached index ===" , cached_index )
	cached_current := circular_set.Current( s.DB , R_KEY_STATE_TWITCH_FOLLOWING_LIVE )
	fmt.Println( "cached current ===" , cached_current )

	// 5.) The actuall point of this function
	// 5.1) if the cache is empty , use the new list ( result ) , return
	// if len( cached_followers ) == 0 { return }
	// 5.2) if there are people in the cache that ARE NOT in the new list , we need to eject them
	// https://pkg.go.dev/github.com/redis/go-redis/v9#Z
	cached_list_with_offline_removed := []string{}
	cached_list_with_offline_removed_map := make( map[ string ]int )
	for i := range cached_followers {
		i_username := cached_followers[ i ].Member.( string )
		if _ , exists := live_followers_currated_map[ i_username ]; exists {
			cached_list_with_offline_removed = append( cached_list_with_offline_removed , i_username )
			cached_list_with_offline_removed_map[ i_username ] = i
		}
	}
	fmt.Println( "cached list with offline removed ===" , cached_list_with_offline_removed )

	// 5.3) find new online users
	// this is where it gets debatable on what you want to do
	// we could resort , or just add people who are newly online to the end , etc
	new_online_users := []string{}
	new_online_users_map := make( map[ string ]int )
	for _ , user := range live_followers_currated {
		if i , exists := live_followers_currated_map[ user ]; !exists {
			new_online_users = append( new_online_users , user )
			new_online_users_map[ user ] = i
		}
	}
	// i think its already sorted , because its in the same order as 'live_followers_currated'
	// sort.Slice( new_online_users , func( i , j int ) bool {
	// 	return currated_followers_map[ new_online_users[ i ] ] < currated_followers_map[ new_online_users[ j ] ]
	// })
	fmt.Println( "new online currated users since last cache ===" , new_online_users )

	cached_list_with_offline_removed_and_new_online_added := append( cached_list_with_offline_removed , new_online_users... )
	fmt.Println( "cached list with offline removed and new online added ===" , cached_list_with_offline_removed_and_new_online_added )

	// 5.4) circular list package doesn't have a remove
	result = cached_list_with_offline_removed_and_new_online_added
	s.DB.Del( context , R_KEY_STATE_TWITCH_FOLLOWING_LIVE )
	for _ , user := range result {
		circular_set.Add( s.DB , R_KEY_STATE_TWITCH_FOLLOWING_LIVE , user )
	}
	s.Set( ( R_KEY_STATE_TWITCH_FOLLOWING_LIVE + ".INDEX" ) , cached_index )

	return
}

func ( s *Server ) GetTwitchLiveUpdate( c *fiber.Ctx ) ( error ) {
	live := s.TwitchLiveUpdate()
	return c.JSON( fiber.Map{
		"url": "/twitch/live/update" ,
		"currated": live ,
		"result": true ,
	})
}

func ( s *Server ) GetTwitchLiveRefresh( c *fiber.Ctx ) ( error ) {
	live := s.TwitchLiveRefresh()
	return c.JSON( fiber.Map{
		"url": "/twitch/live/refresh" ,
		"currated": live ,
		"result": true ,
	})
}

// literally trolling ,
// we have to get pixel values to know where we are being bullied in the quality selection menu
func ( s *Server ) TwitchLiveSetQualityMax( c *fiber.Ctx ) ( error ) {
	// Right , Down , Right , Enter , Down , Enter
	rand.Seed( time.Now().UnixNano() )
	min_sleep := 200
	max_sleep := 700
	s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_DPAD_LEFT" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_ENTER" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_UP" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_UP" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_UP" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_UP" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_UP" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_DPAD_DOWN" )
	time.Sleep( time.Duration( rand.Intn( max_sleep + 1 ) + min_sleep ) * time.Millisecond )
	s.ADB.PressKeyName( "KEYCODE_ENTER" )
	return c.JSON( fiber.Map{
		"url": "/twitch/live/set/quality/max" ,
		"result": true ,
	})
}

func ( s *Server ) TwitchLiveUser( c *fiber.Ctx ) ( error ) {
	username := c.Params( "username" )
	log.Debug( fmt.Sprintf( "TwitchLiveUser( %s )" , username ) )
	s.TwitchContinuousOpen()
	uri := fmt.Sprintf( "twitch://stream/%s" , username )
	s.ADB.OpenURI( uri )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	s.Set( "STATE.TWITCH.LIVE.NOW_PLAYING" , username )
	s.Set( "active_player_now_playing_id" , username )
	s.Set( "active_player_now_playing_text" , "" )
	return c.JSON( fiber.Map{
		"url": "/twitch/view/:username" ,
		"stream": username ,
		"result": true ,
	})
}

// func ( s *Server ) TwitchLiveCylce( c *fiber.Ctx ) ( error ) {
// 	username := c.Params( "username" )
// 	log.Debug( fmt.Sprintf( "TwitchLiveUser( %s )" , username ) )
// 	s.TwitchContinuousOpen()
// 	uri := fmt.Sprintf( "twitch://stream/%s" , username )
// 	s.ADB.OpenURI( uri )
// 	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
// 	s.Set( "STATE.TWITCH.LIVE.NOW_PLAYING" , username )
// 	s.Set( "active_player_now_playing_id" , username )
// 	s.Set( "active_player_now_playing_text" , "" )
// 	return c.JSON( fiber.Map{
// 		"url": "/twitch/view/:username" ,
// 		"stream": username ,
// 		"result": true ,
// 	})
// }