package configs

import "context"

// Bootstrap 仅以本地文件为入口初始化配置，随后按文件中的 etcd/consul 配置进行增量合并并挂载热更新。
// 返回最终 Config 以及可动态读取/更新的 Store。
func Bootstrap(path string) (Config, *Store) {
	base := NewCfg(path)
	merged := base
	ctx := context.Background()

	// 增量合并：etcd
	if len(base.Etcd.Endpoints) > 0 && base.Etcd.Key != "" {
		etcdSource := EtcdSource{
			Endpoints:          base.Etcd.Endpoints,
			Key:                base.Etcd.Key,
			Format:             base.Etcd.Format,
			Username:           base.Etcd.Username,
			Password:           base.Etcd.Password,
			DialTimeoutSeconds: base.Etcd.DialTimeoutSeconds,
		}
		if etcdCfg, err := etcdSource.Load(ctx); err == nil {
			merged = Merge(merged, etcdCfg)
		}
	}
	// 增量合并：consul
	if base.Consul.Address != "" && base.Consul.Key != "" {
		consulSource := ConsulSource{
			Address: base.Consul.Address,
			Token:   base.Consul.Token,
			Key:     base.Consul.Key,
			Format:  base.Consul.Format,
		}
		if consulCfg, err := consulSource.Load(ctx); err == nil {
			merged = Merge(merged, consulCfg)
		}
	}

	store := NewStore(merged)
	// 文件热更新（作为默认入口）
	_ = FileSource{Path: path}.Watch(ctx, func(nc Config) {
		store.Update(Merge(store.Current(), nc))
	})
	// etcd 热更新（如果配置了）
	if len(base.Etcd.Endpoints) > 0 && base.Etcd.Key != "" {
		etcdSource := EtcdSource{
			Endpoints:          base.Etcd.Endpoints,
			Key:                base.Etcd.Key,
			Format:             base.Etcd.Format,
			Username:           base.Etcd.Username,
			Password:           base.Etcd.Password,
			DialTimeoutSeconds: base.Etcd.DialTimeoutSeconds,
		}
		_ = etcdSource.Watch(ctx, func(nc Config) {
			store.Update(Merge(store.Current(), nc))
		})
	}
	// consul 热更新（如果配置了）
	if base.Consul.Address != "" && base.Consul.Key != "" {
		consulSource := ConsulSource{
			Address: base.Consul.Address,
			Token:   base.Consul.Token,
			Key:     base.Consul.Key,
			Format:  base.Consul.Format,
		}
		_ = consulSource.Watch(ctx, func(nc Config) {
			store.Update(Merge(store.Current(), nc))
		})
	}

	return merged, store
}
