package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kamilsk/grafaman/internal/cmd"
)

var _ = Describe("fetch queries", func() {
	BeforeEach(func() {
		buffer.Reset()

		root = New()
		root.SetErr(buffer)
		root.SetOut(buffer)
	})

	When("invalid usage", func() {
		It("returns an error if a Graphite API endpoint is omitted", func() {
			root.SetArgs([]string{"metrics", "-m", "apps.services.awesome-service"})
			Expect(root.Execute()).To(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("please provide Graphite API endpoint"))
		})

		It("returns an error if a subset of metrics is omitted", func() {
			root.SetArgs([]string{"metrics", "--graphite", graphite.URL})
			Expect(root.Execute()).To(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("please provide metric prefix"))
		})

		It("returns an error if a subset of metrics is invalid", func() {
			root.SetArgs([]string{
				"metrics",
				"--graphite", graphite.URL,
				"-m", "$invalid.name",
			})
			Expect(root.Execute()).To(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("invalid metric prefix: $invalid.name"))
		})
	})

	When("correct usage", func() {})
})
