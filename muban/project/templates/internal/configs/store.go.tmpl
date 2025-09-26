package configs

import (
	"sync"
	"sync/atomic"
)

// Store 提供原子读取/更新配置的能力，用于热更新场景
type Store struct {
	v    atomic.Value // holds Config
	mu   sync.RWMutex
	subs map[string][]chan Config // key="*" 表示全量
}

func NewStore(initial Config) *Store {
	s := &Store{subs: make(map[string][]chan Config)}
	s.v.Store(initial)
	return s
}

func (s *Store) Current() Config {
	c, _ := s.v.Load().(Config)
	return c
}

func (s *Store) Update(c Config) {
	s.v.Store(c)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ch := range s.subs["*"] {
		select {
		case ch <- c:
		default:
		}
	}
}

// Subscribe 订阅配置更新；key 目前支持 "*"（全量），预留子树订阅扩展
func (s *Store) Subscribe(key string) <-chan Config {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan Config, 1)
	s.subs[key] = append(s.subs[key], ch)
	return ch
}
