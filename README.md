# noport

Easily setup subdomains for you local servers.

The UI configures a local nginx server listening on ports 80 and 433 to proxy
traffic to the ports where you run your projects, mapping them by subdomain.
E.g.:

* http://foo.localhost -> http://localhost:8080
* https://bar.localhost -> http://localhost:1313


## Configuration

Your configuration is saved to `~/.noport.json` by the server.


## Setup

This setup assumes Arch Linux. Adapt it for your OS if needed.

Install [nginx](https://wiki.archlinux.org/title/Nginx).
```
# pacman -Syu nginx
# systemctl start nginx   # start nginx running on port 80 (http://localhost)
# systemctl enable nginx  # start nginx on boot
# systemctl status nginx  # check the service's status
```

`noport` works by running a server that generates and changes `/etc/nginx/nginx.conf`.
For that reason, it needs permission to write that file. One way to arrange that
is to create a group called `noport` for that purpose.

After installing nginx, usually:
```
$ ls -l /etc/nginx/nginx.conf
-rw-r--r-- 1 root root ... /etc/nginx/nginx.conf
```

Set the owner group to `noport` and grant write access by members of the group:
```
# chown root:noport /etc/nginx/nginx.conf
# chmod g+w /etc/nginx/nginx.conf
-rw-rw-r-- 1 root noport ... /etc/nginx/nginx.conf
```

Add yourself to the `noport` group. You will need to log out and log back in,
or restart the machine for the change to take effect. Check running `groups`.
```
# gpasswd -a spelufo noport
```

You should now be able to write to the file.
```
$ /etc/nginx/nginx.conf
```

## Install and run

`noport` is a self-contained static binary. Download it and put it somewhere in
your PATH to install it. Run `noport` without arguments to start the server.

Open http://localhost:8765 for the web interface.


## Development

The entrypoint for running development tasks is `./dev.sh`.


## TODO

* [ ] Validate no duplicate ports, no duplicate domains.
* [ ] Reorder items.
* [ ] Run only noport's server, not figwheel's. Get rid of allow CORS.
* [ ] ./dev.sh build_release
* [ ] SSL certificates with mkcert
* [ ] Installer script.
* [ ] AUR package.
* [ ] Leverage ports.json to suggest servers and avoid collisions.
* [ ] Other OSes.
