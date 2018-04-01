FROM golang:1.10

WORKDIR /tmp

# Update dependencies
RUN apt-get update -y && \
    apt-get install vim htop unzip -y

RUN wget https://releases.hashicorp.com/serf/0.8.1/serf_0.8.1_linux_amd64.zip && \
    unzip serf_0.8.1_linux_amd64.zip && \
    mkdir -p /opt/serf-0.8.1 && \
    mv serf /opt/serf-0.8.1 && \
    ln -s /opt/serf-0.8.1/serf /usr/bin/serf && \
    rm -rf serf_0.8.1_linux_amd64.zip

ADD ./init.sh /init.sh 
ENTRYPOINT /init.sh
