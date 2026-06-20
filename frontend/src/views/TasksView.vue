<template>
  <div class="page-layout">
    <AppNavbar />
    <main class="page-main">

      <PageHeader title="Tareas" :subtitle="`${completedCount} de ${tasks.length} completadas`">
        <template #actions>
          <BaseButton @click="openCreate">+ Nueva tarea</BaseButton>
        </template>
      </PageHeader>

      <!-- Stats -->
      <div class="stats-strip">
        <div class="stat-card">
          <div class="stat-num">{{ tasks.length }}</div>
          <div class="stat-label">Total</div>
        </div>
        <div class="stat-card">
          <div class="stat-num accent">{{ completedCount }}</div>
          <div class="stat-label">Completadas</div>
        </div>
        <div class="stat-card">
          <div class="stat-num">{{ pendingCount }}</div>
          <div class="stat-label">Pendientes</div>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="loading-row">Cargando tareas...</div>

      <!-- Empty -->
      <EmptyState
        v-else-if="tasks.length === 0"
        icon="✅"
        title="Sin tareas todavia"
        description="Crea tu primera tarea usando el boton de arriba."
      >
        <template #action>
          <BaseButton @click="openCreate">Crear primera tarea</BaseButton>
        </template>
      </EmptyState>

      <!-- Lista -->
      <div v-else class="task-list">
        <template v-for="task in tasks" :key="task.id">

          <!-- Fila edicion inline -->
          <div v-if="editingId === task.id" class="task-edit">
            <BaseAlert :message="editError" type="error" />
            <div class="task-edit__fields">
              <BaseInput v-model="editForm.title"       label="Titulo"       required />
              <BaseInput v-model="editForm.description" label="Descripcion" />
              <label class="check-row">
                <input type="checkbox" v-model="editForm.completed" />
                Marcar como completada
              </label>
            </div>
            <div class="task-edit__actions">
              <BaseButton :loading="saving" @click="handleSaveEdit(task)">Guardar</BaseButton>
              <BaseButton variant="ghost" @click="cancelEdit">Cancelar</BaseButton>
            </div>
          </div>

          <!-- Fila normal -->
          <div v-else class="task-row" :class="{ 'task-row--done': task.completed }">
            <button
              class="task-check"
              :class="{ 'task-check--done': task.completed }"
              :title="task.completed ? 'Marcar pendiente' : 'Marcar completada'"
              @click="handleToggle(task)"
            />
            <div class="task-body">
              <div class="task-id">#{{ String(task.id).padStart(3, '0') }}</div>
              <div class="task-title" :class="{ 'task-title--done': task.completed }">{{ task.title }}</div>
              <div v-if="task.description" class="task-desc">{{ task.description }}</div>
            </div>
            <div class="task-actions">
              <BaseButton variant="ghost" size="sm" @click="startEdit(task)">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </BaseButton>
              <BaseButton variant="danger" size="sm" @click="handleDelete(task.id)">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/>
                  <path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4h6v2"/>
                </svg>
              </BaseButton>
            </div>
          </div>

        </template>
      </div>

    </main>

    <!-- Modal crear -->
    <BaseModal :open="showCreate" title="Nueva tarea" @close="showCreate = false">
      <BaseAlert :message="createError" type="error" />
      <BaseInput v-model="createForm.title"       label="Titulo"       placeholder="Que hay que hacer?" required />
      <BaseInput v-model="createForm.description" label="Descripcion"  placeholder="Mas detalles..." />
      <template #footer>
        <BaseButton variant="ghost" @click="showCreate = false">Cancelar</BaseButton>
        <BaseButton :loading="creating" @click="handleCreate">Crear tarea</BaseButton>
      </template>
    </BaseModal>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AppNavbar  from '@/components/layout/AppNavbar.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput  from '@/components/common/BaseInput.vue'
import BaseAlert  from '@/components/common/BaseAlert.vue'
import BaseModal  from '@/components/common/BaseModal.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { tasksApi } from '@/api/tasks'

const tasks       = ref([])
const loading     = ref(true)
const showCreate  = ref(false)
const creating    = ref(false)
const createError = ref('')
const createForm  = ref({ title: '', description: '' })

const editingId = ref(null)
const saving    = ref(false)
const editError = ref('')
const editForm  = ref({ title: '', description: '', completed: false })

