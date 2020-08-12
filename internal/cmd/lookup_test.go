package cmd_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	. "github.com/kamilsk/grafaman/internal/cmd"
)

var _ = Describe("lookup cache", func() {
	BeforeEach(func() {
		buffer.Reset()
		viper.Reset()

		root = New()
		root.SetOut(buffer)
	})

	When("invalid usage", func() {
		It("returns an error if a subset of metrics is omitted", func() {
			root.SetArgs([]string{"cache-lookup"})
			Expect(root.Execute()).To(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("please provide metric prefix"))
		})

		It("returns an error if a subset of metrics is invalid", func() {
			root.SetArgs([]string{"cache-lookup", "-m", "$invalid.name"})
			Expect(root.Execute()).To(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("invalid metric prefix: $invalid.name"))
		})
	})

	When("when correct usage", func() {
		It("contains temp dir and a subset of metrics", func() {
			root.SetArgs([]string{"cache-lookup", "-m", "apps.services.awesome-service"})
			Expect(root.Execute()).ToNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(os.TempDir()))
			Expect(buffer.String()).To(ContainSubstring("apps.services.awesome-service.grafaman.json"))
		})
	})
})
