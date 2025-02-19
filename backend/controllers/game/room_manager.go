package game

var GlobalRoomManager = NewRoomManager()

// NewRoomManager 創建新的房間管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

// GetRoom 獲取房間
func (rm *RoomManager) GetRoom(id string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	return rm.rooms[id]
}

// AddRoom 添加房間
func (rm *RoomManager) AddRoom(room *Room) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.rooms[room.ID] = room
}

// RemoveRoom 移除房間
func (rm *RoomManager) RemoveRoom(id string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, id)
}

// GetAllRooms 獲取所有房間
func (rm *RoomManager) GetAllRooms() []*Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rooms := make([]*Room, 0, len(rm.rooms))
	for _, room := range rm.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

// CleanEmptyRooms 清理空房間
func (rm *RoomManager) CleanEmptyRooms() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	for id, room := range rm.rooms {
		room.mu.Lock()
		if len(room.Players) == 0 {
			delete(rm.rooms, id)
		}
		room.mu.Unlock()
	}
}

// GetRoomCount 獲取房間數量
func (rm *RoomManager) GetRoomCount() int {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	return len(rm.rooms)
}
