package csptenant_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEcho(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "cspTenant Suite")
}
