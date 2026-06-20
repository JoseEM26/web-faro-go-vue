<template>
  <div class="input-wrap">
    <label v-if="label" :for="uid" class="input-label">
      {{ label }}
      <span v-if="hint" class="input-hint">{{ hint }}</span>
    </label>
    <input
      :id="uid"
      v-bind="$attrs"
      :value="modelValue"
      :class="['input-field', { 'input-field--error': error }]"
      @input="$emit('update:modelValue', $event.target.value)"
    />
    <span v-if="error" class="input-error">{{ error }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  label:      { type: String, default: '' },
  hint:       { type: String, default: '' },
  error:      { type: String, default: '' },
})

defineEmits(['update:modelValue'])

let _n = 0
const uid = computed(() => `input-${++_n}`)
</script>

<style scoped>
.input-wrap { display: flex; flex-direction: column; gap: 5px; }

.input-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
  display: flex;
  align-items: center;
  gap: 6px;
}

.input-hint { font-weight: 400; color: var(--text-light); }

.input-field {
  width: 100%;
  padding: 9px 12px;
  border: 1px solid var(--border);
  border-radius: var(--r);
  font-size: 14px;
  color: var(--text);
  background: var(--surface);
  outline: none;
  transition: border-color var(--t), box-shadow var(--t);
}

.input-field:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(0,122,110,.12);
}

.input-field--error {
  border-color: var(--danger);
}

.input-field--error:focus {
  box-shadow: 0 0 0 3px rgba(220,38,38,.12);
}

.input-error { font-size: 12px; color: var(--danger); }
</style>
