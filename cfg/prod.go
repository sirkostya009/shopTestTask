//go:build prod

package cfg

import "os"

var (
	Username    = os.Getenv("USERNAME")
	Password    = os.Getenv("PASSWORD")
	Address     = os.Getenv("ADDR")
	DatabaseURL = os.Getenv("DATABASE_URL")
)

func init() {
	if Username == "" {
		panic("USERNAME env var is not set")
	}
	if Password == "" {
		panic("PASSWORD env var is not set")
	}
	if Address == "" {
		panic("ADDR env var is not set")
	}
	if DatabaseURL == "" {
		panic("DATABASE_URL env var is not set")
	}
}
