package dolly

import "net"

// randomPort returns a random port number that is not in use.
func randomPort() int {
	addr, _ := net.Listen("tcp", ":0")
	addr.Close()
	return addr.Addr().(*net.TCPAddr).Port
}
