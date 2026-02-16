package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------

// newTestRouter builds the same mux + middleware stack used in production so
// that every test exercises the real routing and handler code.
func newTestRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/employees", handleGetEmployees)
	mux.HandleFunc("GET /api/employees/{id}", handleGetEmployeeByID)
	mux.HandleFunc("GET /health", handleHealth)
	return withLogging(withMethodNotAllowed(mux))
}

// doRequest is a convenience wrapper that sends a GET request to the given
// path on a test server and returns the response.  It calls t.Fatal on any
// transport-level error so callers do not need to handle err themselves.
func doRequest(t *testing.T, srv *httptest.Server, method, path string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, srv.URL+path, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("executing request %s %s: %v", method, path, err)
	}
	return resp
}

// readBody reads and returns the full response body, closing it afterward.
func readBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	return body
}

// assertStatusCode fails the test if the response status code does not match
// the expected value.
func assertStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		t.Errorf("expected status %d, got %d", expected, resp.StatusCode)
	}
}

// assertContentType fails the test if the Content-Type header does not match
// the expected value.
func assertContentType(t *testing.T, resp *http.Response, expected string) {
	t.Helper()
	ct := resp.Header.Get("Content-Type")
	if ct != expected {
		t.Errorf("expected Content-Type %q, got %q", expected, ct)
	}
}

// assertJSONArray unmarshals body into a slice of Employee and returns it.
// It fails the test if the JSON is malformed.
func assertJSONArray(t *testing.T, body []byte) []Employee {
	t.Helper()
	var emps []Employee
	if err := json.Unmarshal(body, &emps); err != nil {
		t.Fatalf("body is not a valid JSON array of Employee: %v\nbody: %s", err, string(body))
	}
	return emps
}

// assertJSONObject unmarshals body into a single Employee and returns it.
// It fails the test if the JSON is malformed.
func assertJSONObject(t *testing.T, body []byte) Employee {
	t.Helper()
	var emp Employee
	if err := json.Unmarshal(body, &emp); err != nil {
		t.Fatalf("body is not a valid Employee JSON object: %v\nbody: %s", err, string(body))
	}
	return emp
}

// --------------------------------------------------------------------------
// Tests
// --------------------------------------------------------------------------

// TestGetAllEmployees verifies that GET /api/employees returns HTTP 200, a
// valid JSON array, and a non-empty list of employees.
func TestGetAllEmployees(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/employees")
	body := readBody(t, resp)

	assertStatusCode(t, resp, http.StatusOK)

	emps := assertJSONArray(t, body)
	if len(emps) == 0 {
		t.Fatal("expected non-empty employee list, got empty array")
	}

	// Sanity-check that the count matches the package-level seed data.
	if len(emps) != len(employees) {
		t.Errorf("expected %d employees, got %d", len(employees), len(emps))
	}
}

// TestGetEmployeeByID verifies that GET /api/employees/1 returns HTTP 200
// with a valid Employee JSON object whose ID equals 1.
func TestGetEmployeeByID(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/employees/1")
	body := readBody(t, resp)

	assertStatusCode(t, resp, http.StatusOK)

	emp := assertJSONObject(t, body)
	if emp.ID != 1 {
		t.Errorf("expected employee ID 1, got %d", emp.ID)
	}
}

// TestGetEmployeeByID_AllKnownIDs is a table-driven test that checks every
// seed employee can be fetched by ID and that the returned data matches.
func TestGetEmployeeByID_AllKnownIDs(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	for _, want := range employees {
		t.Run(fmt.Sprintf("ID_%d", want.ID), func(t *testing.T) {
			path := fmt.Sprintf("/api/employees/%d", want.ID)
			resp := doRequest(t, srv, http.MethodGet, path)
			body := readBody(t, resp)

			assertStatusCode(t, resp, http.StatusOK)

			got := assertJSONObject(t, body)
			if got.ID != want.ID {
				t.Errorf("ID: want %d, got %d", want.ID, got.ID)
			}
			if got.FirstName != want.FirstName {
				t.Errorf("FirstName: want %q, got %q", want.FirstName, got.FirstName)
			}
			if got.LastName != want.LastName {
				t.Errorf("LastName: want %q, got %q", want.LastName, got.LastName)
			}
			if got.Email != want.Email {
				t.Errorf("Email: want %q, got %q", want.Email, got.Email)
			}
			if got.Department != want.Department {
				t.Errorf("Department: want %q, got %q", want.Department, got.Department)
			}
			if got.Position != want.Position {
				t.Errorf("Position: want %q, got %q", want.Position, got.Position)
			}
		})
	}
}

