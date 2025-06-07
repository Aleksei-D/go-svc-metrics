package storage

type MemStorage struct {
	CounterMetrics map[string]int64
	GaugeMetrics   map[string]float64
}

type Repositories interface {
	UpdateGauge(metricName string, value float64)
	UpdateCounter(metricName string, value int64)
}

func InitMemStorage() *MemStorage {
	return &MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}
}

func (m *MemStorage) UpdateCounter(metricName string, value int64) {
	_, ok := m.CounterMetrics[metricName]
	if !ok {
		m.CounterMetrics[metricName] = value
	} else {
		m.CounterMetrics[metricName] += value
	}
}

func (m *MemStorage) UpdateGauge(metricName string, value float64) {
	_, ok := m.GaugeMetrics[metricName]
	if !ok {
		m.GaugeMetrics[metricName] = value
	} else {
		m.GaugeMetrics[metricName] += value
	}
}
