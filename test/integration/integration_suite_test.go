// +build integration

package integration_test

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	runLocal bool
)

func init() {
	flag.BoolVar(&runLocal, "local", true, "Run the integration tests locally using SAM and DynamoDB emulator")
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Test Suite")
}
