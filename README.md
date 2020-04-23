# webserver
Extremely fast static file server without any bloat features

This webserver serves all files under `/serve` on `:8080`.

To improve performance, on startup it loads all files lying there into memory.

Can run without privileges, as non-root, and without any capabilities.

Environment variable `ROOT` can be set to a filename which should be served under `/`.
Default is set to `/index.html`

Also responds with `204` on `/` for health checks, e.g. for kubernetes, even if there is no file for `/`.

`docker run --rm -it -p 8080:8080 -v $PWD:/serve cwrau/webserver:1.0.0`
