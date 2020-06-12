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

Here's how it works:
- Fetches hourly forecast data with [Puppeteer](https://github.com/puppeteer/puppeteer) (headless Chrome) scheduled in a cron job
- Renders the forecast, along with live cam streams, in kiosk-mode Chrome on a standalone X server
- Runs on a "headless" Raspberry Pi controlled remotely via SSH
- (Optional) Displays on an old, beautiful, hacked Apple iMac G4 17" 1440x900 px

## Getting Started

### Dependencies

If needed, [install Raspbian](https://github.com/noperator/guides/blob/master/install_raspbian.md), along with a display server, browser, and other dependencies.
```
# Install NodeSource Node.js 10.x repo.
curl -sL https://deb.nodesource.com/setup_10.x | sudo bash -

sudo apt install -y \
xserver-xorg-core x11-xserver-utils xinit \
chromium-browser \
nodejs git jq
```

Install [Puppeteer](https://github.com/puppeteer/puppeteer).
```
npm install puppeteer
```

Install [WebSocat](https://github.com/vi/websocat).
```
sudo wget -O /usr/local/bin/websocat $(curl -sk https://api.github.com/repos/vi/websocat/releases | jq -r '.[0] | .assets[] | select(.name == "websocat_arm-linux-static") | .browser_download_url')
sudo chmod +x /usr/local/bin/websocat
```

### Installing and Configuring

Clone this repo.
```
git clone https://github.com/noperator/spindrift
```

Start at boot, schedule hourly forecast refresh, and save power while display is not in use.
```
crontab -e
@reboot    /bin/bash /home/pi/spindrift/startx.sh   # Start spindrift.
0 *  * * * /bin/bash /home/pi/spindrift/refresh.sh  # Refresh forecast.
0 8  * * * /usr/bin/vcgencmd display_power 1        # Turn on display.
0 22 * * * /usr/bin/vcgencmd display_power 0        # Turn off display.
```

Specify surf spot and weather location in `config/.env` config file. For example, pull the string `Ala-Moana-Surf-Report/661/` from Ala Moana's surf forecast URL `https://magicseaweed.com/Ala-Moana-Surf-Report/661/`.
```
export SPOT="Ala-Moana-Surf-Report/661/"
export LOCATION="Waikiki HI"
```

Configure `config/streams.js` with the HLS live cam streams you'd like to display.
```
"priority": "1",
"url":      "<URL>",
"title":    "Ala Moana Bowls",
"source":   "Surfline"
```

### Executing

The kiosk will launch automatically when you boot the Raspberry Pi, but you can also start it manually:
```
./startx.sh
```

You can also manually refresh it:
```
./refresh.sh
```

### Troubleshooting

If needed, fix `startx` error, "Only console users are allowed to run the X server."
```
sudo sed -i -E 's/(allowed_users=)console/\1anybody/' /etc/X11/Xwrapper.config
```

## Todo

- [ ] Turn off display signal while sleeping, rather than blanking it
- [x] Add weather report
- [x] Load backup streams in order of preference
- [ ] Merge screenshot scripts
