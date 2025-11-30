package tokens

import "sync"

type Tokens struct {
	mu           sync.RWMutex
	accessToken  string
	refreshToken string
}

func (t *Tokens) GetAccess() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.accessToken
}

func (t *Tokens) GetRefresh() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.refreshToken
}

func (t *Tokens) SetAccess(access string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.accessToken = access
}

func (t *Tokens) SetRefresh(refresh string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.refreshToken = refresh
}
