<template>
    <div class="player-container">
      <p class="best-time" v-if="bestTime">最佳紀錄 {{ bestTime }} 秒</p>
      <div class="color-grid">
        <div
          v-for="color in colors"
          :key="color"
          :class="['color-box', color]"
          @click="handleColorClick(color)"
        ></div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref } from 'vue'
  import { gameApi } from '@/api/game'
  
  export default {
    name: 'PlayerView',
    setup() {
      const bestTime = ref(null)
      const colors = ['red', 'green', 'blue', 'yellow', 'orange', 'purple']
  
      const handleColorClick = async (color) => {
        try {
          const { data } = await gameApi.submitAnswer(color)
          if (data.isFinished) {
            bestTime.value = data.timeUsed.toFixed(1)
          }
        } catch (error) {
          console.error('提交答案失敗:', error)
        }
      }
  
      return {
        bestTime,
        colors,
        handleColorClick
      }
    }
  }
  </script>
  
  <style scoped>
  .player-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 100vh;
    padding: 20px;
  }
  
  .color-grid {
    display: grid;
    grid-template-columns: repeat(3, 100px);
    gap: 20px;
  }
  
  .color-box {
    width: 100px;
    height: 100px;
    cursor: pointer;
    border: 2px solid #000;
  }
  
  .red { background-color: red; }
  .green { background-color: green; }
  .blue { background-color: blue; }
  .yellow { background-color: yellow; }
  .orange { background-color: orange; }
  .purple { background-color: purple; }
  </style>