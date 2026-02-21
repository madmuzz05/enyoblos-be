package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger untuk inisialisasi zerolog dengan warna
func InitLogger(env string) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
	}

	// Format level dengan warna
	output.FormatLevel = func(i interface{}) string {
		level := strings.ToUpper(fmt.Sprintf("%s", i))
		switch level {
		case "ERROR":
			return "\033[31m" + level + "\033[0m" // merah
		case "WARN":
			return "\033[33m" + level + "\033[0m" // kuning
		case "INFO":
			return "\033[32m" + level + "\033[0m" // hijau
		default:
			return level
		}
	}

	// Format pesan
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// NewLogger middleware Fiber
func NewLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := float64(time.Since(start).Microseconds()) / 1000.0

		statusCode := c.Response().StatusCode()

		// Warna status code
		var statusColor string
		switch {
		case statusCode >= 200 && statusCode < 300:
			statusColor = "\033[32m" // hijau
		case statusCode >= 300 && statusCode < 400:
			statusColor = "\033[36m" // cyan
		case statusCode >= 400 && statusCode < 500:
			statusColor = "\033[33m" // kuning
		default:
			statusColor = "\033[31m" // merah
		}

		msg := "-"
		if c.Response().StatusCode() >= 400 && c.Response().StatusCode() < 500 {
			byteRes := c.Response().Body()
			msg = string(byteRes)
			log.Warn().
				Msgf("%s%d\033[0m | %.2f ms | %s %s | %s | %s | %s",
					statusColor, statusCode, latency, c.Method(), c.OriginalURL(), c.Get("User-Agent"), c.IP(), msg)

		} else if c.Response().StatusCode() >= 500 {
			byteRes := c.Response().Body()
			msg = string(byteRes)
			log.Error().
				Msgf("%s%d\033[0m | %.2f ms | %s %s | %s | %s | %s ",
					statusColor, statusCode, latency, c.Method(), c.OriginalURL(), c.Get("User-Agent"), c.IP(), msg)
		} else {
			log.Info().
				Msgf("%s%d\033[0m | %.2f ms | %s %s | %s | %s | %s",
					statusColor, statusCode, latency, c.Method(), c.OriginalURL(), c.Get("User-Agent"), c.IP(), msg)

		}

		return err
	}
}
