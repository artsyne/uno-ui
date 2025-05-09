module.exports = {
  content: [
    './website/static/src/**/*.html',
    './website/static/src/**/*.js',
    './website/static/src/**/*.jsx',
    './website/static/src/**/*.ts',
    './website/static/src/**/*.tsx',
    './website/static/src/**/*.vue',
    './website/views/**/*.html',
    './website/controllers/**/*.go',
    './website/routes/**/*.go',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Arial', 'Helvetica', 'sans-serif'], // 设置默认字体
      },
    },
  },
  plugins: [],
};
