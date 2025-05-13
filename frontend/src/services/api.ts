import axios from 'axios';

const API_URL = 'https://to-read.lg.gl/api/v1';

// 定义API响应的通用接口
interface ApiResponse<T = any> {
  code: number;
  msg: string | null;
  data: T;
}

// 定义各种数据类型的接口
interface User {
  user_id: number;
  user_name: string;
  role: number;
}

interface LoginResponse {
  user_id: number;
  user_name: string;
  role: number;
  token: string;
  token_expiration_time: number;
}

interface Collection {
  collection_id: number;
  url: string;
  type: 'text' | 'image' | 'video';
  title: string;
  description: string;
  tags: string[];
  created_at: number;
}

interface CollectionListResponse {
  collections: Collection[];
}

interface SummaryResponse {
  summary: string;
}

interface TagsResponse {
  tags: string[];
}

const instance = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器，添加token
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 创建一个包装函数，处理响应并提取数据
async function api<T>(promise: Promise<any>): Promise<ApiResponse<T>> {
  try {
    const response = await promise;
    return response.data as ApiResponse<T>;
  } catch (error: any) {
    if (error.response && error.response.status === 401) {
      // Token过期，清除登录状态
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    throw error;
  }
}

// 使用包装函数来处理所有API调用
export const healthCheck = () => api<string>(instance.get('/health'));

// 用户相关API
export const userAPI = {
  getUser: (params: { user_id?: number; user_name?: string }) => 
    api<User>(instance.get('/user', { params })),
  register: (data: { user_name: string; password: string }) => 
    api<User>(instance.post('/user/register', data)),
  login: (data: { user_name: string; password: string }) => 
    api<LoginResponse>(instance.post('/user/login', data)),
};

// 收藏相关API
export const collectionAPI = {
  getList: (params?: { search?: string; tags?: string[] }) => 
    api<CollectionListResponse>(instance.get('/collection/list', { 
      params: {
        ...params,
        tags: params?.tags?.join(',')
      }
    })),
  addCollection: (data: { url: string }) => 
    api<Collection>(instance.post('/collection/add', data)),
  getSummary: (params?: { search?: string; tags?: string[] }) => 
    api<SummaryResponse>(instance.get('/collection/summary', { 
      params: {
        ...params,
        tags: params?.tags?.join(',')
      }
    })),
  getTags: () => api<TagsResponse>(instance.get('/collection/tag')),
};

export type { ApiResponse, User, Collection, LoginResponse, CollectionListResponse, SummaryResponse, TagsResponse };
export default instance;
