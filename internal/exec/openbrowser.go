package exec

import (
	"os/exec"
	"runtime"
)

// OpenBrowser opens the given URL in the default system browser.
func OpenBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return
	}

	_ = cmd.Start()
}
