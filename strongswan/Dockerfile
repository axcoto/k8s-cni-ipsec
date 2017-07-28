FROM ubuntu

RUN wget http://download.strongswan.org/strongswan-5.5.3.tar.bz2 && \
    tar xvf strongswan-5.5.3.tar.bz2 && \
    cd strongswan-5.5.3 && \
    apt install build-essential libgmp-dev iproute2  && \
    ./configure --prefix=/usr --sysconfdir=/etc --with-piddir=/etc/ipsec.d/run && \
    make && make install

ADD ipsec.secrets /etc/ipsec.secrets
ADD ipsec.conf    /etc/ipsec.conf

RUN mkdir /etc/ipsec.d/run

EXPOSE 500/udp 4500/udp

ENTRYPOINT ["/entrypoint.sh"]
