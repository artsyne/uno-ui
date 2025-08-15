# 主题系统测试指南

本文档提供了 iDiary.ai Hebe 项目主题系统的完整测试方案和实施指南。

## 测试概述

主题系统测试分为以下几个层面：
1. **单元测试** - Kiko UI 组件和主题管理器
2. **集成测试** - 主题系统与应用的集成
3. **端到端测试** - 完整的用户体验测试
4. **视觉回归测试** - 主题切换的视觉一致性
5. **性能测试** - 主题切换的性能影响

## 1. 单元测试

### 1.1 Kiko UI 主题管理器测试

**测试文件**: `tests/unit/kiko-theme.test.js`

```javascript
// 主题管理器单元测试
describe('Kiko Theme Manager', () => {
  let themeManager;
  
  beforeEach(() => {
    // 重置 DOM 和本地存储
    document.body.innerHTML = '';
    localStorage.clear();
    themeManager = new Kiko.ThemeManager();
  });
  
  describe('初始化', () => {
    test('应该正确初始化默认主题', () => {
      expect(themeManager.getCurrentTheme()).toBe('light');
    });
    
    test('应该从本地存储恢复主题', () => {
      localStorage.setItem('theme', 'dark');
      const tm = new Kiko.ThemeManager();
      expect(tm.getCurrentTheme()).toBe('dark');
    });
    
    test('应该检测系统主题偏好', () => {
      // 模拟系统偏好
      Object.defineProperty(window, 'matchMedia', {
        writable: true,
        value: jest.fn().mockImplementation(query => ({
          matches: query === '(prefers-color-scheme: dark)',
          media: query,
          onchange: null,
          addListener: jest.fn(),
          removeListener: jest.fn(),
        })),
      });
      
      const tm = new Kiko.ThemeManager({ autoDetect: true });
      expect(tm.getCurrentTheme()).toBe('dark');
    });
  });
  
  describe('主题切换', () => {
    test('应该正确切换到 dark 主题', () => {
      themeManager.setTheme('dark');
      expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
      expect(themeManager.getCurrentTheme()).toBe('dark');
    });
    
    test('应该正确切换到 light 主题', () => {
      themeManager.setTheme('light');
      expect(document.documentElement.getAttribute('data-theme')).toBe('light');
      expect(themeManager.getCurrentTheme()).toBe('light');
    });
    
    test('应该触发主题变更事件', () => {
      const callback = jest.fn();
      themeManager.on('themeChanged', callback);
      
      themeManager.setTheme('dark');
      expect(callback).toHaveBeenCalledWith('dark', 'light');
    });
    
    test('应该保存主题到本地存储', () => {
      themeManager.setTheme('dark');
      expect(localStorage.getItem('theme')).toBe('dark');
    });
  });
  
  describe('主题切换动画', () => {
    test('应该添加切换动画类', async () => {
      themeManager.setTheme('dark');
      expect(document.documentElement.classList.contains('theme-transitioning')).toBe(true);
      
      // 等待动画完成
      await new Promise(resolve => setTimeout(resolve, 300));
      expect(document.documentElement.classList.contains('theme-transitioning')).toBe(false);
    });
  });
});
```

### 1.2 Kiko UI 组件主题测试

**测试文件**: `tests/unit/kiko-components.test.js`

```javascript
// 组件主题测试
describe('Kiko UI Components Theme Support', () => {
  let container;
  
  beforeEach(() => {
    container = document.createElement('div');
    document.body.appendChild(container);
  });
  
  afterEach(() => {
    document.body.removeChild(container);
  });
  
  describe('Notify 组件', () => {
    test('应该应用正确的主题样式', () => {
      const notify = new Kiko.Notify({
        message: 'Test message',
        type: 'success'
      });
      
      // 测试 light 主题
      document.documentElement.setAttribute('data-theme', 'light');
      notify.show();
      
      const notifyEl = document.querySelector('.kiko-notify');
      expect(notifyEl).toHaveClass('theme-aware');
      
      // 测试 dark 主题
      document.documentElement.setAttribute('data-theme', 'dark');
      // 组件应该自动更新样式
      expect(notifyEl).toHaveClass('theme-aware');
    });
  });
  
  describe('Box 组件', () => {
    test('应该支持主题感知的卡片样式', () => {
      container.innerHTML = '<div class="card">Test content</div>';
      const box = container.querySelector('.card');
      
      // 验证主题类是否正确应用
      expect(box).toHaveClass('card');
      
      // 测试主题切换
      document.documentElement.setAttribute('data-theme', 'dark');
      // CSS 变量应该自动更新
      const styles = getComputedStyle(box);
      expect(styles.getPropertyValue('--card-bg')).toBeTruthy();
    });
  });
});
```

