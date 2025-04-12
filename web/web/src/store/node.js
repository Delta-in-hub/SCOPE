// src/store/node.js
import { defineStore } from 'pinia';
import { ref } from 'vue';
import * as nodeApi from '@/services/nodeApi';

export const useNodeStore = defineStore('node', () => {
  // State
  const nodes = ref([]); // 存储节点列表

  // Actions
  async function fetchNodes() {
    try {
      const response = await nodeApi.listNodes();
      // API 返回的是数组，直接赋值
      // 后端返回的数据可能需要进行一些处理或验证
      if (Array.isArray(response.data)) {
          nodes.value = response.data;
      } else {
          console.warn('Node list API did not return an array:', response.data);
          nodes.value = []; // 或者保持不变，看业务需求
      }

    } catch (error) {
      console.error('Failed to fetch nodes in store:', error);
      nodes.value = []; // 出错时清空列表
      throw error; // 将错误抛出，让组件知道
    }
  }

  return {
    // state
    nodes,
    // actions
    fetchNodes,
  };
});