package types


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

type Library struct {
	Spotify SpotifyLibrary `yaml:"spotify"`
	Twitch TwitchLibrary `yaml:"twitch"`
	Disney DisneyLibrary `yaml:"disney"`
	YouTube YouTubeLibrary `yaml:"youtube"`
	VLC VLCLibrary `yaml:"vlc"`
}

type ConfigFile struct {
	ServerName string `yaml:"server_name"`
	ServerBaseUrl string `yaml:"server_base_url"`
	ServerLiveUrl string `yaml:"server_live_url"`
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
	BoltDBPath string `yaml:"bolt_db_path"`
	BoltDBEncryptionKey string `yaml:"bolt_db_encryption_key"`
	RedisAddress string `yaml:"redis_address"`
	RedisDBNumber int `yaml:"redis_db_number"`
	RedisPassword string `yaml:"redis_password"`
	StreamDeckServerUrl string `yaml:"stream_deck_server_url"`
	StreamDeckServerAPIKey string `yaml:"stream_deck_server_api_key"`
	TVWakeOnLan bool `yaml:"tv_wake_on_lan"`
	TVBrand string `yaml:"tv_brand"`
	TVIP string `yaml:"tv_ip"`
	TVWebSocketPort string `yaml:"tv_websocket_port"`
	TVMAC string `yaml:"tv_mac"`
	TVDefaultVolume int `yaml:"tv_default_volume"`
	TVDefaultInput string `yaml:"tv_default_input"`
	TVLGClientKey string `yaml:"tv_lg_client_key"`
	TVVizioAuthToken string `yaml:"tv_vizio_auth_token"`
	TVTimeoutSeconds int `yaml:"tv_timeout_seconds"`
	ADBPath string `yaml:"adb_path"`
	ADBConnectionType string `yaml:"adb_connection_type"`
	ADBSerial string `yaml:"adb_serial"`
	ADBServerIP string `yaml:"adb_server_ip"`
	ADBServerPort string `yaml:"adb_server_port"`
	YouTubeAPIKeys []string `yaml:"youtube_api_keys"`
	TwitchUserID string `yaml:"twitch_user_id"`
	TwitchClientID string `yaml:"twitch_client_id"`
	TwitchClientSecret string `yaml:"twitch_client_secret"`
	TwitchOAUTHToken string `yaml:"twitch_oauth_token"`
	TwitchAccessToken string `yaml:"twitch_access_token"`
	TwitchRefreshToken string `yaml:"twitch_refresh_token"`
	Library Library `yaml:"library"`
}