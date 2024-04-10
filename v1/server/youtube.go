package server

import (
	fmt "fmt"
	time "time"
	"math/rand"
	"strings"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// const YOUTUBE_ACTIVITY = "com.amazon.firetv.youtube/dev.cobalt.app.MainActivity"
// const YOUTUBE_APP_NAME = "com.amazon.firetv.youtube"

func ( s *Server ) YouTubeReopenApp() {
	log.Debug( "YouTubeReopenApp()" )
	s.ADB.StopAllPackages()
	// s.ADB.SetBrightness( 0 )
	s.ADB.ClosePackage( s.Config.ADB.APKS[ "youtube" ][ s.Config.ADB.DeviceType ].Package )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenPackage( s.Config.ADB.APKS[ "youtube" ][ s.Config.ADB.DeviceType ].Package )
	log.Debug( "Done" )
}

func ( s *Server ) YouTubeContinuousOpen() ( was_open bool ) {
	was_open = false
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "YouTubeContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "youtube" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	s.ADBWakeup()
	windows := s.ADB.GetWindowStack()
	for _ , window := range windows {
		activity_lower := strings.ToLower( window.Activity )
		if strings.Contains( activity_lower , "youtube" ) {
			log.Debug( "youtube was already open" )
			was_open = true
			return
		}
	}
	log.Debug( "youtube was NOT already open" )
	s.YouTubeReopenApp()
	time.Sleep( 1 * time.Second )
	return
}

func parse_youtube_sent_id( sent_id string ) ( uri string ) {
	is_url , _ := utils.IsURL( sent_id )
	if is_url {
		fmt.Println( "is url" )
		uri = sent_id
		return
	}
	return
}

func ( s *Server ) YouTubeOpenID( sent_id string ) {
	log.Debug( fmt.Sprintf( "YouTubeOpenID( %s )" , sent_id ) )
	s.YouTubeContinuousOpen()
	uri := parse_youtube_sent_id( sent_id )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
}

func ( s *Server ) YouTubeID( c *fiber.Ctx ) ( error ) {
	sent_id := c.Params( "*" )
	sent_query := c.Request().URI().QueryArgs().String()
	if sent_query != "" { sent_id += "?" + sent_query }
	log.Debug( fmt.Sprintf( "YouTubeID( %s )" , sent_id ) )
	s.YouTubeOpenID( sent_id )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_uri" , sent_id )
	return c.JSON( fiber.Map{
		"url": "/youtube/:id" ,
		"id": sent_id ,
		"result": true ,
	})
}

func ( s *Server ) YouTubePlaylistGetNextAvailableVideoID( playlist_key string ) ( video_id string ) {
	for {
		video_id = circular_set.Next( s.DB , playlist_key )
		available := s.YouTubeIsVideoIdAvailable( video_id )
		fmt.Println( video_id , available )
		if available == false {
			circular_set.Remove( s.DB , playlist_key , video_id )
			total_items := s.DB.ZCard( context.Background() , playlist_key ).Val()
			if total_items == 0 { return }
			time.Sleep( 3 * time.Second )
		} else {
			return
		}
	}
}

func ( s *Server ) YouTubePlaylistGetPreviousAvailableVideoID( playlist_key string ) ( video_id string ) {
	for {
		video_id = circular_set.Previous( s.DB , playlist_key )
		available := s.YouTubeIsVideoIdAvailable( video_id )
		fmt.Println( video_id , available )
		if available == false {
			circular_set.Remove( s.DB , playlist_key , video_id )
			total_items := s.DB.ZCard( context.Background() , playlist_key ).Val()
			if total_items == 0 { return }
			time.Sleep( 3 * time.Second )
		} else {
			return
		}
	}
}

// This is not a youtube playlist , this is our own "playlist"
// android youtube app still doesn't have intents for playlist loading
func ( s *Server ) YouTubePlaylistNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubePlaylistNext()" )
	playlist_name := c.Params( "name" )
	key := fmt.Sprintf( "LIBRARY.YOUTUBE.PLAYLISTS.%s" , playlist_name )
	key_index := key + ".INDEX"
	playlist_index := s.DB.Get( context.Background() , key_index ).Val()
	video_id := s.YouTubePlaylistGetNextAvailableVideoID( key )
	uri := fmt.Sprintf( "https://www.youtube.com/watch?v=%s" , video_id )
	s.YouTubeOpenID( uri )
	s.Set( "active_player_now_playing_id" , video_id )
	s.Set( "active_player_now_playing_text" , "" )
	return c.JSON( fiber.Map{
		"url": "/youtube/playlist/:name/next" ,
		"playlist_name": playlist_name ,
		"playlist_index": playlist_index ,
		"video_id": video_id ,
		"result": true ,
	})
}

func ( s *Server ) YouTubePlaylistPrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubePlaylistPrevious()" )
	playlist_name := c.Params( "name" )
	key := fmt.Sprintf( "LIBRARY.YOUTUBE.PLAYLISTS.%s" , playlist_name )
	key_index := key + ".INDEX"
	playlist_index := s.DB.Get( context.Background() , key_index ).Val()
	video_id := s.YouTubePlaylistGetPreviousAvailableVideoID( key )
	uri := fmt.Sprintf( "https://www.youtube.com/watch?v=%s" , video_id )
	s.YouTubeOpenID( uri )
	s.Set( "active_player_now_playing_id" , video_id )
	s.Set( "active_player_now_playing_text" , "" )
	return c.JSON( fiber.Map{
		"url": "/youtube/playlist/:name/previous" ,
		"playlist_name": playlist_name ,
		"playlist_index": playlist_index ,
		"video_id": video_id ,
		"result": true ,
	})
}

func ( s *Server ) YouTubePlaylistAddVideo( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubePlaylistAddVideo()" )
	playlist_name := c.Params( "name" )
	video_id := c.Params( "id" )
	available := s.YouTubeIsVideoIdAvailable( video_id )
	if available == false {
		return c.JSON( fiber.Map{
			"url": "/youtube/playlist/:name/add/:id" ,
			"playlist_name": playlist_name ,
			"video_id": video_id ,
			"video_available": available ,
			"result": false ,
		})
	}
	key := fmt.Sprintf( "LIBRARY.YOUTUBE.PLAYLISTS.%s" , playlist_name )
	circular_set.Add( s.DB , key , video_id )
	return c.JSON( fiber.Map{
		"url": "/youtube/playlist/:name/add/:id" ,
		"playlist_name": playlist_name ,
		"video_id": video_id ,
		"video_available": available ,
		"result": true ,
	})
}

func ( s *Server ) YouTubePlaylistAddPlaylist( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubePlaylistAddPlaylist()" )
	playlist_name := c.Params( "name" )
	playlist_id := c.Params( "id" )
	key := fmt.Sprintf( "LIBRARY.YOUTUBE.PLAYLISTS.%s" , playlist_name )
	playlist_videos := s.YouTubeGetPlaylistVideos( playlist_id )
	for _ , video := range playlist_videos {
		circular_set.Add( s.DB , key , video.Id )
	}
	return c.JSON( fiber.Map{
		"url": "/youtube/playlist/:name/add/playlist/:id" ,
		"playlist_name": playlist_name ,
		"playlist_id": playlist_id ,
		"playlist_videos": playlist_videos ,
		"result": true ,
	})
}

func ( s *Server ) YouTubePlaylistDeleteVideo( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubePlaylistDeleteVideo()" )
	playlist_name := c.Params( "name" )
	video_id := c.Params( "id" )
	key := fmt.Sprintf( "LIBRARY.YOUTUBE.PLAYLISTS.%s" , playlist_name )
	circular_set.Remove( s.DB , key , video_id )
	return c.JSON( fiber.Map{
		"url": "/youtube/playlist/:name/delete/:id" ,
		"playlist_name": playlist_name ,
		"video_id": video_id ,
		"result": true ,
	})
}

func ( s *Server ) YouTubePlaylistGet( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubePlaylistGet()" )
	playlist_name := c.Params( "name" )
	var ctx = context.Background()
	key := fmt.Sprintf( "LIBRARY.YOUTUBE.PLAYLISTS.%s" , playlist_name )
	key_index := key + ".INDEX"
	videos := s.DB.ZRange( ctx , key , 0 , -1 ).Val()
	current_index := s.DB.Get( ctx , key_index ).Val()
	return c.JSON( fiber.Map{
		"url": "/youtube/playlist/:name/get" ,
		"playlist_name": playlist_name ,
		"current_index": current_index ,
		"videos": videos ,
		"result": true ,
	})
}

