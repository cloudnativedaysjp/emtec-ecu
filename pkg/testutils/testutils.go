
package testutils

import (
	"os"
	"testing"
)

func Getenv(t *testing.T, key string) string {
	v := os.Getenv(key)
	if v == "" {
		t.Fatalf("envvar %s was not set", key)
	}
	return v
}
