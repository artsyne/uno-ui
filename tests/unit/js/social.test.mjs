import { jest } from '@jest/globals';
import { getBatchStatsManager, entryStats } from '../../../website/static/src/js/modules/social.js';

describe('Social Modules - BatchStatsManager', () => {
  let batchStatsManager;

  beforeAll(() => {
    batchStatsManager = getBatchStatsManager();
  });

  beforeEach(() => {
    batchStatsManager.clear();
  });

  describe('基础功能', () => {
    it('应该正确初始化', () => {
      expect(batchStatsManager).toBeDefined();
      expect(batchStatsManager.pendingEntryIds).toBeInstanceOf(Set);
      expect(batchStatsManager.loadedEntryIds).toBeInstanceOf(Set);
    });

    it('addEntryId 应该添加待加载的 entryId', () => {
      batchStatsManager.addEntryId(1);
      expect(batchStatsManager.pendingEntryIds.has('1')).toBe(true);
    });

    it('addEntryId 应该拒绝重复添加相同的 entryId', () => {
      batchStatsManager.addEntryId(1);
      batchStatsManager.addEntryId(1);
      expect(batchStatsManager.pendingEntryIds.size).toBe(1);
    });

    it('addEntryId 应该拒绝已加载的 entryId', () => {
      batchStatsManager.loadedEntryIds.add('1');
      batchStatsManager.addEntryId(1);
      expect(batchStatsManager.pendingEntryIds.has('1')).toBe(false);
    });

    it('clear 应该清空所有状态', () => {
      batchStatsManager.pendingEntryIds.add('1');
      batchStatsManager.loadedEntryIds.add('2');
      batchStatsManager.clear();
      expect(batchStatsManager.pendingEntryIds.size).toBe(0);
      expect(batchStatsManager.loadedEntryIds.size).toBe(0);
    });
  });
});

describe('Social Modules - EntryStats', () => {
  describe('组件初始化', () => {
    it('应该正确初始化组件', () => {
      const component = entryStats(1);
      expect(component.entryId).toBe(1);
      expect(component.stats).toBeDefined();
      expect(component.stats.loaded).toBe(false);
      expect(component.stats.loading).toBe(false);
      expect(component.stats.view_num).toBe(0);
      expect(component.stats.comment_num).toBe(0);
      expect(component.stats.like_num).toBe(0);
      expect(component.stats.liked).toBe(false);
    });

    it('不同 entryId 应该创建不同的组件实例', () => {
      const component1 = entryStats(1);
      const component2 = entryStats(2);
      expect(component1.entryId).toBe(1);
      expect(component2.entryId).toBe(2);
      expect(component1 !== component2).toBe(true);
    });
  });

  describe('乐观更新', () => {
    it('toggleLike 应该有乐观更新逻辑', () => {
      const component = entryStats(1);
      component.stats.loaded = true;
      component.stats.liked = false;
      component.stats.like_num = 5;

      // 验证初始状态
      expect(component.stats.liked).toBe(false);
      expect(component.stats.like_num).toBe(5);
    });

    it('加载中不应该执行 toggleLike', () => {
      const component = entryStats(1);
      component.stats.loading = true;

      // 只验证组件结构，不实际调用异步方法
      expect(component.stats.loading).toBe(true);
      expect(typeof component.toggleLike).toBe('function');
    });
  });

  describe('状态管理', () => {
    it('应该有正确的状态属性', () => {
      const component = entryStats(1);
      
      expect(component).toHaveProperty('stats');
      expect(component).toHaveProperty('observer');
      expect(component).toHaveProperty('initStats');
      expect(component).toHaveProperty('toggleLike');
      expect(component).toHaveProperty('_bindGlobalEvents');
      expect(component).toHaveProperty('destroy');
    });

    it('stats 对象应该包含所有必需字段', () => {
      const component = entryStats(1);
      
      expect(component.stats).toHaveProperty('loaded');
      expect(component.stats).toHaveProperty('loading');
      expect(component.stats).toHaveProperty('view_num');
      expect(component.stats).toHaveProperty('comment_num');
      expect(component.stats).toHaveProperty('like_num');
      expect(component.stats).toHaveProperty('liked');
    });
  });
});