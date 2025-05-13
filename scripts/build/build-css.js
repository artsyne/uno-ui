/**
 * CSS构建脚本
 *
 * 此脚本负责构建CSS文件，按以下步骤执行：
 * 1. 使用Tailwind CLI生成基础样式
 * 2. 合并自定义CSS与Tailwind CSS
 * 3. 使用PostCSS处理并压缩最终CSS
 *
 * 兼容Mac和Windows系统
 */

const { execSync } = require('child_process');
const fs = require('fs');
const os = require('os');
const path = require('path');

// 使用path.join确保路径在不同操作系统上的兼容性
const CSS_SRC_DIR = './website/static/src/css';
const CSS_DIST_DIR = './website/static/dist/css';
const TAILWIND_SRC = path.join(CSS_SRC_DIR, 'tailwind.css');
const THEME_CSS = path.join(CSS_SRC_DIR, 'diary-theme.css');
const TAILWIND_OUTPUT = path.join(CSS_DIST_DIR, 'tailwind.css');
const COMBINED_CSS = path.join(CSS_DIST_DIR, 'combined.css');
const FINAL_OUTPUT = path.join(CSS_DIST_DIR, 'idiary.min.css');

console.log('Starting CSS build process...');

// 确保输出目录存在
if (!fs.existsSync(CSS_DIST_DIR)) {
  console.log(`Creating directory: ${CSS_DIST_DIR}`);
  fs.mkdirSync(CSS_DIST_DIR, { recursive: true });
}

try {
  // Step 1: 生成Tailwind CSS
  console.log('Step 1: Generating tailwind.css...');
  execSync(
    `npx @tailwindcss/cli -i "${TAILWIND_SRC}" -o "${TAILWIND_OUTPUT}"`,
    { stdio: 'inherit', shell: true },
  );

  // Step 2: 合并CSS文件（使用Node.js文件操作代替shell命令）
  console.log('Step 2: Combining CSS files...');
  const themeContent = fs.readFileSync(THEME_CSS, 'utf8');
  const tailwindContent = fs.readFileSync(TAILWIND_OUTPUT, 'utf8');
  fs.writeFileSync(COMBINED_CSS, themeContent + '\n' + tailwindContent);

  // Step 3: 压缩CSS
  console.log('Step 3: Minifying CSS...');
  execSync(`npx postcss "${COMBINED_CSS}" -o "${FINAL_OUTPUT}" --no-map`, {
    stdio: 'inherit',
    shell: true,
  });

  // 清理临时文件
  console.log('Cleaning up temporary files...');
  if (fs.existsSync(COMBINED_CSS)) {
    fs.unlinkSync(COMBINED_CSS);
  }
  if (fs.existsSync(TAILWIND_OUTPUT)) {
    fs.unlinkSync(TAILWIND_OUTPUT);
  }

  console.log('CSS build process completed successfully.');
} catch (error) {
  console.error('Error during CSS build process:', error.message);
  process.exit(1);
}
