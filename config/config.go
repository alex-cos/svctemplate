package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ConfigMgmt struct {
	v        *viper.Viper
	cfg      *Config
	log      *slog.Logger
	reloadCh chan struct{}
	mutex    sync.RWMutex
}

func New(log *slog.Logger) *ConfigMgmt {
	return &ConfigMgmt{
		v:        viper.New(),
		cfg:      &Config{},
		log:      log,
		reloadCh: make(chan struct{}, 1),
		mutex:    sync.RWMutex{},
	}
}

func (cm *ConfigMgmt) Close() {
	close(cm.reloadCh)
}

func (cm *ConfigMgmt) GetMode() string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	return cm.v.GetString("mode")
}

func (cm *ConfigMgmt) GetConfig() *Config {
	var cfg Config

	cm.mutex.RLock()
	cfg = *cm.cfg
	cm.mutex.RUnlock()

	return &cfg
}

func (cm *ConfigMgmt) SetDefault() {
	SetDefault(cm.v)
}

func (cm *ConfigMgmt) Load(configFile string) (*Config, error) {
	var cfg Config

	if configFile != "" {
		cm.v.SetConfigFile(configFile)
	} else {
		cm.v.SetConfigName("config")
		cm.v.SetConfigType("yaml")
		cm.v.AddConfigPath(".")
		cm.v.AddConfigPath("./conf")
		cm.v.AddConfigPath("../conf")
	}
	cm.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cm.v.SetEnvPrefix("MYAPP")
	err := cm.v.BindEnv("mode")
	if err != nil {
		return nil, fmt.Errorf("failed to bind env: %w", err)
	}
	cm.v.AutomaticEnv()
	cm.SetDefault()

	err = cm.v.ReadInConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "no config file found, using defaults + env")
	}

	err = cm.v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	cm.mutex.Lock()
	*cm.cfg = cfg
	cm.mutex.Unlock()

	return &cfg, nil
}

func (cm *ConfigMgmt) Watch(onChange func(oldCfg, newCfg *Config)) {
	// watch config file and trigger reload
	cm.v.WatchConfig()
	cm.v.OnConfigChange(func(e fsnotify.Event) {
		cm.logInfo("config file changed", slog.String("file", e.Name))
		cm.NotifyChange()
	})

	// reload worker: handles changes
	go cm.reloadConfig(onChange)
}

func (cm *ConfigMgmt) NotifyChange() bool {
	select {
	case cm.reloadCh <- struct{}{}:
		return true
	default:
		return false
	}
}

// reloadConfig listens on reloadCh and applies configuration changes.
// Structural changes (ports, pprof enable/auth, network mode) trigger a full restart of servers.
func (cm *ConfigMgmt) reloadConfig(onChange func(oldCfg, newCfg *Config)) {
	var newCfg Config

	cm.mutex.RLock()
	oldCfg := *cm.cfg
	cm.mutex.RUnlock()

	for range cm.reloadCh {
		if err := cm.v.Unmarshal(&newCfg); err != nil {
			cm.logError("failed to unmarshal configuration", slog.Any("error", err))
			continue
		}
		cm.mutex.Lock()
		*cm.cfg = newCfg
		cm.mutex.Unlock()
		// detect structural changes
		if !cm.cfg.IsEqualTo(&oldCfg) && onChange != nil {
			onChange(&oldCfg, &newCfg)
		}

		cm.mutex.RLock()
		oldCfg = *cm.cfg
		cm.logInfo("configuration reloaded", slog.Any("config", cm.cfg))
		cm.mutex.RUnlock()
	}
}

func (cm *ConfigMgmt) logInfo(msg string, args ...any) {
	if cm.log != nil {
		cm.log.Info(msg, args...)
	}
}

func (cm *ConfigMgmt) logError(msg string, args ...any) {
	if cm.log != nil {
		cm.log.Error(msg, args...)
	}
}
