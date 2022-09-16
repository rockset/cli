package cmd

import (
	"os"
	"testing"
)

func skipUnlessIntegrationTest(t *testing.T) {
	const keyName = "ROCKSET_APIKEY"
	if os.Getenv(keyName) == "" {
		t.Skipf("could not find %s", keyName)
	}
}
