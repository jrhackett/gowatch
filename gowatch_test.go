package main_test

import (
	"testing"

	gowatch "github.com/jrhackett/gowatch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGowatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gowatch")
}

var _ = Describe("Gowatch", func() {
	Describe("NewFileWatcher", func() {
		It("should return a non-nil FileWatcher", func() {
			Expect(gowatch.NewFileWatcher("./")).ToNot(BeNil())
		})
	})
})
