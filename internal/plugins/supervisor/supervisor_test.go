package supervisor

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSupervisorRestartsCrashedProcess(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "crash.sh")
	counter := filepath.Join(dir, "counter")
	content := `#!/usr/bin/env sh
count=0
if [ -f "` + counter + `" ]; then
  count=$(cat "` + counter + `")
fi
count=$((count + 1))
printf '%s' "$count" > "` + counter + `"
exit 1
`
	if err := os.WriteFile(script, []byte(content), 0o700); err != nil {
		t.Fatalf("write script: %v", err)
	}

	sup := NewSupervisor(ProcessOptions{CallTimeout: 50 * time.Millisecond, HealthInterval: time.Hour, RestartBackoff: 10 * time.Millisecond, MaxRestarts: 2})
	if err := sup.Register(ProcessSpec{ID: "crashy", Command: script}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := sup.Start(context.Background(), "crashy"); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer sup.Stop("crashy")

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		status, ok := sup.Status("crashy")
		if ok && status.Restarts >= 1 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	status, _ := sup.Status("crashy")
	t.Fatalf("expected restart after crash, got status=%+v", status)
}

func TestSupervisorCallTimesOutWhenPluginDoesNotRespond(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "silent.sh")
	content := `#!/usr/bin/env sh
while IFS= read -r line; do
  sleep 1
done
`
	if err := os.WriteFile(script, []byte(content), 0o700); err != nil {
		t.Fatalf("write script: %v", err)
	}

	sup := NewSupervisor(ProcessOptions{CallTimeout: 20 * time.Millisecond, HealthInterval: time.Hour, RestartBackoff: time.Millisecond})
	if err := sup.Register(ProcessSpec{ID: "silent", Command: script}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := sup.Start(context.Background(), "silent"); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer sup.Stop("silent")

	_, err := sup.Call(context.Background(), "silent", "plugin.health", nil)
	if err == nil {
		t.Fatal("expected call timeout")
	}
}

func TestSupervisorKillsAndRestartsAfterHealthFailures(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "bad-health.sh")
	content := `#!/usr/bin/env sh
while IFS= read -r line; do
  id=$(printf '%s' "$line" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
  printf '{"jsonrpc":"2.0","id":"%s","error":{"code":-32603,"message":"unhealthy"}}\n' "$id"
done
`
	if err := os.WriteFile(script, []byte(content), 0o700); err != nil {
		t.Fatalf("write script: %v", err)
	}

	sup := NewSupervisor(ProcessOptions{CallTimeout: 50 * time.Millisecond, HealthInterval: 10 * time.Millisecond, RestartBackoff: 10 * time.Millisecond, MaxRestarts: 2, MaxHealthFailures: 2})
	if err := sup.Register(ProcessSpec{ID: "bad-health", Command: script}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := sup.Start(context.Background(), "bad-health"); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer sup.Stop("bad-health")

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		status, ok := sup.Status("bad-health")
		if ok && status.Restarts >= 1 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	status, _ := sup.Status("bad-health")
	t.Fatalf("expected restart after health failures, got status=%+v", status)
}

func TestSupervisorStartsAndCallsHealth(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "plugin.sh")
	content := `#!/usr/bin/env sh
while IFS= read -r line; do
  id=$(printf '%s' "$line" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
  method=$(printf '%s' "$line" | sed -n 's/.*"method":"\([^"]*\)".*/\1/p')
  if [ "$method" = "plugin.health" ]; then
    printf '{"jsonrpc":"2.0","id":"%s","result":{"ok":true}}\n' "$id"
  else
    printf '{"jsonrpc":"2.0","id":"%s","result":{"ok":true}}\n' "$id"
  fi
done
`
	if err := os.WriteFile(script, []byte(content), 0o700); err != nil {
		t.Fatalf("write script: %v", err)
	}

	sup := NewSupervisor(ProcessOptions{CallTimeout: time.Second, HealthInterval: time.Hour})
	if err := sup.Register(ProcessSpec{ID: "fake", Command: script}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := sup.Start(context.Background(), "fake"); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer sup.Stop("fake")

	result, err := sup.Call(context.Background(), "fake", "plugin.health", nil)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	var parsed struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !parsed.OK {
		t.Fatal("expected ok")
	}
}
