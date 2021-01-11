package osystem

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"

	errors "github.com/pkg/errors"
	cmd "nenvoy.com/pkg/utils/cmd"
)

// GetHomeDir - Return the current home directory of the user
func GetHomeDir() (homeDir string) {
	// Get the current user
	userAccount, _ := user.Current()
	// Get Home Directory
	homeDir = userAccount.HomeDir
	return homeDir
}

// PathExists checks to see if the path specified exists
func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil

	} else if os.IsNotExist(err) {
		return false, nil

	} else {
		return false, err
	}
}

// DisableSwapSpace - Disable os swap space, required for kubelet to run
func DisableSwapSpace() (err error) {

	// Run command for non persistent disable of swapspace
	_, stderr, err := cmd.Output("swapoff", "-a")
	if err != nil {
		err = errors.Wrap(err, stderr)
		return err
	}

	// Run command for persistent disable of swapspace on restart
	_, stderr, err = cmd.Output("sed", "-ri", "/\\sswap\\s/s/^#?/#/", "/etc/fstab")
	if err != nil {
		err = errors.Wrap(err, stderr)
		return err
	}

	return nil
}

// ChangeHostname - Update node os's hostname
func ChangeHostname(hostname string) (err error) {

	// Update /etc/hosts
	prevHostname, err := os.Hostname()

	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}

	etcHosts, err := ioutil.ReadFile("/etc/hosts")

	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}

	output := bytes.ReplaceAll(etcHosts, []byte(" "+prevHostname+" "), []byte(" "+hostname+" "))
	output = bytes.ReplaceAll(etcHosts, []byte(" "+prevHostname+"\n"), []byte(" "+hostname+"\n"))

	err = ioutil.WriteFile("/etc/hosts", output, 0777)

	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}

	// Run command for non persistent hostname change
	_, stderr, err := cmd.Output("hostname", hostname)

	if err != nil {
		err = errors.Wrap(err, stderr)
		return err
	}

	// Update /etc/hostname for persistant change after reboot
	config := hostname
	err = ioutil.WriteFile("/etc/hostname", []byte(config), 0777)

	if err != nil {
		return err
	}

	return nil
}

// LoadModules - Load modprobe Linux kernal modules
func LoadModules(modules []string) error {
	for _, module := range modules {

		_, stderr, err := cmd.Output("modprobe", module)

		if err != nil {
			return errors.Wrap(err, stderr)
		}

	}

	return nil
}
