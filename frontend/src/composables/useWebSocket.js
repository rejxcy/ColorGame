import { ref } from 'vue'

// 使用單例模式來保持連接狀態
let wsInstance = null

export const useWebSocket = () => {
  const ws = ref(null)
  const isConnected = ref(false)
  const messageHandlers = new Set()
  
  // 如果已經有實例，直接返回
  if (wsInstance) {
    return wsInstance
  }

  const connect = (roomId, playerName, isHost) => {
    return new Promise((resolve, reject) => {
      try {
        // 使用當前主機名來建立連線
        const wsUrl = `ws://${window.location.hostname}:8080/api/game/ws`
        const params = new URLSearchParams({
          room_id: roomId,
          player_name: playerName,
          is_host: isHost
        })

        console.log('Connecting to WebSocket:', `${wsUrl}?${params}`)
        ws.value = new WebSocket(`${wsUrl}?${params}`)

        ws.value.onopen = () => {
          console.log('WebSocket connected successfully')
          isConnected.value = true
          resolve()
        }

        ws.value.onerror = (error) => {
          console.error('WebSocket connection error:', error)
          isConnected.value = false
          reject(new Error('WebSocket connection failed'))
        }

        ws.value.onclose = (event) => {
          console.log('WebSocket disconnected:', event.code, event.reason)
          isConnected.value = false
          ws.value = null
        }

        ws.value.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data)
            console.log('Received message:', data)
            messageHandlers.forEach(handler => handler(data))
          } catch (err) {
            console.error('Failed to parse message:', err)
          }
        }
      } catch (err) {
        console.error('Failed to create WebSocket:', err)
        isConnected.value = false
        reject(err)
      }
    })
  }

  const send = (message) => {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected, attempting to reconnect...')
      // 可以在這裡添加重連邏輯
      return
    }
    
    try {
      console.log('Sending message:', message)
      ws.value.send(JSON.stringify(message))
    } catch (err) {
      console.error('Failed to send message:', err)
      isConnected.value = false
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    messageHandlers.clear()
    isConnected.value = false
  }

  const on = (handler) => {
    if (typeof handler === 'function') {
      messageHandlers.add(handler)
    }
  }

  const off = (handler) => {
    messageHandlers.delete(handler)
  }

  // 創建實例
  wsInstance = {
    isConnected,
    connect,
    disconnect,
    send,
    on,
    off
  }

  return wsInstance
}

// 添加 WebSocket 實例的引用
const ws = ref(null) 