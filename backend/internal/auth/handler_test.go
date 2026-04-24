package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestHandler() *AuthHandler {
	return NewAuthHandler(&Service{jwtSecret: "test-secret"}, false)
}

func performRequest(method, path string, body any, setup func(*gin.Context), handler gin.HandlerFunc) *httptest.ResponseRecorder {
	var reader io.Reader
	switch b := body.(type) {
	case nil:
		reader = nil
	case string:
		reader = strings.NewReader(b)
	default:
		buf := &bytes.Buffer{}
		_ = json.NewEncoder(buf).Encode(b)
		reader = buf
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	if setup != nil {
		setup(c)
	}
	handler(c)
	return w
}

func decodeError(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	msg, _ := body["error"].(string)
	return msg
}

func TestRegister_BindingErrors(t *testing.T) {
	h := newTestHandler()

	tests := []struct {
		name string
		body any
	}{
		{"malformed JSON", "{not json"},
		{"missing email", map[string]string{"username": "alice", "password": "pass"}},
		{"invalid email format", map[string]string{"email": "not-an-email", "username": "alice", "password": "pass"}},
		{"missing username", map[string]string{"email": "a@b.com", "password": "pass"}},
		{"username too short", map[string]string{"email": "a@b.com", "username": "ab", "password": "pass"}},
		{"missing password", map[string]string{"email": "a@b.com", "username": "alice"}},
		{"password too short", map[string]string{"email": "a@b.com", "username": "alice", "password": "x"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := performRequest(http.MethodPost, "/register", tc.body, nil, h.Register)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status: got %d, want %d (body=%s)", w.Code, http.StatusBadRequest, w.Body.String())
			}
			if msg := decodeError(t, w); msg != "invalid request" {
				t.Errorf("error message: got %q, want %q", msg, "invalid request")
			}
		})
	}
}

func TestLogin_BindingErrors(t *testing.T) {
	h := newTestHandler()

	tests := []struct {
		name string
		body any
	}{
		{"malformed JSON", "{"},
		{"missing email", map[string]string{"password": "pass"}},
		{"invalid email format", map[string]string{"email": "not-an-email", "password": "pass"}},
		{"missing password", map[string]string{"email": "a@b.com"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := performRequest(http.MethodPost, "/login", tc.body, nil, h.Login)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
			}
			if msg := decodeError(t, w); msg != "invalid request" {
				t.Errorf("error message: got %q, want %q", msg, "invalid request")
			}
		})
	}
}

func TestLogout_ClearsCookie(t *testing.T) {
	h := newTestHandler()
	w := performRequest(http.MethodPost, "/logout", nil, nil, h.Logout)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	var authCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == authCookieName {
			authCookie = c
			break
		}
	}
	if authCookie == nil {
		t.Fatal("expected auth_token cookie in response, got none")
	}
	if authCookie.Value != "" {
		t.Errorf("cookie value: got %q, want empty string (cleared)", authCookie.Value)
	}
	if authCookie.MaxAge >= 0 {
		t.Errorf("cookie MaxAge: got %d, want negative (to clear cookie)", authCookie.MaxAge)
	}
}

func TestGetMe_Unauthorized(t *testing.T) {
	h := newTestHandler()
	w := performRequest(http.MethodGet, "/me", nil, nil, h.GetMe)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if msg := decodeError(t, w); msg != "unauthorized" {
		t.Errorf("error message: got %q, want %q", msg, "unauthorized")
	}
}

func TestUpdateMe_Unauthorized(t *testing.T) {
	h := newTestHandler()
	w := performRequest(http.MethodPut, "/me", map[string]string{"email": "a@b.com"}, nil, h.UpdateMe)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUpdateMe_BindingErrorWithUserID(t *testing.T) {
	h := newTestHandler()
	setup := func(c *gin.Context) { c.Set("userID", uint(1)) }

	w := performRequest(http.MethodPut, "/me", "{not json", setup, h.UpdateMe)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func runProtected(svc *Service, req *http.Request) (*httptest.ResponseRecorder, uint) {
	var resolved uint
	r := gin.New()
	r.GET("/protected", AuthMiddleware(svc), func(c *gin.Context) {
		if v, ok := c.Get("userID"); ok {
			if id, ok := v.(uint); ok {
				resolved = id
			}
		}
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, resolved
}

func TestAuthMiddleware_RejectsMissingHeaderAndCookie(t *testing.T) {
	svc := newTestService("test-secret")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	w, _ := runProtected(svc, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_RejectsMalformedAuthorizationHeader(t *testing.T) {
	svc := newTestService("test-secret")

	cases := []string{
		"Token abc.def.ghi", // wrong scheme
		"Bearer",            // only one part
		"Basic dXNlcjpwdw",  // different auth scheme
	}
	for _, header := range cases {
		t.Run(header, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", header)

			w, _ := runProtected(svc, req)
			if w.Code != http.StatusUnauthorized {
				t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestAuthMiddleware_RejectsInvalidToken(t *testing.T) {
	svc := newTestService("test-secret")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")

	w, _ := runProtected(svc, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_AcceptsValidBearerToken(t *testing.T) {
	svc := newTestService("test-secret")
	token, err := svc.generateToken(99, "x@y.com")
	if err != nil {
		t.Fatalf("generateToken: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w, uid := runProtected(svc, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d (body=%s)", w.Code, http.StatusOK, w.Body.String())
	}
	if uid != 99 {
		t.Errorf("userID: got %d, want 99", uid)
	}
}

func TestAuthMiddleware_AcceptsValidCookieToken(t *testing.T) {
	svc := newTestService("test-secret")
	token, err := svc.generateToken(7, "x@y.com")
	if err != nil {
		t.Fatalf("generateToken: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: authCookieName, Value: token})

	w, uid := runProtected(svc, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if uid != 7 {
		t.Errorf("userID: got %d, want 7", uid)
	}
}

func TestAuthMiddleware_EmptyCookieFallsThroughToUnauthorized(t *testing.T) {
	svc := newTestService("test-secret")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: authCookieName, Value: ""})

	w, _ := runProtected(svc, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_HeaderTakesPrecedenceOverCookie(t *testing.T) {
	svc := newTestService("test-secret")
	headerToken, err := svc.generateToken(111, "h@h.com")
	if err != nil {
		t.Fatalf("generateToken header: %v", err)
	}
	cookieToken, err := svc.generateToken(222, "c@c.com")
	if err != nil {
		t.Fatalf("generateToken cookie: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+headerToken)
	req.AddCookie(&http.Cookie{Name: authCookieName, Value: cookieToken})

	w, uid := runProtected(svc, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if uid != 111 {
		t.Errorf("userID: got %d, want 111 (header should win over cookie)", uid)
	}
}
