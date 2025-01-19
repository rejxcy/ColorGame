import { createRouter, createWebHistory } from 'vue-router'
import gameView from '../views/gameView.vue'

const routes = [
  {
    path: '/',
    name: 'game',
    component: gameView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 