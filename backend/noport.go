package noport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"text/template"
)

func Main() {
	if len(os.Args) <= 1 {
		loadUserConfig()
		runServer()
	} else {
		switch os.Args[1] {
		case "server", "serve":
			loadUserConfig()
			runServer()
		case "conf":
			loadUserConfig()
			writeNginxConfTo(os.Stdout)
		case "install":
			install()
		default:
			fmt.Fprintf(os.Stderr, "unrecognized command")
			os.Exit(1)
		}
	}
}

type Config struct {
	Servers []Server `json:"servers" binding:"required"`
}

type Server struct {
	Domain string `json:"domain" binding:"required"`
	Port   int    `json:"port" binding:"required"`
	SSL    bool   `json:"ssl"`
}
 
const (
	// The filename of the config file (~/.noport.json).
	theConfigFileName = ".noport.json"
)

var (
	// The current config.
	theConfig Config

	// The nginx.conf template.
	theTemplate *template.Template

	theNginxConfPath = "/etc/nginx/nginx.conf"
	theNginxRestartCommand = []string{"systemctl", "restart", "nginx"}
)

func init() {
	if runtime.GOOS == "darwin" {
		theNginxConfPath = "/usr/local/etc/nginx/nginx.conf"
		theNginxRestartCommand = []string{"brew", "services", "restart", "nginx"}
	}

	nginxConf := os.Getenv("NOPORT_NGINX_CONF")
	if nginxConf != "" {
		theNginxConfPath = nginxConf
	}

	nginxExe := os.Getenv("NOPORT_NGINX_EXE")
	if nginxExe != "" {
		theNginxRestartCommand = []string{nginxExe, "-s", "restart"}
	}

	theTemplate = template.Must(template.New("nginx.conf").Parse(`
user http;
worker_processes  1;

events {
    worker_connections  1024;
}

http {
    charset utf-8;
    include mime.types;
    default_type application/octet-stream;
    sendfile on;
    keepalive_timeout 65;
{{define "server"}}
    server {
        listen {{if .SSL}}443{{else}}80{{end}};
        server_name {{.Domain}}.localhost;

        {{- if .SSL}}
        ssl_certificate      cert.pem;
        ssl_certificate_key  cert.key;
        ssl_session_cache    shared:SSL:1m;
        ssl_session_timeout  5m;
        ssl_ciphers  HIGH:!aNULL:!MD5;
        ssl_prefer_server_ciphers  on;{{end}}

        location / {
            proxy_set_header   X-Forwarded-For $remote_addr;
            proxy_set_header   Host $http_host;
            proxy_pass         "http{{if .SSL}}s{{end}}://127.0.0.1:{{.Port}}";
        }
    }
{{end}}
{{range .Servers -}}{{template "server" .}}{{end}}
}
`))
}

func install() {
	confPath := theNginxConfPath
	if len(os.Args) > 2 {
		confPath = os.Args[2]
	}
	installNginxConf(confPath)
}

func installNginxConf(confPath string) error {
	err := writeNginxConf(confPath)
	if err != nil {
		return err
	}
	cmd := exec.Command(theNginxRestartCommand[0], theNginxRestartCommand[1:]...)
	return cmd.Run()
}

func writeNginxConf(confPath string) error {
	b := bytes.Buffer{}
	err := writeNginxConfTo(&b)
	if err != nil {
		return err
	}
	return os.WriteFile(confPath, b.Bytes(), 0644)
}

func writeNginxConfTo(w io.Writer) error {
	return theTemplate.Execute(w, theConfig)
}

func loadUserConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return loadConfig(path.Join(home, theConfigFileName))
}

func saveUserConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return saveConfig(path.Join(home, theConfigFileName))
}

func loadConfig(configPath string) error {
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(configBytes, &theConfig)
}

func saveConfig(configPath string) error {
	file, err := json.MarshalIndent(theConfig, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, file, 0644)
}
