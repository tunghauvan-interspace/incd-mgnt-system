import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatusBadge from '@/components/StatusBadge.vue'

describe('StatusBadge', () => {
  it('renders correctly with open status', () => {
    const wrapper = mount(StatusBadge, {
      props: {
        status: 'open'
      }
    })

    expect(wrapper.text()).toBe('Open')
    expect(wrapper.classes()).toContain('status-badge')
    expect(wrapper.classes()).toContain('status-badge--md')
    expect(wrapper.classes()).toContain('status-open')
  })

  it('renders correctly with acknowledged status', () => {
    const wrapper = mount(StatusBadge, {
      props: {
        status: 'acknowledged'
      }
    })

    expect(wrapper.text()).toBe('Acknowledged')
    expect(wrapper.classes()).toContain('status-acknowledged')
  })

  it('renders correctly with resolved status', () => {
    const wrapper = mount(StatusBadge, {
      props: {
        status: 'resolved'
      }
    })

    expect(wrapper.text()).toBe('Resolved')
    expect(wrapper.classes()).toContain('status-resolved')
  })

  it('handles unknown status gracefully', () => {
    const wrapper = mount(StatusBadge, {
      props: {
        status: 'unknown-status'
      }
    })

    expect(wrapper.text()).toBe('Unknown-status')
    expect(wrapper.classes()).toContain('status-default')
  })

  it('applies correct size classes', () => {
    const wrapperSm = mount(StatusBadge, {
      props: {
        status: 'open',
        size: 'sm'
      }
    })

    const wrapperLg = mount(StatusBadge, {
      props: {
        status: 'open',
        size: 'lg'
      }
    })

    expect(wrapperSm.classes()).toContain('status-badge--sm')
    expect(wrapperLg.classes()).toContain('status-badge--lg')
  })

  it('has correct accessibility attributes', () => {
    const wrapper = mount(StatusBadge, {
      props: {
        status: 'acknowledged'
      }
    })

    expect(wrapper.attributes('title')).toBe('Status: Acknowledged')
  })

  it('defaults to medium size when size prop is not provided', () => {
    const wrapper = mount(StatusBadge, {
      props: {
        status: 'open'
      }
    })

    expect(wrapper.classes()).toContain('status-badge--md')
  })
})