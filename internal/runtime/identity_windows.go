//go:build windows

package runtime

import "fmt"

func currentRuntimeIdentityImpl() (int, int, error) {
	return 0, 0, fmt.Errorf("runtime identity is only supported on Linux")
}
