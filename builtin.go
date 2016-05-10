package cli

import (
	"fmt"
)

// Helper is builtin Help flag
type Helper struct {
	Help bool `cli:"!h,help" usage:"display help information"`
}

// AutoHelp implements AutoHelper interface
func (h Helper) AutoHelp() bool {
	return h.Help
}

// Addr is builtin host,port flag
// NOTE:Deprecated
type Addr struct {
	Host string `cli:"host" usage:"specify host" dft:"0.0.0.0"`
	Port uint16 `cli:"port" usage:"specify port" dft:"8080"`
}

// AddrWithShort is builtin host,port flag contains short flag
// NOTE:Deprecated
type AddrWithShort struct {
	Host string `cli:"H,host" usage:"specify host" dft:"0.0.0.0"`
	Port uint16 `cli:"p,port" usage:"specify port" dft:"8080"`
}

// ToString ...
// NOTE:Deprecated
func (addr Addr) ToString() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

// ToString ...
// NOTE:Deprecated
func (addr AddrWithShort) ToString() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}
