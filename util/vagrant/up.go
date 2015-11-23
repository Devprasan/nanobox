//
package vagrant

import "os/exec"
import "runtime"

// Up runs a vagrant up
func Up() error {

	// gain sudo privilages; not handling error here because worst case scenario
	// this fails and just prompts for a password later
	if runtime.GOOS != "windows" {
		exec.Command("sudo", "ls").Run()
	}

	return runInContext(exec.Command("vagrant", "up"))
}
