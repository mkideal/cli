package cli

type Helper struct {
	Help bool `cli:"!h,help" usage:"display help"`
}

type Addr struct {
	Host string `cli:"host" usage:"remote host"`
	Port uint16 `cli:"port" usage:"remote port" dft:"8080"`
}
