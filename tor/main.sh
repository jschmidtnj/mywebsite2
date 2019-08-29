#!/bin/bash

if [ -z "$1" ]
then
    echo '[-] You dont provided any mask, please inform an mask to generate your address'
    exit -1
else
    echo '[+] Generating the address with mask: '$1
    shallot -f /tmp/key $1
    echo '[+] '$(grep Found /tmp/key)
    grep 'BEGIN RSA' -A 99 /tmp/key > /web/private_key
fi

address=$(grep Found /tmp/key | cut -d ':' -f 2 )

echo '[+] Generating nginx configuration for site '$address
echo 'server {' > /web/site.conf
echo '  listen 127.0.0.1:8080;' >> /web/site.conf
echo '  root /web/www/;' >> /web/site.conf
echo '  index index.html index.htm;' >> /web/site.conf
echo '  server_name '$address';' >> /web/site.conf
echo '}' >> /web/site.conf

echo '[+] Creating www folder'
echo '[+] Generating website'
if [ -f /web/src ] && [ -f /web/nuxt ] && [ -f /web/dist ]
then
    echo '[+] src directory found for website'
else
    git clone https://github.com/jschmidtnj/mywebsite2 /web/src
    cd /web/src/nuxt
    yarn
    yarn predeploy
    cd /
fi
rm -rf /web/www
chmod 755 /web/
chmod 755 /web/www
cp -ar /web/src/nuxt/dist /web/www
chown hidden:hidden -R /web/www

echo '[+] Initializing local clock'
ntpdate -B -q 0.debian.pool.ntp.org
echo '[+] Starting tor'
tor -f /etc/tor/torrc &
echo '[+] Starting nginx'
nginx &
# Monitor logs
sleep infinity
