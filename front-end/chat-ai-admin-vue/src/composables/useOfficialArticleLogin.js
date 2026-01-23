import { ref, onMounted } from 'vue'
import {runPlugin} from "@/api/plugins/index.js";
import {useUserStore} from "@/stores/modules/user.js";

export function useOfficialArticleLogin(autoInit = true) {
  const userStore = useUserStore()

  const loginInfo = ref({})
  const loginStatus = ref(false)
  const loading = ref(false)
  const qrcodeUrl = ref('')
  const loginUrl = ref('')

  // 开源版本仅支持admin
  //const getUsername = () => userStore?.userInfo?.user_id
  const getUsername = () => 'admin'

  async function init() {
    await register()
    await getLoginStatus()
    getLoginUrl()
  }

  async function register() {
    const username = getUsername()
    return runPlugin({
      name: 'official_article',
      action: 'default/exec',
      params: JSON.stringify({
        business: 'register',
        arguments: { username }
      })
    })
  }

  async function getLoginStatus() {
    const username = getUsername()
    try {
      const res = await runPlugin({
        name: 'official_article',
        action: 'default/exec',
        params: JSON.stringify({
          business: 'get_login_status',
          arguments: { username }
        })
      })

      loginStatus.value = !!res?.data?.online
      loginInfo.value = res?.data || {}
    } finally {
      setTimeout(() => {
        loading.value = false
      }, 1200)
    }
  }

  function refresh() {
    if (loading.value) return
    loading.value = true
    getLoginStatus()
  }

  // 云版
  async function getLoginQrcode() {
    const username = getUsername()
    const {data} = await runPlugin({
      name: 'official_article',
      action: 'default/exec',
      params: JSON.stringify({
        business: 'wechat_qrcode_login',
        arguments: { username }
      })
    })
    qrcodeUrl.value = data?.qrcode_base64 || ''
  }

  // 开源版
  async function getLoginUrl() {
    const username = getUsername()
    const {data} = await runPlugin({
      name: 'official_article',
      action: 'default/exec',
      params: JSON.stringify({
        business: 'get_login_url',
        arguments: { username }
      })
    })
    loginUrl.value = data?.url
  }

  if (autoInit) {
    onMounted(init)
  }

  async function goLogin() {
    if (!loginUrl.value) await getLoginUrl()
    window.open(loginUrl.value)
  }

  return {
    loginInfo,
    loginStatus,
    loginUrl,
    qrcodeUrl,
    loading,
    init,
    refresh,
    goLogin,
    getLoginUrl,
    getLoginQrcode
  }
}
