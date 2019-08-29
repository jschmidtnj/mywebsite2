#!/bin/bash

if [ -f /web/private_key ]
then
    echo '[-] You already have an private key, delete it if you want to generate a new key'
    exit -1
fi
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
mkdir /web/www
chmod 755 /web/
echo '[+] Generating website'
if [ -f /web/src ]
then
    git clone https://github.com/jschmidtnj/mywebsite2 /web/src
    cd /web/src/nuxt
    yarn
    yarn predeploy
    mv /web/src/nuxt/dist /web/www
    chmod 755 /web/www
fi
chown hidden:hidden -R /web/www
chown "$USER":666 /web/www
chmod 755 /web/

echo '[+] Initializing local clock'
ntpdate -B -q 0.debian.pool.ntp.org
echo '[+] Starting tor'
tor -f /etc/tor/torrc &
echo '[+] Starting nginx'
nginx &
# Monitor logs
sleep infinity
