# go-telegramatorr

Simple bot (in development) for my particular use.

Currently has two functions:
1. Primary: getting current sessions from Jellyfin, as a response to comand /jellystatus in telegram. This returns some information about the currently running streams (filename, subtitle filename, bitrate, username).

    >Example:  
    >Here's an activity report from Jellyfin:  
    USER is playing (in progress): Agatha Christie's Poirot - S10E01 - Mystery of the Blue Train Bluray-1080p h265 AC3 2.0  
    Playback: Transcode  
    Bitrate: 4.43 Mbps  
    Subtitles: English - SUBRIP - External  
    Device: USER’s MacBook Pro



2. Beta(in development - included since image version 2.0.x): so-called botMonitor, which tracks user sessions (every 30s) and reports on them in the chat when they finish, providing duration and what was played. Restarting zeroes all session information. Sessions older than 3hrs are automatically purged.

    >Example:  
    User USER was playing Good People for 1 minutes - finished.

Requires 3 env vars:  
TELEGRAM_APIKEY,  
JELLYFIN_APIKEY,  
JELLYFIN_ADDRESS  
&  
TELEGRAM_CHATID if BOT_MONITOR set to true (enabled botMonitor beta feature)


Still a work in progress - sometimes subtitles aren't correct, no other bugs were found for now.

Image at https://hub.docker.com/repository/docker/januszadlo/go-telegramator/general, otherwise docker build yourself.

Image tags follow semantic versioning.
