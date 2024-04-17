package config

import (
	"fmt"
	"strconv"
)

type ConfigPostgres struct {
	PostgresHostname string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSslMode  string
}

func (p ConfigPostgres) Dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		p.PostgresUser,
		p.PostgresPassword,
		p.PostgresHostname,
		p.PostgresDB,
		p.PostgresSslMode,
	)
}

type ConfigHttpServer struct {
	ServerListenAddress string
	ServerListenPort    int
}

func (s ConfigHttpServer) Host() string {
	return s.ServerListenAddress + ":" + strconv.Itoa(s.ServerListenPort)
}

type CmdRoot struct {
	ConfigPostgres
	ConfigHttpServer
}

var Root = CmdRoot{
	ConfigPostgres: ConfigPostgres{
		PostgresHostname: "127.0.0.1",
		PostgresUser:     "postgres",
		PostgresPassword: "postgres",
		PostgresDB:       "college_event_website",
		PostgresSslMode:  "disable",
	},
	ConfigHttpServer: ConfigHttpServer{
		ServerListenAddress: "0.0.0.0",
		ServerListenPort:    3000,
	},
}
