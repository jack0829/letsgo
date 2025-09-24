package log

// Option 日志选项
type Option func(c *Config)

// Level 日志级别
func Level(v string) Option {
	return func(c *Config) {
		c.Level = v
	}
}

// Path 日志路径
func Path(path ...string) Option {
	return func(c *Config) {
		l := len(path)
		if l > 0 {
			c.Dir = path[0]
		}
		if l > 1 {
			c.FileName = path[1:]
		}
	}
}

// TimeFormat 日志时间格式（不设置此选项为 ts 时间戳）
func TimeFormat(layout string) Option {
	return func(c *Config) {
		c.TimeFormat = layout
	}
}
