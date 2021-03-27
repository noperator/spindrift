#!/bin/bash

errcho() { >&2 echo -e '\n\x01\x1B[0;32m\x02==>' $@ '\x01\x1B[m\x02\n'; }

if [[ "$UID" -ne 0 ]]; then
    errcho 'This script requires elevated privileges to install a few packages. You may need to enter your password.'
    sudo true
fi

# https://github.com/nodesource/distributions#debinstall
errcho 'Installing Node.js APT source...'
curl -sL https://deb.nodesource.com/setup_10.x | sudo bash -

errcho 'Installing display server, browser, Node.js, and other dependencies...'
sudo apt update &&
sudo apt install -y \
    xserver-xorg-core x11-xserver-utils xinit \
    chromium-browser \
    nodejs jq

errcho 'Installing Spindrift forecast checker...'
go install

# https://github.com/vi/websocat
errcho 'Installing WebSocat...'
sudo wget -O /usr/local/bin/websocat $(curl -sk https://api.github.com/repos/vi/websocat/releases |
    jq -r '.[0] | .assets[] | select(.name == "websocat_arm-linux-static") | .browser_download_url')
sudo chmod +x /usr/local/bin/websocat

errcho 'Scheduling cron jobs...'
CRON_TMP=$(mktemp)
crontab -l | grep '^#' > "$CRON_TMP"
(crontab -l | grep -v '^#'; echo '
@reboot    /bin/bash "$HOME/spindrift/startx.sh"   # Start spindrift.
0 *  * * * /bin/bash "$HOME/spindrift/refresh.sh"  # Refresh forecast.
0 8  * * * /usr/bin/vcgencmd display_power 1       # Turn on display.
0 22 * * * /usr/bin/vcgencmd display_power 0       # Turn off display.
') | sort -ru | sed '/^$/d' >> "$CRON_TMP"
crontab "$CRON_TMP"
crontab -l
