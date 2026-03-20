package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"grade-api/config"
	"grade-api/handlers"
	"grade-api/middleware"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setupDB(t *testing.T) {
	t.Helper()
	os.Setenv("DB_PATH", ":memory:")
	os.Setenv("JWT_SECRET", "test-secret")
	config.InitDB()
}

func TestLogin_Success(t *testing.T) {
	setupDB(t)
	r := setupRouter()
	r.POST("/login", handlers.Login)

	body := `{"email":"admin@test.com","password":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Fatal("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	setupDB(t)
	r := setupRouter()
	r.POST("/login", handlers.Login)

	body := `{"email":"admin@test.com","password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	setupDB(t)
	r := setupRouter()
	r.GET("/protected", middleware.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetPerformance_StudentCannotViewOthers(t *testing.T) {
	setupDB(t)

	// Login as student
	loginRouter := setupRouter()
	loginRouter.POST("/login", handlers.Login)
	loginBody := `{"email":"student@test.com","password":"123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	lw := httptest.NewRecorder()
	loginRouter.ServeHTTP(lw, loginReq)

	var loginResp map[string]string
	json.Unmarshal(lw.Body.Bytes(), &loginResp)
	token := loginResp["token"]

	r := setupRouter()
	r.GET("/students/:id/performance",
		middleware.Auth(),
		middleware.RequireRole("admin", "teacher", "student"),
		handlers.GetPerformance,
	)

	// Student tries to view admin's performance
	req := httptest.NewRequest(http.MethodGet, "/students/1/performance", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", w.Code, w.Body.String())
	}
}