package route122

import (
	"testing"
)

func TestBasicRouting(t *testing.T) {
	r := New()

	err := r.Handle("GET /users/{id}", "getUserHandler")
	if err != nil {
		t.Fatalf("Failed to register route: %v", err)
	}

	match, found := r.Match("GET", "", "/users/123")
	if !found {
		t.Fatal("Expected to find match")
	}

	if match.Handler != "getUserHandler" {
		t.Errorf("Expected handler=getUserHandler, got %v", match.Handler)
	}

	if match.Params["id"] != "123" {
		t.Errorf("Expected id=123, got %s", match.Params["id"])
	}
}

func TestWildcardRouting(t *testing.T) {
	r := New()

	err := r.Handle("GET /files/{path...}", "fileHandler")
	if err != nil {
		t.Fatalf("Failed to register route: %v", err)
	}

	match, found := r.Match("GET", "", "/files/docs/readme.txt")
	if !found {
		t.Fatal("Expected to find match")
	}

	expectedPath := "docs/readme.txt"
	if match.Params["path"] != expectedPath {
		t.Errorf("Expected path=%s, got %s", expectedPath, match.Params["path"])
	}
}

func TestHostSpecificRouting(t *testing.T) {
	r := New()

	err := r.Handle("GET api.example.com/users/{id}", "apiHandler")
	if err != nil {
		t.Fatalf("Failed to register route: %v", err)
	}

	// Test with correct host
	match, found := r.Match("GET", "api.example.com", "/users/123")
	if !found {
		t.Fatal("Expected to find match with correct host")
	}

	if match.Handler != "apiHandler" {
		t.Errorf("Expected handler=apiHandler, got %v", match.Handler)
	}

	if match.Params["id"] != "123" {
		t.Errorf("Expected id=123, got %s", match.Params["id"])
	}

	// Test with wrong host (should not match)
	_, found = r.Match("GET", "wrong.example.com", "/users/123")
	if found {
		t.Error("Should not match with wrong host")
	}
}

func TestMultipleRoutes(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		requestURL  string
		method      string
		host        string
		shouldMatch bool
		expectedParams map[string]string
	}{
		{
			name:           "Simple param route",
			pattern:        "GET /users/{id}",
			requestURL:     "/users/42",
			method:         "GET",
			host:           "",
			shouldMatch:    true,
			expectedParams: map[string]string{"id": "42"},
		},
		{
			name:           "POST route without params",
			pattern:        "POST /users",
			requestURL:     "/users",
			method:         "POST",
			host:           "",
			shouldMatch:    true,
			expectedParams: map[string]string{},
		},
		{
			name:           "Multiple params",
			pattern:        "GET /users/{id}/posts/{postId}",
			requestURL:     "/users/123/posts/456",
			method:         "GET",
			host:           "",
			shouldMatch:    true,
			expectedParams: map[string]string{"id": "123", "postId": "456"},
		},
		{
			name:           "Multi-segment wildcard",
			pattern:        "GET /files/{path...}",
			requestURL:     "/files/docs/api/readme.txt",
			method:         "GET",
			host:           "",
			shouldMatch:    true,
			expectedParams: map[string]string{"path": "docs/api/readme.txt"},
		},
		{
			name:           "Root path",
			pattern:        "GET /",
			requestURL:     "/",
			method:         "GET",
			host:           "",
			shouldMatch:    true,
			expectedParams: map[string]string{"": ""},
		},
		{
			name:           "Host-specific route",
			pattern:        "GET api.example.com/v1/users/{id}",
			requestURL:     "/v1/users/789",
			method:         "GET",
			host:           "api.example.com",
			shouldMatch:    true,
			expectedParams: map[string]string{"id": "789"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()

			err := r.Handle(tt.pattern, "handler")
			if err != nil {
				t.Fatalf("Failed to register route %s: %v", tt.pattern, err)
			}

			match, found := r.Match(tt.method, tt.host, tt.requestURL)

			if found != tt.shouldMatch {
				t.Errorf("Expected shouldMatch=%v, got %v", tt.shouldMatch, found)
				return
			}

			if tt.shouldMatch {
				if len(tt.expectedParams) != len(match.Params) {
					t.Errorf("Expected %d params, got %d. Expected: %v, Actual: %v", len(tt.expectedParams), len(match.Params), tt.expectedParams, match.Params)
				}

				for param, expectedValue := range tt.expectedParams {
					if actualValue, exists := match.Params[param]; !exists {
						t.Errorf("Expected param %s to exist", param)
					} else if actualValue != expectedValue {
						t.Errorf("Expected param %s=%s, got %s", param, expectedValue, actualValue)
					}
				}
			}
		})
	}
}

func TestMismatchCases(t *testing.T) {
	mismatchTests := []struct {
		name           string
		registeredRoutes []string
		method         string
		host           string
		path           string
	}{
		{
			name:             "Wrong method",
			registeredRoutes: []string{"GET /users/{id}"},
			method:           "POST",
			host:             "",
			path:             "/users/123",
		},
		{
			name:             "Wrong path",
			registeredRoutes: []string{"GET /users/{id}"},
			method:           "GET",
			host:             "",
			path:             "/posts/123",
		},
		{
			name:             "Nonexistent route",
			registeredRoutes: []string{"GET /users/{id}"},
			method:           "PUT",
			host:             "",
			path:             "/something/else",
		},
	}

	for _, tt := range mismatchTests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()

			for _, route := range tt.registeredRoutes {
				if err := r.Handle(route, "handler"); err != nil {
					t.Fatalf("Failed to register route %s: %v", route, err)
				}
			}

			_, found := r.Match(tt.method, tt.host, tt.path)
			if found {
				t.Errorf("Expected no match for %s %s%s", tt.method, tt.host, tt.path)
			}
		})
	}
}

func TestRouteRegistrationErrors(t *testing.T) {
	r := New()

	// Test nil handler
	err := r.Handle("GET /test", nil)
	if err == nil {
		t.Error("Expected error for nil handler")
	}

	// Test invalid pattern (missing slash)
	err = r.Handle("INVALID", "handler")
	if err == nil {
		t.Error("Expected error for invalid pattern")
	}
}