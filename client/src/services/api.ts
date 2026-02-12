import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
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

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authApi = {
  register: async (data: { name: string; email: string; password: string }) => {
    const response = await api.post('/auth/register', data)
    return response.data
  },

  login: async (data: { email: string; password: string }) => {
    const response = await api.post('/auth/login', data)
    return response.data
  },

  validate: async () => {
    const response = await api.get('/auth/validate')
    return response.data
  },
}

// Task API
export interface Task {
  id: string
  title: string
  description: string
  status: 'PENDING' | 'IN_PROGRESS' | 'DONE'
  user_id: string
  created_at: string
  updated_at: string
}

export interface TaskListResponse {
  tasks: Task[]
  total_count: number
  page: number
  limit: number
  total_pages: number
}

export const taskApi = {
  list: async (params?: { page?: number; limit?: number; status?: string }) => {
    const response = await api.get<{ success: boolean; data: TaskListResponse; message: string }>('/tasks', { params })
    return response.data
  },

  get: async (id: string) => {
    const response = await api.get<{ success: boolean; data: Task; message: string }>(`/tasks/${id}`)
    return response.data
  },

  create: async (data: { title: string; description?: string; status?: string }) => {
    const response = await api.post<{ success: boolean; data: Task; message: string }>('/tasks', data)
    return response.data
  },

  update: async (id: string, data: { title?: string; description?: string; status?: string }) => {
    const response = await api.put<{ success: boolean; data: Task; message: string }>(`/tasks/${id}`, data)
    return response.data
  },

  delete: async (id: string) => {
    const response = await api.delete<{ success: boolean; message: string }>(`/tasks/${id}`)
    return response.data
  },
}

export default api
