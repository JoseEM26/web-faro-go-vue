<template>
  <div class="auth-layout">
    <div class="auth-card">
      <div class="auth-brand">
        <RouterLink to="/login" class="brand">
          <div class="brand-mark">go</div>
          <span class="brand-name">TaskGo</span>
          <span class="brand-tag">v1</span>
        </RouterLink>
      </div>

      <h1 class="auth-title">Iniciar sesion</h1>
      <p class="auth-sub">Accede a tus tareas y notas</p>

      <BaseAlert :message="error" type="error" style="margin-bottom:16px" />

      <form @submit.prevent="handleLogin" style="display:flex;flex-direction:column;gap:14px">
        <BaseInput
          v-model="email"
          label="Email"
          type="email"
          placeholder="tu@email.com"
          autocomplete="email"
          required
        />
        <BaseInput
          v-model="password"
          label="Contrasena"
          type="password"
          placeholder="••••••••"
          autocomplete="current-password"
          required
        />
        <BaseButton type="submit" :loading="loading" :full="true" style="margin-top:4px">
          Iniciar sesion
        </BaseButton>
      </form>

      <div class="auth-footer">
        Sin cuenta? <RouterLink to="/register">Registrate</RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import BaseInput  from '@/components/common/BaseInput.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseAlert  from '@/components/common/BaseAlert.vue'

const router = useRouter()
const auth   = useAuthStore()

const email    = ref('')
const password = ref('')
const loading  = ref(false)
const error    = ref('')

async function handleLogin() {
  error.value   = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    router.push('/tasks')
  } catch (err) {
    error.value = err.response?.data?.error ?? 'Error al iniciar sesion. Intenta de nuevo.'
  } finally {
    loading.value = false
  }
}
</script>
