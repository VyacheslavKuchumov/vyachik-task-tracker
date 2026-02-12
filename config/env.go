package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/lpernett/godotenv"
)

var Envs = initConfig()

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

func initConfig() Config {
	// Load .env file from project root
	loadEnvFromProjectRoot()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", ":8000"),
		DBUser:                 getEnv("DB_USER", "admin"),
		DBPassword:             getEnv("DB_PASSWORD", "admin"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "dev-db"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "CHANGE_ME"),
	}
}

func loadEnvFromProjectRoot() {
	// Try multiple approaches to find project root

	// Method 1: Use runtime.Caller to get this file's location
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		// This file is at: /path/to/project/config/config.go
		// So project root is one directory up
		projectRoot := filepath.Join(filepath.Dir(filename), "..")
		envPath := filepath.Join(projectRoot, ".env")

		if err := godotenv.Load(envPath); err == nil {
			log.Printf("Loaded .env from: %s", envPath)
			return
		}
	}

	// Method 2: Look for go.mod starting from current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting working directory: %v", err)
	} else {
		projectRoot := findGoModRoot(cwd)
		if projectRoot != "" {
			envPath := filepath.Join(projectRoot, ".env")
			if err := godotenv.Load(envPath); err == nil {
				log.Printf("Loaded .env from: %s", envPath)
				return
			}
		}
	}

	// Method 3: Try parent directories (common in test scenarios)
	for i := 0; i < 5; i++ {
		parentPath := filepath.Join(".." + string(filepath.Separator) + "..")
		envPath := filepath.Join(parentPath, ".env")
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("Loaded .env from: %s", envPath)
			return
		}
	}

	// Method 4: Try current directory as last resort
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file from any location")
	}
}

func findGoModRoot(dir string) string {
	// Check if go.mod exists in current directory
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return dir
	}

	// Move to parent directory
	parent := filepath.Dir(dir)
	if parent == dir {
		// Reached filesystem root
		return ""
	}

	// Recursively check parent
	return findGoModRoot(parent)
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)

	if ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	value, ok := os.LookupEnv(key)

	if ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
