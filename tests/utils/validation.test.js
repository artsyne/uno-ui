import { describe, it, expect } from 'vitest'
import { VALIDATORS, validate } from '../../src/utils/validation'

describe('Validation Utils - VALIDATORS', () => {
  describe('email validator', () => {
    it('should validate correct emails', () => {
      expect(VALIDATORS.email.test('user@example.com')).toBe(true)
      expect(VALIDATORS.email.test('test.user@domain.co.uk')).toBe(true)
      expect(VALIDATORS.email.test('user+tag@example.com')).toBe(true)
      expect(VALIDATORS.email.test('user_name@example.com')).toBe(true)
    })

    it('should reject invalid emails', () => {
      expect(VALIDATORS.email.test('invalid')).toBe(false)
      expect(VALIDATORS.email.test('@example.com')).toBe(false)
      expect(VALIDATORS.email.test('user@')).toBe(false)
      expect(VALIDATORS.email.test('user @example.com')).toBe(false)
      expect(VALIDATORS.email.test('user@.com')).toBe(false)
    })
  })

  describe('url validator', () => {
    it('should validate correct URLs', () => {
      expect(VALIDATORS.url.test('https://example.com')).toBe(true)
      expect(VALIDATORS.url.test('http://example.com')).toBe(true)
      expect(VALIDATORS.url.test('example.com')).toBe(true)
      expect(VALIDATORS.url.test('https://sub.domain.com/path')).toBe(true)
      expect(VALIDATORS.url.test('example.com/path/to/page')).toBe(true)
    })

    it('should reject invalid URLs', () => {
      expect(VALIDATORS.url.test('not a url')).toBe(false)
      expect(VALIDATORS.url.test('http://')).toBe(false)
      expect(VALIDATORS.url.test('://example.com')).toBe(false)
    })
  })

  describe('date validator', () => {
    it('should validate correct date formats', () => {
      expect(VALIDATORS.date.test('2024-01-15')).toBe(true) // ISO format
      expect(VALIDATORS.date.test('15/01/2024')).toBe(true) // DD/MM/YYYY
      expect(VALIDATORS.date.test('15.01.2024')).toBe(true) // DD.MM.YYYY
    })

    it('should reject invalid dates', () => {
      expect(VALIDATORS.date.test('2024-1-5')).toBe(false) // Missing leading zeros
      expect(VALIDATORS.date.test('15-01-2024')).toBe(false) // Wrong format
      expect(VALIDATORS.date.test('not a date')).toBe(false)
    })
  })

  describe('int validator', () => {
    it('should validate integers', () => {
      expect(VALIDATORS.int.test('123')).toBe(true)
      expect(VALIDATORS.int.test('-456')).toBe(true)
      expect(VALIDATORS.int.test('0')).toBe(true)
    })

    it('should reject non-integers', () => {
      expect(VALIDATORS.int.test('123.45')).toBe(false)
      expect(VALIDATORS.int.test('12.0')).toBe(false)
      expect(VALIDATORS.int.test('abc')).toBe(false)
      expect(VALIDATORS.int.test('12a')).toBe(false)
    })
  })

  describe('number validator', () => {
    it('should validate numbers', () => {
      expect(VALIDATORS.number.test('123')).toBe(true)
      expect(VALIDATORS.number.test('123.45')).toBe(true)
      expect(VALIDATORS.number.test('-123.45')).toBe(true)
      expect(VALIDATORS.number.test('0.5')).toBe(true)
      // Note: '.5' and '5.' patterns may not match the regex
      // Depending on the actual regex implementation
    })

    it('should reject non-numbers', () => {
      expect(VALIDATORS.number.test('abc')).toBe(false)
      expect(VALIDATORS.number.test('12a')).toBe(false)
      expect(VALIDATORS.number.test('12.34.56')).toBe(false)
    })
  })

  describe('validname validator', () => {
    it('should validate valid names', () => {
      expect(VALIDATORS.validname.test('username123')).toBe(true)
      expect(VALIDATORS.validname.test('user_name')).toBe(true)
      expect(VALIDATORS.validname.test('user-name')).toBe(true)
      expect(VALIDATORS.validname.test('ABC123')).toBe(true)
    })

    it('should reject invalid names', () => {
      expect(VALIDATORS.validname.test('user name')).toBe(false) // Space
      expect(VALIDATORS.validname.test('user@name')).toBe(false) // Special char
      expect(VALIDATORS.validname.test('用户名')).toBe(false) // Non-ASCII
    })
  })

  describe('en (English letters) validator', () => {
    it('should validate English letters only', () => {
      expect(VALIDATORS.en.test('abc')).toBe(true)
      expect(VALIDATORS.en.test('ABC')).toBe(true)
      expect(VALIDATORS.en.test('AbCdEf')).toBe(true)
    })

    it('should reject non-English input', () => {
      expect(VALIDATORS.en.test('abc123')).toBe(false)
      expect(VALIDATORS.en.test('abc-def')).toBe(false)
      expect(VALIDATORS.en.test('abc def')).toBe(false)
      expect(VALIDATORS.en.test('中文')).toBe(false)
    })
  })

  describe('domain validator', () => {
    it('should validate correct domains', () => {
      expect(VALIDATORS.domain.test('example.com')).toBe(true)
      expect(VALIDATORS.domain.test('sub.domain.com')).toBe(true)
      expect(VALIDATORS.domain.test('my-site.co.uk')).toBe(true)
    })

    it('should reject invalid domains', () => {
      expect(VALIDATORS.domain.test('-example.com')).toBe(false) // Leading hyphen
      expect(VALIDATORS.domain.test('example-.com')).toBe(false) // Trailing hyphen
      expect(VALIDATORS.domain.test('example')).toBe(false) // No TLD
      expect(VALIDATORS.domain.test('example..com')).toBe(false) // Double dot
    })
  })

  describe('mobile validator', () => {
    it('should validate phone numbers', () => {
      expect(VALIDATORS.mobile.test('+1234567')).toBe(true)
      expect(VALIDATORS.mobile.test('1234567890')).toBe(true)
      expect(VALIDATORS.mobile.test('+1 234 567 8900')).toBe(true)
      expect(VALIDATORS.mobile.test('123-456-7890')).toBe(true)
    })

    it('should reject invalid phone numbers', () => {
      expect(VALIDATORS.mobile.test('123')).toBe(false) // Too short
      expect(VALIDATORS.mobile.test('abc123')).toBe(false) // Letters
    })
  })

  describe('pwd (password) validator', () => {
    it('should validate passwords with letters and numbers (8-72 chars)', () => {
      expect(VALIDATORS.pwd.test('abcd1234')).toBe(true)
      expect(VALIDATORS.pwd.test('Password1')).toBe(true)
      expect(VALIDATORS.pwd.test('12345abc')).toBe(true)
      expect(VALIDATORS.pwd.test('Pass!word1')).toBe(true) // "!" allowed, aligned with backend pkg.IsPwd
      expect(VALIDATORS.pwd.test('a1'.repeat(36))).toBe(true) // Exactly 72 chars (bcrypt limit) is valid
    })

    it('should reject weak passwords', () => {
      expect(VALIDATORS.pwd.test('12345678')).toBe(false) // No letters
      expect(VALIDATORS.pwd.test('abcdefgh')).toBe(false) // No numbers
      expect(VALIDATORS.pwd.test('abc123')).toBe(false) // Too short
      expect(VALIDATORS.pwd.test('a1'.repeat(36) + 'x')).toBe(false) // Exceeds 72-char bcrypt limit
      expect(VALIDATORS.pwd.test('abcd1234/')).toBe(false) // "/" is not in the allowed charset
      expect(VALIDATORS.pwd.test('')).toBe(false)
    })
  })

  describe('strongpwd (strong password) validator', () => {
    it('should validate strong passwords with letters, numbers, and special chars', () => {
      expect(VALIDATORS.strongpwd.test('Abcd123!')).toBe(true)
      expect(VALIDATORS.strongpwd.test('Pass@word1')).toBe(true)
      expect(VALIDATORS.strongpwd.test('Secure#2024')).toBe(true)
    })

    it('should reject passwords without all requirements', () => {
      expect(VALIDATORS.strongpwd.test('abcd1234')).toBe(false) // No special char
      expect(VALIDATORS.strongpwd.test('Abcd@#$%')).toBe(false) // No numbers
      expect(VALIDATORS.strongpwd.test('1234@#$%')).toBe(false) // No letters
      expect(VALIDATORS.strongpwd.test('Abc1!')).toBe(false) // Too short
    })
  })
})

