# go-telegramatorr

Simple bot (in development) for my particular use.

Currently has three functions:
1. getting current sessions from Jellyfin and/or PLex, as a response to comand /playstatus in telegram chat with bot. This returns some information about the currently running streams (filename, subtitle filename, bitrate, username).

    >Example:  
    Here's an activity report from your player(s):  
    USER is playing(directplay) on Plex: Batman Begins  
    Bitrate: 6.46 Mbps  
    Device: DEVICE  
    Subs: English-SRT  

2. botMonitor, which tracks user sessions (every 30s) and reports on them in the chat when they finish, providing duration and what was played. Restarting zeroes currently tracked session information (in development). Sessions older than 3hrs are automatically purged.

    >Example:  
    User USER (DEVICE) was playing Batman Begins on Plex for 32 minutes  
    method: directplay  
    bitrate: 6.46 Mbps  
    subs: English-SRT  

Requires at least 3 env vars:  
TELEGRAM_APIKEY (required),  
and (to enable Jellyfin gatherer)  
JELLYFIN_APIKEY,  
JELLYFIN_ADDRESS  
and/or (to enable Plex gatherer)  
PLEX_ADDRESS  
PLEX_APIKEY  
and optionally (to enable botMonitor)  
TELEGRAM_CHATID  

EXAMPLE COMPOSE:   
```yaml  
version: "3"
services:
  go-telegramator:
    image: januszadlo/go-telegramator:(version)
    restart: always
    container_name: go-telegramator-test
    environment:
      TELEGRAM_APIKEY: "REPLACE"
      JELLYFIN_ADDRESS: "http://REPLACE:8096"
      JELLYFIN_APIKEY: "REPLACE"
      TELEGRAM_CHATID: "REPLACE"
      PLEX_ADDRESS: "http://REPLACE:32400"
      PLEX_APIKEY: "REPLACE"
      TZ: "Europe/Warsaw"
    volumes:
      - /home/user/telegramator:/data #example hostpath for reporting db
```

Still a work in progress - sometimes subtitles aren't correct, there are some inconsistencies in outputs.

Image at https://hub.docker.com/repository/docker/januszadlo/go-telegramator/general, otherwise docker build yourself.

Image tags follow semantic versioning.
