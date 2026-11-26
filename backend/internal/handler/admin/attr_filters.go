package admin

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// parseAttributeFilters extracts attr[{id}]=value query params.
func parseAttributeFilters(c *gin.Context) map[int64]string {
	out := make(map[int64]string)
	for key, values := range c.Request.URL.Query() {
		if !strings.HasPrefix(key, "attr[") || !strings.HasSuffix(key, "]") {
			continue
		}
		rawID := strings.TrimSuffix(strings.TrimPrefix(key, "attr["), "]")
		id, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		if len(values) == 0 {
			continue
		}
		val := strings.TrimSpace(values[0])
		if val == "" {
			continue
		}
		out[id] = val
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
