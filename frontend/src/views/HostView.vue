<template>
  <div class="host-container">
    <div class="game-header">
      <h1>顏色配對遊戲</h1>
      <div v-if="roomId" class="room-info">
        房間號碼: {{ roomId }}
      </div>
    </div>

    <div class="main-content">
      <!-- 遊戲未開始時顯示 QR Code 與玩家列表 -->
      <div v-if="gameStatus === 'waiting'" class="waiting-screen">
        <div class="qrcode-section">
          <canvas ref="qrcodeRef"></canvas>
          <p class="scan-hint">掃描 QR Code 加入遊戲</p>
          <p class="room-link">
            或複製連結：
            <span class="link">{{ joinUrl }}</span>
          </p>
        </div>

        <div class="player-list">
          <h2>已加入玩家 ({{ players.length }})</h2>
          <div class="players">
            <div v-for="player in players" :key="player.id" class="player-item">
              <span class="player-name">{{ player.name }}</span>
              <span class="player-status" :class="{ ready: player.isReady }">
                {{ player.isReady ? '已準備' : '未準備' }}
              </span>
            </div>
          </div>
          
          <!-- 開始遊戲按鈕僅在所有玩家準備且至少2人時啟用 -->
          <button 
            v-if="canStartGame" 
            class="start-game-button"
            @click="startGame"
          >
            開始遊戲
          </button>
        </div>
      </div>

      <!-- 遊戲進行中顯示排行榜 -->
      <div v-else-if="gameStatus === 'playing'" class="game-progress">
        <h2>遊戲進行中</h2>
        <div class="ranking-list">
          <div v-for="player in sortedPlayers" :key="player.id" class="ranking-item">
            <span class="rank-number">{{ player.rank }}</span>
            <span class="player-name">{{ player.name }}</span>
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: `${(player.progress / totalQuiz) * 100}%` }"></div>
            </div>
            <span class="player-score">{{ player.score }}</span>
          </div>
        </div>
      </div>

      <div v-else-if="gameStatus === 'finished'" class="game-end">
        <h2>遊戲結束</h2>
        <button @click="gameReset" class="restart-button">重新開始</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import QRCode from 'qrcode'
import { useRouter } from 'vue-router'
import { useWebSocket } from '../composables/useWebSocket'

const router = useRouter()
const roomId = ref('')
const players = ref([])
const totalQuiz = ref(0); // 總題數
// gameStatus 可能值包括："waiting", "playing", "finished"
const gameStatus = ref('waiting')
const qrcodeRef = ref(null)
const ws = useWebSocket()

// 產生房間連結
const joinUrl = computed(() => {
  if (!roomId.value) return ''
  return `${window.location.origin}/room/${roomId.value}/join`
})

// 添加類型檢查，確保即使數據異常也不會報錯
const canStartGame = computed(() => {
  if (!Array.isArray(players.value)) {
    console.warn('players.value is not an array:', players.value)
    return false
  }
  return players.value.length >= 2 && players.value.every(player => player.isReady)
})

const sortedPlayers = computed(() => {
  if (!Array.isArray(players.value)) {
    console.warn('players.value is not an array:', players.value)
    return []
  }
  return [...players.value].sort((a, b) => b.score - a.score)
})

// 生成 QR Code
const generateQRCode = async () => {
  if (!qrcodeRef.value || !joinUrl.value) return
  try {
    await QRCode.toCanvas(qrcodeRef.value, joinUrl.value, {
      width: 256,
      margin: 2,
      color: {
        dark: '#2c3e50',
        light: '#ffffff'
      }
    })
  } catch (err) {
    console.error('生成 QR Code 失敗:', err)
  }
}

// 開始遊戲
const startGame = () => {
  if (players.value.length < 2) {
    alert('至少需要 2 名玩家才能開始遊戲')
    return
  }
  ws.send({
    type: 'game_start'
  })
}

// 重置遊戲
const gameReset = () => {
  ws.send({
    type: 'game_reset'
  })
}

// 處理 WebSocket 傳來的訊息
const handleWebSocketMessage = (data) => {
  switch (data.type) {
    case 'player_list':
      players.value = data.payload
      // 當遊戲狀態為 waiting 時，確保 QR Code 存在
      if (gameStatus.value === 'waiting') {
        generateQRCode()
      }
      break
    case 'game_start':
      gameStatus.value = 'playing'
      break
    case 'game_end':
      gameStatus.value = 'finished'
      players.value = data.payload
      break
    case 'game_state':
      totalQuiz.value = data.payload.totalQuiz
      break
    case 'game_reset':
      gameStatus.value = 'waiting'
      break
  }
}

// 連接 WebSocket 並初始化房間
const connectWebSocket = async () => {
  try {
    if (!roomId.value) {
      roomId.value = Math.random().toString(36).substring(2, 8).toUpperCase()
    }
    await ws.connect(roomId.value, 'Host', true)
    ws.on(handleWebSocketMessage)
    await generateQRCode()
  } catch (err) {
    console.error('WebSocket 連接失敗:', err)
  }
}

onMounted(async () => {
  await connectWebSocket()
})

onUnmounted(() => {
  ws.off(handleWebSocketMessage)
})
</script>

<style scoped>
.host-container {
  max-width: 800px;
  margin: 40px auto;
  padding: 20px;
}

.game-header {
  text-align: center;
  margin-bottom: 40px;
}

.room-info {
  margin-top: 10px;
  font-size: 1.2em;
  color: #666;
}

.main-content {
  display: flex;
  justify-content: center;
}

.waiting-screen {
  display: flex;
  gap: 40px;
}

.qrcode-section {
  text-align: center;
  padding: 20px;
  background: #f5f5f5;
  border-radius: 8px;
}

.qrcode-section canvas {
  margin-bottom: 20px;
  background: white;
  padding: 10px;
  border-radius: 4px;
}

.scan-hint {
  color: #666;
  margin: 10px 0;
}

.room-link {
  color: #666;
  margin-top: 10px;
  word-break: break-all;
}

.link {
  color: #2196F3;
  cursor: pointer;
}

.player-list {
  min-width: 300px;
}

.player-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  margin: 8px 0;
  background: #f5f5f5;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.player-item:hover {
  background: #e9e9e9;
}

.player-name {
  font-weight: 500;
}

.player-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.9em;
  background: #f0f0f0;
}

.player-status.ready {
  background: #4CAF50;
  color: white;
}

.start-game-button {
  margin-top: 20px;
  padding: 12px 24px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1.1em;
  width: 100%;
  transition: background-color 0.3s;
}

.start-game-button:hover {
  background-color: #45a049;
}

.start-game-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.players {
  max-height: 300px;
  overflow-y: auto;
  margin-bottom: 20px;
}

.game-progress {
  width: 100%;
  max-width: 800px;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px;
  margin: 5px 0;
  background: #f5f5f5;
  border-radius: 5px;
}

.rank-number {
  width: 30px;
  text-align: center;
  font-weight: bold;
}

.progress-bar {
  flex-grow: 1;
  height: 20px;
  background: #e0e0e0;
  border-radius: 10px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #4CAF50;
  transition: width 0.3s ease;
}

.player-score {
  min-width: 60px;
  text-align: right;
}

.restart-button {
  padding: 12px 24px;
  font-size: 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  background-color: #4CAF50;
  color: white;
  transition: background-color 0.3s;
}

.restart-button:hover {
  background-color: #45a049;
}

@media (max-width: 768px) {
  .waiting-screen {
    flex-direction: column;
    align-items: center;
  }
  
  .qrcode-section {
    margin-bottom: 30px;
    width: 100%;
    max-width: 300px;
  }
}
</style> 