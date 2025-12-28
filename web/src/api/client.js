import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1'
})

export const getSkills = (query) => {
  return api.get('/skills', { params: { q: query } })
}

export const getSkillDetail = (name) => {
  return api.get(`/skills/${name}`)
}
