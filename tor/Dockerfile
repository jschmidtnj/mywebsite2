FROM debian:jessie

LABEL maintainer "joshua@joshuaschmidt.tech"

# Base packages
RUN apt-get update && \
    apt-get -y install \
    nginx \
    tor torsocks ntpdate

# yarn and node install
RUN curl -sL https://deb.nodesource.com/setup_11.x | bash - && \
    apt-get install -y nodejs && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list && \
    sudo apt-get update && sudo apt-get install yarn

# Compile shallot
ADD ./shallot /shallot
RUN apt-get -y install \
    build-essential \
    libssl-dev && \
    cd /shallot && \
    ./configure && \
    make && \
    mv ./shallot /bin && \
    cd / && \
    rm -Rf /shallot && \
    apt-get -y purge build-essential libssl-dev && \
    rm -Rf /var/lib/apt/lists/*

# Security and permissions
RUN useradd --system --uid 666 -M --shell /usr/sbin/nologin hidden

# Configure nginx logs to go to Docker log collection (via stdout/stderr)
RUN ln --symbolic --force /dev/stdout /var/log/nginx/access.log
RUN ln --symbolic --force /dev/stderr /var/log/nginx/error.log

# Main script
ADD ./main.sh /main.sh

# Tor Config
ADD ./torrc /etc/tor/torrc

# Add nginx default configuration 
ADD ./nginx.conf /etc/nginx/nginx.conf

# Configuration files and data output folder
VOLUME /web
WORKDIR /web

ENTRYPOINT ["/main.sh"]
CMD ["serve"]
