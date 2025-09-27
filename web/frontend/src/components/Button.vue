<script setup lang="ts">
import { computed } from 'vue'

export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  type?: 'button' | 'submit' | 'reset'
  block?: boolean
}

interface ButtonEmits {
  (e: 'click', event: MouseEvent): void
}

const props = withDefaults(defineProps<ButtonProps>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
  type: 'button',
  block: false
})

const emit = defineEmits<ButtonEmits>()

const buttonClasses = computed(() => {
  const classes = ['btn', `btn--${props.variant}`, `btn--${props.size}`]

  if (props.disabled || props.loading) {
    classes.push('btn--disabled')
  }

  if (props.block) {
    classes.push('btn--block')
  }

  if (props.loading) {
    classes.push('btn--loading')
  }

  return classes
})

const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<template>
  <button :class="buttonClasses" :type="type" :disabled="disabled || loading" @click="handleClick">
    <span v-if="loading" class="btn__spinner" aria-hidden="true"></span>
    <slot></slot>
  </button>
</template>

<style scoped>
.btn {
  /* Base styles using design tokens */
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);

  font-family: var(--font-family-base);
  font-weight: var(--font-weight-medium);
  text-decoration: none;
  text-align: center;

  border: none;
  border-radius: var(--radius-base);
  cursor: pointer;

  transition: all var(--transition-base);

  /* Prevent text selection */
  user-select: none;
  -webkit-user-select: none;

  /* Focus styles */
  &:focus {
    outline: 2px solid var(--color-primary);
    outline-offset: 2px;
  }

  &:focus:not(:focus-visible) {
    outline: none;
  }
}

/* Size variants */
.btn--sm {
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-tight);
}

.btn--md {
  padding: var(--spacing-sm) var(--spacing-base);
  font-size: var(--font-size-base);
  line-height: var(--line-height-base);
}

.btn--lg {
  padding: var(--spacing-md) var(--spacing-lg);
  font-size: var(--font-size-lg);
  line-height: var(--line-height-base);
}

/* Color variants */
.btn--primary {
  background-color: var(--color-primary);
  color: var(--color-text-white);

  &:hover:not(.btn--disabled) {
    background-color: var(--color-primary-hover);
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }
}

.btn--secondary {
  background-color: var(--color-secondary);
  color: var(--color-text-white);

  &:hover:not(.btn--disabled) {
    background-color: var(--color-secondary-hover);
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }
}

.btn--success {
  background-color: var(--color-success);
  color: var(--color-text-white);

  &:hover:not(.btn--disabled) {
    background-color: var(--color-success-hover);
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }
}

.btn--warning {
  background-color: var(--color-warning);
  color: var(--color-text-white);

  &:hover:not(.btn--disabled) {
    background-color: var(--color-warning-hover);
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }
}

.btn--danger {
  background-color: var(--color-danger);
  color: var(--color-text-white);

  &:hover:not(.btn--disabled) {
    background-color: var(--color-danger-hover);
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }
}

/* Disabled state */
.btn--disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
}

/* Block button */
.btn--block {
  width: 100%;
}

/* Loading state */
.btn--loading {
  cursor: not-allowed;
}

.btn__spinner {
  display: inline-block;
  width: 1em;
  height: 1em;
  border: 2px solid currentColor;
  border-right-color: transparent;
  border-radius: var(--radius-full);
  animation: spin 0.75s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .btn--lg {
    padding: var(--spacing-sm) var(--spacing-base);
    font-size: var(--font-size-base);
  }
}
</style>
