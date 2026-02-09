module.exports = {
  testEnvironment: 'jsdom',
  transform: {},
  moduleFileExtensions: ['js', 'mjs', 'json', 'node'],
  testMatch: [
    '**/tests/unit/**/*.test.[jt]s?(x)',
    '**/tests/unit/**/*.test.mjs',
    '**/tests/js/unit/**/*.test.[jt]s?(x)',
    '**/tests/js/unit/**/*.test.mjs',
    '**/tests/js/integration/**/*.test.[jt]s?(x)',
    '**/tests/js/integration/**/*.test.mjs'
  ]
};
