#!/bin/sh

# get third party dependencies
# this can be commented out after the first run
go get github.com/bwmarrin/discordgo

# build a the bot as a static binary
go build -ldflags '-linkmode external -extldflags "-fno-PIC -static"'
