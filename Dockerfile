FROM debian:stable-slim

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ="US/Eastern"
ARG USERNAME="morphs"
ARG PASSWORD="asdf"

# https://github.com/Pulse-Eight/libcec/blob/master/docs/README.linux.md

RUN apt-get update -y && apt-get install -y \
    apt-transport-https \
    apt-utils \
    gcc \
    g++ \
    nano \
    tar \
    bash \
    sudo \
    openssl \
    git \
    make \
    wget \
    curl \
    net-tools \
    iproute2 \
    bc \
    pkg-config \
    cmake \
    libudev-dev \
    libxrandr-dev \
    libcec-dev \
    build-essential \
    python3-dev \
    python3-pip \
    python3-venv \
    build-essential \
    python3-dev \
    python3-setuptools \
    python3-smbus \
    python3-numpy \
    python3-scipy \
    libncursesw5-dev \
    libgdbm-dev \
    libc6-dev \
    zlib1g-dev \
    libsqlite3-dev \
    tk-dev \
    libssl-dev \
    openssl \
    libffi-dev \
    swig \
    libopencv-dev \
    libgtk2.0-dev \
    libavcodec-dev \
    libavformat-dev \
    libswscale-dev \
    python3-numpy \
    libtbb-dev \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    libsm6 \
    libxrender1 \
    libfontconfig1 \
    python3-opencv \
    python3-h5py \
    yasm \
    ffmpeg \
    libpq-dev \
    libxvidcore-dev \
    libx264-dev \
    libv4l-dev \
    libgtk-3-dev \
    libjpeg62 \
    libopenjp2-7 \
    libilmbase-dev \
    libatlas-base-dev \
    libgstreamer1.0-dev \
    openexr \
    libopenexr-dev \
    android-tools-adb \
    libasound2-dev \
    libhidapi-dev \
    udev \
    alsa-utils \
    && rm -rf /var/lib/apt/lists/*

ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/lib/pkgconfig

RUN groupadd -g 1000 morphs && \
    useradd -u 1000 -g morphs -m $USERNAME -p $PASSWORD -s "/bin/bash" && \
    mkdir -p /home/$USERNAME && \
    chown -R $USERNAME:$USERNAME /home/$USERNAME && \
    usermod -aG sudo $USERNAME && \
    echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers && \
    usermod -a -G dialout $USERNAME && \
    usermod -a -G audio $USERNAME
USER $USERNAME

# install go with specific version and progress
WORKDIR /home/$USERNAME
COPY ./go_install.sh /home/$USERNAME/go_install.sh
RUN sudo chmod +x /home/$USERNAME/go_install.sh
RUN sudo chown $USERNAME:$USERNAME /home/$USERNAME/go_install.sh
RUN /home/$USERNAME/go_install.sh
RUN sudo tar --checkpoint=100 --checkpoint-action=exec='/bin/bash -c "cmd=$(echo ZXhwb3J0IEdPX1RBUl9LSUxPQllURVM9JChwcmludGYgIiUuM2ZcbiIgJChlY2hvICIkKHN0YXQgLS1mb3JtYXQ9IiVzIiAvaG9tZS9tb3JwaHMvZ28udGFyLmd6KSAvIDEwMDAiIHwgYmMgLWwpKSAmJiBlY2hvIEV4dHJhY3RpbmcgWyRUQVJfQ0hFQ0tQT0lOVF0gb2YgJEdPX1RBUl9LSUxPQllURVMga2lsb2J5dGVzIC91c3IvbG9jYWwvZ28= | base64 -d ; echo); eval $cmd"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
RUN echo "PATH=$PATH:/usr/local/go/bin" | tee -a /home/$USERNAME/.bashrc

USER root

RUN git clone https://github.com/Pulse-Eight/platform.git && \
    mkdir platform/build && \
    cd platform/build && \
    cmake .. && \
    make -j4 && \
    make install && \
    ldconfig

RUN git clone https://github.com/Pulse-Eight/libcec.git && \
    mkdir libcec/build && \
    cd libcec/build && \
    cmake -DHAVE_LINUX_API=1 .. && \
    make -j4 && \
    make install && \
    ldconfig

# https://github.com/opencv/opencv/releases
# https://github.com/opencv/opencv/wiki
# https://docs.opencv.org/4.x/d2/de6/tutorial_py_setup_in_ubuntu.html
# https://docs.opencv.org/4.x/d7/d9f/tutorial_linux_install.html
#WORKDIR /build-opencv
#RUN git clone https://github.com/opencv/opencv.git && \
#    cd opencv && \
#    git checkout 4.8.0 && \
#    git clone https://github.com/opencv/opencv_contrib.git && \
#    cd opencv_contrib && \
#    git checkout 4.8.0 && \
#    cd ../.. && \
#    mkdir build && cd build && \
#    cmake -D CMAKE_BUILD_TYPE=RELEASE \
#    -D CMAKE_INSTALL_PREFIX=/usr/local \
#    -D OPENCV_EXTRA_MODULES_PATH=/build-opencv/opencv/opencv_contrib/modules \
#    -D BUILD_EXAMPLES=OFF \
#    /build-opencv/opencv && \
#    make -j$(nproc) && \
#    make install && \
#    ldconfig

WORKDIR /build-opencv
RUN git clone https://github.com/opencv/opencv.git && \
    cd opencv && \
    git checkout 4.8.0 && \
    git clone https://github.com/opencv/opencv_contrib.git && \
    cd opencv_contrib && \
    git checkout 4.8.0 && \
    cd ../.. && \
    mkdir build && cd build && \
    cmake -D CMAKE_BUILD_TYPE=RELEASE \
    -D OPENCV_GENERATE_PKGCONFIG=ON \
    -D WITH_CUDA=OFF \
    -D BUILD_opencv_python=OFF \
    -D BUILD_opencv_python2=OFF \
    -D BUILD_opencv_python3=OFF \
    -D CMAKE_BUILD_TYPE=RELEASE \
    -D BUILD_SHARED_LIBS=ON \
    -D CMAKE_INSTALL_PREFIX=/usr \
    -D INSTALL_C_EXAMPLES=OFF \
    -D INSTALL_PYTHON_EXAMPLES=OFF \
    -D BUILD_PYTHON_SUPPORT=OFF \
    -D BUILD_NEW_PYTHON_SUPPORT=OFF \
    -D WITH_TBB=ON \
    -D WITH_PTHREADS_PF=ON \
    -D WITH_OPENNI=OFF \
    -D WITH_OPENNI2=ON \
    -D WITH_EIGEN=ON \
    -D BUILD_DOCS=OFF \
    -D BUILD_TESTS=OFF \
    -D BUILD_PERF_TESTS=OFF \
    -D BUILD_EXAMPLES=OFF \
    -D WITH_OPENCL=$OPENCL_ENABLED \
    -D USE_GStreamer=ON \
    -D WITH_GDAL=ON \
    -D WITH_CSTRIPES=ON \
    -D ENABLE_FAST_MATH=1 \
    -D WITH_OPENGL=ON \
    -D WITH_QT=OFF \
    -D WITH_IPP=OFF \
    -D WITH_FFMPEG=ON \
    -D WITH_PROTOBUF=ON \
    -D BUILD_PROTOBUF=ON \
    -D CMAKE_SHARED_LINKER_FLAGS=-Wl,-Bsymbolic \
    -D WITH_V4L=ON \
    -D OPENCV_EXTRA_MODULES_PATH=/build-opencv/opencv/opencv_contrib/modules \
    /build-opencv/opencv && \
    make -j$(nproc) && \
    make install && \
    ldconfig


WORKDIR /home/$USERNAME
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER $USERNAME
ENTRYPOINT [ "/entrypoint.sh" ]