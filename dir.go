package dolly

import "os"

func randomDir() string {
	tmp, _ := os.MkdirTemp(os.TempDir(), "dolly")
	return tmp
}
