package metrics

type MockMetrics struct{}

func NewMockMetrics() *MockMetrics {
	return &MockMetrics{}
}

func (m *MockMetrics) IncPVZCreated() {}

func (m *MockMetrics) IncReceptionCreated() {}

func (m *MockMetrics) IncProductAdded() {}

func (m *MockMetrics) ObserveRequestDuration(method, endpoint string, duration float64) {}

func (m *MockMetrics) IncRequestCount(method, endpoint, status string) {}
