# 项目技能指南 (Project Skills Guide)

本文档列出了项目中集成的 AI 辅助技能 (AI Skills)，这些技能旨在增强开发效率和代码质量。

## 技能列表

### 1. Frontend Design (`frontend-design`)

**描述**: 创建具有高设计质量、独特且生产级的前端界面。避免通用的 AI 生成美学。

**使用场景**:
- 构建 Web 组件、页面、应用程序。
- 美化现有 UI。
- 需要独特设计风格时（如极简主义、极繁主义、复古未来主义等）。

**配置位置**: `.trae/skills/frontend-design/SKILL.md`

**核心原则**:
- **Typography**: 使用独特的字体。
- **Color**: 大胆的配色方案。
- **Motion**: 强调 CSS 动画和微交互。
- **Spatial**: 打破常规的布局。

**禁用项**:
- 通用字体 (Arial, Inter)。
- 刻板的配色 (紫色渐变)。
- 模板化的布局。

### 2. Brand Guidelines (`brand-guidelines`)

**描述**: 应用项目官方的极简主义设计系统，严格遵循 Tailwind CSS 4 实用类规范。

**使用场景**:
- 确保新页面或组件符合官方视觉风格。
- 修复 UI 样式不一致的问题。
- 生成符合 Dark Mode 规范的代码。

**配置位置**: `.trae/skills/brand-guidelines/SKILL.md`

**核心原则**:
- **Minimalism**: 内容至上，UI 退后。
- **Grayscale**: 基于灰度构建结构，强调边框和间距。
- **Utility-First**: 禁止自定义 CSS，完全依赖 Tailwind Utility。
- **Dark Mode**: 必须支持深色模式。

## 如何添加新技能

1. 在 `.trae/skills/` 下创建新目录。
2. 添加 `SKILL.md` 定义技能行为。
3. 在 `.trae/skills.yaml` 中注册技能。
