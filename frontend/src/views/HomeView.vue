<template>
  <div class="home-container">
    <h1>顏色配對遊戲</h1>
    <div class="button-group">
      <button class="create-room" @click="createRoom">創建房間</button>
      <div class="join-section">
        <input 
          v-model="roomId" 
          type="text" 
          placeholder="輸入房間號碼"
          @keyup.enter="joinRoom"
        >
        <button class="join-room" @click="joinRoom" :disabled="!roomId">
          加入房間
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const roomId = ref('')

const createRoom = () => {
  router.push({ name: 'host' })
}

const joinRoom = () => {
  if (roomId.value) {
    router.push({
      name: 'join',
      params: { roomId: roomId.value }
    })
  }
}
</script>

<style scoped>
.home-container {
  max-width: 500px;
  margin: 40px auto;
  padding: 20px;
  text-align: center;
}

h1 {
  margin-bottom: 40px;
  color: #2c3e50;
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

button {
  padding: 12px 24px;
  font-size: 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
}

.create-room {
  background-color: #4CAF50;
  color: white;
}

.create-room:hover {
  background-color: #45a049;
}

.join-section {
  display: flex;
  gap: 10px;
}

input {
  flex: 1;
  padding: 12px;
  font-size: 16px;
  border: 2px solid #ddd;
  border-radius: 6px;
  outline: none;
}

input:focus {
  border-color: #4CAF50;
}

.join-room {
  background-color: #2196F3;
  color: white;
}

.join-room:hover {
  background-color: #1e88e5;
}

.join-room:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style> 