import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 添加JWT token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      // 可以在这里跳转到登录页
    }
    return Promise.reject(error)
  }
)

// 公开API
export const getSites = () => api.get('/sites')
export const getContacts = () => api.get('/contacts')
export const getSiteConfig = () => api.get('/config')
export const getFrontendConfig = () => api.get('/frontend-config')

// 认证API
export const login = (username, password) =>
  api.post('/auth/login', { username, password })

// 管理API
export const adminAPI = {
  // 站点管理
  getSites: () => api.get('/admin/sites'),
  createSite: (data) => api.post('/admin/sites', data),
  updateSite: (id, data) => api.put(`/admin/sites/${id}`, data),
  deleteSite: (id) => api.delete(`/admin/sites/${id}`),

  // 联系方式管理
  getContacts: () => api.get('/admin/contacts'),
  createContact: (data) => api.post('/admin/contacts', data),
  updateContact: (id, data) => api.put(`/admin/contacts/${id}`, data),
  deleteContact: (id) => api.delete(`/admin/contacts/${id}`),

  // 站点配置管理
  getSiteConfig: () => api.get('/admin/config'),
  updateSiteConfig: (data) => api.put('/admin/config', data),

  // 文件上传
  uploadFile: (file) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/admin/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
  },

  // 统计API
  getStats: () => api.get('/admin/stats'),
  getChartData: (period) => api.get(`/admin/charts?period=${period}`),
  getRecentVisits: (limit = 5) => api.get(`/admin/recent-visits?limit=${limit}`),
  
  // 热重载通知
  notifyConfigUpdate: () => api.post('/admin/notify-update'),
  
  // 用户管理
  changePassword: (oldPassword, newPassword) => 
    api.put('/admin/change-password', { oldPassword, newPassword }),
  
  // 日志API
  getBackendLogs: (lines = 100) => api.get(`/admin/logs?lines=${lines}`),
  
  // 登录历史API
  getLoginHistory: (limit = 20) => api.get(`/admin/login-history?limit=${limit}`),
  
  // 轮换文本配置API
  getRotatingTexts: () => api.get('/api/rotating-texts'),
  updateRotatingTexts: (texts) => api.put('/admin/rotating-texts', { texts }),
}

export default api
