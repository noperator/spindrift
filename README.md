## Install

If needed, [install Raspbian](https://github.com/noperator/guides/blob/master/install_raspbian.md).

Install display server, browser, and other dependencies.
```
# Install NodeSource Node.js 10.x repo.
curl -sL https://deb.nodesource.com/setup_10.x | sudo bash -

sudo apt install -y \
xserver-xorg-core x11-xserver-utils xinit \
chromium-browser \
nodejs git jq
```

Install Puppeteer.
```
npm install puppeteer
```

Install WebSocat.
```
sudo wget -O /usr/local/bin/websocat $(curl -sk https://api.github.com/repos/vi/websocat/releases | jq '.[0] | .assets[] | select(.name == "websocat_arm-linux-static") | .browser_download_url' -r)
sudo chmod +x /usr/local/bin/websocat
```

Clone this repo.
```
git clone https://github.com/noperator/spindrift && cd spindrift
```

Specify surf spot in `.env` config file. For example, pull the string `Laniakea-Surf-Report/3672/` from Laniakea's surf forecast URL `https://magicseaweed.com/Laniakea-Surf-Report/3672/`.
```
export SPOT="Laniakea-Surf-Report/3672/"
```

Schedule hourly forecast refresh, and save power while display is not in use.
```
crontab -e
0 *  * * * /bin/bash /home/pi/spindrift/refresh.sh  # Refresh forecast.
0 8  * * * /usr/bin/vcgencmd display_power 1        # Turn on display.
0 22 * * * /usr/bin/vcgencmd display_power 0        # Turn off display.
```

If needed, fix `startx` error, "Only console users are allowed to run the X server."
```
sudo sed -i -E 's/(allowed_users=)console/\1anybody/' /etc/X11/Xwrapper.config
```

Launch kiosk (absolute path required).
```
startx /home/pi/spindrift/launch.sh -- -nocursor &
```
