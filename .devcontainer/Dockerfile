FROM golang:1.24.7

RUN apt-get update
RUN apt-get install -y vim net-tools iputils-ping dnsutils zsh less lsof git jq

RUN git clone https://github.com/ohmyzsh/ohmyzsh.git ~/.oh-my-zsh
RUN cp ~/.oh-my-zsh/templates/zshrc.zsh-template ~/.zshrc
RUN chsh -s $(which zsh)

RUN go install golang.org/x/tools/gopls@v0.16.2
RUN go install honnef.co/go/tools/cmd/staticcheck@v0.5.1
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/cweill/gotests/gotests@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

ENV MAINTAINER=roy@koinvote.com