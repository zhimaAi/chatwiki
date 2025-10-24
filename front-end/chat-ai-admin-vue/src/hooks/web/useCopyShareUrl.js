import { message, Modal } from 'ant-design-vue'
import useClipboard from 'vue-clipboard3'
import { usePublicLibraryStore } from '@/stores/modules/public-library'
import { useRouter } from 'vue-router'

export const useCopyShareUrl = () => {
  const libraryStore = usePublicLibraryStore()
  const router = useRouter()
  
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
        title: '提示',
        content,
        okText: '立即配置',
        cancelText: '取消',
        onOk: navigateToPermissions
      })
    }

    // 检查是否有自定义域名
    if (!share_url) {
      showConfigModal('复制链接需要先配置自定义域名，是否立即前往配置？')
      return
    }

    // 检查访问权限是否为公开
    if (access_rights != 1) {
      showConfigModal('仅公开文档支持复制链接，是否立即前往设置访问权限？')
      return
    }
    
    // 执行复制操作
    const { toClipboard } = useClipboard()
    const url = share_url + docUrl

    try {
      await toClipboard(url)
      message.success('复制成功')
    } catch (err) {
      console.error('复制失败:', err)
      message.error('复制失败')
    }
  }

  return {
    copyShareUrl
  }
}
