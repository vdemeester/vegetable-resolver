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

# [Optional] Uncomment the next lines
# to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>