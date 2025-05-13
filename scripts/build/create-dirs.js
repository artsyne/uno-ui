/**
 * 创建JS源码目录
 *
 * 此脚本用于创建JS源码的目录结构，兼容Mac和Windows系统
 */

const fs = require('fs');
const path = require('path');

// 要创建的目录
const directories = [
  'website/static/src/js/tools',
  'website/static/src/js/app',
  'website/static/src/js/vendors',
];

// 确保目录存在
function ensureDir(dirPath) {
  const normalizedPath = path.normalize(dirPath);

  if (!fs.existsSync(normalizedPath)) {
    console.log(`创建目录: ${normalizedPath}`);
    fs.mkdirSync(normalizedPath, { recursive: true });
  } else {
    console.log(`目录已存在: ${normalizedPath}`);
  }
}

// 创建所有目录
function createDirs() {
  console.log('开始创建目录结构...');

  directories.forEach((dir) => {
    ensureDir(dir);
  });

  console.log('目录结构创建完成！');
}

// 执行目录创建
createDirs();
