package wechat

const (
	EventEncryptModePlainText  eventEncryptMode = iota // 明文模式
	EventEncryptModeCompatible                         // 兼容模式
	EventEncryptModeSafe                               // 安全模式
)

type Event struct {
	Token          string           // 必须为英文或数字，长度为3-32字符。
	EncodingAESKey string           // 消息加密密钥由43位字符组成，可随机修改，字符范围为A-Z，a-z，0-9。
	EncryptMode    eventEncryptMode // 0：明文模式；1：兼容模式；2：安全模式；（当前只支持明文模式）
	w              *Wechat
}

type eventEncryptMode uint8

func WithEvent(token, aesKey string) Option {
	return func(w *Wechat) {
		w.e = &Event{
			Token:          token,
			EncodingAESKey: aesKey,
			w:              w,
		}
	}
}
