package snow

import "reflect"

// 挂载hook
type HookMount interface {
	OnMount()
}

// 取消挂载
type HookUnMount interface {
	OnUnMount()
}

// 自定义消息处理器
type HookCall interface {
	OnCall(name string, call func([]reflect.Value) []reflect.Value, args []reflect.Value) []reflect.Value
}
