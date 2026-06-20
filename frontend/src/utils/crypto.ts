import md5 from 'js-md5'

export function hashPassword(password: string): string {
  return md5(password)
}
