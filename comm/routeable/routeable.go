package routeable

import (
	"net"

	"github.com/p9c/pkg/app/slog"
)

// GetInterface returns the address and interface of multicast capable interfaces
func GetInterface() (lanInterface []*net.Interface) {
	var err error
	var interfaces []net.Interface
	if interfaces, err = net.Interfaces(); slog.Check(err) {
	}
	// Traces(interfaces)
	for ifi := range interfaces {
		if interfaces[ifi].Flags&net.FlagLoopback == 0 && interfaces[ifi].
			HardwareAddr != nil {
			// iads, _ := interfaces[ifi].Addrs()
			// for i := range iads {
			//	//Traces(iads[i].Network())
			// }
			// Debug(interfaces[ifi].MulticastAddrs())
			lanInterface = append(lanInterface, &interfaces[ifi])
		}
	}
	// Traces(lanInterface)
	return
}
