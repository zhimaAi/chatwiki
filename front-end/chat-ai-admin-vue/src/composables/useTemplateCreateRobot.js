import {computed, h} from 'vue';
import {useRouter} from 'vue-router';
import {message, Modal} from 'ant-design-vue';
import {ExclamationCircleOutlined} from '@ant-design/icons-vue';
import {useCompanyStore} from "@/stores/modules/company.js";
import {useRobotTemplate} from "@/api/explore/template.js";
import {useI18n} from '@/hooks/web/useI18n';

export function useTemplateCreateRobot() {
  const {t} = useI18n('composables.use-template-create-robot')
  const router = useRouter()
  const companyStore = useCompanyStore()
  const {companyInfo} = companyStore

  const sysVersion = computed(() => companyInfo?.version)

  function checkVersion(sys_v, tpl_v) {
    sys_v = Number(sys_v.replace(/\D/g, ''))
    tpl_v = Number(tpl_v.replace(/\D/g, ''))
    return sys_v >= tpl_v
  }

  function run(tpl, cb = null) {
    useRobotTemplate({template_id: tpl.id, csl_url: tpl.csl_url}).then((res) => {
      typeof cb === "function" && cb(true)
      message.success(t('msg_use_success'))
      const {id, robot_key} = res.data
      const url = router.resolve({path: '/robot/config/workflow', query: {id, robot_key}})
      window.open(url.href, '_blank')
    })
  }

  function exec(tpl, _call = null) {
    const versionLow = !checkVersion(sysVersion.value, tpl.version)
    const {content, okText} = versionLow
      ? {
        content: t('msg_version_low'),
        okText: t('btn_continue_use')
      }
      : {
        content: t('msg_confirm_use_template', {name: tpl.name}),
        okText: t('btn_confirm')
      }
    Modal.confirm({
      title: t('title_prompt'),
      content: content,
      icon: h(ExclamationCircleOutlined),
      okText: okText,
      cancelText: t('btn_cancel'),
      onOk: () => run(tpl, _call)
    })
  }

  return {
    useTpl: exec,
  }
}
