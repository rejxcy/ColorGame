<template>
  <div class="game-container">
    <!-- 等待遊戲開始 -->
    <div v-if="!gameStarted" class="waiting-screen">
      <h2>等待遊戲開始</h2>
      <div class="ready-section">
        <button 
          class="ready-button"
          :class="{ 'is-ready': isReady }"
          @click="toggleReady"
          :disabled="isConnecting">
          {{ isReady ? '取消準備' : '準備' }}
        </button>
      </div>
      <div class="player-list">
        <h3>玩家列表</h3>
        <div v-for="player in players" :key="player.id" class="player-item">
          <span class="player-name">{{ player.name }}</span>
          <span class="player-status" :class="{ ready: player.isReady }">
            {{ player.isReady ? '已準備' : '未準備' }}
          </span>
        </div>
      </div>
    </div>

    <!-- 遊戲進行中 -->
    <template v-else>
      <div class="game-info">
        <div class="score">得分：{{ currentPlayer?.score || 0 }}</div>
        <div class="progress">進度：{{ gameState.progress }}/{{ gameState.totalQuiz }}</div>
        <div class="wrong-count">錯誤：{{ gameState.wrongCount }}</div>
      </div>
      
      <div class="quiz-container">
        <p v-if="gameState.quiz" 
           class="quiz" 
           :style="{ color: gameState.displayColor }">
          {{ gameState.quiz }}
        </p>
      </div>
      
      <div class="color-grid">
        <button
          v-for="color in validColors"
          :key="color"
          :class="['color-button', color]"
          @click="handleAnswer(color)"
          :disabled="gameState.isFinished"
        />
      </div>

      <!-- 遊戲結束 -->
      <div v-if="gameState.isFinished" class="game-over">
        <h2>遊戲結束！</h2>
        <div class="ranking-list">
          <div v-for="rank in rankings" :key="rank.id" class="ranking-item">
            <span class="rank-number">{{ rank.rank }}</span>
            <span class="player-name">{{ rank.name }}</span>
            <span class="player-score">得分：{{ rank.score }}</span>
            <span class="player-time">用時：{{ formatDuration(rank.duration) }}</span>
          </div>
        </div>
        <button v-if="isHost" class="restart-button" @click="handleRestart">
          重新開始
        </button>
      </div>
    </template>

    <div v-if="isConnecting" class="connecting-overlay">
      正在重新連接...
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWebSocket } from '../composables/useWebSocket'

const route = useRoute()
const router = useRouter()
const ws = useWebSocket()

// 狀態
const gameStarted = ref(false)
const isReady = ref(false)
const isConnecting = ref(false)
const players = ref([])
const gameState = ref({
  quiz: '',
  displayColor: '',
  progress: 0,
  percentage: 0,
  wrongCount: 0,
  isFinished: false,
  totalQuiz: 10
})
const rankings = ref([])

// 常量
const validColors = ['red', 'green', 'blue', 'yellow', 'orange', 'purple']

// 計算屬性
const currentPlayer = computed(() => 
  players.value.find(p => p.id === route.query.playerId)
)

const isHost = computed(() => currentPlayer.value?.isHost || false)

// 確保 WebSocket 連接
const ensureConnection = async () => {
  if (!ws.isConnected.value && !isConnecting.value) {
    isConnecting.value = true
    try {
      const playerName = localStorage.getItem('playerName')
      const roomId = route.params.roomId
      
      if (!playerName || !roomId) {
        console.error('Missing player info')
        router.push('/')
        return false
      }
      
      await ws.connect(roomId, playerName, false)
      return true
    } catch (err) {
      console.error('重新連接失敗:', err)
      return false
    } finally {
      isConnecting.value = false
    }
  }
  return ws.isConnected.value
}

// 方法
const handleAnswer = async (color) => {
  ws.send({
    type: 'answer',
    payload: color
  })
}

const toggleReady = async () => {
  try {
    if (!await ensureConnection()) {
      console.error('無法發送準備狀態：WebSocket 未連接')
      return
    }
    
    const newReadyState = !isReady.value
    console.log('Sending ready state:', newReadyState)
    
    ws.send({
      type: 'ready',
      payload: newReadyState
    })
    
    // 立即更新本地狀態
    isReady.value = newReadyState
    console.log('Local ready state updated to:', isReady.value)
  } catch (err) {
    console.error('發送準備狀態失敗:', err)
  }
}

const handleRestart = () => {
  ws.send({
    type: 'restart'
  })
}

const formatDuration = (ms) => {
  return (ms / 1000).toFixed(2) + '秒'
}

const handleWebSocketMessage = (data) => {
  console.log('Game received message:', data)
  switch (data.type) {
    case 'player_list':
      console.log('Updating player list:', data.payload)
      players.value = data.payload
      // 使用 playerId 來確保只更新當前玩家的準備狀態
      const currentPlayerData = data.payload.find(
        p => p.id === route.query.playerId
      )
      if (currentPlayerData) {
        console.log(`Found current player (ID: ${route.query.playerId}), ready state:`, currentPlayerData.isReady)
        // 只有當找到當前玩家時才更新準備狀態
        isReady.value = currentPlayerData.isReady
      } else {
        console.log('Current player not found in updated list')
      }
      break
    case 'game_start':
      gameStarted.value = true
      break
    case 'game_state':
      gameState.value = data.payload
      break
    case 'rankings':
      rankings.value = data.payload
      break
    case 'error':
      console.error('收到錯誤消息:', data.payload)
      break
  }
}

onMounted(async () => {
  await ensureConnection()
  ws.on(handleWebSocketMessage)
})

onUnmounted(() => {
  ws.off(handleWebSocketMessage)
})
</script>

<style scoped>
.game-container {
  position: relative;
  max-width: 800px;
  margin: 40px auto;
  padding: 20px;
}

.waiting-screen {
  text-align: center;
}

.ready-section {
  margin: 20px 0;
}

.player-list {
  margin: 20px 0;
}

.player-item {
  display: flex;
  justify-content: space-between;
  padding: 10px;
  margin: 5px 0;
  background: #f5f5f5;
  border-radius: 4px;
}

.ready-button {
  padding: 12px 24px;
  font-size: 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  background-color: #4CAF50;
  color: white;
  transition: background-color 0.3s;
}

.ready-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.ready-button.is-ready {
  background-color: #f44336;
}

.game-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 30px;
  font-size: 1.2em;
}

.quiz-container {
  text-align: center;
  margin: 40px 0;
}

.quiz {
  font-size: 3em;
  font-weight: bold;
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin: 40px auto;
  max-width: 500px;
}

.color-button {
  aspect-ratio: 1;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
}

.color-button:hover:not(:disabled) {
  transform: scale(1.05);
}

.color-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.game-over {
  text-align: center;
  margin-top: 40px;
}

.ranking-list {
  margin: 20px 0;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 10px;
  margin: 5px 0;
  background: #f5f5f5;
  border-radius: 5px;
}

.restart-button {
  padding: 12px 24px;
  font-size: 1.2em;
  background: #4CAF50;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.restart-button:hover {
  background: #45a049;
}

/* 顏色按鈕樣式 */
.red { background-color: red; }
.green { background-color: green; }
.blue { background-color: blue; }
.yellow { background-color: yellow; }
.orange { background-color: orange; }
.purple { background-color: purple; }

@media (max-width: 600px) {
  .color-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.connecting-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
</style>