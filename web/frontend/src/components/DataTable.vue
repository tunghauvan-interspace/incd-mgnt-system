<script setup lang="ts" generic="T extends Record<string, any>">
import { computed, ref } from 'vue'

export interface TableColumn<T = Record<string, any>> {
  key: keyof T
  label: string
  sortable?: boolean
  width?: string
  align?: 'left' | 'center' | 'right'
  formatter?: (value: any, row: T) => string
}

export interface TableProps<T = Record<string, any>> {
  columns: TableColumn<T>[]
  data: T[]
  loading?: boolean
  emptyMessage?: string
  sortBy?: keyof T | null
  sortOrder?: 'asc' | 'desc'
  hoverable?: boolean
  striped?: boolean
  bordered?: boolean
  compact?: boolean
}

interface TableEmits<T = Record<string, any>> {
  (e: 'sort', column: keyof T, order: 'asc' | 'desc'): void
  (e: 'rowClick', row: T, index: number): void
}

const props = withDefaults(defineProps<TableProps<T>>(), {
  loading: false,
  emptyMessage: 'No data available',
  sortBy: null,
  sortOrder: 'asc',
  hoverable: true,
  striped: false,
  bordered: false,
  compact: false
})

const emit = defineEmits<TableEmits<T>>()

const internalSortBy = ref(props.sortBy)
const internalSortOrder = ref(props.sortOrder)

const tableClasses = computed(() => {
  const classes = ['data-table']

  if (props.hoverable) classes.push('data-table--hoverable')
  if (props.striped) classes.push('data-table--striped')
  if (props.bordered) classes.push('data-table--bordered')
  if (props.compact) classes.push('data-table--compact')

  return classes
})

const sortedData = computed(() => {
  if (!internalSortBy.value) return props.data

  return [...props.data].sort((a, b) => {
    const aVal = a[internalSortBy.value!]
    const bVal = b[internalSortBy.value!]

    let comparison = 0

    if (aVal < bVal) comparison = -1
    else if (aVal > bVal) comparison = 1

    return internalSortOrder.value === 'desc' ? comparison * -1 : comparison
  })
})

const handleSort = (column: TableColumn<T>) => {
  if (!column.sortable) return

  const columnKey = column.key

  if (internalSortBy.value === columnKey) {
    // Toggle sort order
    internalSortOrder.value = internalSortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    // New column
    internalSortBy.value = columnKey
    internalSortOrder.value = 'asc'
  }

  emit('sort', columnKey, internalSortOrder.value)
}

const handleRowClick = (row: T, index: number) => {
  emit('rowClick', row, index)
}

const getSortIcon = (column: TableColumn<T>) => {
  if (!column.sortable) return ''
  if (internalSortBy.value !== column.key) return '↕️'
  return internalSortOrder.value === 'asc' ? '↑' : '↓'
}

const formatCellValue = (column: TableColumn<T>, row: T) => {
  const value = row[column.key]
  return column.formatter ? column.formatter(value, row) : String(value || '')
}
</script>

<template>
  <div class="table-container">
    <table :class="tableClasses">
      <thead class="data-table__head">
        <tr>
          <th
            v-for="column in columns"
            :key="String(column.key)"
            :class="[
              'data-table__header',
              `data-table__header--${column.align || 'left'}`,
              { 'data-table__header--sortable': column.sortable }
            ]"
            :style="{ width: column.width }"
            @click="handleSort(column)"
          >
            <div class="data-table__header-content">
              {{ column.label }}
              <span
                v-if="column.sortable"
                class="data-table__sort-icon"
                :class="{
                  'data-table__sort-icon--active': internalSortBy === column.key
                }"
              >
                {{ getSortIcon(column) }}
              </span>
            </div>
          </th>
        </tr>
      </thead>

      <tbody class="data-table__body">
        <tr v-if="loading">
          <td :colspan="columns.length" class="data-table__loading">
            <div class="loading-spinner"></div>
            Loading...
          </td>
        </tr>

        <tr v-else-if="sortedData.length === 0">
          <td :colspan="columns.length" class="data-table__empty">
            {{ emptyMessage }}
          </td>
        </tr>

        <tr
          v-else
          v-for="(row, index) in sortedData"
          :key="index"
          class="data-table__row"
          @click="handleRowClick(row, index)"
        >
          <td
            v-for="column in columns"
            :key="String(column.key)"
            :class="['data-table__cell', `data-table__cell--${column.align || 'left'}`]"
          >
            <slot
              :name="`cell-${String(column.key)}`"
              :row="row"
              :column="column"
              :value="row[column.key]"
              :index="index"
            >
              {{ formatCellValue(column, row) }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.table-container {
  background: var(--color-bg-primary);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-base);
  overflow: hidden;
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-family: var(--font-family-base);
}

/* Header styles */
.data-table__head {
  background-color: var(--color-bg-muted);
}

.data-table__header {
  padding: var(--spacing-md);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  border-bottom: 2px solid var(--color-border-light);
  white-space: nowrap;
}

.data-table__header--left {
  text-align: left;
}
.data-table__header--center {
  text-align: center;
}
.data-table__header--right {
  text-align: right;
}

.data-table__header--sortable {
  cursor: pointer;
  user-select: none;
  transition: background-color var(--transition-base);
}

.data-table__header--sortable:hover {
  background-color: var(--color-gray-100);
}

.data-table__header-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.data-table__sort-icon {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  transition: color var(--transition-base);
}

.data-table__sort-icon--active {
  color: var(--color-primary);
}

/* Body styles */
.data-table__row {
  transition: background-color var(--transition-fast);
}

.data-table__cell {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--color-border-light);
  color: var(--color-text-primary);
}

.data-table__cell--left {
  text-align: left;
}
.data-table__cell--center {
  text-align: center;
}
.data-table__cell--right {
  text-align: right;
}

/* Loading state */
.data-table__loading {
  padding: var(--spacing-2xl);
  text-align: center;
  color: var(--color-text-secondary);
}

.loading-spinner {
  display: inline-block;
  width: 1rem;
  height: 1rem;
  margin-right: var(--spacing-sm);
  border: 2px solid var(--color-gray-300);
  border-top-color: var(--color-primary);
  border-radius: var(--radius-full);
  animation: spin 0.8s linear infinite;
}

/* Empty state */
.data-table__empty {
  padding: var(--spacing-2xl);
  text-align: center;
  color: var(--color-text-muted);
  font-style: italic;
}

/* Table variants */
.data-table--hoverable .data-table__row:hover {
  background-color: var(--color-bg-muted);
  cursor: pointer;
}

.data-table--striped .data-table__row:nth-child(even) {
  background-color: var(--color-gray-50);
}

.data-table--bordered {
  border: 1px solid var(--color-border-light);
}

.data-table--bordered .data-table__cell {
  border-right: 1px solid var(--color-border-light);
}

.data-table--bordered .data-table__cell:last-child {
  border-right: none;
}

.data-table--compact .data-table__header,
.data-table--compact .data-table__cell {
  padding: var(--spacing-sm);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .data-table__header,
  .data-table__cell {
    padding: var(--spacing-sm);
    font-size: var(--font-size-sm);
  }

  .data-table__header-content {
    gap: var(--spacing-xs);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
