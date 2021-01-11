package network_test

import (
	"net"
	"testing"

	"nenvoy.com/pkg/utils/network"
)

func TestGetInterface(t *testing.T) {

	network.GetInterface(net.ParseIP("20.0.0.1"))

}
