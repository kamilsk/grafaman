package cmd_test

import (
	"strings"

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

	metrics := strings.TrimSpace(`
a.b.c
a.b.d
a.b.e
a.f.g
a.f.h
a.i.j
a.i.k`)

	When("correct usage", func() {
		It("returns metrics if all work well", func() {
			root.SetArgs([]string{
				"metrics",
				"--graphite", graphite.URL,
				"-m", "a",
				"-f", "tsv",
				"--no-cache",
			})
			Expect(root.Execute()).ToNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(metrics))
		})
	})
})
