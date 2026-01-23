import { message, Modal } from 'ant-design-vue'
import useClipboard from 'vue-clipboard3'
import { usePublicLibraryStore } from '@/stores/modules/public-library'
import { useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

export const useCopyShareUrl = () => {
  const libraryStore = usePublicLibraryStore()
  const router = useRouter()
  const { t } = useI18n('hooks.web.use-copy-share-url')

  const copyShareUrl = async (docUrl) => {
    const { access_rights, share_url, id, library_key } = libraryStore.libraryInfo

    // 提取公共的跳转到权限配置页面的方法
    const navigateToPermissions = () => {
      router.push({
        path: '/public-library/permissions',
        query: { library_id: id, library_key }
      })
    }

    // 提取公共的确认弹窗方法
    const showConfigModal = (content) => {
      Modal.confirm({
        title: t('prompt'),
        content,
        okText: t('configure_now'),
        cancelText: t('cancel'),
        onOk: navigateToPermissions
      })
    }

    // 检查是否有自定义域名
    if (!share_url) {
      showConfigModal(t('need_custom_domain'))
      return
    }

    // 检查访问权限是否为公开
    if (access_rights != 1) {
      showConfigModal(t('only_public_doc'))
      return
    }

    // 执行复制操作
    const { toClipboard } = useClipboard()
    const url = share_url + docUrl

    try {
      await toClipboard(url)
      message.success(t('copy_success'))
    } catch (err) {
      console.error('复制失败:', err)
      message.error(t('copy_failed'))
    }
  }

  return {
    copyShareUrl
  }
}
