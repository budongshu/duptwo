import { ref, readonly } from 'vue'
import { AuthApi, type UserInfo } from '@/api/auth'

const userInfo = ref<UserInfo | null>(null)

export function useUser() {
  const refreshUser = async () => {
    try {
      const res = await AuthApi.getCurrentUser()
      if (res.code === 200 && res.data) {
        userInfo.value = res.data
        localStorage.setItem('user', JSON.stringify(res.data))
        return res.data
      }
    } catch {
      // ignore
    }
    return null
  }

  const initUser = () => {
    try {
      const stored = localStorage.getItem('user')
      if (stored) {
        userInfo.value = JSON.parse(stored)
      }
    } catch {
      userInfo.value = null
    }
  }

  const getPermissions = (): string[] => {
    return userInfo.value?.permissions || []
  }

  const hasPermission = (perm: string): boolean => {
    return getPermissions().includes(perm)
  }

  return {
    userInfo: readonly(userInfo),
    refreshUser,
    initUser,
    getPermissions,
    hasPermission,
  }
}