func ( s *Server ) YouTubePlaylistGetIndex( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubePlaylistGetIndex()" )
	playlist_name := c.Params( "name" )
	key := fmt.Sprintf( "LIBRARY.YOUTUBE.PLAYLISTS.%s" , playlist_name )
	key_index := key + ".INDEX"
	index := s.DB.Get( context.Background() , key_index ).Val()
	return c.JSON( fiber.Map{
		"url": "/youtube/playlist/:name/index/get" ,
		"playlist_name": playlist_name ,
		"index": index ,
		"result": true ,
	})
}

func ( s *Server ) YouTubePlaylistSetIndex( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubePlaylistSetIndex()" )
	playlist_name := c.Params( "name" )
	set_index := c.Params( "index" )
	key := fmt.Sprintf( "LIBRARY.YOUTUBE.PLAYLISTS.%s" , playlist_name )
	key_index := key + ".INDEX"
	s.DB.Set( context.Background() , key_index , set_index , 0 )
	return c.JSON( fiber.Map{
		"url": "/youtube/playlist/:name/index/set/:index" ,
		"playlist_name": playlist_name ,
		"set_index": set_index ,
		"result": true ,
	})
}


var RESET_COUNTER = 0;
func ( s *Server ) YouTubeLiveNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubeLiveNext()" )
	video_id := circular_set.Next( s.DB , "STATE.YOUTUBE.LIVE.VIDEOS" )
	available := s.YouTubeIsVideoIdAvailable( video_id )
	// IF None Available 5 Times in A Row
	if available == false {
		RESET_COUNTER += 1
		if RESET_COUNTER > 5 {
			RESET_COUNTER = 0
			return c.JSON( fiber.Map{
				"url": "/youtube/live/next" ,
				"video_id": video_id ,
				"result": false ,
			})
		}
		fmt.Println( "Deleting ," , video_id )
		s.DB.ZRem( context.Background() , "STATE.YOUTUBE.LIVE.VIDEOS" , video_id )
		return s.YouTubeLiveNext( c )
	}
	uri := fmt.Sprintf( "https://www.youtube.com/watch?v=%s" , video_id )
	s.YouTubeOpenID( uri )
	// s.Set( "STATE.YOUTUBE.NOW_PLAYING" , video_id )
	s.Set( "active_player_now_playing_id" , video_id )
	s.Set( "active_player_now_playing_text" , "" )
	return c.JSON( fiber.Map{
		"url": "/youtube/live/next" ,
		"video_id": video_id ,
		"result": true ,
	})
}

func ( s *Server ) YouTubeLivePrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubeLivePrevious()" )
	video_id := circular_set.Previous( s.DB , "STATE.YOUTUBE.LIVE.VIDEOS" )
	available := s.YouTubeIsVideoIdAvailable( video_id )
	if available == false {
		fmt.Println( "Deleting ," , video_id )
		s.DB.ZRem( context.Background() , "STATE.YOUTUBE.LIVE.VIDEOS" , video_id )
		return s.YouTubeLivePrevious( c )
	}
	uri := fmt.Sprintf( "https://www.youtube.com/watch?v=%s" , video_id )
	s.YouTubeOpenID( uri )
	s.Set( "active_player_now_playing_id" , video_id )
	s.Set( "active_player_now_playing_text" , "" )
	return c.JSON( fiber.Map{
		"url": "/youtube/live/previous" ,
		"video_id": video_id ,
		"result": true ,
	})
}

func ( s *Server ) YouTubeVideo( c *fiber.Ctx ) ( error ) {
	video_id := c.Params( "video_id" )
	log.Debug( fmt.Sprintf( "YouTubeVideo( %s )" , video_id ) )
	uri := fmt.Sprintf( "https://www.youtube.com/watch?v=%s" , video_id )
	s.YouTubeOpenID( uri )
	s.Set( "active_player_now_playing_id" , video_id )
	s.Set( "active_player_now_playing_text" , "" )
	return c.JSON( fiber.Map{
		"url": "/youtube/:video_id" ,
		"video_id": video_id ,
		"result": true ,
	})
}

