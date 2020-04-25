FROM scratch

ENV INDEX TRUE
CMD ["/webserver"]
COPY webserver /webserver

