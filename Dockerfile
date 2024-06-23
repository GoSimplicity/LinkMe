FROM ubuntu:20.04
COPY LinkMe /app/LinkMe
WORKDIR /app
EXPOSE 9999
CMD ["/app/LinkMe"]