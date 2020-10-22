package foundation

// CacheAccessor 缓存读取器
type CacheAccessor interface {
	Set(key, value string, expiresIn int) error
	Get(key string) (string, error)
	Delete(key string) error
	IncrBy(key string, value int64) (result int64, err error)
	Exists(key string) (bool, error)
}
