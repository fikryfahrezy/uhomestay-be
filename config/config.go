package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	JwtKey          []byte
	CloudinaryUrl   string
	SentryDsn       string
	LogDnaKey       string
	Port            string
	Argon2Salt      string
	JwtAudiencesStr string
	JwtKeyStr       string
	JwtIssuerUrl    string
	PostgreUrl      string
	MongoUri        string
	RedisUrl        string
	Env             string
	JwtAudiences    []string
}

func LoadConfig() Config {
	var c Config

	cd := os.Getenv("HOMESTAY_CLOUDINARY_URL")
	if cd == "" {
		log.Fatal("$HOMESTAY_CLOUDINARY_URL must be set")
	}
	c.CloudinaryUrl = cd

	sentryDsn := os.Getenv("HOMESTAY_SENTRY_DSN")
	if sentryDsn == "" {
		log.Fatal("$HOMESTAY_SENTRY_DSN must be set")
	}
	c.SentryDsn = sentryDsn

	logDnaKey := os.Getenv("HOMESTAY_LOGDNA_KEY")
	if logDnaKey == "" {
		log.Fatal("$HOMESTAY_LOGDNA_KEY must be set")
	}
	c.LogDnaKey = logDnaKey

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	c.Port = port

	argon2Salt := os.Getenv("HOMESTAY_ARG_SALT")
	if argon2Salt == "" {
		log.Fatal("$HOMESTAY_ARG_SALT must be set")
	}
	c.Argon2Salt = argon2Salt

	jwtAudiencesStr := os.Getenv("HOMESTAY_JWT_AUDIENCES")
	if jwtAudiencesStr == "" {
		log.Fatal("$HOMESTAY_JWT_AUDIENCES must be set")
	}
	c.JwtAudiencesStr = jwtAudiencesStr
	c.JwtAudiences = strings.Split(jwtAudiencesStr, ",")

	jwtKeyStr := os.Getenv("HOMESTAY_JWT_SECRET")
	if jwtKeyStr == "" {
		log.Fatal("$HOMESTAY_JWT_SECRET must be set")
	}
	c.JwtKeyStr = jwtKeyStr
	c.JwtKey = []byte(jwtKeyStr)

	jwtIssuerUrl := os.Getenv("HOMESTAY_JWT_ISSUER")
	if jwtKeyStr == "" {
		log.Fatal("$HOMESTAY_JWT_ISSUER must be set")
	}
	c.JwtIssuerUrl = jwtIssuerUrl

	postgreUrl := os.Getenv("DATABASE_URL")
	if postgreUrl == "" {
		postgreUrl = "postgres://postgres:postgres@localhost:5432/homestay"
	}
	c.PostgreUrl = postgreUrl

	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://mongo:mongo@localhost:27017"
	}
	c.MongoUri = mongoUri

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "redis://localhost:6379"
	}
	c.RedisUrl = redisUrl

	env := os.Getenv("ENVI")
	if env == "" {
		env = "dev"
	}
	c.Env = env

	return c
}
