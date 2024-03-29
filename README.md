# Noport

Easily setup localhost subdomains for you local servers.

The UI configures a local nginx server listening on ports 80 and 433 to proxy
traffic to the ports where you run your projects, mapping them by subdomain.

![Noport UI](resources/public/images/noport_ui.png)

The "Save" button saves your configuration to `~/.noport.json`.
The "Install" button saves and then installs an `/etc/nginx/nginx.conf` file
generated from your configuration.


## Setup

* [Linux](#Linux)
* [OSX](#OSX)
* [Windows](#Windows)


### Linux

Install [nginx](https://wiki.archlinux.org/title/Nginx).
```
# pacman -Syu nginx       # or `apt-get install nginx`, etc.
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
# groupadd noport
# chown root:noport /etc/nginx/nginx.conf
# chmod g+w /etc/nginx/nginx.conf
-rw-rw-r-- 1 root noport ... /etc/nginx/nginx.conf
```

Add yourself to the `noport` group. You will need to log out and log back in,
or restart the machine for the change to take effect. Check running `groups`.
```
# gpasswd -a spelufo noport
```

You should now be able to write to the file. Check with:
```
$ touch /etc/nginx/nginx.conf
```

### OSX

Install nginx with brew.
```
brew install nginx
```

The nginx server will run through launchd as a regular user, and the config
is at `/usr/local/etc/nginx/nginx.conf`. You can (re)start it with:
```
brew services restart nginx
```

### Windows

Install nginx for windows. It comes as a zip file that decompreses to a folder
from where you run nginx.exe manually. I've put it at `D:\Programs\nginx`.

Configure noport by setting the following environment variables (e.g. from windows settings):
```
NOPORT_NGINX_EXE=D:\Programs\nginx\nginx.exe
NOPORT_NGINX_CONF=D:\Programs\nginx\conf\nginx.conf
```

Save [nginx.bat](https://raw.githubusercontent.com/spelufo/noport/main/nginx.bat)
and [nginx_reload.bat](https://raw.githubusercontent.com/spelufo/noport/main/nginx_reload.bat)
to nginx's folder.

Setup nginx.bat to run at startup: Press Win+R and enter "shell:startup" which 
will open explorer to the startup folder. Creating a shortcut here causes the
target of the shortcut to run at startup. Right click, new shortcut, and set
the target to nginx.bat.

Restart or run nginx.bat to start nginx.

nginx_reload.bat reloads the configuration after you make changes with noport.
TODO: Make noport trigger the reload, like on osx and linux.


## Running

`noport` is a self-contained static binary. Download it and put it somewhere in
your PATH to install it. Run `noport` without arguments to start the server.

Open http://localhost:8012 for the web interface.

You can keep it running however you plan on running your other services. One
option is to user systemd user services. Or you can just start it manually when
you want to change the configuration. The changes persist in your nginx.conf.


## Development

### Build it

```
./dev.sh build ~/bin/noport
```


### Run in development mode

Start the server.

```
./dev.sh server
```

Run the frontend figwheel server, that builds and hotloads cljs and css code.

```
./dev.sh front
```


### Missing features

This are some features that noport lacks that it would be nice to add.

* [ ] SSL certificates with mkcert
* [ ] Reorder items.
* [ ] Installer script.
* [ ] AUR package.
* [ ] Leverage ports.json to suggest servers and avoid collisions.
* [ ] Windows, reload automatically.


### TODO

* [ ] Fix gin warnings. Run in gin production mode.
