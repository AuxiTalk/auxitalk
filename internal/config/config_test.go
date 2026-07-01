package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultValidate(t *testing.T) {
	cfg := Default()

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected default config to be valid: %v", err)
	}

	if cfg.Mode != ModeDev {
		t.Fatalf("expected dev mode, got %q", cfg.Mode)
	}
	if cfg.Runtime.RequestTimeout.Std() != 10*time.Second {
		t.Fatalf("expected 10s request timeout, got %s", cfg.Runtime.RequestTimeout.Std())
	}
}

func TestLoadEmptyPathUsesDefault(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("expected load default config: %v", err)
	}
	if cfg.Mode != ModeDev {
		t.Fatalf("expected dev mode, got %q", cfg.Mode)
	}
}

func TestLoadConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auxitalk.json")
	data := []byte(`{
		"mode": "local",
		"runtime": {
			"requestTimeout": "5s",
			"healthTimeout": "1s",
			"maxPayloadSize": 2048,
			"maxEventsPerSecond": 10
		},
		"plugins": []
	}`)

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected config to load: %v", err)
	}

	if cfg.Mode != ModeLocal {
		t.Fatalf("expected local mode, got %q", cfg.Mode)
	}
	if cfg.Runtime.MaxPayloadSize != 2048 {
		t.Fatalf("expected max payload 2048, got %d", cfg.Runtime.MaxPayloadSize)
	}
}

func TestInvalidMode(t *testing.T) {
	cfg := Default()
	cfg.Mode = "unsafe"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid mode error")
	}
}

func TestInvalidDuration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auxitalk.json")
	data := []byte(`{
		"mode": "dev",
		"runtime": {
			"requestTimeout": "soon",
			"healthTimeout": "1s",
			"maxPayloadSize": 2048,
			"maxEventsPerSecond": 10
		},
		"plugins": []
	}`)

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("expected invalid duration error")
	}
}

func TestInlinePluginValidation(t *testing.T) {
	cfg := Default()
	cfg.Plugins = []Plugin{{Enabled: true}}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected missing plugin manifest error")
	}
}

func TestPluginResolvedEnv(t *testing.T) {
	plugin := Plugin{
		Manifest: "plugin.json",
		Enabled:  true,
		Env: map[string]string{
			"PLAIN":  "value",
			"SECRET": "${TOKEN}",
			"MIXED":  "prefix-${TOKEN}",
			"EMPTY":  "${MISSING}",
		},
	}

	resolved := plugin.ResolvedEnv(func(key string) (string, bool) {
		if key == "TOKEN" {
			return "abc", true
		}
		return "", false
	})
	values := map[string]string{}
	for _, item := range resolved {
		parts := strings.SplitN(item, "=", 2)
		values[parts[0]] = parts[1]
	}

	if values["PLAIN"] != "value" || values["SECRET"] != "abc" || values["MIXED"] != "prefix-abc" || values["EMPTY"] != "" {
		t.Fatalf("unexpected env: %+v", values)
	}
}
