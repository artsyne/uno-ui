module.exports = {
    content: [
      './src/**/*.{html,js,ts,jsx,tsx}', // 根据项目结构调整路径
      './views/**/*.{html,tmpl,gohtml}',  // 添加 Go 模板文件路径
    ],
    theme: {
      extend: {
        gridTemplateColumns: {
          'masonry': 'repeat(auto-fill, minmax(250px, 1fr))',
        },
      },
    },
    plugins: [],
  }