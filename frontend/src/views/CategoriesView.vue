<template>
  <div class="page-layout">
    <AppNavbar />
    <main class="page-main">

      <PageHeader title="Categorias" subtitle="Organiza tus tareas por categoria">
        <template #actions>
          <BaseButton @click="openCreate">+ Nueva categoria</BaseButton>
        </template>
      </PageHeader>

      <!-- Loading -->
      <div v-if="loading" class="loading-row">Cargando categorias...</div>

      <!-- Empty -->
      <EmptyState
        v-else-if="categories.length === 0"
        icon="🏷️"
        title="Sin categorias"
        description="Crea categorias para organizar mejor tus tareas."
      >
        <template #action>
          <BaseButton @click="openCreate">Crear primera categoria</BaseButton>
        </template>
      </EmptyState>

      <!-- Grid -->
      <div v-else class="cat-grid">
        <div v-for="cat in categories" :key="cat.id" class="cat-card">
          <div class="cat-card__swatch" :style="{ background: cat.color }" />
          <div class="cat-card__body">
            <span class="cat-card__name">{{ cat.name }}</span>
            <BaseBadge :color="cat.color">{{ cat.color }}</BaseBadge>
          </div>
          <div class="cat-card__actions">
            <BaseButton variant="ghost" size="sm" @click="openEdit(cat)">Editar</BaseButton>
            <BaseButton variant="danger" size="sm" @click="handleDelete(cat.id)">Eliminar</BaseButton>
          </div>
        </div>
      </div>

    </main>

    <!-- Modal crear / editar -->
    <BaseModal :open="showModal" :title="editing ? 'Editar categoria' : 'Nueva categoria'" @close="closeModal">
      <BaseAlert :message="modalError" type="error" />
      <BaseInput v-model="form.name"  label="Nombre" placeholder="Ej: Trabajo" required />
      <div class="color-row">
        <label class="color-label">Color</label>
        <div class="color-options">
          <button
            v-for="c in palette"
            :key="c"
            class="color-swatch"
            :class="{ selected: form.color === c }"
            :style="{ background: c }"
            :title="c"
            type="button"
            @click="form.color = c"
          />
          <input type="color" v-model="form.color" class="color-input" title="Color personalizado" />
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="closeModal">Cancelar</BaseButton>
        <BaseButton :loading="saving" @click="handleSave">
          {{ editing ? 'Guardar cambios' : 'Crear categoria' }}
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
import BaseModal    from '@/components/common/BaseModal.vue'
import BaseBadge    from '@/components/common/BaseBadge.vue'
import EmptyState   from '@/components/common/EmptyState.vue'
import { categoriesApi } from '@/api/categories'

const categories = ref([])
const loading    = ref(true)
const showModal  = ref(false)
const editing    = ref(null)  // category object when editing, null when creating
const saving     = ref(false)
const modalError = ref('')
const form       = ref({ name: '', color: '#007A6E' })

const palette = [
  '#007A6E', '#3B82F6', '#10B981', '#EF4444',
  '#8B5CF6', '#F59E0B', '#EC4899', '#14B8A6',
]

onMounted(fetchCategories)

async function fetchCategories() {
  loading.value = true
  try {
    const res = await categoriesApi.getAll()
    categories.value = res.data ?? []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value    = null
  form.value       = { name: '', color: '#007A6E' }
  modalError.value = ''
  showModal.value  = true
}

function openEdit(cat) {
  editing.value    = cat
  form.value       = { name: cat.name, color: cat.color }
  modalError.value = ''
  showModal.value  = true
}

function closeModal() {
  showModal.value = false
}

async function handleSave() {
  modalError.value = ''
  if (!form.value.name.trim()) { modalError.value = 'El nombre es requerido.'; return }
  saving.value = true
  try {
    if (editing.value) {
      const res = await categoriesApi.update(editing.value.id, form.value)
      const idx = categories.value.findIndex(c => c.id === editing.value.id)
      if (idx !== -1) categories.value[idx] = res.data
    } else {
      const res = await categoriesApi.create(form.value)
      categories.value.push(res.data)
    }
    closeModal()
  } catch (err) {
    modalError.value = err.response?.data?.error ?? 'Error al guardar.'
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  try {
    await categoriesApi.remove(id)
    categories.value = categories.value.filter(c => c.id !== id)
  } catch {
    // silencioso
  }
}
</script>

<style scoped>
.loading-row { color: var(--text-muted); padding: 32px 0; }

.cat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.cat-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--t), transform var(--t);
}
.cat-card:hover { box-shadow: var(--shadow); transform: translateY(-2px); }

.cat-card__swatch { height: 6px; }

.cat-card__body {
  padding: 16px 16px 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cat-card__name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
}

.cat-card__actions {
  padding: 10px 16px 14px;
  display: flex;
  gap: 6px;
}

/* Color picker */
.color-row   { display: flex; flex-direction: column; gap: 8px; }
.color-label { font-size: 13px; font-weight: 500; color: var(--text); }
.color-options {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.color-swatch {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: transform var(--t), border-color var(--t);
}
.color-swatch:hover   { transform: scale(1.15); }
.color-swatch.selected { border-color: var(--text); transform: scale(1.15); }

.color-input {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  padding: 0;
  background: none;
}
</style>
