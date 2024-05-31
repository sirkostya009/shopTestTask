//go:build !prod

package cfg

const (
	Username    = "admin"
	Password    = "admin"
	Address     = ":8080"
	DatabaseURL = "postgres://postgres:postgres@localhost:5432/postgres"
)
