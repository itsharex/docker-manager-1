import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
});

export const dockerApi = {
  // Containers
  getContainers: () => api.get('/containers'),
  startContainer: (id: string) => api.post(`/containers/${id}/start`),
  stopContainer: (id: string) => api.post(`/containers/${id}/stop`),
  removeContainer: (id: string) => api.delete(`/containers/${id}/remove`),
  inspectContainer: (id: string) => api.get(`/containers/${id}/inspect`),

  // Images
  getImages: () => api.get('/images'),
  removeImage: (id: string) => api.delete(`/images/${id}`),

  // Volumes
  getVolumes: () => api.get('/volumes'),
  removeVolume: (id: string) => api.delete(`/volumes/${id}`),

  // Networks
  getNetworks: () => api.get('/networks'),
  removeNetwork: (id: string) => api.delete(`/networks/${id}`),

  // System
  getSystemInfo: () => api.get('/info'),
};

export const getWsUrl = (path: string) => {
  return `ws://localhost:8080/ws${path}`;
};
