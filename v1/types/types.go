package types

import (
	tv_controller_types "github.com/0187773933/TVController/v1/types"
)

type SpotifyItem struct {
	Name string `yaml:"name"`
}
type SpotifyLibrary struct {
	Songs map[string]SpotifyItem `yaml:"songs"`
	Playlists map[string]SpotifyItem `yaml:"playlists"`
}

type TwitchFollowing struct {
	Currated []string `yaml:"currated"`
	All []string `yaml:"all"`
}
type TwitchLibrary struct {
	Following TwitchFollowing `yaml:"following"`
}

type DisneyItem struct {
	Name string `yaml:"name"`
}
type DisneyMovies struct {
	Currated map[string]DisneyItem `yaml:"currated"`
}
type DisneyLibrary struct {
	Movies DisneyMovies `yaml:"movies"`
}

type YoutubeLiveItem struct {
	Name string `yaml:"name"`
	Videos []string `yaml:"videos"`
}
type YoutubeSet struct {
	Live map[string]YoutubeLiveItem `yaml:"live"`
	Normal []string `yaml:"normal"`
	Relaxing []string `yaml:"relaxing"`
}
type YouTubeLibrary struct {
	Videos YoutubeSet `yaml:"movies"`
	Playlists YoutubeSet `yaml:"playlists"`
	Following YoutubeSet `yaml:"following"`
}

type VLCLibrary struct {
	Videos []string `yaml:"videos"`
}

type HuluEpisode struct {
	ID string `yaml:"id"`
	Name string `yaml:"name"`
}

type HuluTVShowSeason struct {
	Number string `yaml:"number"`
	Episodes []HuluEpisode `yaml:"episodes"`
}

type HuluTVShow struct {
	Name string `yaml:"name"`
	Seasons []HuluTVShowSeason `yaml:"seasons"`
}

type HuluMovie struct {
	Name string `yaml:"name"`
}

type HuluLibrary struct {
	Movies map[string]HuluMovie `yaml:"movies"`
	TV map[string]HuluTVShow `yaml:"tv"`
}

type NetflixEpisode struct {
	ID string `yaml:"id"`
	Name string `yaml:"name"`
}

type NetflixTVShowSeason struct {
	Number string `yaml:"number"`
	Episodes []NetflixEpisode `yaml:"episodes"`
}

type NetflixTVShow struct {
	Name string `yaml:"name"`
	Seasons []NetflixTVShowSeason `yaml:"seasons"`
}

type NetflixMovie struct {
	Name string `yaml:"name"`
}

type NetflixLibrary struct {
	Movies map[string]NetflixMovie `yaml:"movies"`
	TV map[string]NetflixTVShow `yaml:"tv"`
}

type Library struct {
	Spotify SpotifyLibrary `yaml:"spotify"`
	Twitch TwitchLibrary `yaml:"twitch"`
	Disney DisneyLibrary `yaml:"disney"`
	YouTube YouTubeLibrary `yaml:"youtube"`
	VLC VLCLibrary `yaml:"vlc"`
	Hulu HuluLibrary `yaml:"hulu"`
	Netflix NetflixLibrary `yaml:"netflix"`
}

type APKInfo struct {
	Package string `yaml:"package"`
	Activity string `yaml:"activity"`
}

type ConfigFile struct {
	ServerName string `yaml:"server_name"`
	ServerBaseUrl string `yaml:"server_base_url"`
	ServerLiveUrl string `yaml:"server_live_url"`
	ServerPrivateUrl string `yaml:"server_private_url"`
	ServerPublicUrl string `yaml:"server_public_url"`
	ServerPort string `yaml:"server_port"`
	ServerAPIKey string `yaml:"server_api_key"`
	ServerLoginUrlPrefix string `yaml:"server_login_url_prefix"`
	ServerCookieName string `yaml:"server_cookie_name"`
	ServerCookieSecret string `yaml:"server_cookie_secret"`
	ServerCookieAdminSecretMessage string `yaml:"server_cookie_admin_secret_message"`
	ServerCookieSecretMessage string `yaml:"server_cookie_secret_message"`
	AdminUsername string `yaml:"admin_username"`
	AdminPassword string `yaml:"admin_password"`
	TimeZone string `yaml:"time_zone"`
	SaveFilesPath string `yaml:"save_files_path"`
	BoltDBPath string `yaml:"bolt_db_path"`
	EncryptionKey string `yaml:"encryption_key"`
	RedisAddress string `yaml:"redis_address"`
	RedisDBNumber int `yaml:"redis_db_number"`
	RedisPassword string `yaml:"redis_password"`
	ReStreamServerUrlLocal string `yaml:"restream_server_url_local"`
	ReStreamServerUrl string `yaml:"restream_server_url"`
	ReStreamServerAPIKey string `yaml:"restream_server_api_key"`
	ReStreamServerHLSURLPrefix string `yaml:"restream_server_hls_url_prefix"`
	StreamDeckServerUrl string `yaml:"stream_deck_server_url"`
	StreamDeckServerAPIKey string `yaml:"stream_deck_server_api_key"`
	TV tv_controller_types.ConfigFile `yaml:"tv"`
	ADBPath string `yaml:"adb_path"`
	ADBConnectionType string `yaml:"adb_connection_type"`
	ADBSerial string `yaml:"adb_serial"`
	ADBServerIP string `yaml:"adb_server_ip"`
	ADBServerPort string `yaml:"adb_server_port"`
	ADBTimeoutSeconds int `yaml:"adb_timeout_seconds"`
	ADBDeviceType string `yaml:"adb_device_type"`
	APKS map[string]map[string]string `yaml:"apks"`
	FireCubeTotalUserProfiles int `yaml:"firecube_total_user_profiles"`
	FireCubeUserProfileIndex int `yaml:"firecube_user_profile_index"`
	HuluTotalUserProfiles int `yaml:"hulu_total_user_profiles"`
	HuluUserProfileIndex int `yaml:"hulu_user_profile_index"`
	NetflixTotalUserProfiles int `yaml:"netflix_total_user_profiles"`
	NetflixUserProfileIndex int `yaml:"netflix_user_profile_index"`
	YouTubeAPIKeys []string `yaml:"youtube_api_keys"`
	TwitchUserID string `yaml:"twitch_user_id"`
	TwitchClientID string `yaml:"twitch_client_id"`
	TwitchClientSecret string `yaml:"twitch_client_secret"`
	TwitchOAUTHToken string `yaml:"twitch_oauth_token"`
	TwitchAccessToken string `yaml:"twitch_access_token"`
	TwitchRefreshToken string `yaml:"twitch_refresh_token"`
	BrowserAPIKey string `yaml:"browser_api_key"` // silk browser audio/video player re-auth
	Library Library `yaml:"library"`
}