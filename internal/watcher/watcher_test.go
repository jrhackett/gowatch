package watcher_test

import (
	"gowatch/internal/watcher"
	"testing"

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
			Expect(watcher.NewFileWatcher("./")).ToNot(BeNil())
		})
	})

	Describe("ShouldIgnoreFile", func() {
		It("should return true if path starts with a period", func() {
			Expect(watcher.ShouldIgnoreFile(".something")).To(BeTrue())
		})

		It("should return true if path starts with an underscore", func() {
			Expect(watcher.ShouldIgnoreFile("_something")).To(BeTrue())
		})

		It("should return true if path starts with vendor", func() {
			Expect(watcher.ShouldIgnoreFile("vendorsomething")).To(BeTrue())
		})

		It("should return false if path is valid file", func() {
			Expect(watcher.ShouldIgnoreFile("yay/what/hi")).To(BeFalse())
		})
	})
})
