import globals from 'globals'
import pluginJs from '@eslint/js'
import eslintConfigPrettier from 'eslint-config-prettier'

export default [
  {
    ignores: ['dist/*']
  },
  {
    languageOptions: {
      globals: {
        ...globals.browser,
        Alpine: 'readonly',
        htmx: 'readonly'
      },
      sourceType: 'module'
    }
  },
  pluginJs.configs.recommended,
  eslintConfigPrettier,
  {
    rules: {
      'no-unused-vars': 'warn',
      'no-undef': 'warn'
    }
  }
]
