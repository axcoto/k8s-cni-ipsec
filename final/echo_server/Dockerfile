FROM alpine

ADD echo /usr/local/bin/echo
RUN chmod +x /usr/local/bin/echo

RUN apk add --update curl && \
    rm -rf /var/cache/apk/*

CMD ["/usr/local/bin/echo"]
