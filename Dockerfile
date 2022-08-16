FROM scratch

COPY bunnyshell-cli /

ENTRYPOINT ["/bunnyshell-cli"]
