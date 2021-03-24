#!/bin/bash

if [[ "$UID" -ne 0 ]]; then
    echo 'This script requires elevated privileges to install a few packages. You may need to enter your password.'
    sudo true
fi

# Install NodeSource Node.js 10.x repo (https://github.com/nodesource/distributions#debinstall).
curl -sL https://deb.nodesource.com/setup_10.x | sudo bash -

# Install display server, browser, Node.js, and other dependencies.
sudo apt install -y \
    xserver-xorg-core x11-xserver-utils xinit \
    chromium-browser \
    nodejs jq

# Install Puppeteer (https://github.com/puppeteer/puppeteer).
# npm install puppeteer

# Install Spindrift forecast checker. This installs Playwright.
go install

# Install WebSocat (https://github.com/vi/websocat).
sudo wget -O /usr/local/bin/websocat $(curl -sk https://api.github.com/repos/vi/websocat/releases |
    jq -r '.[0] | .assets[] | select(.name == "websocat_arm-linux-static") | .browser_download_url')
sudo chmod +x /usr/local/bin/websocat

# Schedule cron jobs.
(crontab -l; echo '
@reboot    /bin/bash /home/pi/spindrift/startx.sh   # Start spindrift.
0 *  * * * /bin/bash /home/pi/spindrift/refresh.sh  # Refresh forecast.
0 8  * * * /usr/bin/vcgencmd display_power 1        # Turn on display.
0 22 * * * /usr/bin/vcgencmd display_power 0        # Turn off display.
') | sort -ru | sed '/^$/d' | crontab -
