import api from './axios'

export const notesApi = {
  getAll:  ()         => api.get('/notes'),
  getById: (id)       => api.get(`/notes/${id}`),
  create:  (data)     => api.post('/notes', data),
  update:  (id, data) => api.put(`/notes/${id}`, data),
  remove:  (id)       => api.delete(`/notes/${id}`),
}
