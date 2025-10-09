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

type Config struct {
	LogLevel      string      `json:"logLevel"`
	APIConfig     WebConfig   `json:"apiConfig"`
	AdminConfig   WebConfig   `json:"adminConfig"`
	MetricsConfig WebConfig   `json:"metricsConfig"`
	PprofConfig   PprofConfig `json:"pprofConfig"`
}

func (c *Config) IsEqualTo(cfg *Config) bool {
	return c.APIConfig.IsEqualTo(&cfg.APIConfig) &&
		c.AdminConfig.IsEqualTo(&cfg.AdminConfig) &&
		c.MetricsConfig.IsEqualTo(&cfg.MetricsConfig) &&
		c.PprofConfig.IsEqualTo(&cfg.PprofConfig) &&
		c.LogLevel == cfg.LogLevel
}
