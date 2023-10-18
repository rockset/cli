package test

import (
	"os"
	"testing"
)

func SkipUnlessIntegrationTest(t *testing.T) {
	const keyName = "ROCKSET_APIKEY"
	t.Helper()
	if os.Getenv(keyName) == "" {
		t.Skipf("could not find %s", keyName)
	}
}
