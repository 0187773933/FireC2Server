 package server

import (
	// "fmt"
	// "strings"
	// fiber "github.com/gofiber/fiber/v2"
	websocket "github.com/gofiber/contrib/websocket"
)

func ( s *Server ) SetupAdminRoutes() {

	// 2 x 3 StreamDeck-1
	streamdeck := s.FiberApp.Group( "/streamdeck" )
	streamdeck.Use( validate_admin_mw )
	streamdeck.Get( "spotify" , s.StreamDeckSpotify )
	streamdeck.Get( "youtube" , s.StreamDeckYouTube )
	streamdeck.Get( "disney" , s.StreamDeckDisney )
	streamdeck.Get( "twitch" , s.StreamDeckTwitch )
	streamdeck.Get( "escape-rope" , s.StreamDeckEscapeRope )
	streamdeck.Get( "heart" , s.StreamDeckHeart )

	// TV
	tv := s.FiberApp.Group( "/tv" )
	tv.Use( validate_admin_mw )
	tv.Get( "/prepare" , s.TVPrepare )
	tv.Get( "/power/on" , s.TVPowerOn )
	tv.Get( "/power/off" , s.TVPowerOff )
	tv.Get( "/power/status" , s.TVPowerStatus )
	tv.Get( "/input" , s.TVGetInput )
	tv.Get( "/input/:input" , s.TVSetInput )
	tv.Get( "/mute/on" , s.TVMuteOn )
	tv.Get( "/mute/off" , s.TVMuteOff )
	tv.Get( "/volume" , s.TVGetVolume )
	tv.Get( "/volume/:volume" , s.TVSetVolume )

	// Generic ADB Media Buttons
	adb := s.FiberApp.Group( "/adb" )
	adb.Use( validate_admin_mw )
	adb.Get( "play" , s.ADBPlay )
	adb.Get( "pause" , s.ADBPause )
	adb.Get( "stop" , s.ADBStop )
	adb.Get( "next" , s.ADBNext )
	adb.Get( "previous" , s.ADBPrevious )
	adb.Get( "status" , s.GetStatusUrl )

	// Responsive Media Buttons
	s.FiberApp.Get( "play" , s.Play )
	s.FiberApp.Get( "pause" , s.Pause )
	s.FiberApp.Get( "resume" , s.Resume )
	s.FiberApp.Get( "stop" , s.Stop )
	s.FiberApp.Get( "next" , s.Next )
	s.FiberApp.Get( "previous" , s.Previous )
	s.FiberApp.Get( "status" , s.GetStatusUrl )

	// Spotify
	spotify := s.FiberApp.Group( "/spotify" )
	spotify.Use( validate_admin_mw )
	spotify.Get( "/shuffle/:state" , s.SpotifySetShuffle )
	spotify.Get( "/song/:song_id" , s.SpotifySong )
	spotify.Get( "/playlist/:playlist_id" , s.SpotifyPlaylist )
	spotify.Get( "/playlist-shuffle/:playlist_id" , s.SpotifyPlaylistWithShuffle )
	// spotify.Get( "/next/song" , SpotifyPlaylistWithShuffle )
	// spotify.Get( "/next/playlist" , SpotifyNextPlaylist )
	spotify.Get( "/next/playlist-shuffle" , s.SpotifyNextPlaylistWithShuffle )
	// spotify.Get( "/previous/song" , SpotifyPlaylistWithShuffle )
	// spotify.Get( "/previous/playlist" , SpotifyPlaylistWithShuffle )
	// spotify.Get( "/previous" , SpotifyPressPreviousButton ) // needs a custom previous , requires 2 clicks if in shuffle-mode

	// Twitch
	twitch := s.FiberApp.Group( "/twitch" )
	twitch.Use( validate_admin_mw )
	twitch.Get( "/next" , s.TwitchLiveNext )
	twitch.Get( "/previous" , s.TwitchLivePrevious )
	twitch.Get( "/update" , s.GetTwitchLiveUpdate )
	twitch.Get( "/refresh" , s.GetTwitchLiveRefresh )
	twitch.Get( "/set/quality/max" , s.TwitchLiveSetQualityMax )
	twitch.Get( "/view/:username" , s.GetTwitchLiveUser )

	// Disney
	disney := s.FiberApp.Group( "/disney" )
	disney.Use( validate_admin_mw )
	disney.Get( "/next" , s.DisneyMovieNext )
	disney.Get( "/previous" , s.DisneyMoviePrevious )
	disney.Get( "/movie/:id" , s.DisneyMovie )

	// YouTube
	youtube := s.FiberApp.Group( "/youtube" )
	youtube.Use( validate_admin_mw )
		// Misc
	// youtube.Get( "/search/:query" , s.YouTubeSearch )
		// Live
	youtube.Get( "/live/next" , s.YouTubeLiveNext )
	youtube.Get( "/live/previous" , s.YouTubeLivePrevious )
	youtube.Get( "/live/update" , s.GetYouTubeLiveUpdate )
		// Custom Playlist Wrapper
	youtube.Get( "/playlist/:name/add/:id" , s.YouTubePlaylistAddVideo )
	youtube.Get( "/playlist/:name/add/playlist/:id" , s.YouTubePlaylistAddPlaylist )
	youtube.Get( "/playlist/:name/get" , s.YouTubePlaylistGet )
	youtube.Get( "/playlist/:name/delete/:video_id" , s.YouTubePlaylistDeleteVideo )
	youtube.Get( "/playlist/:name/index/get" , s.YouTubePlaylistGetIndex )
	youtube.Get( "/playlist/:name/index/set/:index" , s.YouTubePlaylistSetIndex )
	youtube.Get( "/playlist/:name/next" , s.YouTubePlaylistNext )
	youtube.Get( "/playlist/:name/previous" , s.YouTubePlaylistPrevious )

	// VLC
	vlc := s.FiberApp.Group( "/vlc" )
	vlc.Use( validate_admin_mw )
	vlc.Get( "/url/*" , s.VLCPlayURL )
	vlc.Get( "/playlist/:name/add/url/*" , s.VLCPlaylistAddURL )
	// vlc.Get( "/playlist/:name/get" , s.VLCPlaylistGet )
	// vlc.Get( "/playlist/:name/delete/:video_id" , s.VLCPlaylistDeleteVideo )
	// vlc.Get( "/playlist/:name/index/get" , s.VLCPlaylistGetIndex )
	// vlc.Get( "/playlist/:name/index/set/:index" , s.VLCPlaylistSetIndex )
	// vlc.Get( "/playlist/:name/next" , s.VLCPlaylistNext )
	// vlc.Get( "/playlist/:name/previous" , s.VLCPlaylistPrevious )

	// ReStreamer
	restreamer := s.FiberApp.Group( "/restream" )
	restreamer.Use( validate_admin_mw )
	restreamer.Get( "/url/*" , s.ReStreamURL )
	restreamer.Get( "/restart" , s.ReStreamRestart )
	restreamer.Get( "/stop" , s.ReStreamStop )

	// Firefox Focus
	s.FiberApp.Get( "/browser/audio" , s.GetBrowserAudioPlayer )
	s.FiberApp.Get( "/browser/video" , s.GetBrowserVideoPlayer )
	s.FiberApp.Get( "/browser/audio/set/:hash/position/:position" , validate_browser_mw , s.BrowserAudioPlayerSetPosition )
	s.FiberApp.Get( "/browser/ws/:id" , validate_browser_mw , websocket.New( s.BrowserWebSocketHandler ) )
	browser := s.FiberApp.Group( "/browser" )
	browser.Use( validate_admin_mw )
	browser.Get( "/url/*" , s.BrowserOpenURL )
	browser.Get( "/audio/*" , s.BrowserOpenAudioPlayer )
	browser.Get( "/video/*" , s.BrowserOpenVideoPlayer )

}