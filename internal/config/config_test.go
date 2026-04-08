package config

import (
	"testing"
	"time"
)

func TestGlobalConfig_Initialized(t *testing.T) {
	if GlobalConfig == nil {
		t.Fatal("GlobalConfig should not be nil")
	}
	if GlobalConfig.ConcurrencyLimit != 10 {
		t.Errorf("ConcurrencyLimit = %v, want 10", GlobalConfig.ConcurrencyLimit)
	}
	if GlobalConfig.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", GlobalConfig.Timeout)
	}
	if GlobalConfig.MaxFileSize != 100*1024*1024 {
		t.Errorf("MaxFileSize = %v, want 100MB", GlobalConfig.MaxFileSize)
	}
}

func TestInit_Configuration(t *testing.T) {
	// Reset to default before test
	GlobalConfig = &Config{
		ConcurrencyLimit: 10,
		Timeout:          30 * time.Second,
		MaxFileSize:      100 * 1024 * 1024,
	}

	Init(Config{
		ConcurrencyLimit: 20,
		Timeout:          60 * time.Second,
		MaxFileSize:      200 * 1024 * 1024,
	})

	if GlobalConfig.ConcurrencyLimit != 20 {
		t.Errorf("ConcurrencyLimit = %v, want 20", GlobalConfig.ConcurrencyLimit)
	}
	if GlobalConfig.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want 60s", GlobalConfig.Timeout)
	}
	if GlobalConfig.MaxFileSize != 200*1024*1024 {
		t.Errorf("MaxFileSize = %v, want 200MB", GlobalConfig.MaxFileSize)
	}
}

func TestInit_Singleton(t *testing.T) {
	// Reset to default before test
	GlobalConfig = &Config{
		ConcurrencyLimit: 10,
		Timeout:          30 * time.Second,
		MaxFileSize:      100 * 1024 * 1024,
	}
	defer func() {
		// Reset after test
		GlobalConfig = &Config{
			ConcurrencyLimit: 10,
			Timeout:          30 * time.Second,
			MaxFileSize:      100 * 1024 * 1024,
		}
	}()

	config1 := Config{ConcurrencyLimit: 1, Timeout: 100 * time.Millisecond, MaxFileSize: 1024}
	config2 := Config{ConcurrencyLimit: 2, Timeout: 200 * time.Millisecond, MaxFileSize: 2048}

	Init(config1)

	if GlobalConfig.ConcurrencyLimit != 1 {
		t.Errorf("ConcurrencyLimit = %v, want 1", GlobalConfig.ConcurrencyLimit)
	}
	if GlobalConfig.Timeout != 100*time.Millisecond {
		t.Errorf("Timeout = %v, want 100ms", GlobalConfig.Timeout)
	}

	// Second call should update the config (not singleton behavior)
	Init(config2)

	if GlobalConfig.ConcurrencyLimit != 2 {
		t.Errorf("ConcurrencyLimit = %v, want 2", GlobalConfig.ConcurrencyLimit)
	}
	if GlobalConfig.Timeout != 200*time.Millisecond {
		t.Errorf("Timeout = %v, want 200ms", GlobalConfig.Timeout)
	}
}
