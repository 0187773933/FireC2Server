FROM debian:stable-slim

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ="US/Eastern"
ARG USERNAME="morphs"
ARG PASSWORD="asdf"

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
    && rm -rf /var/lib/apt/lists/*

# Setup User
RUN useradd -m $USERNAME -p $PASSWORD -s "/bin/bash" && \
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
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER $USERNAME
ENTRYPOINT [ "/entrypoint.sh" ]