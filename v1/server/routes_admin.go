package server

import (
	// "fmt"
	// "strings"
	// fiber "github.com/gofiber/fiber/v2"
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
	twitch.Get( "/set/quality/max" , s.TwitchLiveSetQualityMax )

	// Disney
	disney := s.FiberApp.Group( "/disney" )
	disney.Use( validate_admin_mw )
	disney.Get( "/next" , s.DisneyMovieNext )
	disney.Get( "/previous" , s.DisneyMoviePrevious )
	disney.Get( "/movie/:movie_id" , s.DisneyMovie )

	// YouTube
	youtube := s.FiberApp.Group( "/youtube" )
	youtube.Use( validate_admin_mw )
	youtube.Get( "/:video_id" , s.YouTubeVideo )
	youtube.Get( "/live/next" , s.YouTubeLiveNext )
	youtube.Get( "/live/previous" , s.YouTubeLivePrevious )
	youtube.Get( "/update/live" , s.GetYouTubeLiveUpdate )
	// s.SetupMediaPlayerRoutes( youtube , "youtube" )

}