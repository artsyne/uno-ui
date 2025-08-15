module.exports = (ctx) => {
  const isProduction = ctx.env === 'production';
  
  return {
    plugins: {
      '@tailwindcss/postcss': {},
      autoprefixer: {},
      ...(isProduction && {
        cssnano: {
          preset: ['default', {
            discardComments: {
              removeAll: true,
            },
            normalizeWhitespace: true,
            minifySelectors: true,
          }],
        },
      }),
    },
  };
};
