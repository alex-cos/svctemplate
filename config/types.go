package config

import "time"

type WebConfig struct {
	Host          string        `json:"host"`
	HTTPPort      int           `json:"httpPort"`
	ShutdownGrace time.Duration `json:"shutdownGrace"`
	EnableTLS     bool          `json:"enableTLS"`
	TLSCertFile   string        `json:"tlsCertFile"`
	TLSKeyFile    string        `json:"tlsKeyFile"`
	TLSClientAuth string        `json:"tlsClientAuth"` // "none", "require", "verify"
	TLSCAFile     string        `json:"tlsCAFile"`
}

func (c *WebConfig) IsEqualTo(cfg *WebConfig) bool {
	return c.Host == cfg.Host &&
		c.HTTPPort == cfg.HTTPPort &&
		c.ShutdownGrace == cfg.ShutdownGrace &&
		c.EnableTLS == cfg.EnableTLS &&
		c.TLSCertFile == cfg.TLSCertFile &&
		c.TLSKeyFile == cfg.TLSKeyFile &&
		c.TLSClientAuth == cfg.TLSClientAuth &&
		c.TLSCAFile == cfg.TLSCAFile
}

type PprofConfig struct {
	Enable        bool          `json:"enable"`
	Port          int           `json:"port"`
	User          string        `json:"user"`
	Pass          string        `json:"pass"`
	ShutdownGrace time.Duration `json:"shutdownGrace"`
}

func (c *PprofConfig) IsEqualTo(cfg *PprofConfig) bool {
	return c.Enable == cfg.Enable &&
		c.Port == cfg.Port &&
		c.User == cfg.User &&
		c.Pass == cfg.Pass &&
		c.ShutdownGrace == cfg.ShutdownGrace
}

type RateLimiterConfig struct {
	Rate  float64 `json:"rate"`
	Burst int     `json:"burst"`
}

func (c *RateLimiterConfig) IsEqualTo(cfg *RateLimiterConfig) bool {
	return c.Rate == cfg.Rate &&
		c.Burst == cfg.Burst
}

type Config struct {
	LogLevel    string            `json:"logLevel"`
	API         WebConfig         `json:"api"`
	Admin       WebConfig         `json:"admin"`
	Metrics     WebConfig         `json:"metrics"`
	RateLimiter RateLimiterConfig `json:"rateLimiter"`
	Pprof       PprofConfig       `json:"pprof"`
}

func (c *Config) IsEqualTo(cfg *Config) bool {
	return c.API.IsEqualTo(&cfg.API) &&
		c.Admin.IsEqualTo(&cfg.Admin) &&
		c.Metrics.IsEqualTo(&cfg.Metrics) &&
		c.Pprof.IsEqualTo(&cfg.Pprof) &&
		c.RateLimiter.IsEqualTo(&cfg.RateLimiter) &&
		c.LogLevel == cfg.LogLevel
}
