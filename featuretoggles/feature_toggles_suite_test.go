package featuretoggles

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestFeatureToggles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FeatureToggles Suite")
}
