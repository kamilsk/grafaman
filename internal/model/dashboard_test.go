package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/kamilsk/grafaman/internal/model"
)

func TestDashboard(t *testing.T) {
	t.Run("skip raw", func(t *testing.T) {
		dashboard := Dashboard{
			Prefix: "apps.services.service",
			RawData: Queries{
				"apps.services.service.api.$source.$method.POST.request_time.$code.count",
				"apps.services.service.api.$source.$method.POST.request_time.$code.percentile.$percentile",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.*.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.200.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.4*.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.4*.percentile.$percentile",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.5*.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.5*.percentile.$percentile",
				"env.$env.apps.*.api.$source.$method.POST.request_time.5*.percentile.$percentile",
			},
		}
		queries, err := dashboard.Queries(Config{SkipRaw: true})
		require.NoError(t, err)
		require.Len(t, queries, len(dashboard.RawData)-1)
	})

	t.Run("with invalid queries", func(t *testing.T) {
		dashboard := Dashboard{
			RawData: Queries{
				"",
			},
		}
		queries, err := dashboard.Queries(Config{})
		require.Error(t, err)
		require.Nil(t, queries)
	})

	t.Run("issue#37, with limits of issue#36", func(t *testing.T) {
		dashboard := Dashboard{
			Prefix: "apps.services.service",
			RawData: Queries{
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.$code.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.$code.percentile.$percentile",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.*.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.200.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.4*.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.4*.percentile.$percentile",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.5*.count",
				"env.$env.apps.services.service.api.$source.$method.POST.request_time.5*.percentile.$percentile",
			},
			Variables: []Variable{
				{Name: "env"},
				{
					Name: "source",
					Options: []Option{
						{Name: "consumer", Value: "consumer"},
					},
				},
				{
					Name: "method",
					Options: []Option{
						{Name: "get", Value: "get"},
						{Name: "info", Value: "info"},
					},
				},
				{
					Name: "code",
					Options: []Option{
						{Name: "200", Value: "200"},
						{Name: "499", Value: "499"},
					},
				},
				{
					Name: "percentile",
					Options: []Option{
						{Name: "75", Value: "75"},
						{Name: "95", Value: "95"},
					},
				},
			},
		}
		queries, err := dashboard.Queries(Config{
			SkipRaw:        false,
			SkipDuplicates: false,
			NeedSorting:    true,
			Unpack:         true,
			TrimPrefixes:   []string{"complex.$env.", "env-staging.", "env.$env."},
		})
		require.NoError(t, err)
		require.Len(t, queries, len(dashboard.RawData)-1)

		reporter := NewCoverageReporter(queries)
		report := reporter.CoverageReport(Metrics{
			"apps.services.service.api.source.get.POST.request_time.200.count",
			"apps.services.service.api.source.get.POST.request_time.200.max",
			"apps.services.service.api.source.get.POST.request_time.200.mean",
			"apps.services.service.api.source.get.POST.request_time.200.median",
			"apps.services.service.api.source.get.POST.request_time.200.min",
			"apps.services.service.api.source.get.POST.request_time.200.percentile.75",
			"apps.services.service.api.source.get.POST.request_time.200.percentile.95",
			"apps.services.service.api.source.get.POST.request_time.200.percentile.98",
			"apps.services.service.api.source.get.POST.request_time.200.percentile.99",
			"apps.services.service.api.source.get.POST.request_time.200.percentile.999",
			"apps.services.service.api.source.get.POST.request_time.200.sum",
			"apps.services.service.api.source.get.POST.request_time.499.count",
			"apps.services.service.api.source.get.POST.request_time.499.max",
			"apps.services.service.api.source.get.POST.request_time.499.mean",
			"apps.services.service.api.source.get.POST.request_time.499.median",
			"apps.services.service.api.source.get.POST.request_time.499.min",
			"apps.services.service.api.source.get.POST.request_time.499.percentile.75",
			"apps.services.service.api.source.get.POST.request_time.499.percentile.95",
			"apps.services.service.api.source.get.POST.request_time.499.percentile.98",
			"apps.services.service.api.source.get.POST.request_time.499.percentile.99",
			"apps.services.service.api.source.get.POST.request_time.499.percentile.999",
			"apps.services.service.api.source.get.POST.request_time.499.sum",
			"apps.services.service.api.source.get.POST.request_time.500.count",
			"apps.services.service.api.source.get.POST.request_time.500.max",
			"apps.services.service.api.source.get.POST.request_time.500.mean",
			"apps.services.service.api.source.get.POST.request_time.500.median",
			"apps.services.service.api.source.get.POST.request_time.500.min",
			"apps.services.service.api.source.get.POST.request_time.500.percentile.75",
			"apps.services.service.api.source.get.POST.request_time.500.percentile.95",
			"apps.services.service.api.source.get.POST.request_time.500.percentile.98",
			"apps.services.service.api.source.get.POST.request_time.500.percentile.99",
			"apps.services.service.api.source.get.POST.request_time.500.percentile.999",
			"apps.services.service.api.source.get.POST.request_time.500.sum",
			"apps.services.service.api.source.get.POST.request_time.503.count",
			"apps.services.service.api.source.get.POST.request_time.503.max",
			"apps.services.service.api.source.get.POST.request_time.503.mean",
			"apps.services.service.api.source.get.POST.request_time.503.median",
			"apps.services.service.api.source.get.POST.request_time.503.min",
			"apps.services.service.api.source.get.POST.request_time.503.percentile.75",
			"apps.services.service.api.source.get.POST.request_time.503.percentile.95",
			"apps.services.service.api.source.get.POST.request_time.503.percentile.98",
			"apps.services.service.api.source.get.POST.request_time.503.percentile.99",
			"apps.services.service.api.source.get.POST.request_time.503.percentile.999",
			"apps.services.service.api.source.get.POST.request_time.503.sum",
		}.Exclude(Queries{"*.max", "*.mean", "*.median", "*.min", "*.sum"}.MustMatchers()...))
		assert.Equal(t, 100.0, report.Total())
	})
}
