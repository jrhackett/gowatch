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

	Describe("ShouldIgnoreFile", func() {
		It("should return true if path starts with a period", func() {
			Expect(gowatch.ShouldIgnoreFile(".something")).To(BeTrue())
		})

		It("should return true if path starts with an underscore", func() {
			Expect(gowatch.ShouldIgnoreFile("_something")).To(BeTrue())
		})

		It("should return true if path starts with vendor", func() {
			Expect(gowatch.ShouldIgnoreFile("vendorsomething")).To(BeTrue())
		})

		It("should return false if path is valid file", func() {
			Expect(gowatch.ShouldIgnoreFile("yay/what/hi")).To(BeFalse())
		})
	})
})
