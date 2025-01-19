<template>
  <div class="game-container">
    <h1>按照文字選擇顏色</h1>

    <!-- 遊戲未開始時顯示開始按鈕 -->
    <div v-if="!gameStarted" class="start-screen">
      <button class="start-button" @click="startGame">開始遊戲</button>
    </div>

    <!-- 遊戲開始後顯示遊戲內容 -->
    <template v-else>
      <div class="game-info">
        <p class="timer">用時：{{ formatTime(elapsedTime) }} 秒</p>
        <p class="progress">進度：{{ gameState.progress || 0 }}/{{ gameState.totalQuiz || 10 }}</p>
        <p class="wrong-count">錯誤次數：{{ gameState.wrongCount || 0 }}</p>
      </div>
      
      <p v-if="gameState.quiz" 
         class="question" 
         :style="{ color: gameState.displayColor }">
        {{ gameState.quiz }}
      </p>
      
      <div class="color-grid">
        <button
          v-for="color in colors"
          :key="color"
          :class="['color-button', color]"
          @click="handleColorClick(color)"
        />
      </div>

      <div v-if="gameState.isFinished" class="game-end">
        <div class="alert">遊戲結束！</div>
        <p>總用時：{{ formatTime(elapsedTime) }} 秒</p>
        <p>錯誤次數：{{ gameState.wrongCount }}</p>
        <button class="restart-button" @click="handleRestart">重新開始</button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { GameWebSocket } from '@/api/game'

const gameState = ref({})
const gameStarted = ref(false)
const startTime = ref(0)
const elapsedTime = ref(0)
const timerInterval = ref(null)
const colors = ['red', 'green', 'blue', 'yellow', 'orange', 'purple']
const gameWs = new GameWebSocket()

const formatTime = (time) => (time || 0).toFixed(1)

const startTimer = () => {
  startTime.value = Date.now()
  if (timerInterval.value) {
    clearInterval(timerInterval.value)
  }
  timerInterval.value = setInterval(() => {
    elapsedTime.value = (Date.now() - startTime.value) / 1000
  }, 100)
}

const stopTimer = () => {
  if (timerInterval.value) {
    clearInterval(timerInterval.value)
    timerInterval.value = null
  }
}

const startGame = async () => {
  await setupWebSocket()
  gameStarted.value = true
  startTimer()
}

const handleColorClick = (color) => {
  gameWs.sendAnswer(color)
}

const handleRestart = () => {
  gameWs.restart()
  startTime.value = Date.now()
  elapsedTime.value = 0
  stopTimer()
  startTimer()
}

// WebSocket 消息處理
const setupWebSocket = async () => {
  try {
    await gameWs.connect()

    gameWs.on('game_state', (state) => {
      gameState.value = state
      if (state.progress === 0 && !timerInterval.value) {
        startTimer()
      }
    })

    gameWs.on('game_over', (finalState) => {
      gameState.value = finalState
      stopTimer()  // 遊戲結束時停止計時
    })

    gameWs.on('error', (error) => {
      console.error('Game error:', error)
    })
  } catch (error) {
    console.error('WebSocket connection failed:', error)
  }
}

onMounted(() => {
  setupWebSocket()
})

onUnmounted(() => {
  stopTimer()
  gameWs.close()
})
</script>

<style scoped>
.game-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}

.game-info {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.question {
  font-size: 32px;
  font-weight: bold;
  margin: 20px 0;
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 15px;
  margin: 20px 0;
}

.color-button {
  width: 80px;
  height: 80px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.1s;
}

.color-button:hover {
  transform: scale(1.05);
}

.red { background-color: red; }
.green { background-color: green; }
.blue { background-color: blue; }
.yellow { background-color: yellow; }
.orange { background-color: orange; }
.purple { background-color: purple; }

.game-end {
  text-align: center;
  margin-top: 20px;
}

.restart-button {
  padding: 10px 20px;
  font-size: 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-top: 10px;
}

.restart-button:hover {
  background-color: #45a049;
}

.start-screen {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.start-button {
  padding: 15px 30px;
  font-size: 24px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.start-button:hover {
  background-color: #45a049;
}
</style>