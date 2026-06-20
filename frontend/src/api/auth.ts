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

export async function checkHasUsers(): Promise<{ has_users: boolean }> {
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
