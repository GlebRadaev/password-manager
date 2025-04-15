package state

import (
	"sync"
	"testing"
)

func TestAccessToken(t *testing.T) {
	testToken := "test_token_123"
	SetAccessToken(testToken)

	if got := GetAccessToken(); got != testToken {
		t.Errorf("GetAccessToken() = %v, want %v", got, testToken)
	}

	if !IsAuthenticated() {
		t.Error("IsAuthenticated() = false, want true")
	}

	SetAccessToken("")
	if IsAuthenticated() {
		t.Error("IsAuthenticated() = true for empty token, want false")
	}
}

func TestConcurrentAccess(t *testing.T) {
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			SetAccessToken("token_" + string(rune(i)))
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = GetAccessToken()
			_ = IsAuthenticated()
		}()
	}

	wg.Wait()

	if GetAccessToken() == "" {
		t.Error("Expected some token after concurrent access")
	}
}
