## Install

If needed, [install Raspbian](https://github.com/noperator/guides/blob/master/install_raspbian.md).

Install display server, browser, and other dependencies.
```
sudo apt install -y \
xserver-xorg-core x11-xserver-utils xinit \
chromium-browser \
git jq
```

Download `websocat`.
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

Schedule hourly forecast refresh.
```
crontab -e
# 0 * * * * /bin/bash /home/pi/spindrift/refresh.sh
```

If needed, fix `startx` error, "Only console users are allowed to run the X server."
```
sudo sed -i -E 's/(allowed_users=)console/\1anybody/' /etc/X11/Xwrapper.config
```

Launch kiosk (absolute path required).
```
startx /home/pi/spindrift/launch.sh -- -nocursor &
```

## Install Puppeteer

Install Node.js and Puppeteer.
```
curl -sL https://deb.nodesource.com/setup_10.x | sudo bash -
sudo apt install -y nodejs
npm install puppeteer
```

Get MSW screenshots.
```
node screenshot.js
```