// am start -a android.intent.action.VIEW -d "vnd.youtube:PLcW8xNfZoh7cgn3QNojv1ly5NZfg9xVN9?listType=playlist"
// https://www.youtube.com/playlist?list=PLcW8xNfZoh7cgn3QNojv1ly5NZfg9xVN9

type YoutubeVideoInfo struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}
func ( s *Server ) YouTubeIsVideoIdAvailable( video_id string ) ( result bool ) {
	result = false
	next_api_key := circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
	base_url := "https://youtube.googleapis.com/youtube/v3/videos"
	params := url.Values{}
	params.Add( "part" , "id" )
	params.Add( "id" , video_id )
	params.Add( "key" , next_api_key )
	full_url := fmt.Sprintf( "%s?%s" , base_url , params.Encode() )
	req, _ := http.NewRequest( "GET" , full_url , nil )
	req.Header.Add( "Accept" , "application/json" )
	resp , _ := http.DefaultClient.Do( req )
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		fmt.Println( "api key banned ???" , next_api_key )
		return s.YouTubeIsVideoIdAvailable( video_id )
	}
	var video_info YoutubeVideoInfo
	body , _ := ioutil.ReadAll( resp.Body )
	// fmt.Println( string( body ) )
	json.Unmarshal( body , &video_info )
	if len( video_info.Items ) >= 1 {
		result = true
	}
	return
}

// 1.) Get Channels channelID
// curl \
// 'https://youtube.googleapis.com/youtube/v3/channels' \
// --header 'Accept: application/json' \
// --compressed \
// --get \
// --data-urlencode 'part=snippet' \
// --data-urlencode 'forUsername=MontereyBayAquarium' \
// --data-urlencode 'key=asdf'

// 2.) Get Channels Live Videos
// curl \
// 'https://youtube.googleapis.com/youtube/v3/search' \
// --header 'Accept: application/json' \
// --compressed \
// --get \
// --data-urlencode 'part=snippet' \
// --data-urlencode 'channelId=UCnM5iMGiKsZg-iOlIO2ZkdQ' \
// --data-urlencode 'eventType=live' \
// --data-urlencode 'maxResults=50' \
// --data-urlencode 'type=video' \
// --data-urlencode 'key=asdf'

// or just goto dev console on youtube channel's /stream page
// and run `ytInitialData.metadata.channelMetadataRenderer.externalId`
func ( s *Server ) YouTubeGetChannelId( channel_name string ) ( result string ) {
	next_api_key := circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
	headers := map[string]string{
		"Accept": "application/json" ,
	}
	params := map[string]string{
		"part": "snippet" ,
		"forUsername": channel_name ,
		"key": next_api_key ,
	}
	response_json := utils.GetJSON( "https://youtube.googleapis.com/youtube/v3/channels" , headers , params )
	items , _ := response_json.( map[string]interface{} )[ "items" ].( []interface{} )
	if len( items ) < 1 { fmt.Println( response_json ); return }
	first_result , _ := items[ 0 ].( map[string]interface{} )
	result , _ = first_result[ "id" ].( string )
	fmt.Println( result )
	return
}

type YouTubeResponseItem struct {
	Id struct {
		VideoId string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		Title string `json:"title"`
	} `json:"snippet"`
}
type YoutubeResponse struct {
	Items []YouTubeResponseItem `json:"items"`
}
// https://developers.google.com/youtube/v3/docs/search/list
func ( s *Server ) YouTubeGetChannelsLiveVideos( channel_id string ) ( result []YoutubeVideo ) {
	next_api_key := circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
	base_url := "https://youtube.googleapis.com/youtube/v3/search"
	params := url.Values{}
	params.Add( "part" , "snippet" )
	params.Add( "channelId" , channel_id )
	params.Add( "eventType" , "live" )
	params.Add( "maxResults" , "50" )
	params.Add( "type" , "video" )
	params.Add( "key" , next_api_key )
	full_url := fmt.Sprintf( "%s?%s" , base_url , params.Encode() )
	max_retries := 1
	for i := 0; i < max_retries; i++ {
		req , _ := http.NewRequest( "GET" , full_url , nil )
		req.Header.Add( "Accept" , "application/json" )
		resp , _ := http.DefaultClient.Do( req )
		defer resp.Body.Close()

		if resp.StatusCode == 403 {
			fmt.Println( "api key banned ???" , next_api_key )
			return s.YouTubeGetChannelsLiveVideos( channel_id )
		}

		body , _ := ioutil.ReadAll( resp.Body )
		var live_videos YoutubeResponse
		json.Unmarshal( body, &live_videos )
		if len( live_videos.Items ) >= 1 {
			for _ , item := range live_videos.Items {
				video := YoutubeVideo{
					Id: item.Id.VideoId,
					Name: item.Snippet.Title,
				}
				result = append( result , video )
			}
			return result
		} else {
			fmt.Println( "failed. Retrying..." )
			next_api_key = circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
			// fmt.Println( string( body ) )
			time.Sleep( 2 * time.Second )
		}
	}
	fmt.Println( "Max retries reached. No live videos found." )
	return
}

