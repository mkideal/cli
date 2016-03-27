package cli

type Helper struct {
	Help bool `cli:"!h,help" usage:"display help"`
}

type Addr struct {
	Host string `cli:"host" usage:"remote host" dft:"127.0.0.1"`
	Port uint16 `cli:"port" usage:"remote port" dft:"8080"`
}
