package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterHeroSubResource_RegistersCorrectRoutes(t *testing.T) {
	router := gin.New()
	group := router.Group("/heroes")

	noop := func(c *gin.Context) { c.Status(200) }
	registerHeroSubResource(group, "attributes", noop, noop, noop)

	expected := map[string]string{
		"GET":    "/heroes/:id/attributes",
		"POST":   "/heroes/:id/attributes",
		"DELETE": "/heroes/:id/attributes/:subId",
	}

	found := make(map[string]bool)
	for _, r := range router.Routes() {
		if path, ok := expected[r.Method]; ok && r.Path == path {
			found[r.Method] = true
		}
	}

	for method, path := range expected {
		if !found[method] {
			t.Errorf("route %s %s not registered", method, path)
		}
	}
}

func TestRegisterHeroSubResource_MultipleResources(t *testing.T) {
	router := gin.New()
	group := router.Group("/heroes")

	noop := func(c *gin.Context) { c.Status(200) }
	registerHeroSubResource(group, "skills", noop, noop, noop)
	registerHeroSubResource(group, "talents", noop, noop, noop)

	routes := router.Routes()
	routeSet := make(map[string]bool)
	for _, r := range routes {
		routeSet[r.Method+" "+r.Path] = true
	}

	expectedRoutes := []string{
		"GET /heroes/:id/skills",
		"POST /heroes/:id/skills",
		"DELETE /heroes/:id/skills/:subId",
		"GET /heroes/:id/talents",
		"POST /heroes/:id/talents",
		"DELETE /heroes/:id/talents/:subId",
	}

	for _, route := range expectedRoutes {
		if !routeSet[route] {
			t.Errorf("route %s not registered", route)
		}
	}
}
