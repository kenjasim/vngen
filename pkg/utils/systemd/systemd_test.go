package systemd_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"nenvoy.com/pkg/utils/printing"
	"nenvoy.com/pkg/utils/systemd"
)

var (
	TestDir         = "/tmp/nenvoy/test/systemd"
	SystemdFileName = "/tmp/nenvoy/test/systemd/systemd-test.service"
	Description     = "systemd-description-test"
	Documentation   = "www.documentation-test.url"
	BinaryPath      = "/bin/path/test"
)

// TestSystemdFileCreation
func TestSystemdFileCreation(t *testing.T) {

	// Create test directory
	err := os.MkdirAll(TestDir, os.ModePerm)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, fmt.Sprintf("failed to create test directory: %s", TestDir)))
	}

	// Create basic config file
	config := systemd.CreateBasicConfig(Description, Documentation, BinaryPath)

	// Add options
	config.AddOption("flag1", "value1", "=")
	config.AddOption("flag2", "value2", " ")
	config.AddOption("flag3", "value3", "=")
	config.AddOption("flag4", "", "")
	config.AddOption("flag5", "value5", " ")

	// Add After
	config.AddExecAfter("service1-test.service")
	config.AddExecAfter("service2-test.service")

	// Add Requires
	config.AddExecRequires("service1-test.service")
	config.AddExecRequires("service2-test.service")

	// Add Pre
	config.AddExecPre("service1-test.service")
	config.AddExecPre("service2-test.service")

	// Add Vars
	config.AddExecVar("KillMode", "process")
	config.AddExecVar("Restart", "always")

	// Write to file
	err = config.WriteToFile(SystemdFileName)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, fmt.Sprintf("failed to write file to: %s", SystemdFileName)))
	}

	// // Verify syntax
	// stdout, stderr, err := cmd.Output("systemd-analyze", "verify", SystemdFileName)
	// if err != nil {
	// 	t.Fatalf("%s", errors.Wrap(err, fmt.Sprintf("failed verification: %s", stderr)))
	// }

	// t.Log(printing.SprintSuccess(fmt.Sprintf("Verified file: %s", stdout)))

	// Read file
	dat, err := ioutil.ReadFile(SystemdFileName)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, fmt.Sprintf("failed to read file from: %s", SystemdFileName)))
	}

	t.Log(printing.SprintSuccess(fmt.Sprintf("systemd service file generated:\n%s", string(dat))))

}
