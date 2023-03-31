package noport

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"text/template"
)

// The filename of the config file (~/.noport.json).
const theConfigFileName = ".noport.json"

// The path to the config file.
var theConfigFilePath string

// The current config.
var theConfig Config

// The nginx.conf template.
var nginxConfTemplate *template.Template

type Config struct {
	Servers []Server `json:"servers" binding:"required"`
}

type Server struct {
	Domain string `json:"domain" binding:"required"`
	Port   int    `json:"port" binding:"required"`
	SSL    bool   `json:"ssl"`
}

func init() {
	nginxConfTemplate = template.Must(template.New("nginx.conf").Parse(`
user http;
worker_processes  1;

events {
    worker_connections  1024;
}

http {
    charset utf-8;

    include       mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    #tcp_nopush     on;

    #keepalive_timeout  0;
    keepalive_timeout  65;

    #gzip  on;

		{{define "server"}}
    server {
        listen {{if .SSL}}443{{else}}80{{end}};
        server_name {{.Domain}}.localhost;

        {{if .SSL -}}
        ssl_certificate      cert.pem;
        ssl_certificate_key  cert.key;

        ssl_session_cache    shared:SSL:1m;
        ssl_session_timeout  5m;

        ssl_ciphers  HIGH:!aNULL:!MD5;
        ssl_prefer_server_ciphers  on;
        {{- end}}

        location / {
            proxy_set_header   X-Forwarded-For $remote_addr;
            proxy_set_header   Host $http_host;
            proxy_pass         "http{{if .SSL}}s{{end}}://127.0.0.1:{{.Port}}";
        }
    }
		{{end}}

    {{range .Servers}}
    {{template "server" .}}
    {{end}}
}
`))
}

func Main() {
	loadConfig()
	if len(os.Args) <= 1 {
		runServer()
	} else {
		switch os.Args[1] {
		case "server", "serve":
			runServer()
		case "conf":
			nginxConfTemplate.Execute(os.Stdout, theConfig)
		case "install":
			install()
		default:
			fmt.Fprintf(os.Stderr, "unrecognized command")
			os.Exit(1)
		}
	}
}

func saveConfig(config Config) error {
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(theConfigFilePath, file, 0644)
}

func loadConfig() {
	home, err := os.UserHomeDir()
	die(err)

	theConfigFilePath = path.Join(home, theConfigFileName)
	configBytes, err := os.ReadFile(theConfigFilePath)
	if errors.Is(err, os.ErrNotExist) {
		err = saveConfig(Config{})
		if err == nil {
			fmt.Println("Created ~/" + theConfigFileName)
		}
	}
	if err != nil {
		die(err)
	}
	config := Config{}
	json.Unmarshal(configBytes, &config)
}

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s", err)
		os.Exit(1)
	}
}
