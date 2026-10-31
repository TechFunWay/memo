import api from '@/utils/request'
import { hashPassword } from '@/utils/crypto'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    id: number
    username: string
    nickname: string
    email: string
    avatar: string
    is_admin: boolean
    is_active: boolean
  }
}

export async function login(username: string, password: string, keepLogin = false): Promise<LoginResponse> {
  return api.post('/auth/login', { username, password: hashPassword(password), keep_login: keepLogin })
}

export async function register(username: string, password: string, nickname: string, email: string, securityQuestions?: string[], securityAnswers?: string[]) {
  return api.post('/auth/register', {
    username,
    password: hashPassword(password),
    nickname,
    email,
    security_questions: securityQuestions || [],
    security_answers: securityAnswers || [],
  })
}

export async function fnosLogin(): Promise<LoginResponse & { binding_required?: boolean; fnos_username?: string }> {
  return api.post('/auth/fnos/login')
}

export async function bindFnOSAccount(mode: 'register' | 'bind', username: string, password: string): Promise<LoginResponse> {
  return api.post('/auth/fnos/bind', { mode, username, password: hashPassword(password) })
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export async function changePassword(oldPassword: string, newPassword: string) {
  return api.put('/auth/password', { old_password: hashPassword(oldPassword), new_password: hashPassword(newPassword) })
}

export async function logout() {
  return Promise.resolve()
}

export async function checkAuth() {
  return api.get('/auth/check')
}

export async function getCurrentUser() {
  return api.get('/auth/me')
}

export interface HasUsersResult {
  has_users: boolean
  /** 本次启动时数据库是新建的空库（正常仅首次安装出现） */
  fresh_db?: boolean
  /** 历史上已有账号但当前数据库为空，疑似数据目录未挂载或被更换 */
  suspect_data_loss?: boolean
}

export async function checkHasUsers(): Promise<HasUsersResult> {
  return api.get('/auth/has-users')
}

export async function getSecurityQuestions(username: string): Promise<{ questions: string[] }> {
  return api.get('/auth/security-question', { params: { username } })
}

export async function verifySecurityAnswers(username: string, answers: string[], newPassword: string) {
  return api.post('/auth/verify-security', {
    username,
    answers,
    new_password: hashPassword(newPassword),
  })
}

export async function setSecurityQuestions(questions: string[], answers: string[]) {
  return api.put('/auth/security-question', { questions, answers })
}

export async function checkSecurityQuestion(): Promise<{ has_security_question: boolean }> {
  return api.get('/auth/security-question/check')
}
