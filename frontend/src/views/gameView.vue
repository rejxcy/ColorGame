<template>
  <div class="game-container">
    <h1>按照文字選擇顏色</h1>
    <p class="timer">{{ formatTime(elapsedTime) }} 秒</p>
    <p class="question" :style="{ color: currentColor }">{{ currentQuiz }}</p>
    
    <div v-if="isGameEnd" class="game-end">
      <div class="alert">遊戲結束！</div>
      <button class="restart-button" @click="handleRestart">重新開始</button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { gameApi } from '@/api/game'

export default {
  name: 'GameView',
  setup() {
    const elapsedTime = ref(0)
    const currentQuiz = ref('')
    const currentColor = ref('')
    const isGameEnd = ref(false)
    let timer = null

    const startTimer = () => {
      timer = setInterval(() => {
        elapsedTime.value += 0.1
      }, 100)
    }

    const stopTimer = () => {
      if (timer) {
        clearInterval(timer)
        timer = null
      }
    }

    const formatTime = (time) => time.toFixed(1)

    const handleRestart = async () => {
      try {
        const { data } = await gameApi.restart()
        currentQuiz.value = data.quiz
        currentColor.value = data.displayColor
        elapsedTime.value = 0
        isGameEnd.value = false
        startTimer()
      } catch (error) {
        console.error('重新開始失敗:', error)
      }
    }

    onMounted(async () => {
      try {
        const { data } = await gameApi.newGame()
        currentQuiz.value = data.quiz
        currentColor.value = data.displayColor
        startTimer()
      } catch (error) {
        console.error('開始遊戲失敗:', error)
      }
    })

    onUnmounted(() => {
      stopTimer()
    })

    return {
      elapsedTime,
      currentQuiz,
      currentColor,
      isGameEnd,
      formatTime,
      handleRestart
    }
  }
}
</script>

<style scoped>
.game-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.timer {
  font-size: 20px;
}

.question {
  font-size: 24px;
  font-weight: bold;
}

.restart-button {
  margin-top: 20px;
  padding: 10px 20px;
  font-size: 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
</style>