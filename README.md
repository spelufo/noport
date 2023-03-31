# noport

Easily setup subdomains for you local servers.

The UI configures a local nginx server listening on ports 80 and 433 to proxy
traffic to the ports where you run your projects, mapping them by subdomain.
E.g.:

* http://foo.localhost -> localhost:8080
* https://bar.localhost -> localhost:1313


## Features

* [ ] Config file at ~/.noport.json
* [ ] Installation / bootstrapping process (single binary distributed?)
* [ ] SSL certificates with mkcert
* [ ] Works on Windows
