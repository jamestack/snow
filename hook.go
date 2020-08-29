package snow

// 挂载hook
type HookMount interface {
	OnMount()
}

// 取消挂载
type HookUnMount interface {
	OnUnMount()
}

// 消息处理器
