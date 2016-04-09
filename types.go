package cli

import (
	"fmt"
)

type Helper struct {
	Help bool `cli:"!h,help" usage:"display help"`
}

type Addr struct {
	Host string `cli:"host" usage:"specify host" dft:"0.0.0.0"`
	Port uint16 `cli:"port" usage:"specify port" dft:"8080"`
}

type AddrWithShort struct {
	Host string `cli:"H,host" usage:"specify host" dft:"0.0.0.0"`
	Port uint16 `cli:"p,port" usage:"specify port" dft:"8080"`
}

func (addr Addr) ToString() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

func (addr AddrWithShort) ToString() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}
