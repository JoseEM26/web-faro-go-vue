<template>
  <div class="textarea-wrap">
    <label v-if="label" :for="uid" class="textarea-label">{{ label }}</label>
    <textarea
      :id="uid"
      v-bind="$attrs"
      :value="modelValue"
      class="textarea-field"
      @input="$emit('update:modelValue', $event.target.value)"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'

defineProps({
  modelValue: { type: String, default: '' },
  label:      { type: String, default: '' },
})

defineEmits(['update:modelValue'])

let _n = 0
const uid = computed(() => `ta-${++_n}`)
</script>

<style scoped>
.textarea-wrap  { display: flex; flex-direction: column; gap: 5px; }

.textarea-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
}

.textarea-field {
  width: 100%;
  padding: 9px 12px;
  border: 1px solid var(--border);
  border-radius: var(--r);
  font-size: 14px;
  color: var(--text);
  background: var(--surface);
  outline: none;
  resize: vertical;
  min-height: 90px;
  transition: border-color var(--t), box-shadow var(--t);
  line-height: 1.55;
}

.textarea-field:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(0,122,110,.12);
}
</style>
