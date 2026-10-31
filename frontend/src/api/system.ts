import api from '@/utils/request'

export async function getVersion() {
  return api.get('/version')
}

export async function getSystemConfig() {
  return api.get('/config/system')
}

// 「已支持」上报一次匿名支持计数（金额仅接收端统计用，不做支付核验）
export async function donateSupport(amount?: number) {
  return api.post('/donate/support', { amount: amount || 0 })
}
