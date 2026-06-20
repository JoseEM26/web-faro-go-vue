<template>
  <span class="badge" :style="{ background: lightColor, color: color, borderColor: color + '44' }">
    <span class="badge__dot" :style="{ background: color }" />
    <slot />
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  color: { type: String, default: '#007A6E' },
})

function hexToRgba(hex, alpha) {
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return `rgba(${r},${g},${b},${alpha})`
}

const lightColor = computed(() => hexToRgba(props.color, 0.1))
</script>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 2px 10px 2px 7px;
  border-radius: var(--r-full);
  font-size: 12px;
  font-weight: 500;
  border: 1px solid;
  line-height: 1.6;
}

.badge__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}
</style>
