package constructor

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	structs "nenvoy.com/pkg/constants"
)

//ConvertYAML ...
func ConvertYAML(filename string) (*structs.VirtualNetworkDefinition, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	netDef := &structs.VirtualNetworkDefinition{}
	err = yaml.Unmarshal(buf, netDef)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return netDef, nil
}
