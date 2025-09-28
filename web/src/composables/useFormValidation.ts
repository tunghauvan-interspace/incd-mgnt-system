import { ref, computed, type Ref } from 'vue'

export interface ValidationRule {
  required?: boolean
  minLength?: number
  maxLength?: number
  pattern?: RegExp
  email?: boolean
  custom?: (value: string) => string | null
}

export interface ValidationRules {
  [key: string]: ValidationRule
}

export interface UseFormValidationReturn {
  errors: Ref<Record<string, string>>
  isValid: Ref<boolean>
  validate: (formData: Record<string, any>, rules: ValidationRules) => boolean
  validateField: (fieldName: string, value: any, rule: ValidationRule) => string | null
  clearError: (fieldName: string) => void
  clearAllErrors: () => void
}

export function useFormValidation(): UseFormValidationReturn {
  const errors = ref<Record<string, string>>({})

  const isValid = computed(() => {
    return Object.keys(errors.value).length === 0
  })

  const validateField = (fieldName: string, value: any, rule: ValidationRule): string | null => {
    const stringValue = String(value || '').trim()

    // Required validation
    if (rule.required && !stringValue) {
      return `${fieldName} is required`
    }

    // Skip other validations if field is empty and not required
    if (!stringValue && !rule.required) {
      return null
    }

    // Email validation
    if (rule.email) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(stringValue)) {
        return 'Please enter a valid email address'
      }
    }

    // Min length validation
    if (rule.minLength && stringValue.length < rule.minLength) {
      return `${fieldName} must be at least ${rule.minLength} characters`
    }

    // Max length validation
    if (rule.maxLength && stringValue.length > rule.maxLength) {
      return `${fieldName} must not exceed ${rule.maxLength} characters`
    }

    // Pattern validation
    if (rule.pattern && !rule.pattern.test(stringValue)) {
      return `${fieldName} format is invalid`
    }

    // Custom validation
    if (rule.custom) {
      return rule.custom(stringValue)
    }

    return null
  }

  const validate = (formData: Record<string, any>, rules: ValidationRules): boolean => {
    const newErrors: Record<string, string> = {}

    Object.entries(rules).forEach(([fieldName, rule]) => {
      const error = validateField(fieldName, formData[fieldName], rule)
      if (error) {
        newErrors[fieldName] = error
      }
    })

    errors.value = newErrors
    return Object.keys(newErrors).length === 0
  }

  const clearError = (fieldName: string) => {
    const newErrors = { ...errors.value }
    delete newErrors[fieldName]
    errors.value = newErrors
  }

  const clearAllErrors = () => {
    errors.value = {}
  }

  return {
    errors,
    isValid,
    validate,
    validateField,
    clearError,
    clearAllErrors
  }
}