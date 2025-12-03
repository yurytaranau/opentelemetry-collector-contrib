// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kubelet // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/kubelet"

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	stats "k8s.io/kubelet/pkg/apis/stats/v1alpha1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/metadata"
)

func addCPUMetrics(
	mb *metadata.MetricsBuilder,
	cpuMetrics metadata.CPUMetrics,
	s *stats.CPUStats,
	currentTime pcommon.Timestamp,
	r resources,
	nodeCPULimit float64,
) {
	if s == nil {
		return
	}
	if s.UsageNanoCores != nil {
		usageCores := float64(*s.UsageNanoCores) / 1_000_000_000
		cpuMetrics.Usage(mb, currentTime, usageCores)
		addCPUUtilizationMetrics(mb, cpuMetrics, usageCores, currentTime, r, nodeCPULimit)
	}
	addCPUTimeMetric(mb, cpuMetrics.Time, s, currentTime)
	if s.PSI != nil {
		if s.PSI.Full.Total != 0 {
			cpuMetrics.PressureStalled(mb, currentTime, float64(s.PSI.Full.Total)/1_000_000)
		}
		if s.PSI.Some.Total != 0 {
			cpuMetrics.PressureWaiting(mb, currentTime, float64(s.PSI.Some.Total)/1_000_000)
		}
	}
}

func addCPUUtilizationMetrics(
	mb *metadata.MetricsBuilder,
	cpuMetrics metadata.CPUMetrics,
	usageCores float64,
	currentTime pcommon.Timestamp,
	r resources,
	nodeCPULimit float64,
) {
	if nodeCPULimit > 0 {
		cpuMetrics.NodeUtilization(mb, currentTime, usageCores/nodeCPULimit)
	}
	if r.cpuLimit > 0 {
		cpuMetrics.LimitUtilization(mb, currentTime, usageCores/r.cpuLimit)
	}
	if r.cpuRequest > 0 {
		cpuMetrics.RequestUtilization(mb, currentTime, usageCores/r.cpuRequest)
	}
}

func addCPUTimeMetric(mb *metadata.MetricsBuilder, recordDataPoint metadata.RecordDoubleDataPointFunc, s *stats.CPUStats, currentTime pcommon.Timestamp) {
	if s.UsageCoreNanoSeconds == nil {
		return
	}
	value := float64(*s.UsageCoreNanoSeconds) / 1_000_000_000
	recordDataPoint(mb, currentTime, value)
}
