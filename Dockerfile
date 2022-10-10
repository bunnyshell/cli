FROM alpine as slim

# promoting json to main format along with jq for access
RUN apk --no-cache add jq

# autocomplete support
RUN apk --no-cache add bash bash-completion
RUN echo 'source <(bunnyshell-cli completion bash)' >> ~/.bashrc

# common tools
RUN apk --no-cache add curl sed

# autocomplete
RUN echo 'source <(bunnyshell-cli completion bash)' >> /root/.bashrc

# binaries
COPY bunnyshell-cli /usr/bin
# @deprecated but kept for backwards compatibility
RUN ln -sf /usr/bin/bunnyshell-cli /bunnyshell-cli

# main config file
RUN mkdir /root/.bunnyshell
COPY config.sample.yaml /root/.bunnyshell/config.yaml

# default to the cli
ENTRYPOINT ["bunnyshell-cli"]
