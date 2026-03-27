import request from './index'
import type { ResultData } from './index'

// AD登录请求
export interface ADLoginReq {
  username: string
  password: string
  captchaId?: string
  captcha?: string
}
export interface LoginReq {
  username: string
  password: string
  captchaId?: string
  captcha?: string
}

// 验证码响应
export interface CaptchaResp {
  captchaId: string
  captchaImage: string
  enabled: boolean
}

// 注册请求
export interface RegisterReq {
  username: string
  password: string
  nickname?: string
  email?: string
}

// 用户信息
export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  roleId: number
  roleName: string
  roleCode: string
  status: string
  mfaEnabled: boolean
  source?: string  // LOCAL / AD
  permissions: string[]  // 用户权限列表，用于前端菜单控制
}

// 登录响应
export interface LoginResp {
  token: string
  expireAt: number
  user: UserInfo
  mfaRequired: boolean
}

// MFA验证请求
export interface MFAVerifyReq {
  userId: number
  code: string
  tmpToken?: string
}

// MFA启用请求
export interface MFAEnableReq {
  code: string
}

// MFA状态响应
export interface MFAStatusResp {
  enabled: boolean
  secret?: string
  backupCodes?: string[]
}

// 修改密码请求
export interface ChangePasswordReq {
  oldPassword: string
  newPassword: string
}

// 个人资料更新请求
export interface ProfileUpdateReq {
  nickname?: string
  email?: string
  phone?: string
}

export namespace AuthApi {
  // 获取登录验证码
  export const getCaptcha = () => {
    return request.get<CaptchaResp>('/auth/captcha')
  }

  // 用户登录
  export const login = (data: LoginReq) => {
    return request.post<LoginResp>('/auth/login', data)
  }

  // AD用户登录
  export const adLogin = (data: ADLoginReq) => {
    return request.post<LoginResp>('/auth/ad-login', data)
  }

  // MFA验证
  export const mfaVerify = (data: MFAVerifyReq) => {
    return request.post<LoginResp>('/auth/mfa/verify', data)
  }

  // 用户注册
  export const register = (data: RegisterReq) => {
    return request.post<LoginResp>('/auth/register', data)
  }

  // 获取当前用户信息
  export const getCurrentUser = () => {
    return request.get<UserInfo>('/auth/current')
  }

  // 更新个人资料
  export const updateProfile = (data: ProfileUpdateReq) => {
    return request.put<UserInfo>('/auth/profile', data)
  }

  // 修改密码
  export const changePassword = (data: ChangePasswordReq) => {
    return request.post<null>('/auth/change-password', data)
  }

  // 获取MFA状态
  export const getMFAStatus = () => {
    return request.get<MFAStatusResp>('/auth/mfa/status')
  }

  // 启用MFA
  export const enableMFA = (data: MFAEnableReq) => {
    return request.post<MFAStatusResp>('/auth/mfa/enable', data)
  }

  // 禁用MFA
  export const disableMFA = () => {
    return request.post<null>('/auth/mfa/disable', {})
  }
}
