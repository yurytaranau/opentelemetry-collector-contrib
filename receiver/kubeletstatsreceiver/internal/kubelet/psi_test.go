package kubelet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/receiver/receivertest"
	stats "k8s.io/kubelet/pkg/apis/stats/v1alpha1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/metadata"
)

func TestAddCPUMetrics_PSI(t *testing.T) {
	mb := metadata.NewMetricsBuilder(metadata.DefaultMetricsBuilderConfig(), receivertest.NewNopSettings(metadata.Type))
	currentTime := pcommon.NewTimestampFromTime(time.Now())

	psiTotal := uint64(1000000) // 1 second in microseconds
	s := &stats.CPUStats{
		PSI: &stats.PSIStats{
			Full: stats.PSIData{Total: psiTotal},
			Some: stats.PSIData{Total: psiTotal},
		},
	}

	addCPUMetrics(mb, metadata.ContainerCPUMetrics, s, currentTime, resources{}, 0)

	m := mb.Emit()
	assert.Equal(t, 1, m.ResourceMetrics().Len())
	sm := m.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()

	// Check for PSI metrics
	foundStalled := false
	foundWaiting := false
	for i := 0; i < sm.Len(); i++ {
		metric := sm.At(i)
		if metric.Name() == "container.cpu.pressure.stalled" {
			foundStalled = true
			assert.Equal(t, 1.0, metric.Sum().DataPoints().At(0).DoubleValue())
		}
		if metric.Name() == "container.cpu.pressure.waiting" {
			foundWaiting = true
			assert.Equal(t, 1.0, metric.Sum().DataPoints().At(0).DoubleValue())
		}
	}
	assert.True(t, foundStalled, "container.cpu.pressure.stalled not found")
	assert.True(t, foundWaiting, "container.cpu.pressure.waiting not found")
}

func TestAddMemoryMetrics_PSI(t *testing.T) {
	mb := metadata.NewMetricsBuilder(metadata.DefaultMetricsBuilderConfig(), receivertest.NewNopSettings(metadata.Type))
	currentTime := pcommon.NewTimestampFromTime(time.Now())

	psiTotal := uint64(2000000) // 2 seconds in microseconds
	s := &stats.MemoryStats{
		PSI: &stats.PSIStats{
			Full: stats.PSIData{Total: psiTotal},
			Some: stats.PSIData{Total: psiTotal},
		},
	}

	addMemoryMetrics(mb, metadata.ContainerMemoryMetrics, s, currentTime, resources{}, 0)

	m := mb.Emit()
	assert.Equal(t, 1, m.ResourceMetrics().Len())
	sm := m.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()

	// Check for PSI metrics
	foundStalled := false
	foundWaiting := false
	for i := 0; i < sm.Len(); i++ {
		metric := sm.At(i)
		if metric.Name() == "container.memory.pressure.stalled" {
			foundStalled = true
			assert.Equal(t, 2.0, metric.Sum().DataPoints().At(0).DoubleValue())
		}
		if metric.Name() == "container.memory.pressure.waiting" {
			foundWaiting = true
			assert.Equal(t, 2.0, metric.Sum().DataPoints().At(0).DoubleValue())
		}
	}
	assert.True(t, foundStalled, "container.memory.pressure.stalled not found")
	assert.True(t, foundWaiting, "container.memory.pressure.waiting not found")
}

func TestAddIOMetrics_PSI(t *testing.T) {
	mb := metadata.NewMetricsBuilder(metadata.DefaultMetricsBuilderConfig(), receivertest.NewNopSettings(metadata.Type))
	currentTime := pcommon.NewTimestampFromTime(time.Now())

	psiTotal := uint64(3000000) // 3 seconds in microseconds
	s := &stats.IOStats{
		PSI: &stats.PSIStats{
			Full: stats.PSIData{Total: psiTotal},
			Some: stats.PSIData{Total: psiTotal},
		},
	}

	addIOMetrics(mb, metadata.ContainerIoMetrics, s, currentTime)

	m := mb.Emit()
	assert.Equal(t, 1, m.ResourceMetrics().Len())
	sm := m.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()

	// Check for PSI metrics
	foundStalled := false
	foundWaiting := false
	for i := 0; i < sm.Len(); i++ {
		metric := sm.At(i)
		if metric.Name() == "container.io.pressure.stalled" {
			foundStalled = true
			assert.Equal(t, 3.0, metric.Sum().DataPoints().At(0).DoubleValue())
		}
		if metric.Name() == "container.io.pressure.waiting" {
			foundWaiting = true
			assert.Equal(t, 3.0, metric.Sum().DataPoints().At(0).DoubleValue())
		}
	}
	assert.True(t, foundStalled, "container.io.pressure.stalled not found")
	assert.True(t, foundWaiting, "container.io.pressure.waiting not found")
}
