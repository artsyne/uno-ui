/**
 * JS构建脚本
 *
 * 此脚本负责构建JavaScript文件，按以下步骤执行：
 * 1. 分别处理工具、应用和第三方库JS文件
 * 2. 按正确顺序合并文件
 * 3. 使用terser压缩并混淆
 *
 * 兼容Mac和Windows系统
 *
 * 使用方法：
 * - 标准构建（压缩并移除调试语句）: node scripts/build/build-js.js
 * - 开发模式构建（保留console.log）: node scripts/build/build-js.js --dev
 */

const { execSync } = require('child_process');
const fs = require('fs');
const os = require('os');

// 检查是否为开发模式
const isDevMode = process.argv.includes('--dev');
console.log(isDevMode ? 
  'Starting JS build process in DEVELOPMENT mode (debug info will be preserved)...':
  'Starting JS build process in PRODUCTION mode (debug info will be removed)...');

// 创建输出目录（如果不存在）
console.log('Ensuring output directory exists...');
if (!fs.existsSync('./website/static/dist/js')) {
  fs.mkdirSync('./website/static/dist/js', { recursive: true });
}

// 确定是Windows还是Mac/Linux
const isWindows = os.platform() === 'win32';

// 第1步：合并vendor库文件
console.log('Step 1: Processing vendor files...');
const vendorFiles = isWindows
  ? '.\\website\\static\\src\\js\\vendors\\*.js'
  : './website/static/src/js/vendors/*.js';

// 第2步：合并工具文件
console.log('Step 2: Processing tool files...');
const toolFiles = isWindows
  ? '.\\website\\static\\src\\js\\tools\\*.js'
  : './website/static/src/js/tools/*.js';

// 第3步：合并应用文件
console.log('Step 3: Processing app files...');
const appFiles = isWindows
  ? '.\\website\\static\\src\\js\\app\\*.js'
  : './website/static/src/js/app/*.js';

// 使用terser压缩所有文件
console.log('Step 4: Minifying all JS files...');
try {
  // 为所有目录中的文件创建Glob模式列表
  const globPatterns = [vendorFiles, toolFiles, appFiles];

  // 设置压缩选项，根据模式决定是否保留console和debugger
const compressOptions = isDevMode
    ? 'drop_console=false,drop_debugger=false'
    : 'drop_console=true,drop_debugger=true';

  // 构建terser命令
  const terserCommand = `npx terser ${globPatterns.join(' ')} -o ${
    isWindows
      ? '.\\website\\static\\dist\\js\\idiary.min.js'
      : './website/static/dist/js/idiary.min.js'
  } --compress ${compressOptions} --mangle --comments false`;

  // 执行命令
  execSync(terserCommand, { stdio: 'inherit', shell: true });

  console.log('JS build process completed successfully.');
} catch (error) {
  console.error('Error during JS build:', error.message);
  process.exit(1);
}
