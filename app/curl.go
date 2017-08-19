package app

import (
	"os/exec"
	"time"

	"golang.org/x/net/context"
)

func GetWebPage(url string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output := make(chan string)
	errch := make(chan error)
	var cmd *exec.Cmd
	go func() {
		cmd = exec.CommandContext(ctx, "curl", "-s", "-S", "-L", "-c", "/dev/null", "-m", "3", url)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			errch <- err
			return
		}
		output <- string(stdout)
	}()

	select {
	case result := <-output:
		return result
	case <-ctx.Done():
		cmd.Process.Kill()
		return ""
	case <-errch:
		return ""
	}
}
