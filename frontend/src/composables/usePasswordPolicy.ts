import { ref, computed } from 'vue'
import { AdminApi, type SecuritySettings } from '@/api/admin'

// 全局单例 - 缓存密码策略设置
const cachedPolicy = ref<SecuritySettings | null>(null)
const loading = ref(false)

export function usePasswordPolicy() {
  // 加载密码策略（带缓存）
  const loadPolicy = async () => {
    if (cachedPolicy.value) return cachedPolicy.value
    loading.value = true
    try {
      const res = await AdminApi.getSecuritySettings()
      if (res.code === 200 && res.data) {
        cachedPolicy.value = res.data
        return res.data
      }
    } finally {
      loading.value = false
    }
    return null
  }

  // 清除缓存（当安全设置保存后调用）
  const clearCache = () => {
    cachedPolicy.value = null
  }

  // 默认策略（未加载时使用）
  const defaultPolicy: SecuritySettings = {
    id: 0,
    captchaEnabled: true,
    captchaMinLen: 3,
    inactiveAutoDisable: false,
    inactiveDaysThreshold: 90,
    userLoginMaxAttempts: 5,
    userLoginLockMinutes: 30,
    ipLoginMaxAttempts: 20,
    ipLoginLockMinutes: 60,
    ipWhitelist: '',
    ipBlacklist: '',
    passwordExpiryDays: 0,
    passwordMinLength: 8,
    passwordRequireUppercase: false,
    passwordRequireLowercase: false,
    passwordRequireDigit: false,
    passwordRequireSpecial: false,
    sessionTimeoutHours: 24,
    registrationEnabled: true
  }

  const policy = computed(() => cachedPolicy.value || defaultPolicy)

  // 特殊字符检测：直接列出常用特殊字符，明确覆盖 %
  const hasSpecialChar = (pwd: string): boolean => {
    const specialChars = '!@#$%^&*()_+-=[]{}|;:\'",.<>?/`~\\'
    for (const c of pwd) {
      if (specialChars.includes(c)) return true
    }
    return false
  }

  // 密码强度检查（实时）
  const checkPassword = (password: string) => {
    const p = policy.value
    return {
      length: password.length >= p.passwordMinLength,
      upper: !p.passwordRequireUppercase || /[A-Z]/.test(password),
      lower: !p.passwordRequireLowercase || /[a-z]/.test(password),
      number: !p.passwordRequireDigit || /[0-9]/.test(password),
      special: !p.passwordRequireSpecial || hasSpecialChar(password)
    }
  }

  // 计算密码强度分数 (0-5)
  const passwordStrength = (password: string) => {
    const rules = checkPassword(password)
    let score = 0
    if (rules.length) score++
    if (rules.upper) score++
    if (rules.lower) score++
    if (rules.number) score++
    if (rules.special) score++
    return score
  }

  // 生成密码验证规则（用于 el-form）
  const passwordValidators = (t: (key: string) => string) => {
    const p = policy.value
    const rules: any[] = [
      { required: true, message: t('profile.newPasswordRequired'), trigger: 'blur' },
      { min: p.passwordMinLength, message: t('profile.passwordMinLength', { min: p.passwordMinLength }), trigger: 'blur' }
    ]
    if (p.passwordRequireUppercase) {
      rules.push({ pattern: /[A-Z]/, message: t('profile.passwordRequireUppercase'), trigger: 'blur' })
    }
    if (p.passwordRequireLowercase) {
      rules.push({ pattern: /[a-z]/, message: t('profile.passwordRequireLowercase'), trigger: 'blur' })
    }
    if (p.passwordRequireDigit) {
      rules.push({ pattern: /[0-9]/, message: t('profile.passwordRequireDigit'), trigger: 'blur' })
    }
    if (p.passwordRequireSpecial) {
      rules.push({ validator: (_: any, val: string, cb: any) => {
        if (!hasSpecialChar(val)) cb(new Error(t('profile.passwordRequireSpecial')))
        else cb()
      }, trigger: 'blur' })
    }
    return rules
  }

  return {
    loadPolicy,
    clearCache,
    policy,
    checkPassword,
    passwordStrength,
    passwordValidators,
    loading
  }
}
