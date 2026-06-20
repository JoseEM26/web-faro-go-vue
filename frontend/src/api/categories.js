import api from './axios'

export const categoriesApi = {
  getAll:       ()         => api.get('/categories'),
  create:       (data)     => api.post('/categories', data),
  update:       (id, data) => api.put(`/categories/${id}`, data),
  remove:       (id)       => api.delete(`/categories/${id}`),
}
