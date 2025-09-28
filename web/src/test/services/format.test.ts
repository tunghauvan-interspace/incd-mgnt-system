import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { formatDate, formatDuration, escapeHtml, calculateDuration } from '@/utils/format'

describe('Format Utilities', () => {
  describe('formatDate', () => {
    it('formats valid date string correctly', () => {
      const dateString = '2023-12-25T10:30:00.000Z'
      const result = formatDate(dateString)
      
      // Result will vary by locale, but should contain date elements
      expect(result).toBeTruthy()
      expect(result.length).toBeGreaterThan(0)
    })

    it('handles invalid date string', () => {
      const result = formatDate('invalid-date')
      expect(result).toBe('Invalid Date')
    })

    it('formats ISO date string', () => {
      const isoDate = '2023-01-01T00:00:00.000Z'
      const result = formatDate(isoDate)
      expect(result).toBeTruthy()
    })
  })

  describe('formatDuration', () => {
    it('returns dash for zero nanoseconds', () => {
      expect(formatDuration(0)).toBe('-')
    })

    it('returns dash for undefined/null', () => {
      expect(formatDuration(null as any)).toBe('-')
      expect(formatDuration(undefined as any)).toBe('-')
    })

    it('formats seconds correctly', () => {
      expect(formatDuration(1000000000)).toBe('1s') // 1 second in ns
      expect(formatDuration(30000000000)).toBe('30s') // 30 seconds in ns
    })

    it('formats minutes correctly', () => {
      expect(formatDuration(60000000000)).toBe('1m') // 1 minute in ns
      expect(formatDuration(1800000000000)).toBe('30m') // 30 minutes in ns
    })

    it('formats hours correctly', () => {
      expect(formatDuration(3600000000000)).toBe('1h') // 1 hour in ns
      expect(formatDuration(7200000000000)).toBe('2h') // 2 hours in ns
    })

    it('formats days correctly', () => {
      expect(formatDuration(86400000000000)).toBe('1d') // 1 day in ns
      expect(formatDuration(259200000000000)).toBe('3d') // 3 days in ns
    })

    it('rounds to nearest unit', () => {
      expect(formatDuration(1500000000)).toBe('2s') // 1.5 seconds rounds to 2s
      expect(formatDuration(90000000000)).toBe('2m') // 1.5 minutes rounds to 2m
    })
  })

  describe('escapeHtml', () => {
    it('escapes HTML entities correctly', () => {
      expect(escapeHtml('<script>alert("xss")</script>')).toBe('&lt;script&gt;alert("xss")&lt;/script&gt;')
      expect(escapeHtml('Tom & Jerry')).toBe('Tom &amp; Jerry')
      expect(escapeHtml('"quoted text"')).toBe('"quoted text"')
    })

    it('handles empty string', () => {
      expect(escapeHtml('')).toBe('')
    })

    it('handles regular text without HTML', () => {
      expect(escapeHtml('regular text')).toBe('regular text')
    })

    it('handles multiple HTML entities', () => {
      const input = '<div>Tom & Jerry "cartoon"</div>'
      const expected = '&lt;div&gt;Tom &amp; Jerry "cartoon"&lt;/div&gt;'
      expect(escapeHtml(input)).toBe(expected)
    })
  })

  describe('calculateDuration', () => {
    beforeEach(() => {
      // Mock Date.now() for consistent testing
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2023-12-25T12:00:00.000Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('calculates duration with end time', () => {
      const startTime = '2023-12-25T10:00:00.000Z'
      const endTime = '2023-12-25T11:00:00.000Z'
      const result = calculateDuration(startTime, endTime)
      expect(result).toBe('1h')
    })

    it('calculates duration without end time (uses current time)', () => {
      const startTime = '2023-12-25T11:30:00.000Z' // 30 minutes ago
      const result = calculateDuration(startTime)
      expect(result).toBe('30m')
    })

    it('handles same start and end time', () => {
      const time = '2023-12-25T10:00:00.000Z'
      const result = calculateDuration(time, time)
      expect(result).toBe('-')
    })

    it('calculates duration for different time periods', () => {
      const startTime = '2023-12-24T12:00:00.000Z' // 1 day ago
      const result = calculateDuration(startTime)
      expect(result).toBe('1d')
    })

    it('handles future start time gracefully', () => {
      const futureTime = '2023-12-25T13:00:00.000Z' // 1 hour in future
      const result = calculateDuration(futureTime)
      // Should handle negative duration gracefully
      expect(typeof result).toBe('string')
    })
  })
})