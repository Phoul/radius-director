//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package runtime

import "os"

func currentRuntimeIdentityImpl() (int, int, error) {
	return os.Getuid(), os.Getgid(), nil
}
