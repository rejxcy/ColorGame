<template>
  <div class="join-container">
    <h1>加入遊戲</h1>
    <div class="join-form">
      <div class="input-group">
        <label for="playerName">請輸入您的名字：</label>
        <input 
          id="playerName"
          v-model="playerName"
          type="text"
          placeholder="您的名字"
          :maxlength="10"
          :disabled="isConnecting"
          @keyup.enter="handleJoin"
        >
      </div>
      <button 
        class="join-button"
        @click="handleJoin"
        :disabled="!playerName.trim() || isConnecting">
        {{ isConnecting ? '連接中...' : '加入遊戲' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useWebSocket } from '../composables/useWebSocket' // 假設你有這個 composable

const router = useRouter()
const route = useRoute()
const playerName = ref('')
const ws = useWebSocket()
const isConnecting = ref(false) // 添加連接狀態標記

onMounted(() => {
  // 檢查是否有 roomId
  if (!route.params.roomId) {
    console.error('沒有房間 ID')
    router.push('/')
    return
  }
})

const handleJoin = async () => {
  if (!playerName.value.trim() || isConnecting.value) return
  
  try {
    isConnecting.value = true
    await ws.connect(route.params.roomId, playerName.value, false)
    
    // 保存玩家信息
    localStorage.setItem('playerName', playerName.value)
    
    router.push({
      name: 'game',
      params: { 
        roomId: route.params.roomId 
      }
    })
  } catch (err) {
    console.error('加入遊戲失敗:', err)
    alert('無法連接到遊戲服務器，請稍後再試')
  } finally {
    isConnecting.value = false
  }
}
</script>

<style scoped>
.join-container {
  max-width: 400px;
  margin: 40px auto;
  padding: 20px;
  text-align: center;
}

.join-form {
  margin-top: 30px;
}

.input-group {
  margin-bottom: 20px;
  text-align: left;
}

label {
  display: block;
  margin-bottom: 8px;
  color: #333;
}

input {
  width: 100%;
  padding: 12px;
  font-size: 16px;
  border: 2px solid #ddd;
  border-radius: 6px;
  outline: none;
  transition: border-color 0.3s;
}

input:focus {
  border-color: #4CAF50;
}

.join-button {
  width: 100%;
  padding: 12px;
  font-size: 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.join-button:hover {
  background-color: #45a049;
}

.join-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style> 