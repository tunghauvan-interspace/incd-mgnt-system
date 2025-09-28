<script setup lang="ts">
export interface SeverityBadgeProps {
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info' | string
  size?: 'sm' | 'md' | 'lg'
  showIcon?: boolean
}

const props = withDefaults(defineProps<SeverityBadgeProps>(), {
  size: 'md',
  showIcon: false
})

const getSeverityConfig = (severity: string) => {
  const configs = {
    critical: {
      label: 'Critical',
      colorClass: 'severity-critical',
      icon: '🔴'
    },
    high: {
      label: 'High',
      colorClass: 'severity-high',
      icon: '🟠'
    },
    medium: {
      label: 'Medium',
      colorClass: 'severity-medium',
      icon: '🟡'
    },
    low: {
      label: 'Low',
      colorClass: 'severity-low',
      icon: '🔵'
    },
    info: {
      label: 'Info',
      colorClass: 'severity-info',
      icon: '⚪'
    }
  }

  return (
    configs[severity as keyof typeof configs] || {
      label: severity.charAt(0).toUpperCase() + severity.slice(1),
      colorClass: 'severity-default',
      icon: '⚫'
    }
  )
}

const severityConfig = getSeverityConfig(props.severity)
</script>

<template>
  <span
    :class="['severity-badge', `severity-badge--${size}`, severityConfig.colorClass]"
    :title="`Severity: ${severityConfig.label}`"
  >
    <span v-if="showIcon" class="severity-badge__icon" aria-hidden="true">
      {{ severityConfig.icon }}
    </span>
    {{ severityConfig.label }}
  </span>
</template>

<style scoped>
.severity-badge {
  /* Base styles using design tokens */
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);

  font-family: var(--font-family-base);
  font-weight: var(--font-weight-semibold);
  text-transform: uppercase;
  letter-spacing: 0.05em;

  border-radius: var(--radius-base);
  white-space: nowrap;

  /* Size variants */
  &.severity-badge--sm {
    padding: var(--spacing-xs) var(--spacing-sm);
    font-size: var(--font-size-xs);
    line-height: var(--line-height-tight);
  }

  &.severity-badge--md {
    padding: var(--spacing-xs) var(--spacing-md);
    font-size: var(--font-size-sm);
    line-height: var(--line-height-tight);
  }

  &.severity-badge--lg {
    padding: var(--spacing-sm) var(--spacing-base);
    font-size: var(--font-size-base);
    line-height: var(--line-height-tight);
  }
}

.severity-badge__icon {
  font-size: 0.875em;
}

/* Severity color variants */
.severity-critical {
  background-color: var(--color-severity-critical);
  color: var(--color-text-white);
}

.severity-high {
  background-color: var(--color-severity-high);
  color: var(--color-text-white);
}

.severity-medium {
  background-color: var(--color-severity-medium);
  color: var(--color-text-white);
}

.severity-low {
  background-color: var(--color-severity-low);
  color: var(--color-text-white);
}

.severity-info {
  background-color: var(--color-severity-info);
  color: var(--color-text-white);
}

.severity-default {
  background-color: var(--color-gray-500);
  color: var(--color-text-white);
}

/* Hover effect for interactive contexts */
.severity-badge:hover {
  transform: scale(1.05);
  transition: transform var(--transition-fast);
  box-shadow: var(--shadow-base);
}

/* Accessibility improvements */
@media (prefers-reduced-motion: reduce) {
  .severity-badge:hover {
    transform: none;
  }
}
</style>
