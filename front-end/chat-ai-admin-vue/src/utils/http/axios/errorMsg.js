import { useI18n } from '@/hooks/web/useI18n'
export function getErrorMsg(error) {
  const { t } = useI18n()
  let status = error.response.status
  let msg = error.response.data.message || error.message

  let errMessage = ''

  switch (status) {
    case 400:
      errMessage = `${msg}`
      break
    case 401:
      errMessage = t('common.errMsg401')
      break
    case 403:
      errMessage = t('common.errMsg403')
      break
    // 404请求不存在
    case 404:
      errMessage = t('common.errMsg404')
      break
    case 405:
      errMessage = t('common.errMsg405')
      break
    case 408:
      errMessage = t('common.errMsg408')
      break
    case 500:
      errMessage = t('common.errMsg500')
      break
    case 501:
      errMessage = t('common.errMsg501')
      break
    case 502:
      errMessage = t('common.errMsg502')
      break
    case 503:
      errMessage = t('common.errMsg503')
      break
    case 504:
      errMessage = t('common.errMsg504')
      break
    case 505:
      errMessage = t('common.errMsg505')
      break
    default:
      errMessage = `${msg}`
  }

  return errMessage
}
