import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import HostView from '../views/HostView.vue'
import JoinView from '../views/JoinView.vue'
import GameView from '../views/GameView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/host',
      name: 'host',
      component: HostView
    },
    {
      path: '/room/:roomId/join',
      name: 'join',
      component: JoinView,
      props: true 
    },
    {
      path: '/room/:roomId/game',
      name: 'game',
      component: GameView,
      props: true
    }
  ]
})

export default router 