# FROM ghcr.io/vdemeester/test/vegetable-resolver
ARG VARIANT="1.20-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# Install additional packages
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        bat \
        fd-find \
        fzf \
        gnupg2 \
        jq \
        libcrypt-ssleay-perl \
        libnet-ssleay-perl \
        vim \
        make \
        ;
ARG KO_VERSION=0.14.1
RUN ARCH=$(uname -m) OS=$(uname) && \
    wget -O ko.tar.gz https://github.com/ko-build/ko/releases/download/v${KO_VERSION}/ko_${KO_VERSION}_${OS}_${ARCH}.tar.gz && \
    tar xvzf ko.tar.gz -C /usr/local/bin

# [Optional] Uncomment the next lines
# to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>