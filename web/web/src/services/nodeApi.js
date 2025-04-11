// src/services/nodeApi.js
import axiosInstance from './axiosInstance';

export const listNodes = () => {
  // GET 请求不需要 body，拦截器会自动添加 Authorization header
  return axiosInstance.get('/node/list');
};

// node/up 和 node/down 通常由 agent 调用，前端可能不需要直接调用
// 如果需要（例如管理功能），可以添加相应方法
// export const nodeUp = (nodeInfo) => {
//   return axiosInstance.post('/node/up', nodeInfo);
// };
// export const nodeDown = (nodeInfo) => {
//   return axiosInstance.post('/node/down', nodeInfo);
// };