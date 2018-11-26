package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/internal"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const VerboseAuth = "RELINT_VERBOSE_AUTH"
const CFAuthRetries = 2

func CfAuth(cmdStarter internal.Starter, reporter internal.Reporter, user string, password string, timeout time.Duration) error {
	var auth *gexec.Session
	var err error
	var failures []string

	args := []string{"auth", user, password}
	if os.Getenv(VerboseAuth) == "true" {
		args = append(args, "-v")
	}

	for i := 0; i < CFAuthRetries; i++ {
		auth, err = cmdStarter.Start(reporter, "cf", args...)
		if err != nil {
			return err
		}

		failures = InterceptGomegaFailures(func() {
			auth.Wait(timeout)
		})

		if len(failures) == 0 && auth.ExitCode() == 0 {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	if len(failures) != 0  {
		return fmt.Errorf("cf auth command timed out: %s", failures)
	}

	if auth.ExitCode() != 0 {
		return fmt.Errorf("cf auth command exited with %d", auth.ExitCode())
	}

	return nil
}
