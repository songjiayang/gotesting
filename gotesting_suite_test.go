package gotesting

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGotesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gotesting Suite")
}
