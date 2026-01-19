package test

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func PluginTestDataDir(skip int) string {
	_, file, _, _ := runtime.Caller(skip)
	return filepath.Join(filepath.Dir(file), "testdata")
}

func BuildTestPlugins(t *testing.T) {
	t.Helper()
	t.Logf("[BeforeSuite] Current working directory: %s", PluginTestDataDir(2))
	cmd := exec.Command("make", "-C", PluginTestDataDir(2))
	out, err := cmd.CombinedOutput()
	t.Logf("[BeforeSuite] Make output: %s", string(out))

	if err != nil {
		t.Fatalf("Failed to build test plugins: %v", err)
	}
}
