package ginkgotesting_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGinkgotesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ginkgotesting Suite")
}
