package app

import "os/exec"

func GetWebPage(url string) string {
	cmd := exec.Command("curl", "-s", "-S", "-L", "-c", "/dev/null", url)
	stdout, _ := cmd.CombinedOutput()
	return string(stdout)
}
