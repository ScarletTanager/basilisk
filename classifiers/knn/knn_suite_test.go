package knn_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKnn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Knn Suite")
}
