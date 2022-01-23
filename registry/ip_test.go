package registry

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	t.Log(GetLocalIP())
}
