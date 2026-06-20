<template>
  <div class="page-layout">
    <AppNavbar />
    <main class="page-main">

      <PageHeader title="Notas" subtitle="Tus apuntes personales">
        <template #actions>
          <BaseButton @click="openCreate">+ Nueva nota</BaseButton>
        </template>
      </PageHeader>

      <div v-if="loading" class="loading-row">Cargando notas...</div>

      <EmptyState
        v-else-if="notes.length === 0"
        icon="📝"
        title="Sin notas todavia"
        description="Guarda ideas, apuntes o cualquier cosa que no quieras olvidar."
      >
        <template #action>
          <BaseButton @click="openCreate">Crear primera nota</BaseButton>
        </template>
      </EmptyState>

      <div v-else class="notes-grid">
        <div v-for="note in notes" :key="note.id" class="note-card">
          <div class="note-card__header">
            <span class="note-id">#{{ String(note.id).padStart(3, '0') }}</span>
            <div class="note-card__actions">
              <BaseButton variant="ghost" size="sm" @click="openEdit(note)">Editar</BaseButton>
              <BaseButton variant="danger" size="sm" @click="handleDelete(note.id)">Eliminar</BaseButton>
            </div>
          </div>
          <h3 class="note-title">{{ note.title }}</h3>
          <p v-if="note.content" class="note-content">{{ truncate(note.content) }}</p>
          <div class="note-date">{{ formatDate(note.updated_at) }}</div>
        </div>
      </div>

    </main>

    <!-- Modal -->
    <BaseModal :open="showModal" :title="editing ? 'Editar nota' : 'Nueva nota'" @close="closeModal">
      <BaseAlert :message="modalError" type="error" />
      <BaseInput    v-model="form.title"   label="Titulo" placeholder="Titulo de la nota" required />
      <BaseTextarea v-model="form.content" label="Contenido" placeholder="Escribe aqui..." />
      <template #footer>
        <BaseButton variant="ghost" @click="closeModal">Cancelar</BaseButton>
        <BaseButton :loading="saving" @click="handleSave">
          {{ editing ? 'Guardar cambios' : 'Crear nota' }}
        </BaseButton>
      </template>
    </BaseModal>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppNavbar    from '@/components/layout/AppNavbar.vue'
import PageHeader   from '@/components/common/PageHeader.vue'
import BaseButton   from '@/components/common/BaseButton.vue'
import BaseAlert    from '@/components/common/BaseAlert.vue'
import BaseInput    from '@/components/common/BaseInput.vue'
import BaseTextarea from '@/components/common/BaseTextarea.vue'
import BaseModal    from '@/components/common/BaseModal.vue'
import EmptyState   from '@/components/common/EmptyState.vue'
import { notesApi } from '@/api/notes'

const notes      = ref([])
const loading    = ref(true)
const showModal  = ref(false)
const editing    = ref(null)
const saving     = ref(false)
const modalError = ref('')
const form       = ref({ title: '', content: '' })

onMounted(fetchNotes)

async function fetchNotes() {
  loading.value = true
  try {
    const res = await notesApi.getAll()
    notes.value = res.data ?? []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value    = null
  form.value       = { title: '', content: '' }
  modalError.value = ''
  showModal.value  = true
}

function openEdit(note) {
  editing.value    = note
  form.value       = { title: note.title, content: note.content }
  modalError.value = ''
  showModal.value  = true
}

function closeModal() { showModal.value = false }

async function handleSave() {
  modalError.value = ''
  if (!form.value.title.trim()) { modalError.value = 'El titulo es requerido.'; return }
  saving.value = true
  try {
    if (editing.value) {
      const res = await notesApi.update(editing.value.id, form.value)
      const idx = notes.value.findIndex(n => n.id === editing.value.id)
      if (idx !== -1) notes.value[idx] = res.data
    } else {
      const res = await notesApi.create(form.value)
      notes.value.unshift(res.data)
    }
    closeModal()
  } catch (err) {
    modalError.value = err.response?.data?.error ?? 'Error al guardar la nota.'
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  try {
    await notesApi.remove(id)
    notes.value = notes.value.filter(n => n.id !== id)
  } catch {
    // silencioso
  }
}

function truncate(text, max = 120) {
  return text.length > max ? text.slice(0, max) + '...' : text
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('es-PE', { day: '2-digit', month: 'short', year: 'numeric' })
}
</script>

<style scoped>
.loading-row { color: var(--text-muted); padding: 32px 0; }

.notes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}

.note-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--t), transform var(--t);
}
.note-card:hover { box-shadow: var(--shadow); transform: translateY(-2px); }

.note-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.note-card__actions { display: flex; gap: 4px; opacity: 0; transition: opacity var(--t); }
.note-card:hover .note-card__actions { opacity: 1; }

.note-id    { font-family: var(--font-mono); font-size: 11px; color: var(--text-light); }
.note-title { font-size: 15px; font-weight: 600; color: var(--text); line-height: 1.3; }
.note-content { font-size: 13px; color: var(--text-muted); line-height: 1.5; flex: 1; }
.note-date  { font-size: 11px; color: var(--text-light); margin-top: 4px; }
</style>
