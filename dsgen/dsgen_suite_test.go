package dsgen_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDsgen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dsgen Suite")
}