// TestGetEmployeeNotFound verifies that fetching an ID that does not exist
// returns HTTP 404 with a JSON error body.
func TestGetEmployeeNotFound(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/employees/999")
	body := readBody(t, resp)

	assertStatusCode(t, resp, http.StatusNotFound)

	var errBody map[string]string
	if err := json.Unmarshal(body, &errBody); err != nil {
		t.Fatalf("expected JSON error body, got parse error: %v\nbody: %s", err, string(body))
	}
	if _, ok := errBody["error"]; !ok {
		t.Errorf("expected error key in response body, got: %v", errBody)
	}
}

// TestHealthEndpoint verifies that GET /health returns HTTP 200 with
// {"status":"ok"}.
func TestHealthEndpoint(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/health")
	body := readBody(t, resp)

	assertStatusCode(t, resp, http.StatusOK)

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("expected JSON body, got parse error: %v\nbody: %s", err, string(body))
	}
	if result["status"] != "ok" {
		t.Errorf("expected status \"ok\", got %q", result["status"])
	}
}

// TestContentTypeHeader uses a table-driven approach to verify that every
// endpoint returns the Content-Type: application/json header.
func TestContentTypeHeader(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	tests := []struct {
		name string
		path string
	}{
		{"GetAllEmployees", "/api/employees"},
		{"GetEmployeeByID", "/api/employees/1"},
		{"GetEmployeeNotFound", "/api/employees/999"},
		{"Health", "/health"},
		{"InvalidID", "/api/employees/abc"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, srv, http.MethodGet, tc.path)
			defer resp.Body.Close()
			// Drain the body so the connection can be reused.
			io.ReadAll(resp.Body)

			assertContentType(t, resp, "application/json")
		})
	}
}

// TestEmployeeFields verifies that the employee returned for a known ID has
// all fields populated (non-zero / non-empty).
func TestEmployeeFields(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/employees/1")
	body := readBody(t, resp)

	assertStatusCode(t, resp, http.StatusOK)

	emp := assertJSONObject(t, body)

	// Check that every field has a non-zero value.
	checks := []struct {
		field string
		value string
	}{
		{"FirstName", emp.FirstName},
		{"LastName", emp.LastName},
		{"Email", emp.Email},
		{"Department", emp.Department},
		{"Position", emp.Position},
	}
	for _, c := range checks {
		if c.value == "" {
			t.Errorf("expected non-empty %s, got empty string", c.field)
		}
	}
	if emp.ID == 0 {
		t.Error("expected non-zero ID, got 0")
	}
}

// TestEmployeeFields_JSONKeys verifies that the raw JSON keys use the
// expected snake_case names defined by the struct tags.
func TestEmployeeFields_JSONKeys(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/employees/1")
	body := readBody(t, resp)

	assertStatusCode(t, resp, http.StatusOK)

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("failed to unmarshal response as generic map: %v", err)
	}

	expectedKeys := []string{"id", "first_name", "last_name", "email", "department", "position"}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q to be present in response, keys found: %v", key, keysOf(raw))
		}
	}
}

// keysOf returns the keys of a map for diagnostic output.
func keysOf(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// TestInvalidEmployeeID verifies that a non-numeric ID in the path results
// in an error response (the implementation returns 404 for invalid IDs).
func TestInvalidEmployeeID(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	tests := []struct {
		name string
		path string
	}{
		{"Alphabetic", "/api/employees/abc"},
		{"SpecialChars", "/api/employees/@!"},
		{"Float", "/api/employees/1.5"},
		{"Negative", "/api/employees/-1"},
		{"Empty", "/api/employees/"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, srv, http.MethodGet, tc.path)
			body := readBody(t, resp)

			// The implementation returns 404 for non-numeric IDs.
			// We accept either 400 or 404 as both are reasonable.
			if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusNotFound {
				t.Errorf("expected status 400 or 404, got %d", resp.StatusCode)
			}

			// Ensure the error body is still valid JSON.
			var errBody map[string]interface{}
			if err := json.Unmarshal(body, &errBody); err != nil {
				// The /api/employees/ (empty id) case may route to the list
				// handler instead, which returns an array.  In that case just
				// verify it is valid JSON at all.
				var anything interface{}
				if jsonErr := json.Unmarshal(body, &anything); jsonErr != nil {
					t.Errorf("response is not valid JSON: %v\nbody: %s", jsonErr, string(body))
				}
			}
		})
	}
}