## 2. 集成测试

### 2.1 主题系统集成测试

**测试文件**: `tests/integration/theme-integration.test.js`

```javascript
// 主题系统集成测试
describe('Theme System Integration', () => {
  beforeEach(() => {
    // 加载完整的应用环境
    document.head.innerHTML = `
      <link rel="stylesheet" href="/static/dist/css/idiary.min.css">
    `;
    document.body.innerHTML = `
      <div id="app">
        <header class="card">
          <h1>iDiary.ai</h1>
          <button id="theme-toggle" class="btn-secondary">切换主题</button>
        </header>
        <main class="card">
          <p>测试内容</p>
        </main>
      </div>
    `;
    
    // 初始化应用
    window.App.init();
  });
  
  test('应用启动时应该正确初始化主题', () => {
    expect(document.documentElement.hasAttribute('data-theme')).toBe(true);
    expect(['light', 'dark']).toContain(
      document.documentElement.getAttribute('data-theme')
    );
  });
  
  test('主题切换按钮应该正常工作', () => {
    const toggleBtn = document.getElementById('theme-toggle');
    const initialTheme = document.documentElement.getAttribute('data-theme');
    
    // 点击切换按钮
    toggleBtn.click();
    
    const newTheme = document.documentElement.getAttribute('data-theme');
    expect(newTheme).not.toBe(initialTheme);
    expect(['light', 'dark']).toContain(newTheme);
  });
  
  test('所有卡片组件应该响应主题变更', () => {
    const cards = document.querySelectorAll('.card');
    expect(cards.length).toBeGreaterThan(0);
    
    // 切换主题
    window.Kiko.ThemeManager.getInstance().setTheme('dark');
    
    // 验证所有卡片都有正确的主题属性
    cards.forEach(card => {
      const styles = getComputedStyle(card);
      expect(styles.getPropertyValue('--card-bg')).toBeTruthy();
    });
  });
  
  test('主题偏好应该持久化', () => {
    window.Kiko.ThemeManager.getInstance().setTheme('dark');
    expect(localStorage.getItem('theme')).toBe('dark');
    
    // 模拟页面重新加载
    window.App.init();
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
  });
});
```

## 3. 端到端测试

### 3.1 用户体验测试

**测试文件**: `tests/e2e/theme-e2e.test.js`

