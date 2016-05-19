package metar_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoMetar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoMetar Suite")
}
