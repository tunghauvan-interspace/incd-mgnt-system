// Format date strings to locale format
export function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString()
}

// Format nanoseconds to human readable duration
export function formatDuration(nanoseconds: number): string {
  if (!nanoseconds || nanoseconds === 0) {
    return '-'
  }

  const seconds = nanoseconds / 1000000000
  
  if (seconds < 60) {
    return `${Math.round(seconds)}s`
  } else if (seconds < 3600) {
    const minutes = Math.round(seconds / 60)
    return `${minutes}m`
  } else if (seconds < 86400) {
    const hours = Math.round(seconds / 3600)
    return `${hours}h`
  } else {
    const days = Math.round(seconds / 86400)
    return `${days}d`
  }
}

// Escape HTML to prevent XSS
export function escapeHtml(text: string): string {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

// Calculate duration from start to end time
export function calculateDuration(startTime: string, endTime?: string): string {
  const start = new Date(startTime)
  const end = endTime ? new Date(endTime) : new Date()
  const durationMs = end.getTime() - start.getTime()
  const durationNs = durationMs * 1000000 // Convert ms to ns
  return formatDuration(durationNs)
}