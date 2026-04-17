package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// assertRoutesMounted fails the test if any expected "METHOD /path" entry is
// missing from the router's route table. Used to lock in the public contract
// of each register*Routes helper.
func assertRoutesMounted(t *testing.T, router *gin.Engine, expected []string) {
	t.Helper()
	got := make(map[string]bool, len(router.Routes()))
	for _, r := range router.Routes() {
		got[r.Method+" "+r.Path] = true
	}
	for _, want := range expected {
		if !got[want] {
			t.Errorf("route not registered: %s", want)
		}
	}
}
