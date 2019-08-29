# Easily create and run hidden services

Easily run a hidden service inside the Tor network with this container

Generate the skeleton configuration for you hidden service, replace pattern
for your hidden service pattern name. Example, if you want to your hidden
service contain the word 'boss', just use this word as argument. You can use
regular expressions, like ```^boss```, will generate an address wich will start
with 'boss'. Be aware that bigger the pattern, more time it will take to
generate it.

```sh
sudo dockerd
sudo docker system prune -a
sudo docker rmi -f jschmidtnj/mywebsite2_tor:v0.0.2
sudo docker stop /hiddensite
sudo docker rm /hiddensite
sudo docker build -t jschmidtnj/mywebsite2_tor:v0.0.2 .
rm -rf web
mkdir web
sudo chown -R "$USER":666 web
sudo docker run -d --restart=always --name hiddensite -v $(pwd)/web:/web jschmidtnj/mywebsite2_tor:v0.0.2 ^josh
# after a long time
sudo chmod 755 web
sudo docker login
sudo docker push jschmidtnj/mywebsite2_tor:v0.0.2
```

add to compute engine: `docker.io/jschmidtnj/mywebsite2_tor:v0.0.2`

in the cloud:
```sh
rm -rf web
mkdir web
sudo chown -R "$USER":666 web
docker run -d --restart=always --name hiddensite -v $(pwd)/web:/web jschmidtnj/mywebsite2_tor:v0.0.2 ^josh
# after a long time
sudo chmod 755 web
```

copy .env file to cloud
clone git repo and build frontend. then send to nginx directory

make sure to open all ports

```sh
docker run -it --rm -v $(pwd)/web:/web \
       strm/tor-hiddenservice-nginx generate <pattern>
```

Create an container named 'hiddensite' to serve your generated hidden service

```sh
docker run -d --restart=always --name hiddensite -v $(pwd)/web:/web \
       strm/tor-hiddenservice-nginx 
```

## Example

Let's create a hidden service with the name beginning with strm.

```sh
docker pull strm/tor-hiddenservice-nginx
```

Wait to the container image be downloaded. And them we can generate our site
skeleton:

```sh
$docker run -it --rm -v $(pwd)/web:/web strm/tor-hiddenservice-nginx generate ^strm
[+] Generating the address with mask: ^strm
[+] Found matching domain after 137072 tries: strmfyygjp5st54g.onion
[+] Generating nginx configuration for site  strmfyygjp5st54g.onion
[+] Creating www folder
[+] Generating index.html template
```

Now we have our skeleton generated, we can run the container with:

```sh
docker run -d --restart=always --name hiddensite \
       -v $(pwd)/web:/web strm/tor-hiddenservice-nginx
```

And you have the service running ! :)
