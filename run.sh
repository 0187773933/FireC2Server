#!/bin/bash

# brew install opencv
# brew info opencv
# brew --prefix opencv
# ls /usr/local/Cellar/opencv/4.9.0_8.reinstall/lib
# rm /usr/local/opt/opencv
# ln -s /usr/local/Cellar/opencv/4.9.0_8.reinstall /usr/local/opt/opencv
# otool -l /usr/local/Cellar/opencv/4.9.0_8.reinstall/lib/libopencv_gapi.4.9.0.dylib | grep -A 2 LC_RPATH

# fix paths ??? brew ???
# OPENCV_LIB_DIR="/usr/local/Cellar/opencv/4.9.0_8.reinstall/lib"
# CORRECT_RPATH="/usr/local/Cellar/opencv/4.9.0_8.reinstall/lib"
# OLD_RPATH="/usr/local/Cellar/opencv/4.9.0_8/lib"

# for LIB in $(find $OPENCV_LIB_DIR -name "*.dylib"); do
#     install_name_tool -rpath $OLD_RPATH $CORRECT_RPATH $LIB
# done

# brew install ffmpeg@6
# export DYLD_LIBRARY_PATH=/usr/local/opt/opencv/lib:$DYLD_LIBRARY_PATH
#export DYLD_LIBRARY_PATH=/usr/local/opt/opencv/lib:/usr/local/opt/ffmpeg@6/lib:$DYLD_LIBRARY_PATH
#export DYLD_FALLBACK_LIBRARY_PATH=/usr/local/opt/opencv/lib:/usr/local/opt/ffmpeg@6/lib

export DYLD_LIBRARY_PATH=$(brew --prefix opencv)/lib:$DYLD_LIBRARY_PATH
export PKG_CONFIG_PATH=$(brew --prefix opencv)/lib/pkgconfig:$PKG_CONFIG_PATH

# CGO_ENABLED=1 LOG_LEVEL=debug go run main.go
LOG_LEVEL=debug go run main.go
#LOG_LEVEL=debug air