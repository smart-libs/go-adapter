package httpadpt

// Helper functions and types for testing

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

// mockRequest is a test implementation of Request
type mockRequest struct {
	query QueryParams
}

func (m *mockRequest) Query() QueryParams {
	return m.query
}

// mockQueryParams is a test implementation of QueryParams
type mockQueryParams struct {
	values map[string][]string
}

func (m *mockQueryParams) GetValue(name string) ([]string, bool) {
	if m.values == nil {
		return nil, false
	}
	val, found := m.values[name]
	return val, found
}
