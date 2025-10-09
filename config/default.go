package config

import (
	"time"

	"github.com/spf13/viper"
)

func SetDefault(v *viper.Viper) {
	v.SetDefault("LogLevel", "info")
	v.SetDefault("ShutdownGrace", 5*time.Second)

	v.SetDefault("APIConfig.Host", "")
	v.SetDefault("APIConfig.HTTPPort", 8080)
	v.SetDefault("APIConfig.EnableTLS", false)
	v.SetDefault("APIConfig.TLSCertFile", "")
	v.SetDefault("APIConfig.TLSKeyFile", "")
	v.SetDefault("APIConfig.TLSClientAuth", "")
	v.SetDefault("APIConfig.TLSCAFile", "")

	v.SetDefault("AdminConfig.Host", "")
	v.SetDefault("AdminConfig.HTTPPort", 8070)
	v.SetDefault("AdminConfig.EnableTLS", false)
	v.SetDefault("AdminConfig.TLSCertFile", "")
	v.SetDefault("AdminConfig.TLSKeyFile", "")
	v.SetDefault("AdminConfig.TLSClientAuth", "")
	v.SetDefault("AdminConfig.TLSCAFile", "")

	v.SetDefault("MetricsConfig.Host", "")
	v.SetDefault("MetricsConfig.HTTPPort", 8090)
	v.SetDefault("MetricsConfig.EnableTLS", false)
	v.SetDefault("MetricsConfig.TLSCertFile", "")
	v.SetDefault("MetricsConfig.TLSKeyFile", "")
	v.SetDefault("MetricsConfig.TLSClientAuth", "")
	v.SetDefault("MetricsConfig.TLSCAFile", "")

	v.SetDefault("PprofConfig.Enable", false)
	v.SetDefault("PprofConfig.Port", 8060)
	v.SetDefault("PprofConfig.User", "")
	v.SetDefault("PprofConfig.Pass", "")
}
