package kubelet // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/kubelet"

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	stats "k8s.io/kubelet/pkg/apis/stats/v1alpha1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/metadata"
)

func addIOMetrics(
	mb *metadata.MetricsBuilder,
	ioMetrics metadata.IoMetrics,
	s *stats.IOStats,
	currentTime pcommon.Timestamp,
) {
	if s == nil {
		return
	}

	// PSI stats
	if s.PSI != nil {
		if s.PSI.Full.Total != 0 {
			ioMetrics.PressureStalled(mb, currentTime, float64(s.PSI.Full.Total)/1_000_000) // Microseconds to seconds
		}
		if s.PSI.Some.Total != 0 {
			ioMetrics.PressureWaiting(mb, currentTime, float64(s.PSI.Some.Total)/1_000_000) // Microseconds to seconds
		}
	}
}
