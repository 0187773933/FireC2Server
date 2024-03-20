# Fire C2 Server

[![Go Reference](https://pkg.go.dev/badge/github.com/0187773933/FireC2Server.svg)](https://pkg.go.dev/github.com/0187773933/FireC2Server)

## Example

- http://localhost:5954/spotify/playlist-shuffle/3UMDmO2YJb8DgUjpSBu8y9?k=asdf


## TODO

- all "fresh" / "wake-up" commands need to first check if we are on the amazon firecube profile selection screen ...
	`adb shell pm list users`
- switching back to redis
	- https://github.com/0187773933/RedisCircularList
- Find some new way to get the config and db reference into logger package.
	- apparently there is init() in go? that gets called before main starts ...
- Audio playing locally through speakers here instead of on streamdeck
- HTML Control Panels
- VLC
	- Random M3U8 Music Tracks
- Spotify
	- Detect if media control overlay is already open
		- add just generic get status pre-call
			- don't press anything , get 1 screenshot , detect multiple things
				- now playing , which key index , shuffle status , media_controls open?
- YouTube
	- adhoc playlist support
	- fix s.YouTubeIsVideoIdAvailable( video_id )
		- stream could just be ended , but video_id still valid
			- delete these too
- Disney
	- detect spining circle , stall-out
	- force app reload with same uuid
- Twitch
	- take screenshots of stream ui
	- have to detect where we are in the quality selection menu
	- "weird" quality selection menu trap
- start storing better state details to improve status
- fix sleep ?
	- `adb shell settings get secure sleep_timeout`
	- `adb shell settings put secure sleep_timeout 2147483647`
- PushOver Notifications
- Automatic Staging on HDMI of Twitch Live Streams
- SSH ?
	- https://arachnoid.com/android/SSHelper/index.html
		- https://arachnoid.com/android/SSHelper/resources/SSHelper.apk
		- https://arachnoid.com/android/SSHelper/resources/SSHelper_source.tar.bz2

## ADB First Time Connection

1. Enter new Docker Container and Run :

	```
	adb connect $FireCubeIP:5555
	failed to authenticate to 192.168.4.193:5555
	```

2. Accept the connection on the Fire Cube TV
3. Confirm Connection :
	`adb devices`
4. `sudo docker cp public-fire-c2-server:/home/morphs/.android ADB_KEYS`
5. `sudo chown -R 1000:1000 ADB_KEYS/`
6. Change dockerRun.sh to mount the ADB_KEYS folder

## Inventory

| Item                           | Quantity | Price  | Link                                                         |
| ------------------------------ | -------- | ------ | ------------------------------------------------------------ |
| Firecube                       | 1        | 139.99 | https://www.amazon.com/gp/product/B09BZZ3MM7                 |
| Computer                       | 1        | 105.99 | https://www.amazon.com/gp/product/B07QY8LDGX                 |
| Pulse 8 HDMI CEC Adapter       | 1        | 43.61  | https://www.pulse-eight.com/p/104/usb-hdmi-cec-adapter       |
| IR Droid USB IR Transceiver    | 1        | 37.00  | https://irdroid.com/irdroid-usb-ir-transceiver               |
| TV Stand                       | 1        | 65.99  | https://www.amazon.com/gp/product/B07V4PK8D5                 |
| Stream Deck Mini               | 1        | 59.99  | https://www.amazon.com/Elgato-Stream-Deck-Mini-customizable/dp/B07DYRS1WH |
| Ethernet Switch                | 1        | 18.99  | https://www.amazon.com/gp/product/B00A121WN6                 |
| USB Extension Cable 15 ft      | 2        | 12.97  | https://www.amazon.com/gp/product/B08M68HMJZ                 |
| HDMI CEC Less Adapter          | 1        | 23.40  | https://www.amazon.com/gp/product/B00DL48KVI                 |
| 128 GB Flash Drive             | 1        | 12.99  | https://www.amazon.com/gp/product/B015CH1PJU                 |
| Power Strip Tower              | 1        | 37.99  | https://www.amazon.com/gp/product/B01HPB7E9Q                 |
| Crate                          | 1        | 22.50  | https://www.amazon.com/gp/product/B0054029SM                 |
| Zip Ties                       | 1        | 5.99   | https://www.amazon.com/HAVE-ME-TD-Cable-Ties/dp/B08TVLYB3Q   |
| HDMI Cable 10 ft               | 1        | 7.19   | https://www.amazon.com/PowerBear-Cable-Braided-Nylon-Connectors/dp/B07X37CG9V |
| USB Sound Card Adapter         | 1        | 13.99  | https://www.amazon.com/gp/product/B072BMG9TB                 |
| USB Sound Bar                  | 1        | 20.69  | https://www.amazon.com/gp/product/B085HZPNRJ                 |
| Anti-Slip Rubber Pad           | 1        | 6.49   | https://www.amazon.com/gp/product/B09GLTTTTW                 |
| 4 Port USB Hub                 | 1        | 17.99  | https://www.amazon.com/gp/product/B083XTKV8V                 |
| Microphone                     | 1        | 11.99  | https://www.amazon.com/gp/product/B075VQ7VG7                 |
| Microphone Mute Switch Adapter | 1        | 17.88  | https://www.amazon.com/gp/product/B08DRCRNRQ                 |
| Ethernet Cords - 5 Foot        | 1        | 12.99  | https://www.amazon.com/gp/product/B00E5I7VJG                 |
| Cord Protector Sleeve          | 1        | 22.99  | https://www.amazon.com/gp/product/B07FW3MKGH                 |