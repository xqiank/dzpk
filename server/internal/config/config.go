package config

import (
	"os"
	"strconv"
	"strings"
)

// Config 运行时配置。默认值按本地靶场部署填好，环境变量可覆盖。
type Config struct {
	Addr        string
	DSN         string
	DBName      string
	JWTSecret   string
	Env         string
	Origins     []string
	AdminOrigins []string
	ClientDist  string
	AdminDist   string
	RobotFill   bool
	RobotMin    int
	RobotMax    int
	MaxConnsPerIP int
	RateLimitPerSec int
	// QA 打开后启用仅限联调的自检指令（qa.*），生产环境保持关闭
	QA bool
}

func env(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

func envInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func Load() *Config {
	dsn := env("U5_DSN", "root:@tcp(127.0.0.1:3306)/u5_poker?parseTime=true&charset=utf8mb4&loc=Local&multiStatements=true")
	return &Config{
		Addr:            env("U5_ADDR", "127.0.0.1:8720"),
		DSN:             dsn,
		DBName:          env("U5_DB", "u5_poker"),
		JWTSecret:       env("U5_SECRET", "u5-poker-cft-sandbox-secret-2026-11-09-rotate-me"),
		Env:             env("U5_ENV", "dev"),
		Origins:         strings.Split(env("U5_ORIGINS", "http://127.0.0.1:8721,http://localhost:8721,http://127.0.0.1:8722,http://localhost:8722,http://127.0.0.1:8720"), ","),
		AdminOrigins:    strings.Split(env("U5_ADMIN_ORIGINS", "http://127.0.0.1:8722,http://localhost:8722"), ","),
		ClientDist:      env("U5_CLIENT_DIST", "client/dist"),
		AdminDist:       env("U5_ADMIN_DIST", "admin/dist"),
		RobotFill:       env("U5_ROBOT_FILL", "1") == "1",
		RobotMin:        envInt("U5_ROBOT_MIN", 3),
		RobotMax:        envInt("U5_ROBOT_MAX", 22),
		MaxConnsPerIP:   envInt("U5_MAX_CONN_IP", 12),
		RateLimitPerSec: envInt("U5_RATE", 40),
		QA:              env("U5_QA", "") == "1",
	}
}
