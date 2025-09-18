/**
 * JS构建脚本
 *
 * 此脚本负责构建JavaScript文件，按以下步骤执行：
 * 1. 处理网站的JS文件（按指定顺序合并）
 * 2. 生成两个独立的JS文件:
 *    - kiko.min.js: Kiko UI框架文件
 *    - idiary.min.js: 业务逻辑文件
 * 3. 使用terser压缩并混淆
 *
 * 兼容Mac和Windows系统
 *
 * 使用方法：
 * - 生产模式构建（压缩并移除调试语句）: node scripts/build/build-js.js
 * - 开发模式构建（保留console.log并使用时间戳）: node scripts/build/build-js.js --dev
 */

const { execSync } = require('child_process');
const fs = require('fs');
const os = require('os');
const path = require('path');

// 检查构建模式
const isDevMode = process.argv.includes('--dev');
const isProdMode = process.argv.includes('--prod');

// 如果没有指定模式，默认为开发模式
const buildMode = isProdMode ? 'PRODUCTION' : 'DEVELOPMENT';
const shouldPreserveDebug = !isProdMode;

console.log(
  `Starting JS build process in ${buildMode} mode (debug info will be ${shouldPreserveDebug ? 'preserved' : 'removed'})...`,
);

// 创建输出目录（如果不存在）
console.log('Ensuring output directory exists...');
if (!fs.existsSync('./website/static/dist/js')) {
  fs.mkdirSync('./website/static/dist/js', { recursive: true });
}

// 确定是Windows还是Mac/Linux
const isWindows = os.platform() === 'win32';

// 语言文件目录
const i18nDir = isWindows
  ? '.\\website\\static\\src\\js\\locales'
  : './website/static/src/js/locales';

// 获取所有语言文件（已废弃，语言包现在通过 JSON 动态加载）
const getLanguageFiles = () => {
  // 不再打包语言文件到 JS 中，改为动态加载 JSON
  return [];
};

// 定义JS文件合并顺序
const jsFiles = [
  'kiko-init.js', // 首先加载Kiko初始化
  'i18n.js', // 网站多语言初始化
  'theme.js', // 主题切换功能
  'user-state.js', // 用户状态管理模块
  'auth.js', // 认证模块
  'diary.js', // 日记本模块
  'settings.js', // 设置模块
  'app.js', // 最后是主应用
];

// 源文件和临时文件路径
const srcJsPath = isWindows
  ? '.\\website\\static\\src\\js'
  : './website/static/src/js';

// 临时idiary.js文件路径
const tempIdiaryPath = isWindows
  ? '.\\website\\static\\dist\\js\\idiary.temp.js'
  : './website/static/dist/js/idiary.temp.js';

// 最终输出文件路径
const idiaryOutputPath = path.join(
  __dirname,
  '..',
  '..',
  'website',
  'static',
  'dist',
  'js',
  'idiary.min.js',
);

const kikoOutputPath = path.join(
  __dirname,
  '..',
  '..',
  'website',
  'static',
  'dist',
  'js',
  'kiko.min.js',
);

// 使用terser压缩所有文件
console.log('Processing and minifying website JS files...');
try {
  // 设置压缩选项，根据模式决定是否保留console和debugger
  const compressOptions = shouldPreserveDebug
    ? 'drop_console=false,drop_debugger=false'
    : 'drop_console=true,drop_debugger=true';

  // 首先合并所有JS文件
  let combinedJs = '';
  for (const file of jsFiles) {
    const filePath = path.join(srcJsPath, file);
    if (fs.existsSync(filePath)) {
      console.log(`Adding ${file}...`);
      const content = fs.readFileSync(filePath, 'utf8');
      combinedJs += `\n// ---- ${file} ----\n${content}\n`;
    } else {
      console.warn(`警告: 文件不存在: ${filePath}`);
    }
  }

  // 写入临时合并文件
  fs.writeFileSync(tempIdiaryPath, combinedJs, 'utf8');

  // 构建terser命令压缩网站JS文件
  const terserIdiaryCommand = `npx terser ${tempIdiaryPath} -o ${idiaryOutputPath} --compress ${compressOptions} --mangle --comments false`;

  // 执行网站JS压缩命令
  console.log('Minifying idiary JS...');
  execSync(terserIdiaryCommand, { stdio: 'inherit', shell: true });

  // 获取Kiko UI的JS文件路径
  const kikoSource = path.join(
    __dirname,
    '..',
    '..',
    'kiko',
    'dist',
    'kiko.min.js', 
  );

  // 拷贝Kiko UI文件（不再合并）
  console.log('Copying Kiko UI framework file...');
  fs.copyFileSync(kikoSource, kikoOutputPath);

  // 读取package.json获取版本号（用于添加注释）
  const packageJsonPath = path.join(__dirname, '..', '..', 'package.json');
  let appVersion = '1.1.0'; // 默认版本号

  if (fs.existsSync(packageJsonPath)) {
    try {
      const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
      if (packageJson.version) {
        appVersion = packageJson.version;
      }
    } catch (err) {
      console.warn(
        'Could not parse package.json, using default version: ' + appVersion,
      );
    }
  }

  // 版本号：开发模式使用时间戳，生产模式使用读取的版本号
  const version = shouldPreserveDebug ? 'dev-' + new Date().getTime() : appVersion;
  const buildDate = new Date().toISOString();

  // 添加版本注释到文件
  const addVersionComment = (filePath, fileType) => {
    const content = fs.readFileSync(filePath, 'utf8');
    const commentedContent = `/* ${fileType} - Version ${version} - Built ${buildDate} */\n${content}`;
    fs.writeFileSync(filePath, commentedContent, 'utf8');
  };

  // 添加版本注释
  addVersionComment(idiaryOutputPath, 'iDiary JS');
  addVersionComment(kikoOutputPath, 'Kiko UI');

  console.log(`Created files: 
  - ${kikoOutputPath} (framework)
  - ${idiaryOutputPath} (application)`);

  // 删除临时文件
  fs.unlinkSync(tempIdiaryPath);
  console.log('Removed temporary files');

  // 处理语言包文件 - 合并 JS 语言包到 idiary.min.js 中
  console.log('Processing and merging locale JS files...');
  const localesSourceDir = path.join(__dirname, '..', '..', 'website', 'static', 'src', 'js', 'locales');
  
  let localeJsCode = '';
  
  // 读取并合并语言包 JS 文件
  if (fs.existsSync(localesSourceDir)) {
    const localeFiles = fs.readdirSync(localesSourceDir).filter(file => file.endsWith('.js'));
    
    for (const file of localeFiles) {
      const sourcePath = path.join(localesSourceDir, file);
      
      console.log(`Processing locale file: ${file}`);
      
      // 读取 JS 文件内容
      const content = fs.readFileSync(sourcePath, 'utf8');
      localeJsCode += `\n\n/* ${file} */\n${content}`;
    }
    
    console.log(`Processed and merged ${localeFiles.length} locale JS files`);
  } else {
    console.warn('Locales source directory not found, skipping locale processing');
  }
  
  // 将语言包 JS 代码追加到 idiary.min.js
  if (localeJsCode) {
    const currentContent = fs.readFileSync(idiaryOutputPath, 'utf8');
    fs.writeFileSync(idiaryOutputPath, currentContent + localeJsCode, 'utf8');
    console.log('Language pack JS files merged into idiary.min.js');
  }

  console.log('JS build process completed successfully.');
} catch (error) {
  console.error('Error during JS build:', error.message);
  process.exit(1);
}
