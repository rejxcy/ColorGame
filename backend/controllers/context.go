package controllers

import (
	"sync"
)

// Context 提供控制器所需的共享資源和功能
type Context struct {
	GameRooms sync.Map
}

// NewContext 創建新的 Context
func NewContext() *Context {
	return &Context{}
}
