package autho

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	mux := &testMux{
		registered: make(map[string]http.Handler),
	}

	providerApp := NewProviderApp(
		"callbackURL",
		"loginURL",
		&noopHandler{},
		&noopHandler{},
	)
	app := NewApp(
		providerApp,
	)
	app.Register(mux)
	mux.validate(t, providerApp)
}

func TestDefaultFailureHandler(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	expectedErr := errors.New("expected error")
	errCtx := ContextWithError(r.Context(), expectedErr)
	PassError(expectedErr, DefaultFailureHandle, w, r.WithContext(errCtx))

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code 400 but got %d", w.Result().StatusCode)
	}
	if w.Body.String() != expectedErr.Error()+"\n" {
		t.Fatalf("expected expectedErr: expected error but got %s", w.Body.String())
	}
}

type noopHandler struct{}

func (n *noopHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

type testMux struct {
	registered map[string]http.Handler
}

func (m *testMux) Handle(path string, handler http.Handler) {
	m.registered[path] = handler
}

func (m *testMux) validate(t *testing.T, r Registerer) {
	t.Helper()

	app := r.(*app)
	if hand, ok := m.registered[app.callbackURL]; hand != app.callbackHandler || !ok {
		t.Fatal("didnt find expected handler or path")
	}
	if hand, ok := m.registered[app.loginURL]; hand != app.loginHandler || !ok {
		t.Fatal("didnt find expected handler or path")
	}
}
