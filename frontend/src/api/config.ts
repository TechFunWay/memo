import api from '@/utils/request'

export async function getConfigs() {
  return api.get('/configs')
}

export async function updateConfig(key: string, value: string) {
  return api.post('/configs', { key, value })
}

export async function getPublicConfig() {
  return api.get('/config/public')
}
