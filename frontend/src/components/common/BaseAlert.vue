<template>
  <div v-if="message" :class="['alert', `alert--${type}`]" role="alert">
    <span class="alert__icon" aria-hidden="true">{{ icons[type] }}</span>
    <span class="alert__msg">{{ message }}</span>
  </div>
</template>

<script setup>
defineProps({
  message: { type: String, default: '' },
  type:    { type: String, default: 'error' }, // error | success | warning | info
})

const icons = { error: '✕', success: '✓', warning: '!', info: 'i' }
</script>

<style scoped>
.alert {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 14px;
  border-radius: var(--r);
  font-size: 13px;
  line-height: 1.45;
  border-left: 3px solid currentColor;
}

.alert__icon {
  font-size: 11px;
  font-weight: 700;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: currentColor;
  color: #fff;
  flex-shrink: 0;
  margin-top: 1px;
}

.alert--error   { background: var(--danger-light);  color: var(--danger); }
.alert--success { background: var(--success-light); color: var(--success); }
.alert--warning { background: var(--warning-light); color: var(--warning); }
.alert--info    { background: var(--info-light);    color: var(--info); }
</style>
