package systemd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"

	cmd "nenvoy.com/pkg/utils/cmd"
	"nenvoy.com/pkg/utils/printing"
	root "nenvoy.com/pkg/utils/root"
)

// Config -
type Config struct {
	Description      string
	Documentation    string
	BinaryPath       string
	Vars             []Flag
	RestartOnFailure string
	RestartSec       int
	Options          []Flag
	WantedBy         string
	FilePermissions  os.FileMode
	Requires         []string
	After            []string
	Pre              []string
}

// Flag -
type Flag struct {
	Flag      string
	Value     string
	Delimiter string
}

// CreateBasicConfig -
func CreateBasicConfig(description, documentation, binaryPath string) (config Config) {
	config = Config{
		Description:      description,
		Documentation:    documentation,
		BinaryPath:       binaryPath,
		RestartOnFailure: "on-failure",
		RestartSec:       5,
		WantedBy:         "multi-user.target",
		FilePermissions:  0644,
	}

	return config
}

// AddOption -
func (c *Config) AddOption(flag, value, delimiter string) {
	option := Flag{Flag: flag, Value: value, Delimiter: delimiter}
	c.Options = append(c.Options, option)
}

// AddExecAfter -
func (c *Config) AddExecAfter(serviceName string) {
	c.After = append(c.After, serviceName)
}

// AddExecRequires -
func (c *Config) AddExecRequires(serviceName string) {
	c.Requires = append(c.Requires, serviceName)
}

// AddExecPre -
func (c *Config) AddExecPre(serviceName string) {
	c.Pre = append(c.Pre, serviceName)
}

// AddExecVar -
func (c *Config) AddExecVar(flag, value string) {
	v := Flag{Flag: flag, Value: value, Delimiter: "="}
	c.Vars = append(c.Vars, v)
}

// FormatFlag -
func (o Flag) FormatFlag(flush bool) (entry string) {
	if !flush {
		return fmt.Sprintf("  --%s%s%s", o.Flag, o.Delimiter, o.Value)
	}
	return fmt.Sprintf("%s%s%s", o.Flag, o.Delimiter, o.Value)
}

// WriteToFile -
func (c Config) WriteToFile(filePath string) (err error) {
	wd := "\n\n[Unit]"
	wd = fmt.Sprintf("\nDescription=%s", c.Description)
	wd += fmt.Sprintf("\nDocumentation=%s", c.Documentation)
	if c.After != nil {
		wd += fmt.Sprintf("\nAfter=%s", strings.Join(c.After, " "))
	}
	if c.Requires != nil {
		wd += fmt.Sprintf("\nRequires=%s", strings.Join(c.Requires, " "))
	}
	wd += "\n\n[Service]"
	if c.Pre != nil {
		wd += fmt.Sprintf("\nExecStartPre=%s", strings.Join(c.Pre, " "))
	}
	wd += fmt.Sprintf("\nExecStart=%s", c.BinaryPath)
	for _, o := range c.Options {
		wd += fmt.Sprintf(" \\\n%s", o.FormatFlag(false))
	}
	wd += "\n\n# Vars"
	for _, v := range c.Vars {
		wd += fmt.Sprintf(" \n%s", v.FormatFlag(true))
	}
	wd += fmt.Sprintf("\nRestart=%s", c.RestartOnFailure)
	wd += fmt.Sprintf("\nRestartSec=%d", c.RestartSec)
	wd += "\n\n[Install]"
	wd += fmt.Sprintf("\nWantedBy=%s", c.WantedBy)

	err = ioutil.WriteFile(filePath, []byte(wd), c.FilePermissions)
	if err != nil {
		return err
	}

	return nil
}

// StartService - Start systemd service, required root permissions
func StartService(name string) (err error) {
	// Check running as root
	root, err := root.AsRoot()

	if err != nil {
		return err
	}

	if !root {
		err := errors.New("permission error: root required")
		return errors.Wrap(err, "failed on StartService (1):")
	}

	// Start systemd process
	_, stderr, err := cmd.Output("systemctl", "daemon-reload")

	if err != nil {
		err = errors.Wrap(err, stderr)
		return errors.Wrap(err, "failed on StartService (2):")
	}

	_, stderr, err = cmd.Output("systemctl", "enable", name)

	if err != nil {
		err = errors.Wrap(err, stderr)
		return errors.Wrap(err, "failed on StartService (3):")
	}

	_, stderr, err = cmd.Output("systemctl", "start", name)

	if err != nil {
		err = errors.Wrap(err, stderr)
		return errors.Wrap(err, "failed on StartService (4):")
	}

	return nil
}

// StopService - Stop systemd service, required root permissions
func StopService(name string) (err error) {
	// Check running as root
	root, err := root.AsRoot()

	if err != nil {
		return err
	}

	if !root {
		err := errors.New("permission error: root required")
		return errors.Wrap(err, "failed on StopService (1):")
	}

	// Stop systemd process
	_, stderr, err := cmd.Output("systemctl", "stop", name)

	if err != nil {
		err = errors.Wrap(err, stderr)
		return errors.Wrap(err, "failed on StopService (2):")
	}

	_, stderr, err = cmd.Output("systemctl", "disable", name)

	if err != nil {
		err = errors.Wrap(err, stderr)
		return errors.Wrap(err, "failed on StopService (3):")
	}

	_, stderr, err = cmd.Output("systemctl", "daemon-reload")

	if err != nil {
		err = errors.Wrap(err, stderr)
		return errors.Wrap(err, "failed on StopService (4):")
	}

	_, stderr, err = cmd.Output("systemctl", "reset-failed")

	if err != nil {
		err = errors.Wrap(err, stderr)
		return errors.Wrap(err, "failed on StopService (5):")
	}

	return nil
}

// ServiceExists - Check if systemd service exists, required root permissions
func ServiceExists(name string) (exists bool, err error) {
	// Check running as root
	root, err := root.AsRoot()

	if err != nil {
		return false, err
	}

	if !root {
		err := errors.New("permission error: root required")
		return false, errors.Wrap(err, "failed on ServiceExists (1):")
	}

	// Search for systemd service
	_, stderr, err := cmd.Output("systemctl", "status", name+".service")
	if err != nil && len(stderr) > 0 {
		err = errors.Wrap(err, stderr)
		return false, nil
	}

	return true, nil
}

// StartExistingServices - Start all existing services passed
func StartExistingServices(services []string) (err error) {

	for _, service := range services {
		// Check each service exists. If so stop it
		exists, err := ServiceExists(service)
		if err != nil {
			return err
		} else if exists == true {
			err := StartService(service)
			if err != nil {
				return err
			}
			printing.PrintSuccess("Started " + service)
		} else {
			return errors.New("Service " + service + " not found")
		}

	}
	return nil
}

// StopExistingServices - Stops any existing services which may interfere with setup\
func StopExistingServices(services []string) (err error) {

	for _, service := range services {
		// Check each service exists. If so stop it
		exists, err := ServiceExists(service)
		if err != nil {
			return err
		} else if exists == true {
			err := StopService(service)
			if err != nil {
				return err
			}
		} else {
			printing.PrintWarning("Service " + service + " not found")
		}

	}
	return nil
}
