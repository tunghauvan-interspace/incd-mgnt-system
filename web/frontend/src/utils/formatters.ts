import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

export const formatDate = (date: string | Date): string => {
  return dayjs(date).format('MMM DD, YYYY HH:mm')
}

export const formatRelativeTime = (date: string | Date): string => {
  return dayjs(date).fromNow()
}

export const formatDuration = (start: string | Date, end?: string | Date): string => {
  const startTime = dayjs(start)
  const endTime = end ? dayjs(end) : dayjs()
  
  const diff = endTime.diff(startTime, 'minute')
  
  if (diff < 60) {
    return `${diff}m`
  } else if (diff < 1440) { // 24 hours
    return `${Math.floor(diff / 60)}h ${diff % 60}m`
  } else {
    return `${Math.floor(diff / 1440)}d ${Math.floor((diff % 1440) / 60)}h`
  }
}

export const formatNumber = (num: number): string => {
  return new Intl.NumberFormat().format(num)
}

export const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

export const formatPercentage = (value: number, total: number): string => {
  if (total === 0) return '0%'
  return `${Math.round((value / total) * 100)}%`
}

export const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text
  return text.slice(0, maxLength) + '...'
}