// Update DB With List of Currated Live Followers
// fucking idiots with this god damn quota. bro
func ( s *Server ) YouTubeLiveUpdate() ( result []string ) {
	// s.DB.Del( context.Background() , "STATE.YOUTUBE.LIVE.VIDEOS" )
	for channel_id , _ := range s.Config.Library.YouTube.Following.Live {
		fmt.Println( "\n" , channel_id , s.Config.Library.YouTube.Following.Live[ channel_id ].Name )
		live_videos := s.YouTubeGetChannelsLiveVideos( channel_id )
		for _ , video_item := range live_videos {
			for _ , video_name := range s.Config.Library.YouTube.Following.Live[ channel_id ].Videos {
				if strings.Contains( strings.ToLower( video_item.Name ) , video_name ) {
					fmt.Println( "adding" , video_item.Id , video_item.Name )
					result = append( result , video_item.Id )
				}
			}
		}
		time.Sleep( 1 * time.Second )
	}
	rand.Shuffle( len( result ) , func( i , j int ) {
		result[ i ] , result[ j ] = result[ j ] , result[ i ]
	})
	for _ , video_id := range result {
		circular_set.Add( s.DB , "STATE.YOUTUBE.LIVE.VIDEOS" , video_id )
	}
	return
}

func ( s *Server ) GetYouTubeLiveUpdate( c *fiber.Ctx ) ( error ) {
	live := s.YouTubeLiveUpdate()
	return c.JSON( fiber.Map{
		"url": "/youtube/live/update" ,
		"live": live ,
		"result": true ,
	})
}

type YoutubePlaylistResponse struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		Snippet struct {
			ResourceId struct {
				VideoId string `json:"videoId"`
			} `json:"resourceId"`
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

type YoutubeVideo struct {
	Id   string
	Name string
}

func ( s *Server ) YouTubeGetPlaylistVideos( playlist_id string ) ( result []YoutubeVideo ) {
	next_api_key := circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
	base_url := "https://www.googleapis.com/youtube/v3/playlistItems"
	params := url.Values{}
	params.Add( "part" , "snippet" )
	params.Add( "playlistId" , playlist_id )
	params.Add( "maxResults" , "50" )
	params.Add( "key" , next_api_key )

	var nextPageToken string
	max_retries := 5 // Consider adjusting based on your error handling strategy

	for i := 0; i < max_retries; i++ {
		if nextPageToken != "" {
			params.Set("pageToken", nextPageToken)
		}

		full_url := fmt.Sprintf("%s?%s", base_url, params.Encode())
		req, _ := http.NewRequest("GET", full_url, nil)
		req.Header.Add("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			time.Sleep(2 * time.Second) // Backoff before retry
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 403 {
			fmt.Println("API key banned ???", next_api_key)
			next_api_key = circular_set.Next(s.DB, "CONFIG.YOUTUBE.API_KEYS")
			params.Set("key", next_api_key) // Update API key and retry
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		var playlistResponse YoutubePlaylistResponse
		json.Unmarshal(body, &playlistResponse)

		for _, item := range playlistResponse.Items {
			video := YoutubeVideo{
				Id:   item.Snippet.ResourceId.VideoId,
				Name: item.Snippet.Title,
			}
			fmt.Println( video )
			result = append(result, video)
		}

		nextPageToken = playlistResponse.NextPageToken
		if nextPageToken == "" {
			break // Exit the loop if there's no next page
		}
	}

	if len(result) == 0 {
		fmt.Println("No videos found in playlist or max retries reached.")
	}
	return
}