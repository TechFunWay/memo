import api from '@/utils/request'

export async function getVersion() {
  return api.get('/version')
}

export async function getSystemConfig() {
  return api.get('/config/system')
}
