/**
 * 组件级 axe 无障碍扫描助手（组件门禁，docs/a11y-design.md §3.5）
 *
 * 设计取舍：
 * - 以真实 Alpine 运行时驱动 uno 组件（jsdom 环境），扫描"水合后"的真实 DOM，
 *   而非静态模板字符串——`x-t.aria-label`、`:aria-modal`、`role=tab`（指令在
 *   运行时 setAttribute）在静态扫描中不可见，会制造 button-name 等误报。
 * - 本地实现 `toHaveNoViolations` 匹配器（与 jest-axe/vitest-axe 语义一致），
 *   不引入仅适配旧版 vitest 的包装包，避免版本兼容风险。
 * - 组件片段非完整页面：文档级规则（region/landmark/bypass/html-has-lang/
 *   document-title/page-has-heading-one）属宿主页职责，组件扫描关闭；
 *   color-contrast 依赖真实布局引擎（jsdom 无），由《对比度专项》人工取样覆盖。
 */
import { expect } from 'vitest'
import axe from 'axe-core'

export const COMPONENT_AXE_RULES = {
  'color-contrast': { enabled: false },
  region: { enabled: false },
  'landmark-one-main': { enabled: false },
  'page-has-heading-one': { enabled: false },
  bypass: { enabled: false },
  'html-has-lang': { enabled: false },
  'document-title': { enabled: false }
}

expect.extend({
  toHaveNoViolations(results) {
    const violations = results?.violations ?? []
    const pass = violations.length === 0
    return {
      pass,
      message: () =>
        pass
          ? '期望存在 axe 违规，但扫描结果为零违规'
          : `${violations.length} 个 axe 违规：\n` +
            violations
              .map(
                (v) =>
                  `- ${v.id} [${v.impact}] ${v.help}\n  ` +
                  v.nodes.map((n) => n.target.join(' ')).join('\n  ')
              )
              .join('\n')
    }
  }
})

/** 等待 Alpine 效应与过渡定时器落定 */
export const tick = (ms = 50) => new Promise((resolve) => setTimeout(resolve, ms))

/**
 * 启动真实 Alpine + uno（须在 jsdom 环境的测试文件中使用）。
 * vitest 默认按文件隔离环境，各测试文件互不污染。
 *
 * @returns {Promise<{Alpine: Object, Uno: Object}>}
 */
export const bootUno = async () => {
  const Alpine = (await import('alpinejs')).default
  window.Alpine = Alpine
  const { Uno } = await import('../../src/index')
  Alpine.start()
  await tick()
  return { Alpine, Uno }
}

/**
 * 对给定容器执行 axe 扫描（组件级规则集）。
 * 内置有效性守卫：引擎未分析到任何节点时抛错，防"扫描空 DOM 而假绿"。
 *
 * @param {Element|Document} [context] - 扫描上下文，默认整个文档体
 * @returns {Promise<Object>} axe 原始结果，供 toHaveNoViolations 断言
 */
export const scanA11y = async (context = document.body) => {
  const results = await axe.run(context, {
    resultTypes: ['violations', 'passes'],
    rules: COMPONENT_AXE_RULES
  })
  if (results.violations.length === 0 && results.passes.length === 0) {
    throw new Error('axe 引擎未分析到任何节点，请检查 jsdom 环境与扫描上下文')
  }
  return results
}
