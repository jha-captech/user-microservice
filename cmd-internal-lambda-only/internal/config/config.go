package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Env      string `env:"ENV"`
	Database struct {
		Name            string `env:"DATABASE_NAME,required"`
		User            string `env:"DATABASE_USER,required"`
		Password        string `env:"DATABASE_PASSWORD,required"`
		Host            string `env:"DATABASE_HOST,required"`
		Port            string `env:"DATABASE_PORT,required"`
		ConnectionRetry int    `env:"DATABASE_CONNECTION_RETRY,required"`
	}
}

// MustNewConfiguration returns a Configuration struct from environment variables or panics if this fails.
func MustNewConfiguration() Configuration {
	loadEnv()

	config := Configuration{}

	if err := ParseStructFromEnv(&config); err != nil {
		panic(fmt.Sprintf("Error unmarshaling ENV vars: %v", err))
	}

	return config
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Info("DOTENV: No .env file found or error reading file.")
		return
	}
	slog.Info("DOTENV: .env found.")
}

// ── Env Parser ───────────────────────────────────────────────────────────────────────────────────

// ParseStructFromEnv takes a struct as an input and recursively loops tough all fields on the
// struct. If a field is not another struct and has a `env` tag, the environment variable associated
// with that tag will be retrieved and added to the struct.
//
// If the `errOnMissingValue` flag is set to `true`, any tag that is missing an environment variable
// will result in an error being returned.
func ParseStructFromEnv(obj any) error {
	val := reflect.ValueOf(obj)

	// if pointer, get value
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Iterate through the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the field is a struct
		if field.Kind() == reflect.Struct {
			if err := ParseStructFromEnv(field.Addr().Interface()); err != nil {
				return fmt.Errorf("in ParseStructFromEnv: %w", err)
			}
			continue
		}

		// Get and then set env value based on tag if present
		fieldType := val.Type().Field(i)
		tagsString := fieldType.Tag.Get("env")

		if field.CanSet() && tagsString != "" {
			tags := strings.Split(tagsString, ",")
			key := tags[0]
			required := false
			if len(tags) >= 2 && tags[1] == "required" {
				required = true
			}
			envValue := os.Getenv(key)

			// An empty string will always be an issue if `required` is set.
			if required && envValue == "" {
				return newEnvVarMissingErr(key)
			}

			switch field.Kind() {
			case reflect.String:
				field.SetString(envValue)

			case reflect.Int:
				convertedInt, err := strconv.Atoi(envValue)
				if err != nil && required {
					return newEnvVarParsingErr(key, err)
				}
				field.SetInt(int64(convertedInt))

			case reflect.Bool:
				value, err := strconv.ParseBool(envValue)
				if err != nil && required {
					return newEnvVarParsingErr(key, err)
				}
				field.SetBool(value)

			default:
				continue
			}
		}
	}
	return nil
}

func newEnvVarMissingErr(key string) error {
	errMsg := fmt.Sprintf("enviroment variable '%s' is missing or blank", key)
	return errors.New(errMsg)
}

func newEnvVarParsingErr(key string, err error) error {
	errMsg := fmt.Sprintf("error parsing enviroment variable '%s': %v", key, err)
	return errors.New(errMsg)
}
