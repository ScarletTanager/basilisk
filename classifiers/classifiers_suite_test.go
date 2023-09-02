package classifiers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClassifiers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Classifiers Suite")
}
