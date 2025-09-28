import { describe, it, expect } from 'vitest'
import { useFormValidation } from '@/composables/useFormValidation'

describe('useFormValidation', () => {
  it('initializes with empty errors and valid state', () => {
    const { errors, isValid } = useFormValidation()
    
    expect(errors.value).toEqual({})
    expect(isValid.value).toBe(true)
  })

  describe('validateField', () => {
    it('validates required fields', () => {
      const { validateField } = useFormValidation()
      
      // Test empty required field
      expect(validateField('username', '', { required: true }))
        .toBe('username is required')
      
      // Test valid required field
      expect(validateField('username', 'testuser', { required: true }))
        .toBe(null)
    })

    it('validates minimum length', () => {
      const { validateField } = useFormValidation()
      
      // Test too short
      expect(validateField('password', 'abc', { minLength: 8 }))
        .toBe('password must be at least 8 characters')
      
      // Test valid length
      expect(validateField('password', 'password123', { minLength: 8 }))
        .toBe(null)
    })

    it('validates maximum length', () => {
      const { validateField } = useFormValidation()
      
      // Test too long
      expect(validateField('username', 'a'.repeat(51), { maxLength: 50 }))
        .toBe('username must not exceed 50 characters')
      
      // Test valid length
      expect(validateField('username', 'testuser', { maxLength: 50 }))
        .toBe(null)
    })

    it('validates email format', () => {
      const { validateField } = useFormValidation()
      
      // Test invalid email
      expect(validateField('email', 'invalid-email', { email: true }))
        .toBe('Please enter a valid email address')
      
      // Test valid email
      expect(validateField('email', 'test@example.com', { email: true }))
        .toBe(null)
    })

    it('validates using regex pattern', () => {
      const { validateField } = useFormValidation()
      const phonePattern = /^\d{10}$/
      
      // Test invalid pattern
      expect(validateField('phone', '123-456-7890', { pattern: phonePattern }))
        .toBe('phone format is invalid')
      
      // Test valid pattern
      expect(validateField('phone', '1234567890', { pattern: phonePattern }))
        .toBe(null)
    })

    it('validates using custom function', () => {
      const { validateField } = useFormValidation()
      const customValidator = (value: string) => 
        value === 'forbidden' ? 'Value is forbidden' : null
      
      // Test custom validation failure
      expect(validateField('field', 'forbidden', { custom: customValidator }))
        .toBe('Value is forbidden')
      
      // Test custom validation success
      expect(validateField('field', 'allowed', { custom: customValidator }))
        .toBe(null)
    })

    it('skips validation for empty non-required fields', () => {
      const { validateField } = useFormValidation()
      
      // Empty non-required field should pass all validations
      expect(validateField('optional', '', { minLength: 8, email: true }))
        .toBe(null)
    })
  })

  describe('validate', () => {
    it('validates entire form and returns boolean', () => {
      const { validate, errors } = useFormValidation()
      
      const formData = {
        username: 'test',
        email: 'invalid-email',
        password: '123'
      }
      
      const rules = {
        username: { required: true, minLength: 5 },
        email: { required: true, email: true },
        password: { required: true, minLength: 8 }
      }
      
      // Should fail validation
      const isValid = validate(formData, rules)
      expect(isValid).toBe(false)
      
      // Should have errors for all fields
      expect(errors.value.username).toContain('at least 5 characters')
      expect(errors.value.email).toContain('valid email')
      expect(errors.value.password).toContain('at least 8 characters')
    })

    it('returns true for valid form', () => {
      const { validate, errors } = useFormValidation()
      
      const formData = {
        username: 'testuser',
        email: 'test@example.com',
        password: 'password123'
      }
      
      const rules = {
        username: { required: true, minLength: 5 },
        email: { required: true, email: true },
        password: { required: true, minLength: 8 }
      }
      
      const isValid = validate(formData, rules)
      expect(isValid).toBe(true)
      expect(errors.value).toEqual({})
    })
  })

  describe('error management', () => {
    it('clears individual field errors', () => {
      const { validate, clearError, errors } = useFormValidation()
      
      // Create some errors
      validate({ username: '' }, { username: { required: true } })
      expect(errors.value.username).toBeDefined()
      
      // Clear specific error
      clearError('username')
      expect(errors.value.username).toBeUndefined()
    })

    it('clears all errors', () => {
      const { validate, clearAllErrors, errors } = useFormValidation()
      
      // Create some errors
      validate(
        { username: '', email: '' }, 
        { 
          username: { required: true }, 
          email: { required: true } 
        }
      )
      
      expect(Object.keys(errors.value)).toHaveLength(2)
      
      // Clear all errors
      clearAllErrors()
      expect(errors.value).toEqual({})
    })
  })

  describe('reactive state', () => {
    it('updates isValid when errors change', () => {
      const { validate, clearAllErrors, isValid } = useFormValidation()
      
      // Initially valid
      expect(isValid.value).toBe(true)
      
      // Add errors - should become invalid
      validate({ username: '' }, { username: { required: true } })
      expect(isValid.value).toBe(false)
      
      // Clear errors - should become valid again
      clearAllErrors()
      expect(isValid.value).toBe(true)
    })
  })
})