```javascript
// 使用 Playwright 或 Cypress 进行端到端测试
describe('Theme System E2E Tests', () => {
  beforeEach(async () => {
    await page.goto('http://localhost:8080');
  });
  
  test('用户可以通过界面切换主题', async () => {
    // 等待页面加载
    await page.waitForSelector('[data-theme]');
    
    // 获取初始主题
    const initialTheme = await page.getAttribute('html', 'data-theme');
    
    // 点击主题切换按钮
    await page.click('#theme-toggle');
    
    // 等待主题切换完成
    await page.waitForTimeout(300);
    
    // 验证主题已切换
    const newTheme = await page.getAttribute('html', 'data-theme');
    expect(newTheme).not.toBe(initialTheme);
  });
  
  test('主题切换应该有平滑的动画效果', async () => {
    // 点击主题切换
    await page.click('#theme-toggle');
    
    // 验证动画类被添加
    const hasTransitionClass = await page.evaluate(() => {
      return document.documentElement.classList.contains('theme-transitioning');
    });
    expect(hasTransitionClass).toBe(true);
    
    // 等待动画完成
    await page.waitForTimeout(300);
    
    // 验证动画类被移除
    const stillHasTransitionClass = await page.evaluate(() => {
      return document.documentElement.classList.contains('theme-transitioning');
    });
    expect(stillHasTransitionClass).toBe(false);
  });
  
  test('所有页面都应该支持主题切换', async () => {
    const pages = [
      '/',
      '/diary',
      '/member/settings',
      '/member'
    ];
    
    for (const pagePath of pages) {
      await page.goto(`http://localhost:8080${pagePath}`);
      await page.waitForSelector('[data-theme]');
      
      // 切换主题
      await page.click('#theme-toggle');
      await page.waitForTimeout(300);
      
      // 验证主题已应用
      const theme = await page.getAttribute('html', 'data-theme');
      expect(['light', 'dark']).toContain(theme);
      
      // 验证页面样式正确
      const cards = await page.$$('.card');
      expect(cards.length).toBeGreaterThan(0);
    }
  });
  
  test('主题偏好在页面刷新后应该保持', async () => {
    // 切换到 dark 主题
    await page.click('#theme-toggle');
    await page.waitForTimeout(300);
    
    const themeBeforeReload = await page.getAttribute('html', 'data-theme');
    
    // 刷新页面
    await page.reload();
    await page.waitForSelector('[data-theme]');
    
    // 验证主题保持不变
    const themeAfterReload = await page.getAttribute('html', 'data-theme');
    expect(themeAfterReload).toBe(themeBeforeReload);
  });
});
```

## 4. 视觉回归测试

### 4.1 截图对比测试

**测试文件**: `tests/visual/theme-visual.test.js`

```javascript
// 视觉回归测试
describe('Theme Visual Regression Tests', () => {
  const pages = [
    { name: 'home', path: '/' },
    { name: 'diary', path: '/diary' },
    { name: 'settings', path: '/member/settings' },
    { name: 'profile', path: '/member' }
  ];
  
  const themes = ['light', 'dark'];
  
  themes.forEach(theme => {
    pages.forEach(({ name, path }) => {
      test(`${name} page should look correct in ${theme} theme`, async () => {
        await page.goto(`http://localhost:8080${path}`);
        await page.waitForSelector('[data-theme]');
        
        // 设置主题
        await page.evaluate((theme) => {
          window.Kiko.ThemeManager.getInstance().setTheme(theme);
        }, theme);
        
        await page.waitForTimeout(300); // 等待主题切换完成
        
        // 截图对比
        const screenshot = await page.screenshot({
          fullPage: true,
          clip: { x: 0, y: 0, width: 1200, height: 800 }
        });
        
        expect(screenshot).toMatchImageSnapshot({
          customSnapshotIdentifier: `${name}-${theme}`,
          threshold: 0.2,
          customDiffConfig: {
            threshold: 0.1
          }
        });
      });
    });
  });
  
  test('主题切换动画应该平滑', async () => {
    await page.goto('http://localhost:8080');
    await page.waitForSelector('[data-theme]');
    
    // 录制主题切换过程
    await page.startScreencast({ path: 'theme-transition.webm' });
    
    await page.click('#theme-toggle');
    await page.waitForTimeout(500);
    
    await page.stopScreencast();
    
    // 验证动画文件存在
    const fs = require('fs');
    expect(fs.existsSync('theme-transition.webm')).toBe(true);
  });
});
```

## 5. 性能测试

### 5.1 主题切换性能测试

**测试文件**: `tests/performance/theme-performance.test.js`

```javascript
// 性能测试
describe('Theme Performance Tests', () => {
  test('主题切换应该在 100ms 内完成', async () => {
    await page.goto('http://localhost:8080');
    await page.waitForSelector('[data-theme]');
    
    // 测量主题切换时间
    const startTime = Date.now();
    
    await page.evaluate(() => {
      window.Kiko.ThemeManager.getInstance().setTheme('dark');
    });
    
    await page.waitForFunction(() => {
      return document.documentElement.getAttribute('data-theme') === 'dark';
    });
    
    const endTime = Date.now();
    const switchTime = endTime - startTime;
    
    expect(switchTime).toBeLessThan(100);
  });
  
  test('主题切换不应该导致布局抖动', async () => {
    await page.goto('http://localhost:8080');
    await page.waitForSelector('[data-theme]');
    
    // 获取初始布局
    const initialLayout = await page.evaluate(() => {
      const elements = document.querySelectorAll('.card');
      return Array.from(elements).map(el => ({
        width: el.offsetWidth,
        height: el.offsetHeight,
        top: el.offsetTop,
        left: el.offsetLeft
      }));
    });
    
    // 切换主题
    await page.click('#theme-toggle');
    await page.waitForTimeout(300);
    
    // 获取切换后的布局
    const newLayout = await page.evaluate(() => {
      const elements = document.querySelectorAll('.card');
      return Array.from(elements).map(el => ({
        width: el.offsetWidth,
        height: el.offsetHeight,
        top: el.offsetTop,
        left: el.offsetLeft
      }));
    });
    
    // 验证布局没有变化
    expect(newLayout).toEqual(initialLayout);
  });
  
  test('大量元素的主题切换性能', async () => {
    // 创建包含大量元素的测试页面
    await page.goto('http://localhost:8080/test/performance');
    await page.waitForSelector('[data-theme]');
    
    // 确保页面有足够多的元素
    const elementCount = await page.$$eval('.card', els => els.length);
    expect(elementCount).toBeGreaterThan(100);
    
    // 测量切换时间
    const metrics = await page.metrics();
    const startTime = metrics.Timestamp;
    
    await page.click('#theme-toggle');
    await page.waitForTimeout(300);
    
    const endMetrics = await page.metrics();
    const endTime = endMetrics.Timestamp;
    
    const switchTime = (endTime - startTime) * 1000; // 转换为毫秒
    expect(switchTime).toBeLessThan(200); // 大量元素下也应该在 200ms 内完成
  });
});
```

## 6. 测试环境配置

### 6.1 Jest 配置

**文件**: `jest.config.js`

```javascript
module.exports = {
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/tests/setup.js'],
  testMatch: [
    '<rootDir>/tests/unit/**/*.test.js',
    '<rootDir>/tests/integration/**/*.test.js'
  ],
  collectCoverageFrom: [
    'kiko/src/**/*.js',
    'website/static/src/js/**/*.js',
    '!**/*.min.js'
  ],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80
    }
  },
  moduleNameMapping: {
    '^@/(.*)$': '<rootDir>/website/static/src/$1',
    '^@kiko/(.*)$': '<rootDir>/kiko/src/$1'
  }
};
```

### 6.2 测试设置文件

**文件**: `tests/setup.js`

```javascript
// 测试环境设置
import 'jest-dom/extend-expect';

