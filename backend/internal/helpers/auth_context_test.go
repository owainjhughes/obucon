package helpers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}

func TestUserIDFromContext_Present(t *testing.T) {
	c := newTestContext()
	c.Set("userID", uint(42))

	id, ok := UserIDFromContext(c)
	if !ok {
		t.Fatal("expected ok=true when userID is set to a uint")
	}
	if id != 42 {
		t.Errorf("userID: got %d, want 42", id)
	}
}

func TestUserIDFromContext_Absent(t *testing.T) {
	c := newTestContext()

	id, ok := UserIDFromContext(c)
	if ok {
		t.Error("expected ok=false when userID is not set")
	}
	if id != 0 {
		t.Errorf("userID: got %d, want 0 (zero value) when unset", id)
	}
}

func TestUserIDFromContext_WrongType(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"string", "not-a-uint"},
		{"int", int(42)},
		{"int64", int64(42)},
		{"float64", float64(42)},
		{"nil", nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newTestContext()
			c.Set("userID", tc.value)

			id, ok := UserIDFromContext(c)
			if ok {
				t.Errorf("expected ok=false for value of type %T, got ok=true with id=%d", tc.value, id)
			}
			if id != 0 {
				t.Errorf("userID: got %d, want 0 (zero value)", id)
			}
		})
	}
}

func TestUserIDFromContext_ZeroValueUint(t *testing.T) {
	c := newTestContext()
	c.Set("userID", uint(0))

	id, ok := UserIDFromContext(c)
	if !ok {
		t.Error("expected ok=true when userID is uint(0) — a valid stored value")
	}
	if id != 0 {
		t.Errorf("userID: got %d, want 0", id)
	}
}
