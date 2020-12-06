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
		It("returns an error if a Grafana API endpoint is omitted", func() {
			root.SetArgs([]string{"queries", "-d", "uid"})
			Expect(root.Execute()).To(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("please provide Grafana API endpoint"))
		})

		It("returns an error if a dashboard unique identifier is omitted", func() {
			root.SetArgs([]string{"queries", "--grafana", grafana.URL})
			Expect(root.Execute()).To(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("please provide a dashboard unique identifier"))
		})

		It("returns an error if a subset of metrics is invalid", func() {
			root.SetArgs([]string{
				"queries",
				"--grafana", grafana.URL,
				"-d", "uid",
				"-m", "$invalid.name",
			})
			Expect(root.Execute()).To(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("invalid metric prefix: $invalid.name"))
		})
	})

	When("correct usage", func() {})
})
