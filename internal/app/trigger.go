package app

import (
	"os"
	"os/exec"
)

// TriggerReconcile spawns a one-shot background reconcile process after a mutation.
// It is non-blocking: the subprocess runs asynchronously without waiting for its
// completion, so the CLI command returns immediately to the user.
//
// This avoids the need to run `synctl daemon` manually to apply state changes to
// the runtime. `synctl daemon` remains the tool for real-time monitoring only.
func TriggerReconcile() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	cmd := exec.Command(exe, "reconcile")
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Start() // fire and forget — we do not call cmd.Wait()
}
