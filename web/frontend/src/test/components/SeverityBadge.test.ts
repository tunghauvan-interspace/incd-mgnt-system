import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SeverityBadge from '@/components/SeverityBadge.vue'

describe('SeverityBadge', () => {
  it('renders correctly with critical severity', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'critical'
      }
    })

    expect(wrapper.text()).toBe('Critical')
    expect(wrapper.classes()).toContain('severity-badge')
    expect(wrapper.classes()).toContain('severity-badge--md')
    expect(wrapper.classes()).toContain('severity-critical')
  })

  it('renders correctly with high severity', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'high'
      }
    })

    expect(wrapper.text()).toBe('High')
    expect(wrapper.classes()).toContain('severity-high')
  })

  it('renders correctly with medium severity', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'medium'
      }
    })

    expect(wrapper.text()).toBe('Medium')
    expect(wrapper.classes()).toContain('severity-medium')
  })

  it('renders correctly with low severity', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'low'
      }
    })

    expect(wrapper.text()).toBe('Low')
    expect(wrapper.classes()).toContain('severity-low')
  })

  it('renders correctly with info severity', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'info'
      }
    })

    expect(wrapper.text()).toBe('Info')
    expect(wrapper.classes()).toContain('severity-info')
  })

  it('handles unknown severity gracefully', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'unknown-severity'
      }
    })

    expect(wrapper.text()).toBe('Unknown-severity')
    expect(wrapper.classes()).toContain('severity-default')
  })

  it('applies correct size classes', () => {
    const wrapperSm = mount(SeverityBadge, {
      props: {
        severity: 'critical',
        size: 'sm'
      }
    })

    const wrapperLg = mount(SeverityBadge, {
      props: {
        severity: 'critical',
        size: 'lg'
      }
    })

    expect(wrapperSm.classes()).toContain('severity-badge--sm')
    expect(wrapperLg.classes()).toContain('severity-badge--lg')
  })

  it('shows icon when showIcon prop is true', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'critical',
        showIcon: true
      }
    })

    expect(wrapper.text()).toContain('🔴')
    expect(wrapper.text()).toContain('Critical')
  })

  it('does not show icon by default', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'critical'
      }
    })

    expect(wrapper.text()).not.toContain('🔴')
    expect(wrapper.text()).toBe('Critical')
  })

  it('has correct accessibility attributes', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'high'
      }
    })

    expect(wrapper.attributes('title')).toBe('Severity: High')
  })

  it('defaults to medium size when size prop is not provided', () => {
    const wrapper = mount(SeverityBadge, {
      props: {
        severity: 'critical'
      }
    })

    expect(wrapper.classes()).toContain('severity-badge--md')
  })
})