package noport

type Config struct {
	Servers []Server `json:"servers" binding:"required"`
}

type Server struct {
	Domain string `json:"domain" binding:"required"`
	Port   int    `json:"port" binding:"required"`
}
