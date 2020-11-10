package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureServiceID(t *testing.T) {
	tests := []struct {
		name      string
		promql    string
		serviceID string
		want      string
	}{
		{
			name:      "apple",
			promql:    "container_memory_rss{name=~\"k8s_1aa695008a85f6aacf3b9ed6342279c5.*\"}/1024/1024",
			serviceID: "foobar",
			want:      "container_memory_rss{name=~\"k8s_1aa695008a85f6aacf3b9ed6342279c5.*\",service_id=\"foobar\"} / 1024 / 1024",
		},
		{
			name:      "banana",
			promql:    "rate(container_network_transmit_bytes_total{name=~\"k8s_POD_1aa695008a85f6aacf3b9ed6342279c5.*\"}[1m])/1024",
			serviceID: "foobar",
			want:      "rate(container_network_transmit_bytes_total{name=~\"k8s_POD_1aa695008a85f6aacf3b9ed6342279c5.*\",service_id=\"foobar\"}[1m]) / 1024",
		},
		{
			name:      "cat",
			promql:    "sum(rate(container_cpu_usage_seconds_total{name=~\"k8s_1aa695008a85f6aacf3b9ed6342279c5.*\"}[1m])) by (pod, namespace) / (sum(container_spec_cpu_quota{name=~\"k8s_1aa695008a85f6aacf3b9ed6342279c5.*\"}/container_spec_cpu_period{name=~\"k8s_1aa695008a85f6aacf3b9ed6342279c5.*\"}) by (pod,namespace)) * 100",
			serviceID: "foobar",
			want:      "sum by(pod, namespace) (rate(container_cpu_usage_seconds_total{name=~\"k8s_1aa695008a85f6aacf3b9ed6342279c5.*\",service_id=\"foobar\"}[1m])) / (sum by(pod, namespace) (container_spec_cpu_quota{name=~\"k8s_1aa695008a85f6aacf3b9ed6342279c5.*\",service_id=\"foobar\"} / container_spec_cpu_period{name=~\"k8s_1aa695008a85f6aacf3b9ed6342279c5.*\",service_id=\"foobar\"})) * 100",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			got, err := ensureServiceID(tc.promql, tc.serviceID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				t.FailNow()
			}

			assert.EqualValues(t, tc.want, got)
		})
	}
}
