ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.0

FROM golang:$GO_VERSION-alpine${ALPINE_VERSION}
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000
ARG GOPLS_VERSION=v0.16.1

RUN adduser $USERNAME -s /bin/sh -D -u $USER_UID $USER_GID && \
    mkdir -p /etc/sudoers.d && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

RUN apk add -q --update --progress --no-cache util-linux

RUN apk add -q --update --progress --no-cache \
    git sudo openssh-client bash zsh curl zsh-vcs make gpg graphviz \
    python3 py3-pip yamllint jq curl unzip git protobuf 
    
# RUN go get github.com/golang/protobuf/protoc-gen-go \
# Create a symlink for python
RUN ln -sf python3 /usr/bin/python

# RUN python3 -m ensurepip --upgrade
# RUN pip3 install --no-cache --upgrade pip setuptools

RUN go install golang.org/x/tools/gopls@${GOPLS_VERSION}
RUN for tool in tools/gopls tools/cmd/goimports lint/golint; \
    do go install golang.org/x/${tool}@latest; \
    done

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.26.1
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.26.1

RUN sudo cp /go/bin/protoc-gen-go /usr/bin/

# Detect shadowing bugs
RUN go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@v0.30.0

# Visual Studio Code tools
RUN go install github.com/cweill/gotests/gotests@latest
RUN go install github.com/fatih/gomodifytags@latest
RUN go install github.com/josharian/impl@latest
RUN go install github.com/haya14busa/goplay/cmd/goplay@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@2023.1.7
# RUN go get google.golang.org/grpc

# Setup protobuf
# RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v28.0/protoc-28.0-linux-x86_64.zip
# RUN unzip protoc-28.0-linux-x86_64.zip -d /tmp_protoc
# RUN cp /tmp_protoc/bin/protoc /usr/local/bin/ \
#     && mv /tmp_protoc/include /usr/local/ \
#     && mv /usr/bin/protoc /usr/bin/protoc_old \
#     && rm -r -f /tmp_protoc \
#     && rm -f protoc-28.0-linux-x86_64.zip

# Setup shell
USER $USERNAME
RUN sh -c "$(wget -O- https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" "" --unattended &> /dev/null
ENV ENV="/home/$USERNAME/.ashrc" \
    ZSH=/home/$USERNAME/.oh-my-zsh \
    EDITOR=vi \
    LANG=en_US.UTF-8 \
    PATH=/home/vscode/.local/bin:$PATH
RUN printf 'ZSH_THEME="agnoster"\nENABLE_CORRECTION="false"\nplugins=(git copyfile extract colorize dotenv encode64 golang)\nsource $ZSH/oh-my-zsh.sh' > "/home/$USERNAME/.zshrc"
RUN echo "exec `which zsh`" > "/home/$USERNAME/.ashrc"

# To run mkdocs server as per Eitri standards
# COPY config/* .
# RUN python3 -m pip install -r requirements.txt

USER root

# Allow the developer to download third-party modules
RUN chown -R $USER_UID:$USER_GID /go/pkg