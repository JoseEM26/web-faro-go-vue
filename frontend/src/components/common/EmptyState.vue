<template>
  <div class="empty">
    <div class="empty__icon" aria-hidden="true">{{ icon }}</div>
    <h3 class="empty__title">{{ title }}</h3>
    <p class="empty__desc">{{ description }}</p>
    <div v-if="$slots.action" class="empty__action">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup>
defineProps({
  icon:        { type: String, default: '📭' },
  title:       { type: String, required: true },
  description: { type: String, default: '' },
})
</script>

<style scoped>
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 56px 32px;
  border: 1px dashed var(--border-strong);
  border-radius: var(--r-lg);
  background: var(--surface);
}

.empty__icon  { font-size: 36px; margin-bottom: 14px; line-height: 1; }
.empty__title { font-size: 16px; font-weight: 600; color: var(--text); margin-bottom: 6px; }
.empty__desc  { font-size: 13px; color: var(--text-muted); max-width: 300px; line-height: 1.5; }
.empty__action { margin-top: 20px; }
</style>
