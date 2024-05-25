package server

import (
	"fmt"
	// "wg"
	// "sync"
	"time"
	utils "github.com/0187773933/FireC2Server/v1/utils"
)

// find device :
	// lsusb
	// v4l2-ctl --list-devices

// get modes :
	// ffmpeg -f avfoundation -i "0" -t 1
	// v4l2-ctl --list-formats-ext --device /dev/video0

// 640x480
// 720x480
// 1024x768
// 1280x720
// 1280x960
// 1280x1024
// 1360x768
// 1600x1200
// 1920x1080
// 1920x1200
// 1920x2160
// 2048x1080
// 2048x1152
// 2560x1440
// 2560x1600

// run :
// ffmpeg -f avfoundation -video_size 1280x720 -framerate 30 -i "0:0" -c:v libx264 -preset ultrafast -tune zerolatency -f mpegts - | /Applications/IINA.app/Contents/MacOS/iina-cli --stdin --mpv-cache-secs=30 --mpv-cache

// ffmpeg -f avfoundation -video_size 1920x1080 -framerate 30 -pixel_format nv12 -i "0" -c:v libx264 -preset ultrafast -tune zerolatency -f mpegts - | /Applications/IINA.app/Contents/MacOS/iina-cli --stdin

// ffmpeg -f avfoundation -video_size 1920x1080 -framerate 30 -pixel_format nv12 -i "0:0" -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 192k -ac 2 -f mpegts - | /Applications/IINA.app/Contents/MacOS/iina-cli --stdin

// ffplay -f avfoundation -video_size 1280x720 -framerate 20 -pixel_format nv12 -i "0:0"

// ffmpeg -f avfoundation -video_size 1920x1080 -framerate 30 -pixel_format nv12 -probesize 50M -analyzeduration 100M -i "0:0" -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 192k -ac 2 -ar 48000 -f mpegts - | /Applications/IINA.app/Contents/MacOS/iina-cli --stdin

// ffmpeg -f avfoundation -video_size 1920x1080 -framerate 30 -pixel_format nv12 -i "0:0" -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 192k -ac 2 -ar 48000 -f mpegts - | ffplay -


// ffmpeg -f v4l2 -input_format mjpeg -framerate 30 -video_size 1920x1080 -i /dev/video0        -f alsa -i hw:1,0        -c:v libx264 -preset veryfast -tune zerolatency -b:v 2000k        -c:a aac -b:a 192k -ac 2        -f mpegts udp://239.255.255.250:5004?pkt_size=1316


// just run on linux
// ffmpeg -f v4l2 -input_format mjpeg -framerate 30 -video_size 1920x1080 -thread_queue_size 512 -i /dev/video0 \
//        -f alsa -thread_queue_size 512 -i hw:1,0 \
//        -c:v libx264 -preset superfast -tune zerolatency -b:v 2500k \
//        -c:a aac -b:a 128k -ac 2 \
//        -f mpegts udp://239.255.255.250:5004?pkt_size=1316

// stream back to iina on mac
// udp://239.255.255.250:5004


// https://gstreamer.freedesktop.org/documentation/installing/on-linux.html?gi-language=c
// apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio
// https://github.com/go-gst/go-gst?tab=readme-ov-file

func ( s *Server ) Governor() {
	ticker := time.NewTicker( 30 * time.Second )
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			fmt.Println( "Mutex acquired by governor, performing checks" )

			// TODO
			tv_status := s.TV.Status()
			utils.PrettyPrint( tv_status )


			fmt.Println( "Mutex released by governor" )
		}
	}
}