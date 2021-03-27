# Spindrift

Surf report kiosk built on an [Apple Pi](https://imgur.com/gallery/4I8jm).

<div align="center">
  <kbd>
    <img src="screenshot.png" />
  </kbd>
</div>

## Description

Spindrift displays a dashboard of the following components, refreshed hourly:
- Surf forecast from Magicseaweed
- Weather forecast from Google
- Live cam streams from user-specified sources (e.g., Surfline)

As a bit of background on the `.m3u8` HLS live cam streams that Spindrift displays, here's a condensed explanation from [Wikipedia](https://en.wikipedia.org/wiki/M3U):
> M3U is a plain text computer file format for a multimedia playlist that specifies the locations of one or more media files. The M3U file can also include comments, prefaced by the "#" character.
> - In extended M3U, "#" also introduces extended M3U directives which are terminated by a colon ":" if they support parameters.
> - The Unicode version of M3U is M3U8, which uses UTF-8-encoded characters.
>
> Extended M3U8 files are the basis for HTTP Live Streaming (HLS), a format originally developed by Apple to stream video and radio to iOS devices, and later standardized by the IETF. In HLS, a master playlist references segment playlists which usually contain URLs for short parts of the media stream.

### Features

- Fetches hourly forecast data with Playwright scheduled in a cron job
- Renders the forecast, along with live cam streams, in Chromium running in kiosk mode and controlled by DevTools Protocol on a standalone X server
- Runs on a headless Raspberry Pi controlled remotely via SSH
- (Optional) Displays on an old, beautiful, hacked Apple iMac G4 17" at 1440 × 900 px

### Built with

- Chromium
- [Playwright](https://playwright.dev)
- [WebSocat](https://github.com/vi/websocat)
- [jq](https://github.com/stedolan/jq)
- X server

## Getting started

### Prerequisites

Install [Raspberry Pi OS](https://www.raspberrypi.org/software/) and [Golang](https://golang.org/doc/install). Note that when installing Golang, it's important to download the latest release available on the official Go downloads page, since APT's version is generally out of date.

### Install

Clone this repo to your home directory, and run the installer script.

```
git -C "$HOME" clone https://github.com/noperator/spindrift && cd "$HOME/spindrift"
./install.sh
```

This'll do a few things:
- Install [dependencies](#built-with)
- Schedule cron jobs
  - Start at boot
  - Schedule hourly forecast refresh
  - Turn off display from 10 PM–8 AM (presumably while not in use)

### Configure

Specify your surf spot and weather location in the `config/params.toml` config file. For example, pull the string `Ala-Moana-Surf-Report/661/` from Ala Moana's surf forecast URL `https://magicseaweed.com/Ala-Moana-Surf-Report/661/`.

```
spot = "Ala-Moana-Surf-Report/661/"
locaion = "Waikiki HI"
```

Configure `config/streams.js` with the HLS live cam streams you'd like to display.

```
"priority": "1",
"url":      "<URL>",
"title":    "Ala Moana Bowls",
"source":   "Surfline"
```

Note that this was developed with a 1440 × 900 px display in mind. If you need to make any display adjustments, you may do so in `check-forecast.go`, `launch.sh`, and `style.css`.

### Usage

The kiosk will launch automatically when you boot the Raspberry Pi, but you can also start it manually:

```
./startx.sh
```

You can also manually refresh it:

```
./refresh.sh
```

### Troubleshooting

When first running `check-forecast`, it'll try to download x86 browsers even though we're BYOB (i.e., we're specifying our own ARM-compatible Chromium via Playwright's `ExecutablePath` option); there's currently [no way to stop](https://github.com/mxschmitt/playwright-go/issues/52) browsers from downloading. You'll probably see an error like this:

```
2021/03/27 16:12:12 Downloading driver to /home/pi/.cache/ms-playwright-go/1.10.0-next-1615230258000
2021/03/27 16:12:19 Downloaded driver successfully
2021/03/27 16:12:19 Downloading browsers...
/home/pi/.cache/ms-playwright-go/1.10.0-next-1615230258000/playwright.sh: 3: /home/pi/.cache/ms-playwright-go/1.10.0-next-1615230258000/playwright.sh: /home/pi/.cache/ms-playwright-go/1.10.0-next-1615230258000/node: Exec format error
2021/03/27 16:12:20 could not start playwright: could not install driver: could not install browsers: exit status 2
```

That's because the Playwright driver brought its own x86 version of Node.js, and that won't run on this ARM-based Raspberry Pi. We can fix this by replacing that `node` binary with a link to a different one that'll work on this platform:

```
PW_NODE="$(find $HOME/.cache/ms-playwright-go -maxdepth 2 -name node)"
mv "$PW_NODE" "$PW_NODE.bu.$(date +%s)"
ln -s "$(which node)" "$PW_NODE"
```

Run `check-forecast` again, and you'll see:

```
2021/03/27 16:29:32 Downloading browsers...
Downloading chromium v857950 - 126.3 Mb [====================] 100% 0.0s
chromium v857950 downloaded to /home/pi/.cache/ms-playwright/chromium-857950
Downloading webkit v1443 - 78.2 Mb [====================] 100% 0.0s
webkit v1443 downloaded to /home/pi/.cache/ms-playwright/webkit-1443
Downloading firefox v1234 - 72.2 Mb [====================] 100% 0.0s
firefox v1234 downloaded to /home/pi/.cache/ms-playwright/firefox-1234
Downloading ffmpeg v1005 - 2.6 Mb [====================] 100% 0.0s
ffmpeg v1005 downloaded to /home/pi/.cache/ms-playwright/ffmpeg-1005
2021/03/27 16:33:38 Downloaded browsers successfully
```

If needed, fix `startx` error, "Only console users are allowed to run the X server."

```
sudo sed -i -E 's/(allowed_users=)console/\1anybody/' /etc/X11/Xwrapper.config
```

## Back matter

### See also

- [TunaSurf/ShouldIShred: Web app to check surf conditions at your local spots](https://github.com/TunaSurf/ShouldIShred)

### To-do

- [ ] Turn off display signal while sleeping, rather than blanking it
- [x] Add weather report
- [x] Load backup streams in order of preference
- [x] Merge screenshot scripts
- [ ] ~~Consolidate config files into a single JavaScript file, if possible~~
- [ ] Manually rotate through streams
- [x] Describe installing Playwright with Chromium on Raspberry Pi
- [x] Move display dimensions to config file
