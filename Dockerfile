FROM scratch

ENV ROOT /index.html
CMD ["/webserver"]
COPY webserver /webserver

