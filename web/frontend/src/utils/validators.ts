export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

export const validatePassword = (password: string): { valid: boolean; errors: string[] } => {
  const errors: string[] = []
  
  if (password.length < 8) {
    errors.push('Password must be at least 8 characters long')
  }
  
  if (!/[A-Z]/.test(password)) {
    errors.push('Password must contain at least one uppercase letter')
  }
  
  if (!/[a-z]/.test(password)) {
    errors.push('Password must contain at least one lowercase letter')
  }
  
  if (!/\d/.test(password)) {
    errors.push('Password must contain at least one number')
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

export const validateRequired = (value: any): boolean => {
  if (value === null || value === undefined) return false
  if (typeof value === 'string') return value.trim().length > 0
  if (Array.isArray(value)) return value.length > 0
  return true
}

export const validateUrl = (url: string): boolean => {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

export const validatePhoneNumber = (phone: string): boolean => {
  const phoneRegex = /^\+?[\d\s\-\(\)]+$/
  return phoneRegex.test(phone) && phone.replace(/\D/g, '').length >= 10
}

export const validateIncidentTitle = (title: string): { valid: boolean; error?: string } => {
  if (!validateRequired(title)) {
    return { valid: false, error: 'Title is required' }
  }
  
  if (title.length < 5) {
    return { valid: false, error: 'Title must be at least 5 characters long' }
  }
  
  if (title.length > 100) {
    return { valid: false, error: 'Title must be less than 100 characters' }
  }
  
  return { valid: true }
}

export const validateIncidentDescription = (description: string): { valid: boolean; error?: string } => {
  if (!validateRequired(description)) {
    return { valid: false, error: 'Description is required' }
  }
  
  if (description.length < 10) {
    return { valid: false, error: 'Description must be at least 10 characters long' }
  }
  
  if (description.length > 1000) {
    return { valid: false, error: 'Description must be less than 1000 characters' }
  }
  
  return { valid: true }
}