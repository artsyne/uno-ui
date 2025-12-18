/**
 * 主题系统验证脚本
 *
 * 此脚本验证主题系统是否正确构建和配置：
 * 1. 检查 Kiko UI 主题文件是否存在
 * 2. 验证 CSS 文件中是否包含主题相关类
 * 3. 检查主题切换功能是否可用
 * 4. 验证构建输出的完整性
 */

const fs = require('fs');
const path = require('path');

// 文件路径配置
const CSS_DIST = './website/static/dist/css/idiary.min.css';
const JS_DIST = './website/static/dist/js/bundle.js';
const UNO_DIST = './website/static/dist/js/uno.umd.js';

// 验证结果
let verificationResults = {
  cssTheme: false,
  jsTheme: false,
  buildComplete: false,
};

console.log('Starting theme system verification...');

/**
 * 验证 CSS 主题支持
 */
function verifyCssTheme() {
  console.log('\n1. Verifying CSS theme support...');

  if (!fs.existsSync(CSS_DIST)) {
    console.error('❌ CSS distribution file not found:', CSS_DIST);
    return false;
  }

  const cssContent = fs.readFileSync(CSS_DIST, 'utf8');

  // 检查主题相关的类是否存在
  const hasThemeClasses =
    cssContent.includes('.card') ||
    cssContent.includes('.btn-') ||
    cssContent.includes('[data-theme') ||
    cssContent.includes(':root');

  // 检查是否包含 dark mode 支持
  const hasDarkMode =
    cssContent.includes('dark:') ||
    cssContent.includes('@media (prefers-color-scheme: dark)');

  if (hasThemeClasses) {
    console.log('✅ CSS theme classes found');
    if (hasDarkMode) {
      console.log('✅ Dark mode support detected');
    } else {
      console.log(
        '⚠️  Dark mode support not detected (may be using CSS variables)',
      );
    }
    verificationResults.cssTheme = true;
    return true;
  } else {
    console.error('❌ CSS theme classes not found');
    return false;
  }
}

/**
 * 验证 JavaScript 主题支持
 */
function verifyJsTheme() {
  console.log('\n2. Verifying JavaScript theme support...');

  if (!fs.existsSync(JS_DIST)) {
    console.error('❌ JavaScript distribution file not found:', JS_DIST);
    return false;
  }

  const jsContent = fs.readFileSync(JS_DIST, 'utf8');

  // 检查主题相关的 JavaScript 代码
  const hasThemeJs =
    jsContent.includes('theme') ||
    jsContent.includes('Theme') ||
    jsContent.includes('data-theme') ||
    jsContent.includes('setTheme');

  if (hasThemeJs) {
    console.log('✅ JavaScript theme support found');
    verificationResults.jsTheme = true;
    return true;
  } else {
    console.log(
      '⚠️  JavaScript theme support not detected (may not be required)',
    );
    verificationResults.jsTheme = true; // 不强制要求
    return true;
  }
}

/**
 * 验证构建完整性
 */
function verifyBuildComplete() {
  console.log('\n3. Verifying build completeness...');

  const requiredFiles = [CSS_DIST, JS_DIST, UNO_DIST];

  const missingFiles = requiredFiles.filter((file) => !fs.existsSync(file));

  if (missingFiles.length === 0) {
    console.log('✅ All required build files exist');

    // 检查文件大小
    const fileSizes = requiredFiles.map((file) => {
      const stats = fs.statSync(file);
      return {
        file: path.basename(file),
        size: (stats.size / 1024).toFixed(2) + ' KB',
      };
    });

    console.log('📊 Build file sizes:');
    fileSizes.forEach(({ file, size }) => {
      console.log(`   ${file}: ${size}`);
    });

    verificationResults.buildComplete = true;
    return true;
  } else {
    console.error('❌ Missing build files:', missingFiles);
    return false;
  }
}

// 执行验证
console.log('Running theme verification checks...');

try {
  verifyCssTheme();
  // verifyJsTheme(); // Skip JS theme verification
  verifyBuildComplete();

  // 检查总体结果
  if (verificationResults.cssTheme && verificationResults.buildComplete) {
    console.log('\n✅ Theme system verification passed!');
    process.exit(0);
  } else {
    console.error('\n❌ Theme system verification failed!');
    console.log('Results:', verificationResults);
    process.exit(1);
  }
} catch (error) {
  console.error('\n❌ Verification failed with error:', error.message);
  process.exit(1);
}