describe('Validation Utils - validate()', () => {
  it('should validate using built-in validator', () => {
    expect(validate('user@example.com', 'email')).toBe(true)
    expect(validate('invalid', 'email')).toBe(false)
  })

  it('should return true for unknown validator type', () => {
    expect(validate('anything', 'unknown')).toBe(true)
  })

  it('should return true for null/undefined type', () => {
    expect(validate('test', null)).toBe(true)
    expect(validate('test', undefined)).toBe(true)
    expect(validate('test', '')).toBe(true)
  })

  it('should validate different types correctly', () => {
    expect(validate('123', 'int')).toBe(true)
    expect(validate('123.45', 'number')).toBe(true)
    expect(validate('https://example.com', 'url')).toBe(true)
    expect(validate('2024-01-15', 'date')).toBe(true)
  })

  it('should reject invalid values', () => {
    expect(validate('abc', 'int')).toBe(false)
    expect(validate('invalid', 'url')).toBe(false)
    expect(validate('not-a-date', 'date')).toBe(false)
  })
})

describe('Validation Utils - edge cases', () => {
  it('should handle empty strings', () => {
    expect(VALIDATORS.email.test('')).toBe(false)
    expect(VALIDATORS.url.test('')).toBe(false)
    expect(VALIDATORS.int.test('')).toBe(false)
  })

  it('should handle whitespace', () => {
    expect(VALIDATORS.email.test('  user@example.com  ')).toBe(false)
    expect(VALIDATORS.int.test(' 123 ')).toBe(false)
  })

  it('should be case-insensitive for URLs', () => {
    expect(VALIDATORS.url.test('HTTPS://EXAMPLE.COM')).toBe(true)
    expect(VALIDATORS.url.test('HTTP://example.com')).toBe(true)
  })
})
