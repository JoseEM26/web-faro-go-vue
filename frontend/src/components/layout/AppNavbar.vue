<template>
  <header class="navbar">
    <div class="navbar__inner">

      <RouterLink to="/tasks" class="brand">
        <div class="brand-mark">go</div>
        <span class="brand-name">TaskGo</span>
      </RouterLink>

      <nav class="navbar__nav" aria-label="Navegacion principal">
        <RouterLink to="/tasks"      class="nav-link">Tareas</RouterLink>
        <RouterLink to="/notes"      class="nav-link">Notas</RouterLink>
        <RouterLink to="/categories" class="nav-link">Categorias</RouterLink>
      </nav>

      <div class="navbar__user">
        <span class="navbar__email">{{ auth.user?.email }}</span>
        <button class="navbar__logout" @click="handleLogout">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
          Salir
        </button>
      </div>

    </div>
  </header>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth   = useAuthStore()

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  height: var(--navbar-h);
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 40;
  box-shadow: var(--shadow-sm);
}

.navbar__inner {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  padding: 0 24px;
  display: flex;
  align-items: center;
  gap: 32px;
}

/* Nav links */
.navbar__nav {
  display: flex;
  align-items: center;
  gap: 2px;
  flex: 1;
}

.nav-link {
  padding: 5px 12px;
  border-radius: var(--r);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-muted);
  transition: background var(--t), color var(--t);
  text-decoration: none;
}

.nav-link:hover { background: var(--surface-2); color: var(--text); }

.nav-link.router-link-active {
  background: var(--accent-light);
  color: var(--accent);
}

/* User */
.navbar__user {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.navbar__email {
  font-size: 12px;
  color: var(--text-muted);
  font-family: var(--font-mono);
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.navbar__logout {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 11px;
  border-radius: var(--r);
  font-size: 13px;
  font-weight: 500;
  color: var(--text-muted);
  transition: background var(--t), color var(--t);
}
.navbar__logout:hover { background: var(--danger-light); color: var(--danger); }
</style>
