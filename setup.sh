#!/usr/bin/env sh

go install github.com/githubnemo/CompileDaemon@latest
CompileDaemon -log-prefix=false -build "go build -o bin/plain.do ./main.go" -command "./bin/plain.do" -exclude-dir=".git"
