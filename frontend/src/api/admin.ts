import request from './index'

// AD配置
export interface ADConfig {
  enabled: boolean
  server: string
  port: number
  use_ssl: boolean
  base_dn: string
  bind_dn: string
  bind_password: string
  user_filter: string
  auto_register: boolean
  default_role_id: number
}

export interface ADConfigResp {
  code: number
  message: string
  data: ADConfig
}

export interface TestConnectionReq {
  server: string
  port: number
  use_ssl: boolean
  base_dn: string
  bind_dn: string
  bind_password: string
  user_filter: string
}

// ============ 安全设置 ============

export interface SecuritySettings {
  id: number
  // 登录验证码
  captchaEnabled: boolean
  captchaMinLen: number
  // 注册功能开关
  registrationEnabled: boolean
  // 不活跃自动禁用
  inactiveAutoDisable: boolean
  inactiveDaysThreshold: number
  // 用户级别登录限制
  userLoginMaxAttempts: number
  userLoginLockMinutes: number
  // IP级别登录限制
  ipLoginMaxAttempts: number
  ipLoginLockMinutes: number
  ipWhitelist: string
  ipBlacklist: string
  // 密码安全
  passwordExpiryDays: number
  passwordMinLength: number
  passwordRequireUppercase: boolean
  passwordRequireLowercase: boolean
  passwordRequireDigit: boolean
  passwordRequireSpecial: boolean
  // 会话配置
  sessionTimeoutHours: number
}

export interface LockoutRecord {
  id: number
  target: string   // 用户名或IP
  type: string     // 'user' | 'ip'
  failCount: number
  locked: boolean
  lockedAt: string | null
  unlockedAt: string | null
  reason: string
  createdAt: string
  updatedAt: string
}

export interface SecurityOverview {
  totalLockedUsers: number
  totalLockedIPs: number
  whitelistCount: number
  blacklistCount: number
  passwordExpiryWarn: number
  inactiveUsersWarn: number
}

export namespace AdminApi {
  // 获取AD配置
  export const getADConfig = () => {
    return request.get<ADConfig>('/admin/ad-config')
  }

  // 更新AD配置
  export const updateADConfig = (data: Partial<ADConfig>) => {
    return request.post<null>('/admin/ad-config', data)
  }

  // 测试AD连接
  export const testADConnection = (data: TestConnectionReq) => {
    return request.post<{ message: string }>('/admin/ad-config/test', data)
  }

  // 获取安全设置
  export const getSecuritySettings = () => {
    return request.get<SecuritySettings>('/admin/security-settings')
  }

  // 更新安全设置
  export const updateSecuritySettings = (data: Partial<SecuritySettings>) => {
    return request.put<null>('/admin/security-settings', data)
  }

  // 获取安全总览
  export const getSecurityOverview = () => {
    return request.get<SecurityOverview>('/admin/security/overview')
  }

  // 获取被锁定的用户
  export const getLockedUsers = () => {
    return request.get<LockoutRecord[]>('/admin/security/locked-users')
  }

  // 解锁用户
  export const unlockUser = (username: string) => {
    return request.post<null>('/admin/security/unlock-user', undefined, { params: { username } })
  }

  // 获取被锁定的IP
  export const getLockedIPs = () => {
    return request.get<LockoutRecord[]>('/admin/security/locked-ips')
  }

  // 解锁IP
  export const unlockIP = (ip: string) => {
    return request.post<null>('/admin/security/unlock-ip', undefined, { params: { ip } })
  }
}
