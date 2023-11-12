FROM nexus.uptimezeus.com/ubuntu:20.04

COPY  . .

ENTRYPOINT ["/project/app"]
