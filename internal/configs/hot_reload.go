package configs

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// HotReloader 配置热重载器
type HotReloader struct {
	watcher    *fsnotify.Watcher
	store      *Store
	callbacks  []func(*Config)
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	configPath string
}

// NewHotReloader 创建配置热重载器
func NewHotReloader(store *Store, configPath string) (*HotReloader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &HotReloader{
		watcher:    watcher,
		store:      store,
		callbacks:  make([]func(*Config), 0),
		ctx:        ctx,
		cancel:     cancel,
		configPath: configPath,
	}, nil
}

// Watch 开始监听配置文件变化
func (hr *HotReloader) Watch(configPath string) error {
	if err := hr.watcher.Add(configPath); err != nil {
		return err
	}

	go hr.watchLoop()
	return nil
}

// AddCallback 添加配置变化回调
func (hr *HotReloader) AddCallback(callback func(*Config)) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	hr.callbacks = append(hr.callbacks, callback)
}

// watchLoop 监听循环
func (hr *HotReloader) watchLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-hr.ctx.Done():
			return
		case event, ok := <-hr.watcher.Events:
			if !ok {
				return
			}
			hr.handleEvent(event)
		case err, ok := <-hr.watcher.Errors:
			if !ok {
				return
			}
			slog.Error("config watcher error", slog.Any("error", err))
		case <-ticker.C:
			// 定期检查，防止事件丢失
		}
	}
}

// handleEvent 处理文件变化事件
func (hr *HotReloader) handleEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Write == fsnotify.Write {
		slog.Info("config file changed, reloading...", slog.String("file", event.Name))

		// 延迟重载，避免文件写入未完成
		time.Sleep(100 * time.Millisecond)

		if err := hr.reloadConfig(); err != nil {
			slog.Error("failed to reload config", slog.Any("error", err))
			return
		}

		slog.Info("config reloaded successfully")
	}
}

// reloadConfig 重新加载配置
func (hr *HotReloader) reloadConfig() error {
	// 重新加载配置
	config := NewCfg(hr.configPath)

	// 更新存储
	hr.store.Update(config)

	// 执行回调
	hr.mu.RLock()
	callbacks := make([]func(*Config), len(hr.callbacks))
	copy(callbacks, hr.callbacks)
	hr.mu.RUnlock()

	for _, callback := range callbacks {
		go func(cb func(*Config)) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("config callback panic", slog.Any("panic", r))
				}
			}()
			cb(&config)
		}(callback)
	}

	return nil
}

// Close 关闭热重载器
func (hr *HotReloader) Close() error {
	hr.cancel()
	return hr.watcher.Close()
}

// ConfigReloadCallback 配置重载回调函数类型
type ConfigReloadCallback func(*Config)

// RegisterConfigReloadCallback 注册配置重载回调
func (s *Store) RegisterConfigReloadCallback(callback ConfigReloadCallback) {
	// 通过订阅机制实现配置重载回调
	ch := s.Subscribe("*")
	go func() {
		for config := range ch {
			callback(&config)
		}
	}()
}
