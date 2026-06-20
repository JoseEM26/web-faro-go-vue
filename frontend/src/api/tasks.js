import api from './axios'

export const tasksApi = {
  getAll:  ()           => api.get('/tasks/'),
  getById: (id)         => api.get(`/tasks/${id}`),
  create:  (data)       => api.post('/tasks/', data),
  update:  (id, data)   => api.put(`/tasks/${id}`, data),
  remove:  (id)         => api.delete(`/tasks/${id}`),
}