// --------------------------------------------------------------------------
// httptest.NewRecorder tests (unit-level, no real TCP server)
// --------------------------------------------------------------------------

// TestGetAllEmployees_Recorder tests handleGetEmployees using a ResponseRecorder
// instead of a full test server, demonstrating the recorder approach.
func TestGetAllEmployees_Recorder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/employees", nil)
	rec := httptest.NewRecorder()

	handler := newTestRouter()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var emps []Employee
	if err := json.Unmarshal(rec.Body.Bytes(), &emps); err != nil {
		t.Fatalf("response is not a valid JSON array: %v\nbody: %s", err, rec.Body.String())
	}
	if len(emps) == 0 {
		t.Error("expected non-empty employee list")
	}
}

// TestGetEmployeeByID_Recorder tests handleGetEmployeeByID using a
// ResponseRecorder.
func TestGetEmployeeByID_Recorder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/employees/2", nil)
	rec := httptest.NewRecorder()

	handler := newTestRouter()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	emp := assertJSONObject(t, rec.Body.Bytes())
	if emp.ID != 2 {
		t.Errorf("expected employee ID 2, got %d", emp.ID)
	}
	if emp.FirstName != "Brian" {
		t.Errorf("expected FirstName %q, got %q", "Brian", emp.FirstName)
	}
}

// TestHealthEndpoint_Recorder tests handleHealth using a ResponseRecorder.
func TestHealthEndpoint_Recorder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler := newTestRouter()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var result map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("response is not valid JSON: %v\nbody: %s", err, rec.Body.String())
	}
	if result["status"] != "ok" {
		t.Errorf("expected status \"ok\", got %q", result["status"])
	}
}

// TestGetEmployeeNotFound_Recorder tests the 404 path using a
// ResponseRecorder.
func TestGetEmployeeNotFound_Recorder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/employees/9999", nil)
	rec := httptest.NewRecorder()

	handler := newTestRouter()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	var errBody map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &errBody); err != nil {
		t.Fatalf("response is not valid JSON: %v\nbody: %s", err, rec.Body.String())
	}
	if _, ok := errBody["error"]; !ok {
		t.Error("expected 'error' key in JSON response")
	}
}

// --------------------------------------------------------------------------
// Edge cases and additional coverage
// --------------------------------------------------------------------------

// TestMethodNotAllowed verifies that using POST on a GET-only endpoint
// returns an appropriate error status.
func TestMethodNotAllowed(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	paths := []string{"/api/employees", "/api/employees/1", "/health"}

	for _, method := range methods {
		for _, path := range paths {
			name := fmt.Sprintf("%s_%s", method, path)
			t.Run(name, func(t *testing.T) {
				resp := doRequest(t, srv, method, path)
				defer resp.Body.Close()
				io.ReadAll(resp.Body)

				if resp.StatusCode == http.StatusOK {
					t.Errorf("%s %s should not return 200", method, path)
				}
			})
		}
	}
}

// TestUnknownRoute verifies that hitting an unregistered path returns a
// non-200 status (typically 404).
func TestUnknownRoute(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	paths := []string{
		"/",
		"/api",
		"/api/departments",
		"/nonexistent",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			resp := doRequest(t, srv, http.MethodGet, path)
			defer resp.Body.Close()
			io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				t.Errorf("GET %s should not return 200 for unknown route", path)
			}
		})
	}
}

// TestEmployeesListOrder verifies that the returned employee list preserves
// the order of the seed data and that IDs are sequential.
func TestEmployeesListOrder(t *testing.T) {
	srv := httptest.NewServer(newTestRouter())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/employees")
	body := readBody(t, resp)

	assertStatusCode(t, resp, http.StatusOK)

	emps := assertJSONArray(t, body)

	for i, emp := range emps {
		expected := employees[i]
		if emp.ID != expected.ID {
			t.Errorf("employee[%d]: expected ID %d, got %d", i, expected.ID, emp.ID)
		}
		if emp.Email != expected.Email {
			t.Errorf("employee[%d]: expected Email %q, got %q", i, expected.Email, emp.Email)
		}
	}
}
