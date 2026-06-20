<template>
  <button
    v-bind="$attrs"
    :class="['btn', `btn--${variant}`, `btn--${size}`, { 'btn--loading': loading, 'btn--full': full }]"
    :disabled="disabled || loading"
  >
    <span v-if="loading" class="btn__spinner" aria-hidden="true" />
    <slot />
  </button>
</template>

<script setup>
defineProps({
  variant:  { type: String, default: 'primary' },  // primary | secondary | ghost | danger
  size:     { type: String, default: 'md' },        // sm | md | lg
  loading:  { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  full:     { type: Boolean, default: false },
})
</script>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border-radius: var(--r);
  font-weight: 500;
  border: 1px solid transparent;
  transition: background var(--t), color var(--t), border-color var(--t), opacity var(--t);
  line-height: 1;
  white-space: nowrap;
  cursor: pointer;
}

.btn:focus-visible { outline: 2px solid var(--accent); outline-offset: 2px; }
.btn:disabled      { opacity: .5; cursor: not-allowed; }
.btn--full         { width: 100%; }

/* Sizes */
.btn--sm { padding: 5px 12px;  font-size: 12px; }
.btn--md { padding: 8px 16px;  font-size: 14px; }
.btn--lg { padding: 11px 22px; font-size: 15px; }

/* Variants */
.btn--primary {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}
.btn--primary:not(:disabled):hover {
  background: var(--accent-hover);
  border-color: var(--accent-hover);
}

.btn--secondary {
  background: transparent;
  color: var(--accent);
  border-color: var(--accent);
}
.btn--secondary:not(:disabled):hover {
  background: var(--accent-light);
}

.btn--ghost {
  background: transparent;
  color: var(--text-muted);
}
.btn--ghost:not(:disabled):hover {
  background: var(--accent-light);
  color: var(--accent);
}

.btn--danger {
  background: var(--danger-light);
  color: var(--danger);
  border-color: transparent;
}
.btn--danger:not(:disabled):hover {
  background: #fca5a5;
}

/* Spinner */
.btn__spinner {
  width: 13px;
  height: 13px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin .7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }
</style>
