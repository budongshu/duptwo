import { ElMessage } from 'element-plus'

// 用户友好的错误消息映射
const ERROR_MESSAGES: Record<number, { title: string; message: string }> = {
  400: { title: '请求参数错误', message: '请检查输入的数据格式是否正确' },
  401: { title: '登录已过期', message: '请重新登录后再操作' },
  403: { title: '权限不足', message: '您没有执行此操作的权限，请联系管理员开通' },
  404: { title: '资源不存在', message: '请求的资源可能已被删除' },
  500: { title: '服务器错误', message: '服务器处理失败，请稍后重试' },
  502: { title: '服务维护中', message: '服务暂时不可用，请稍后重试' },
  503: { title: '服务繁忙', message: '服务器繁忙，请稍后重试' },
}

// HTTP状态码到权限标识的映射
const HTTP_TO_PERMISSION: Record<string, string> = {
  '/users': 'user',
  '/upload-records': 'upload',
  '/projects': 'project',
  '/personnels': 'personnel',
  '/roles': 'role',
  '/user-groups': 'userGroup',
  '/field-configs': 'field-config',
  '/admin': 'config',
}

interface ErrorInfo {
  title: string
  message: string
  detail?: string
}

function parseError(error: any): ErrorInfo {
  // axios错误
  if (error?.response) {
    const status = error.response.status
    const apiPath = error.config?.url || ''

    // 403权限错误，尝试给出具体权限建议
    if (status === 403) {
      let requiredPermission = ''
      for (const [path, perm] of Object.entries(HTTP_TO_PERMISSION)) {
        if (apiPath.includes(path)) {
          requiredPermission = perm
          break
        }
      }
      return {
        title: '权限不足',
        message: '您没有执行此操作的权限',
        detail: requiredPermission ? `需要权限标识: ${requiredPermission}:create/update/delete` : undefined
      }
    }

    // 其他HTTP错误
    const config = ERROR_MESSAGES[status]
    if (config) {
      return {
        title: config.title,
        message: config.message
      }
    }
  }

  // 网络错误
  if (error?.message?.includes('Network Error') || error?.code === 'ERR_NETWORK') {
    return {
      title: '网络错误',
      message: '无法连接服务器，请检查网络后重试'
    }
  }

  // 超时错误
  if (error?.message?.includes('timeout') || error?.code === 'ECONNABORTED') {
    return {
      title: '请求超时',
      message: '服务器响应超时，请稍后重试'
    }
  }

  // 后端返回的业务错误
  if (error?.message) {
    return {
      title: '操作失败',
      message: error.message
    }
  }

  // 未知错误
  return {
    title: '操作失败',
    message: '发生了未知错误，请稍后重试'
  }
}

export function showErrorMessage(error: any, fallbackMessage?: string) {
  const errInfo = parseError(error)

  let finalMessage = errInfo.message
  if (!finalMessage || finalMessage === '操作失败') {
    finalMessage = fallbackMessage || '操作失败，请稍后重试'
  }

  if (errInfo.detail) {
    ElMessage.error({
      message: finalMessage,
      duration: 5000
    })
    console.warn(`[权限提示] ${errInfo.detail}`)
  } else {
    ElMessage.error({
      message: finalMessage,
      duration: 4000
    })
  }
}

export function showSuccessMessage(message: string) {
  ElMessage.success({
    message,
    duration: 3000
  })
}

export function showWarningMessage(message: string) {
  ElMessage.warning({
    message,
    duration: 4000
  })
}

export { parseError }
