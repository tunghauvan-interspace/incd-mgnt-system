import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Modal from '@/components/Modal.vue'

describe('Modal', () => {
  let originalBodyStyle: string

  beforeEach(() => {
    originalBodyStyle = document.body.style.overflow
    // Mock document methods
    document.addEventListener = vi.fn()
    document.removeEventListener = vi.fn()

    // Create a div to act as the teleport target
    const div = document.createElement('div')
    document.body.appendChild(div)
  })

  afterEach(() => {
    document.body.style.overflow = originalBodyStyle
    document.body.innerHTML = '' // Clear teleported content
    vi.clearAllMocks()
  })

  it('renders when show prop is true', () => {
    const wrapper = mount(Modal, {
      props: {
        show: true,
        title: 'Test Modal'
      },
      slots: {
        default: '<p>Modal content</p>',
        footer: '<button>Close</button>'
      },
      attachTo: document.body
    })

    // Check if modal is rendered in document body (due to Teleport)
    expect(document.querySelector('.modal-overlay')).toBeTruthy()
    expect(document.querySelector('.modal h3')?.textContent).toBe('Test Modal')
    expect(document.body.textContent).toContain('Modal content')
  })

  it('does not render when show prop is false', () => {
    mount(Modal, {
      props: {
        show: false,
        title: 'Test Modal'
      },
      attachTo: document.body
    })

    expect(document.querySelector('.modal-overlay')).toBeFalsy()
  })

  it('emits close event when close button is clicked', async () => {
    const wrapper = mount(Modal, {
      props: {
        show: true,
        title: 'Test Modal'
      },
      attachTo: document.body
    })

    const closeButton = document.querySelector('.modal-close') as HTMLButtonElement
    expect(closeButton).toBeTruthy()

    closeButton?.click()

    expect(wrapper.emitted('close')).toBeTruthy()
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('emits close event when overlay is clicked', async () => {
    const wrapper = mount(Modal, {
      props: {
        show: true,
        title: 'Test Modal'
      },
      attachTo: document.body
    })

    // Get overlay element and simulate click
    const overlayElement = document.querySelector('.modal-overlay') as HTMLElement
    expect(overlayElement).toBeTruthy()

    // Create and dispatch click event
    const clickEvent = new Event('click', { bubbles: true })
    Object.defineProperty(clickEvent, 'target', { value: overlayElement })
    overlayElement?.dispatchEvent(clickEvent)

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('does not emit close when modal content is clicked', async () => {
    const wrapper = mount(Modal, {
      props: {
        show: true,
        title: 'Test Modal'
      },
      attachTo: document.body
    })

    // Click on modal content instead of overlay
    const modalElement = document.querySelector('.modal') as HTMLElement
    expect(modalElement).toBeTruthy()

    modalElement?.click()

    expect(wrapper.emitted('close')).toBeFalsy()
  })

  it('renders without title when title prop is not provided', () => {
    mount(Modal, {
      props: {
        show: true
      },
      attachTo: document.body
    })

    expect(document.querySelector('.modal h3')).toBeFalsy()
  })

  it('renders default slot content', () => {
    mount(Modal, {
      props: {
        show: true
      },
      slots: {
        default: '<div class="test-content">Test content</div>'
      },
      attachTo: document.body
    })

    expect(document.querySelector('.test-content')).toBeTruthy()
    expect(document.body.textContent).toContain('Test content')
  })

  it('renders footer slot content', () => {
    mount(Modal, {
      props: {
        show: true
      },
      slots: {
        footer: '<div class="test-footer">Test footer</div>'
      },
      attachTo: document.body
    })

    expect(document.querySelector('.test-footer')).toBeTruthy()
    expect(document.body.textContent).toContain('Test footer')
  })

  it('applies correct CSS classes', () => {
    mount(Modal, {
      props: {
        show: true,
        title: 'Test Modal'
      },
      attachTo: document.body
    })

    expect(document.querySelector('.modal-overlay')).toBeTruthy()
    expect(document.querySelector('.modal')).toBeTruthy()
    expect(document.querySelector('.modal-header')).toBeTruthy()
    expect(document.querySelector('.modal-body')).toBeTruthy()
  })

  it('shows close button with correct aria-label', () => {
    mount(Modal, {
      props: {
        show: true,
        title: 'Test Modal'
      },
      attachTo: document.body
    })

    const closeButton = document.querySelector('.modal-close')
    expect(closeButton).toBeTruthy()
    expect(closeButton?.getAttribute('aria-label')).toBe('Close modal')
  })
})
