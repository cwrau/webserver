# webserver
Extremely fast static file server without any bloat features

This webserver serves all files under `/serve`.

To improve performance, on startup it loads all files lying there into memory.

Can run without privileges, as non-root, and without any capabilities.

Also responds with `204` on `/` for health checks, e.g. for kubernetes.

`docker run --rm -it -v $PWD:/serve cwrau/webserver:1.0.0`
