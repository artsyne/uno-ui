/**
 * CSS构建脚本
 *
 * 此脚本负责构建CSS文件，按以下步骤执行：
 * 1. 使用Tailwind CLI生成基础样式（支持主题系统）
 * 2. 合并自定义CSS与Tailwind CSS
 * 3. 使用PostCSS处理并压缩最终CSS
 * 4. 支持 Kiko UI 主题架构
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
const MARKDOWN_CSS = path.join(CSS_SRC_DIR, 'markdown.css');
const TAILWIND_OUTPUT = path.join(CSS_DIST_DIR, 'tailwind.css');
const COMBINED_CSS = path.join(CSS_DIST_DIR, 'combined.css');
const FINAL_OUTPUT = path.join(CSS_DIST_DIR, 'idiary.min.css');

// 检查构建模式
const isDevMode = process.argv.includes('--dev');
const buildMode = isDevMode ? 'DEVELOPMENT' : 'PRODUCTION';

console.log(`Building CSS in ${buildMode} mode...`);

console.log('Starting CSS build process...');

// 确保输出目录存在
if (!fs.existsSync(CSS_DIST_DIR)) {
  console.log(`Creating directory: ${CSS_DIST_DIR}`);
  fs.mkdirSync(CSS_DIST_DIR, { recursive: true });
}

try {
  // Step 1: 合并源CSS文件 (Tailwind + Markdown)
  // 注意：Theme CSS 不在此处合并，而是待 Tailwind 生成后再追加，以确保最高优先级
  console.log('Step 1: Combining source CSS files (Tailwind + Markdown)...');
  const themeContent = fs.readFileSync(THEME_CSS, 'utf8');
  const markdownContent = fs.readFileSync(MARKDOWN_CSS, 'utf8');
  const tailwindContent = fs.readFileSync(TAILWIND_SRC, 'utf8');

  // 仅合并 Tailwind 和 Markdown 供 CLI 处理
  fs.writeFileSync(COMBINED_CSS, tailwindContent + '\n' + markdownContent);

  // Step 2: 生成Tailwind CSS
  console.log('Step 2: Generating tailwind.css...');
  const contentPaths = [
    './website/views/**/*.html',
    './website/static/src/js/**/*.js',
    './uno/src/**/*.js',
  ];

  const contentFlags = contentPaths.map((p) => `--content "${p}"`).join(' ');

  const tailwindCmd = isDevMode
    ? `npx @tailwindcss/cli -i "${COMBINED_CSS}" -o "${TAILWIND_OUTPUT}" ${contentFlags} --watch=false`
    : `npx @tailwindcss/cli -i "${COMBINED_CSS}" -o "${TAILWIND_OUTPUT}" ${contentFlags} --minify`;

  execSync(tailwindCmd, { stdio: 'inherit', shell: true });

  // Step 2.5: 追加 Theme CSS
  // 强制将 Theme CSS 追加到 Tailwind 生成结果的末尾，物理上保证优先级
  console.log('Step 2.5: Appending Theme CSS to ensure highest priority...');
  const generatedTailwind = fs.readFileSync(TAILWIND_OUTPUT, 'utf8');
  fs.writeFileSync(TAILWIND_OUTPUT, generatedTailwind + '\n' + themeContent);

  // Step 3: 处理和压缩CSS
  console.log('Step 3: Processing CSS with PostCSS...');
  const postcssCmd = isDevMode
    ? `npx postcss "${TAILWIND_OUTPUT}" -o "${FINAL_OUTPUT}" --no-map --env development`
    : `npx postcss "${TAILWIND_OUTPUT}" -o "${FINAL_OUTPUT}" --no-map --env production`;

  execSync(postcssCmd, {
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