// 模拟 localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};
global.localStorage = localStorageMock;

// 模拟 matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(),
    removeListener: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

// 模拟 CSS 变量
Object.defineProperty(document.documentElement.style, 'setProperty', {
  value: jest.fn(),
});

Object.defineProperty(document.documentElement.style, 'getPropertyValue', {
  value: jest.fn().mockReturnValue(''),
});

// 加载 Kiko UI 用于测试
require('../kiko/dist/kiko.js');
```

### 6.3 Playwright 配置

**文件**: `playwright.config.js`

```javascript
module.exports = {
  testDir: './tests/e2e',
  timeout: 30000,
  expect: {
    timeout: 5000
  },
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:8080',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure'
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    }
  ],
  webServer: {
    command: 'npm run dev',
    port: 8080,
    reuseExistingServer: !process.env.CI
  }
};
```

## 7. 测试执行指南

### 7.1 运行所有测试

```bash
# 运行单元测试
npm run test:unit

# 运行集成测试
npm run test:integration

# 运行端到端测试
npm run test:e2e

# 运行视觉回归测试
npm run test:visual

# 运行性能测试
npm run test:performance

# 运行所有测试
npm test
```

### 7.2 测试覆盖率

```bash
# 生成覆盖率报告
npm run test:coverage