const completedCount = computed(() => tasks.value.filter(t => t.completed).length)
const pendingCount   = computed(() => tasks.value.filter(t => !t.completed).length)

onMounted(fetchTasks)

async function fetchTasks() {
  loading.value = true
  try {
    const res = await tasksApi.getAll()
    tasks.value = res.data ?? []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  createForm.value  = { title: '', description: '' }
  createError.value = ''
  showCreate.value  = true
}

async function handleCreate() {
  createError.value = ''
  creating.value    = true
  try {
    const res = await tasksApi.create({ title: createForm.value.title, description: createForm.value.description })
    tasks.value.push(res.data)
    showCreate.value = false
    createForm.value = { title: '', description: '' }
  } catch (err) {
    createError.value = err.response?.data?.error ?? 'No se pudo crear la tarea.'
  } finally {
    creating.value = false
  }
}

function startEdit(task) {
  editingId.value            = task.id
  editForm.value.title       = task.title
  editForm.value.description = task.description
  editForm.value.completed   = task.completed
  editError.value            = ''
}

function cancelEdit() { editingId.value = null }

async function handleSaveEdit(task) {
  editError.value = ''
  saving.value    = true
  try {
    const res = await tasksApi.update(task.id, editForm.value)
    const idx = tasks.value.findIndex(t => t.id === task.id)
    if (idx !== -1) tasks.value[idx] = res.data
    cancelEdit()
  } catch (err) {
    editError.value = err.response?.data?.error ?? 'No se pudo actualizar.'
  } finally {
    saving.value = false
  }
}

async function handleToggle(task) {
  try {
    const res = await tasksApi.update(task.id, { title: task.title, description: task.description, completed: !task.completed })
    const idx = tasks.value.findIndex(t => t.id === task.id)
    if (idx !== -1) tasks.value[idx] = res.data
  } catch { /* silencioso */ }
}

async function handleDelete(id) {
  try {
    await tasksApi.remove(id)
    tasks.value = tasks.value.filter(t => t.id !== id)
  } catch { /* silencioso */ }
}
</script>

<style scoped>
.loading-row { color: var(--text-muted); padding: 32px 0; }

.stats-strip {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 28px;
}

.stat-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  padding: 16px 20px;
  box-shadow: var(--shadow-sm);
}

.stat-num {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -1px;
  font-family: var(--font-mono);
  line-height: 1;
  color: var(--text);
}
.stat-num.accent { color: var(--accent); }

.stat-label {
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: .6px;
  margin-top: 5px;
}

.task-list {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.task-row {
  display: flex;
  align-items: flex-start;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
  gap: 12px;
  transition: background var(--t);
}
.task-row:last-child { border-bottom: none; }
.task-row:hover      { background: var(--surface-2); }
.task-row--done      { opacity: .75; }

.task-check {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid var(--border-strong);
  flex-shrink: 0;
  margin-top: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color var(--t), background var(--t);
}
.task-check:hover { border-color: var(--accent); }
.task-check--done { background: var(--accent); border-color: var(--accent); }
.task-check--done::after {
  content: '';
  width: 4px;
  height: 8px;
  border: 2px solid #fff;
  border-top: none;
  border-left: none;
  transform: rotate(45deg) translateY(-1px);
  display: block;
}

.task-body  { flex: 1; min-width: 0; }
.task-id    { font-family: var(--font-mono); font-size: 11px; color: var(--text-light); margin-bottom: 2px; }
.task-title { font-size: 14px; font-weight: 500; color: var(--text); }
.task-title--done {
  text-decoration: line-through;
  text-decoration-color: var(--accent);
  text-decoration-thickness: 2px;
  color: var(--text-muted);
}
.task-desc { font-size: 13px; color: var(--text-muted); margin-top: 2px; }

.task-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity var(--t);
  flex-shrink: 0;
}
.task-row:hover .task-actions { opacity: 1; }

.task-edit {
  padding: 16px;
  border-bottom: 1px solid var(--border);
  background: var(--accent-light);
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.task-edit__fields  { display: flex; flex-direction: column; gap: 10px; }
.task-edit__actions { display: flex; gap: 8px; }

.check-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text);
  cursor: pointer;
}
</style>
