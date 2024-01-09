FROM alpine as slim

# promoting json to main format along with jq for access
RUN apk --no-cache add jq

# autocomplete support
RUN apk --no-cache add bash bash-completion
RUN echo 'source <(bns completion bash)' >> ~/.bashrc

# common tools
RUN apk --no-cache add curl sed

# autocomplete
RUN echo 'source <(bns completion bash)' >> /root/.bashrc

# binaries
COPY bns /usr/bin

# main config file
RUN mkdir /root/.bunnyshell
COPY config.sample.yaml /root/.bunnyshell/config.yaml

# default to the cli
ENTRYPOINT ["bns"]
