import axios from 'axios'

const BASE_URL = 'http://localhost:8080/api'

export const gameApi = {
  // 開始新遊戲
  newGame() {
    return axios.post(`${BASE_URL}/game/new`)
  },

  // 提交答案
  submitAnswer(color) {
    return axios.post(`${BASE_URL}/game/answer`, { color })
  },

  // 獲取遊戲進度
  getProgress() {
    return axios.get(`${BASE_URL}/game/progress`)
  },

  // 重新開始遊戲
  restart() {
    return axios.post(`${BASE_URL}/game/restart`)
  },

  // 設置玩家名稱
  setPlayerName(name) {
    return axios.post(`${BASE_URL}/player/name`, { name })
  }
}