# 查看覆盖率报告
open coverage/lcov-report/index.html
```

### 7.3 持续集成

**文件**: `.github/workflows/theme-tests.yml`

```yaml
name: Theme System Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    strategy:
      matrix:
        node-version: [18.x, 20.x]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Use Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v3
      with:
        node-version: ${{ matrix.node-version }}
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Build project
      run: npm run build:prod
    
    - name: Verify theme system
      run: npm run verify:theme
    
    - name: Run unit tests
      run: npm run test:unit
    
    - name: Run integration tests
      run: npm run test:integration
    
    - name: Install Playwright
      run: npx playwright install
    
    - name: Run E2E tests
      run: npm run test:e2e
    
    - name: Upload test results
      uses: actions/upload-artifact@v3
      if: failure()
      with:
        name: test-results
        path: test-results/
```

## 8. 测试最佳实践

### 8.1 测试编写原则

1. **测试应该独立** - 每个测试都应该能够独立运行
2. **测试应该快速** - 单元测试应该在毫秒级完成
3. **测试应该可重复** - 相同的测试应该产生相同的结果
4. **测试应该有意义** - 测试应该验证实际的业务需求
5. **测试应该易于维护** - 测试代码应该清晰易懂

### 8.2 测试数据管理

```javascript
// 使用测试工厂函数
const createTestTheme = (overrides = {}) => ({
  name: 'test-theme',
  colors: {
    primary: '#007bff',
    secondary: '#6c757d',
    ...overrides.colors
  },
  ...overrides
});

// 使用测试夹具
const testFixtures = {
  lightTheme: createTestTheme({ name: 'light' }),
  darkTheme: createTestTheme({ 
    name: 'dark',
    colors: {
      primary: '#0d6efd',
      secondary: '#6c757d'
    }
  })
};
```

### 8.3 错误处理测试

```javascript
// 测试错误情况
test('应该处理无效的主题名称', () => {
  const themeManager = new Kiko.ThemeManager();
  
  expect(() => {
    themeManager.setTheme('invalid-theme');
  }).toThrow('Invalid theme name: invalid-theme');
});

// 测试网络错误
test('应该处理主题资源加载失败', async () => {
  // 模拟网络错误
  global.fetch = jest.fn().mockRejectedValue(new Error('Network error'));
  
  const themeManager = new Kiko.ThemeManager();
  await expect(themeManager.loadTheme('custom-theme')).rejects.toThrow('Network error');
});
```

## 9. 故障排除

### 9.1 常见测试问题

1. **测试超时**
   - 检查异步操作是否正确等待
   - 增加超时时间
   - 使用 `waitFor` 等待条件满足

2. **样式测试失败**
   - 确保 CSS 文件已正确加载
   - 检查 CSS 变量是否正确设置
   - 验证浏览器兼容性

3. **主题切换测试不稳定**
   - 添加适当的等待时间
   - 检查动画完成状态
   - 使用事件监听而非固定延时

### 9.2 调试技巧

```javascript
// 在测试中添加调试信息
test('主题切换调试', async () => {
  console.log('Initial theme:', document.documentElement.getAttribute('data-theme'));
  
  await page.click('#theme-toggle');
  
  console.log('After click:', document.documentElement.getAttribute('data-theme'));
  
  // 截图用于调试
  await page.screenshot({ path: 'debug-theme-switch.png' });
});
```

## 10. 测试报告

### 10.1 生成测试报告

```bash
# 生成 HTML 测试报告
npm run test:report

# 生成性能报告
npm run test:performance:report

# 生成视觉回归报告
npm run test:visual:report
```

### 10.2 测试指标

- **代码覆盖率**: 目标 > 80%
- **测试通过率**: 目标 > 95%
- **性能指标**: 主题切换 < 100ms
- **视觉一致性**: 像素差异 < 0.2%

---

通过遵循这个测试指南，可以确保主题系统的稳定性、性能和用户体验质量。定期运行这些测试，并根据项目发展持续更新测试用例。