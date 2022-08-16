FROM alpine

COPY bunnyshell-cli /

ENTRYPOINT ["/bunnyshell-cli"]
