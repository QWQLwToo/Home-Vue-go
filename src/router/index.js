import { createRouter, createWebHistory } from 'vue-router'
import Home from '../components/Home.vue'
import Admin from '../components/Admin.vue'
import Login from '../components/Login.vue'

const routes = [
  {
    path: '/',
    component: Home,
  },
  {
    path: '/admin',
    component: Admin,
    beforeEnter: (to, from, next) => {
      const token = localStorage.getItem('token')
      if (!token) {
        next('/login')
      } else {
        next()
      }
    },
  },
  {
    path: '/login',
    component: Login,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
