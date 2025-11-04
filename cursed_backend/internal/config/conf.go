package config

type Config struct {
	Port        string `env:"PORT" envDefault:"8080"`
	Env         string `env:"ENV" envDefault:"dev"`
	DBHost      string `env:"DB_HOST" envDefault:"localhost"`
	DBUsername  string `env:"DB_USERNAME" envDefault:"postgres"`
	DBPassword  string `env:"DB_PASSWORD" envDefault:"qazaq001"`
	DBName      string `env:"DB_NAME" envDefault:"petStore"`
	DBPort      string `env:"DB_PORT" envDefault:"5432"`
	SSLMode     string `env:"SSL_MODE" envDefault:"disable"`
	JWTSecret   string `env:"JWT_SECRET"`
	CSRFKey     string `env:"CSRF_KEY"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	CORSOrigins string `env:"CORS_ORIGINS" envDefault:"http://localhost:3000,http://localhost:5173"`
}
