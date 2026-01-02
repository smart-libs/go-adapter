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
	query  QueryParams
	header HeaderParams
	path   PathParams
}

func (m *mockRequest) Query() QueryParams {
	return m.query
}

func (m *mockRequest) Header() HeaderParams {
	return m.header
}

func (m *mockRequest) Path() PathParams {
	return m.path
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

// mockPathParams is a test implementation of PathParams
type mockPathParams struct {
	values map[string]string
}

func (m *mockPathParams) GetValue(name string) (string, bool) {
	if m.values == nil {
		return "", false
	}
	val, found := m.values[name]
	return val, found
}

// mockHeaderParams is a test implementation of HeaderParams
type mockHeaderParams struct {
	values map[string][]string
}

func (m *mockHeaderParams) GetValue(name string) ([]string, bool) {
	if m.values == nil {
		return nil, false
	}
	val, found := m.values[name]
	return val, found
}
