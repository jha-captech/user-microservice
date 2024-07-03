package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

type configuration struct {
	Env      string `env:"ENV"`
	Database struct {
		Name            string `env:"DATABASE_NAME"`
		User            string `env:"DATABASE_USER"`
		Password        string `env:"DATABASE_PASSWORD"`
		Host            string `env:"DATABASE_HOST"`
		Port            string `env:"DATABASE_PORT"`
		ConnectionRetry int    `env:"DATABASE_CONNECTION_RETRY"`
	}
	HTTP struct {
		Domain              string `env:"HTTP_DOMAIN"`
		Port                string `env:"HTTP_PORT"`
		ShutdownGracePeriod int    `env:"HTTP_SHUTDOWN_GRACE_PERIOD"`
	}
}

// mustNewConfiguration returns a configuration struct from environment variables or panics if this fails.
func mustNewConfiguration() configuration {
	loadEnv()

	config := configuration{}
	if err := ParseStructFromEnv(&config, true); err != nil {
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
func ParseStructFromEnv(obj any, errOnMissingValue bool) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("in ParseStructFromEnv: %w", err)
		}
	}()
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
			if err := ParseStructFromEnv(field.Addr().Interface(), errOnMissingValue); err != nil {
				return err
			}
			continue
		}

		// Get and then set env value based on tag if present
		fieldType := val.Type().Field(i)
		envTag := fieldType.Tag.Get("env")

		if field.CanSet() && envTag != "" {
			switch field.Kind() {
			case reflect.String:
				value, err := getEnvString(envTag, errOnMissingValue)
				if err != nil {
					return err
				}
				field.SetString(value)
			case reflect.Int:
				value, err := getEnvInt64(envTag, errOnMissingValue)
				if err != nil {
					return err
				}
				field.SetInt(value)
			case reflect.Bool:
				value, err := getEnvBool(envTag, errOnMissingValue)
				if err != nil {
					return err
				}
				field.SetBool(value)
			default:
				continue
			}
		}
	}
	return nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────────────────────────

func newEnvVarMissingErr[T any](key string) (T, error) {
	var blank T
	errMsg := fmt.Sprintf("enviroment variable '%s' is missing or blank", key)
	return blank, errors.New(errMsg)
}

func newEnvVarParsingErr[T any](key string, err error) (T, error) {
	var blank T
	errMsg := fmt.Sprintf(
		"error parsing enviroment variable '%s' to type '%T': %v",
		key,
		blank,
		err,
	)
	return blank, errors.New(errMsg)
}

func getEnvString(key string, errIfMissing bool) (string, error) {
	value := os.Getenv(key)
	if errIfMissing && value == "" {
		return newEnvVarMissingErr[string](key)
	}
	return value, nil
}

func getEnvInt64(key string, errIfMissing bool) (int64, error) {
	value := os.Getenv(key)
	if errIfMissing && value == "" {
		return newEnvVarMissingErr[int64](key)
	}
	convertedInt, err := strconv.Atoi(value)
	if err != nil {
		return newEnvVarParsingErr[int64](key, err)
	}
	return int64(convertedInt), nil
}

func getEnvBool(key string, errIfMissing bool) (bool, error) {
	value := os.Getenv(key)
	if errIfMissing && value == "" {
		return newEnvVarMissingErr[bool](key)
	}
	convertedBool, err := strconv.ParseBool(value)
	if err != nil {
		return newEnvVarParsingErr[bool](key, err)
	}
	return convertedBool, nil
}

func getEnvFloat64(key string, errIfMissing bool) (float64, error) {
	value := os.Getenv(key)
	if errIfMissing && value == "" {
		return newEnvVarMissingErr[float64](key)
	}
	convertedFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return newEnvVarParsingErr[float64](key, err)
	}
	return convertedFloat, nil
}
