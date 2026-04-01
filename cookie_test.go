package lightning

import (
	"testing"
)

func TestCookiesMapSet(t *testing.T) {
	cookies := make(cookiesMap)
	cookies.set("key1", "value1")

	if cookies["key1"] != "value1" {
		t.Errorf("Expected 'value1', got '%s'", cookies["key1"])
	}
}

func TestCookiesMapDel(t *testing.T) {
	cookies := make(cookiesMap)
	cookies.set("key1", "value1")
	cookies.del("key1")

	if _, exists := cookies["key1"]; exists {
		t.Error("Expected key to be deleted")
	}
}

func TestCookiesMapSetMultiple(t *testing.T) {
	cookies := make(cookiesMap)
	cookies.set("key1", "value1")
	cookies.set("key2", "value2")

	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookies, got %d", len(cookies))
	}

	cookies.del("key1")
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie after deletion, got %d", len(cookies))
	}
}
