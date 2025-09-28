<script setup lang="ts">
export interface StatusBadgeProps {
  status: 'open' | 'acknowledged' | 'resolved' | string
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<StatusBadgeProps>(), {
  size: 'md'
})

const getStatusConfig = (status: string) => {
  const configs = {
    open: {
      label: 'Open',
      colorClass: 'status-open'
    },
    acknowledged: {
      label: 'Acknowledged',
      colorClass: 'status-acknowledged'
    },
    resolved: {
      label: 'Resolved',
      colorClass: 'status-resolved'
    }
  }

  return (
    configs[status as keyof typeof configs] || {
      label: status.charAt(0).toUpperCase() + status.slice(1),
      colorClass: 'status-default'
    }
  )
}

const statusConfig = getStatusConfig(props.status)
</script>

<template>
  <span
    :class="['status-badge', `status-badge--${size}`, statusConfig.colorClass]"
    :title="`Status: ${statusConfig.label}`"
  >
    {{ statusConfig.label }}
  </span>
</template>

<style scoped>
.status-badge {
  /* Base styles using design tokens */
  display: inline-flex;
  align-items: center;
  justify-content: center;

  font-family: var(--font-family-base);
  font-weight: var(--font-weight-medium);
  text-transform: uppercase;
  letter-spacing: 0.05em;

  border-radius: var(--radius-full);
  white-space: nowrap;

  /* Size variants */
  &.status-badge--sm {
    padding: var(--spacing-xs) var(--spacing-sm);
    font-size: var(--font-size-xs);
    line-height: var(--line-height-tight);
  }

  &.status-badge--md {
    padding: var(--spacing-xs) var(--spacing-md);
    font-size: var(--font-size-sm);
    line-height: var(--line-height-tight);
  }

  &.status-badge--lg {
    padding: var(--spacing-sm) var(--spacing-base);
    font-size: var(--font-size-base);
    line-height: var(--line-height-tight);
  }
}

/* Status color variants */
.status-open {
  background-color: var(--color-status-open-bg);
  color: var(--color-status-open);
}

.status-acknowledged {
  background-color: var(--color-status-acknowledged-bg);
  color: var(--color-status-acknowledged);
}

.status-resolved {
  background-color: var(--color-status-resolved-bg);
  color: var(--color-status-resolved);
}

.status-default {
  background-color: var(--color-gray-200);
  color: var(--color-gray-700);
}

/* Hover effect for interactive contexts */
.status-badge:hover {
  transform: scale(1.05);
  transition: transform var(--transition-fast);
}

/* Accessibility improvements */
@media (prefers-reduced-motion: reduce) {
  .status-badge:hover {
    transform: none;
  }
}
</style>
