FROM ubuntu:20.04
COPY LinkMe /app/LinkMe
WORKDIR /app
CMD ["/app/LinkMe"]