import { createRouter, createWebHistory } from 'vue-router'
import LoginView      from '@/views/LoginView.vue'
import RegisterView   from '@/views/RegisterView.vue'
import TasksView      from '@/views/TasksView.vue'
import NotesView      from '@/views/NotesView.vue'
import CategoriesView from '@/views/CategoriesView.vue'

const routes = [
  { path: '/',            redirect: '/tasks' },
  { path: '/login',       component: LoginView },
  { path: '/register',    component: RegisterView },
  { path: '/tasks',       component: TasksView,      meta: { requiresAuth: true } },
  { path: '/notes',       component: NotesView,      meta: { requiresAuth: true } },
  { path: '/categories',  component: CategoriesView, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if ((to.path === '/login' || to.path === '/register') && token) {
    next('/tasks')
  } else {
    next()
  }
})

export default router
