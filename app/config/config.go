package config

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/galrub/go/jobSearch/internal/logger"
)

func DevMode() bool {
	return os.Getenv("ENV") == "development"
}

type PostgresData struct {
	Username string
	Pw       string
	DbName   string
	Host     string
}

func (d PostgresData) CreateDSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s sslmode=disable", d.Username, d.Pw, d.Host, d.DbName)
}

func clearString(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, text)
}

// DB_URL=postgresql://user:pass@localhost:5432/postgres?sslmode=disable
func GetPostgresData() PostgresData {
	filename := os.Getenv("POSTGRES_PASSWORD_FILE")
	pw := os.Getenv("POSTGRES_PW")
	if _, err := os.Stat(filename); err != nil {
		logger.LOG.Info().Msg("fialed to find docker secret file, will use env. variables instead")
	} else {
		b, err := os.ReadFile(filename)
		if err != nil {
			logger.LOG.Info().Msg("fialed to read docker secret file, will use env. variables instead")
		} else {
			pw = clearString(string(b))
			logger.LOG.Info().Msg("reading docker secret from file")
		}
	}
	return PostgresData{
		Username: os.Getenv("POSTGRES_USER"),
		Pw:       pw,
		Host:     os.Getenv("POSTGRES_HOST"),
		DbName:   os.Getenv("POSTGRES_DB"),
	}

}
