<!-- 
	This is the sign in page, it uses the dashboard layout in: 
	"./layouts/Default.vue" .
 -->

<template>
  <div class="login-box" :class="['lang-' + lang]">
    <div class="header">
      <div class="header-left">
        <div class="logo">
          <LayoutLogo />
        </div>
        <div class="header-item" @click="toHome">{{ t('home') }}</div>
        <div class="header-item">{{ t('helpDocument') }}</div>
        <div class="header-item">{{ t('githubUrl') }}</div>
      </div>
      <div class="header-right">
        <LocaleDropdown />
      </div>
    </div>
    <div class="content">
      <div class="banner"></div>
      <div class="sign-in">
        <h2 class="project-name">{{ companyName }}</h2>
        <div class="sign-in-content">
          <!-- Sign In Form Column -->
          <div class="col-form">
            <h2 class="login-title">{{ t('accountLogin') }}</h2>

            <!-- Sign In Form -->
            <a-form
              class="login-form"
              :model="formState"
              name="basic"
              autocomplete="off"
              @finish="onFinish"
              @finishFailed="onFinishFailed"
            >
              <a-form-item
                name="username"
                class="usernames"
                :rules="[{ required: true, message: t('pleaseNumber') }]"
              >
                <a-input class="login-item login-username" v-model:value="formState.username" autocomplete="off" :placeholder="t('pleaseNumber')">
                  <template #prefix><UserOutlined style="color: rgba(0, 0, 0, 0.25)" /></template>
                </a-input>
              </a-form-item>

              <a-form-item
                name="password"
                :rules="[{ required: true, message: t('pleasePassword') }]"
              >
                <a-input-password class="login-item login-password" v-model:value="formState.password" autocomplete="off" type="password" :placeholder="t('pleasePassword')">
                  <template #prefix><LockOutlined style="color: rgba(0, 0, 0, 0.25)" /></template>
                </a-input-password>
              </a-form-item>

              <!-- <a-form-item name="remember">
                <a-checkbox v-model:checked="formState.remember">Remember Me</a-checkbox>
              </a-form-item> -->

              <a-form-item>
                <a-button class="login-btn" type="primary" block html-type="submit">{{ t('login') }}</a-button>
              </a-form-item>
            </a-form>
            <!-- / Sign In Form -->
            <!-- 
            <p class="font-semibold text-muted">
              Don't have an account?
              <router-link to="/login" class="font-bold text-dark">Sign Up</router-link>
            </p> -->
          </div>
   
        </div>
      </div>
      <div class="main-box">
        
        <!-- Instructions -->
        <div class="instructions-box">
          <div class="instructions-item one-deployment">{{ t('oneDeployment') }}</div>
          <div class="instructions-item out-the-box">{{ t('outTheBox') }}</div>
          <div class="instructions-item safe-controllable">{{ t('safeControllable') }}</div>
        </div>

        <div class="profession-box">
          <div class="profession-title"><span class="profession-title-stress">{{ t('professionTitleOne') }}</span>{{ t('professionTitleTwo') }}</div>
          <div class="profession-info">{{ t('professionInfo') }}</div>
          <div class="profession-item-box">
            <div class="profession-item" v-for="item in profession" :key="item.id">
              <div class="profession-item-top">
                <div class="profession-item-icon">
                  <svg-icon
                    class="profession-item-icon-svg"
                    :name="item.icon"
                  ></svg-icon>
                </div>
                <div class="profession-item-title">{{  item.title }}</div>
              </div>
              <div class="profession-item-middle"></div>
              <div class="profession-item-bottom">{{ item.info }}</div>
            </div>
          </div>
        </div>

        <div class="ai-model-box">
          <div class="ai-content">
            <div class="ai-model-title">{{ t('aiModelTitleOne') }}<span class="ai-model-title-stress">{{ t('aiModelTitleTwo') }}</span>{{ t('aiModelTitleThree') }}</div>
            <div class="ai-model-item-box">
              <div class="ai-model-item" v-for="item in aiModel" :key="item.id">
                <svg-icon :class="{'ai-model-item-svg-gemini': item.icon == 'gemini', 'ai-model-item-svg-anthrop': item.icon == 'anthrop'}" class="ai-model-item-svg" :name="item.icon"></svg-icon>
                <div class="ai-model-item-title">{{ item.title }}</div>
              </div>
            </div>
          </div>
        </div>

        <div class="open-source-box">
          <div class="open-source-title">{{ t('openSourceTitle') }}</div>
          <div class="open-source-item-box">
            <div class="open-source-item" v-for="item in openSource" :key="item.id">
              <div class="open-source-item-title">{{ item.title }}</div>
              <div class="open-source-item-info">{{ item.info }}</div>
              <div class="open-source-item-list" v-for="list in item.list" :key="list.listId">
                <svg-icon class="open-source-item-list-svg" :name="list.icon"></svg-icon>
                <div class="open-source-item-list-title">{{ list.title }}</div>
              </div>
              <svg-icon class="open-source-item-svg" :name="item.icon"></svg-icon>
            </div>
          </div>
        </div>
      </div>

    </div>
    <div class="layout-footer-wrapper">
      <div class="layout-footer">
        <div class="copyright-text-box">
          <div class="copyright-text-info">Powered by</div>
          <div class="copyright-text" @click="toChatWiki">  ChatWiki</div>
        </div>
        <div class="footer-line"></div>
        <div class="copyright-text" @click="toXiaokefu">{{ t('copyrightTextTwo') }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useLocaleStore } from '@/stores/modules/locale'
import { useI18n } from '@/hooks/web/useI18n'
import { reactive, computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import LocaleDropdown from '@/layouts/AdminLayout/compoents/locale-dropdown.vue'
import LayoutLogo from '@/layouts/AdminLayout/compoents/layout-logo.vue'
import { useUserStore } from '@/stores/modules/user'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue';
import { useCompanyStore } from '@/stores/modules/company'

const localeStore = useLocaleStore()
const lang = computed(() => localeStore.getCurrentLocale.lang)

const companyStore = useCompanyStore()
!companyStore.name && companyStore.getCompanyInfo();
const companyName = companyStore.name || 'ChatWiki'

const { t } = useI18n('views.login.login')

// const selectedKeys = ref(['2']);

const profession = computed(() => [
  {
    title: t('professionOneTitle'),
    info: t('professionOneInfo'),
    id: 1,
    icon: 'finance'
  },
  {
    title: t('professionTwoTitle'),
    info: t('professionTwoInfo'),
    id: 2,
    icon: 'government'
  },
  {
    title: t('professionThreeTitle'),
    info: t('professionThreeInfo'),
    id: 3,
    icon: 'law'
  },
  {
    title: t('professionFourTitle'),
    info: t('professionFourInfo'),
    id: 4,
    icon: 'medical-treatment'
  },
  {
    title: t('professionFiveTitle'),
    info: t('professionFiveInfo'),
    id: 5,
    icon: 'manufacture'
  },
  {
    title: t('professionSixTitle'),
    info: t('professionSixInfo'),
    id: 6,
    icon: 'source-energy'
  },
  {
    title: t('professionSevenTitle'),
    info: t('professionSevenInfo'),
    id: 7,
    icon: 'retail'
  },
  {
    title: t('professionEightTitle'),
    info: t('professionEightInfo'),
    id: 8,
    icon: 'unit'
  }
])

const aiModel = computed(() => [
  {
    title: 'OpenAI',
    id: 1,
    icon: 'open-ai'
  },
  {
    title: 'Azure OpenAI',
    id: 2,
    icon: 'azure'
  },
  {
    title: '',
    id: 3,
    icon: 'gemini'
  },
  {
    title: '',
    id: 4,
    icon: 'anthrop'
  },
  {
    title: t('aiModelWenxinyiyan'),
    id: 5,
    icon: 'wenxinyiyan'
  },
  {
    title: t('aiModelYuezhianmian'),
    id: 6,
    icon: 'yuezhianmian'
  },
  {
    title: t('aiModelTongyiqianwen'),
    id: 7,
    icon: 'tongyiqianwen'
  },
  {
    title: t('aiModelXunfeixinghuo'),
    id: 8,
    icon: 'xunfeixinghuo'
  },
  {
    title: t('aiModelHuoshanyinqin'),
    id: 9,
    icon: 'huoshanyinqin'
  },
  {
    title: t('aiModelBaichuan'),
    id: 10,
    icon: 'baichuan'
  },
  {
    title: t('aiModelLingyiwanwu'),
    id: 11,
    icon: 'lingyiwanwu'
  },
  {
    title: t('aiModelZhipuai'),
    id: 12,
    icon: 'zhipuai'
  },
  {
    title: t('aiModelTengxunhunyuan'),
    id: 13,
    icon: 'tengxunhunyuan'
  },
  {
    title: 'Cohere',
    id: 14,
    icon: 'cohere'
  },
  {
    title: 'Deepseek',
    id: 15,
    icon: 'deepseek'
  },
  {
    title: 'Jina',
    id: 16,
    icon: 'jina'
  }
])

const openSource = computed(() => [
  {
    title: t('sourceOneTitle'),
    id: 1,
    icon: 'out-the-box',
    info: t('sourceOneInfo'),
    list: [
      {
        title: t('sourceOneListOneTitle'),
        icon: 'login-success',
        listId: '1-1'
      },
      {
        title: t('sourceOneListTwoTitle'),
        icon: 'login-success',
        listId: '1-2'
      },
      {
        title: t('sourceOneListThreeTitle'),
        icon: 'login-success',
        listId: '1-3'
      }
    ]
  },
  {
    title: t('sourceTwoTitle'),
    id: 2,
    icon: 'data-security',
    info: t('sourceTwoInfo'),
    list: [
      {
        title: t('sourceTwoListOneTitle'),
        icon: 'login-success',
        listId: '2-1'
      },
      {
        title: t('sourceTwoListTwoTitle'),
        icon: 'login-success',
        listId: '2-2'
      },
      {
        title: t('sourceTwoListThreeTitle'),
        icon: 'login-success',
        listId: '2-3'
      }
    ]
  },
  {
    title: t('sourceThreeTitle'),
    id: 3,
    icon: 'multimode-box',
    info: t('sourceThreeInfo'),
    list: [
      {
        title: t('sourceThreeListOneTitle'),
        icon: 'login-success',
        listId: '3-1'
      },
      {
        title: t('sourceThreeListTwoTitle'),
        icon: 'login-success',
        listId: '3-2'
      },
      {
        title: t('sourceThreeListThreeTitle'),
        icon: 'login-success',
        listId: '3-3'
      },
      {
        title: t('sourceThreeListFourTitle'),
        icon: 'login-success',
        listId: '3-4'
      },
      {
        title: t('sourceThreeListFiveTitle'),
        icon: 'login-success',
        listId: '3-5'
      }
    ]
  },
  {
    title: t('sourceFourTitle'),
    id: 4,
    icon: 'more-compatible',
    info: t('sourceFourInfo'),
    list: [
      {
        title: t('sourceFourListOneTitle'),
        icon: 'login-success',
        listId: '4-1'
      },
      {
        title: t('sourceFourListTwoTitle'),
        icon: 'login-success',
        listId: '4-2'
      }
    ]
  }
])

const { replace } = useRouter()
const route = useRoute()

const redirect = computed(() => {
  return route.query.redirect ? decodeURIComponent(route.query.redirect) : '/'
})

const formState = reactive({
  username: '',
  password: '',
  remember: true
})

const userStore = useUserStore()

const onFinish = () => {
  handleLogin()
}
const onFinishFailed = (errorInfo) => {
  console.log('Failed:', errorInfo)
}

const handleLogin = () => {
  userStore
    .login({
      toHome: true,
      username: formState.username,
      password: formState.password
    })
    .then(() => {
      replace(redirect.value || '/')
    })
    .catch((err) => {
      console.log(err.message)
    })
}

const toChatWiki = () => {
  let currentUrl = btoa(window.location.href);
  
  window.open(`https://chatwiki.com?source=${currentUrl}`, '_blank');
}

const toXiaokefu = () => {
  let currentUrl = btoa(window.location.href);
  
  window.open(`https://xiaokefu.com.cn?source=${currentUrl}`, '_blank');
}

const toHome = () => {
  window.location.reload();
}

const bgColor = ref('rgba(255, 255, 255, 0)');
 
function handleScroll() {
  const threshold = 50; // 滚动多少像素时改变颜色
  const scrollTop = window.scrollY;
  const newBgColor = `rgba(255, 255, 255, ${scrollTop > threshold ? 1 : scrollTop / threshold})`;
  bgColor.value = newBgColor;
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
});
 
onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});

</script>

<style lang="less" scoped>

.main-box {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
}

.project-name {
  position: absolute;
  top: -53px;
  width: calc(100% - 40px);
  text-align: center;
  font-weight: bolder;
  color: #2475FC;
  font-size: 24px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.sign-in {
  position: absolute;
  top: 200px;
  right: 178px;
  background-color: #ffffff;
  padding: 20px;
  width: 382px;
  height: 334px;
  border-radius: 16px;
  padding: 0px 20px;
  box-shadow: 0 4px 32px 0 #1e44703d;

  .sign-in-content {
    height: 100%;
  }
  .col-img img {
    width: 100%;
    max-width: 500px;
    margin: auto;
    display: block;

    @media (min-width: 992px) {
      margin: 0;
    }
  }
  .login-title {
    font-weight: bolder;
    color: #000;
    margin: 48px 0 28px 0;
    text-align: center;
  }
  .text-muted {
    color: #8c8c8c !important;
  }
  h5 {
    font-size: 20px;
    margin-bottom: 0.5em;
  }
  .login-form {
    font-weight: 700;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .login-item {
    display: flex;
    width: 294px;
    padding: 8px 12px;
    height: 40px;
    box-sizing: border-box;
    align-items: center;
    border-radius: 6px;
    background: #FFF;
  }

  .login-btn {
    height: 40px;
    box-sizing: border-box;
    display: flex;
    width: 294px;
    padding: 8px 12px;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
  }
}

body {
  background-color: #ffffff;
}

.header {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1;
  display: flex;
  width: 100%;
  height: 52px;
  padding: 0 24px 0 20px;
  transition: background-color 0.3s; /* 平滑过渡效果 */
  background-color: v-bind(bgColor); /* 绑定计算出的背景颜色 */
  opacity: 0.95;

  .header-left {
    width: 80%;
    display: flex;

    .logo {
      display: flex;
      align-items: center;
      margin-right: 37px;
    }

    .header-item {
      font-family: "PingFang SC";
      font-size: 16px;
      font-style: normal;
      font-weight: 400;
      cursor: pointer;
      width: auto;
      padding: 0 16px;
      display: flex;      
      justify-content: center;
      align-items: center;
      color: #3A4559;

      &:hover {
        color: #2475fc;
      }
    }

  }

  .header-right {
    display: flex;
    align-items: center;
    justify-content: end;
    width: 20%;

    .langs {
      padding: 5px 16px;
      justify-content: center;
      align-items: center;
      border-radius: 16px;
      border: 1px solid #C3CBD9;
      width: 108px;
      height: 32px;
    }
  }
}

.banner {
  width: 100%;
  height: 652px;
  background-image: url('@/assets/img/login/banner.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.instructions-box {
  // width: 100%;
  width: 1200px;
  height: 80px;
  display: flex;
  justify-content: space-between;
  padding: 20px;
  align-items: center;
  margin-top: -40px;

  .instructions-item {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    width: 370px;
    height: 80px;
    color: #164799;
    font-family: "PingFang SC";
    font-size: 20px;
    font-weight: 600;
    padding-left: 80px;
    background-size: contain;
    background-position: center;
    background-repeat: no-repeat;
    box-shadow: 0 8px 32px 0 #2042741f;
    border-radius: 16px;
  }

  .one-deployment {
    background-image: url('@/assets/img/login/one-deployment.png');
  }

  .out-the-box {
    background-image: url('@/assets/img/login/out-the-box.png');
  }

  .safe-controllable {
    background-image: url('@/assets/img/login/safe-controllable.png');
  }
}

.profession-item-box {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px;
  height: 514px;
  width: 1200px;
}

.profession-title {
  margin-top: 60px;
  color: #262626;
  text-align: center;
  font-family: "PingFang SC";
  font-size: 28px;
  font-style: normal;
  font-weight: 600;
  line-height: 36px;
  width: 1200px;

  .profession-title-stress {
    color: #2475fc;
    text-align: center;
    font-family: "PingFang SC";
    font-size: 28px;
    font-style: normal;
    font-weight: 600;
    line-height: 36px;
  }
}

.profession-info {
  margin: 6px 0 20px 0;
  color: #595959;
  text-align: center;
  font-family: "PingFang SC";
  font-size: 16px;
  font-style: normal;
  font-weight: 400;
  line-height: 24px;
  width: 1200px;
}

.profession-item {
  display: flex;
  flex-direction: column;
  flex-basis: calc((100% - 60px)/ 4);
  height: 230px;
  width: 250px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 0 16px 0 #1e457314;
  padding: 24px;
  border: 1px solid #E5EFFF;

  .profession-item-top {
    display: flex;
    flex-direction: row;
    width: 100%;
    align-items: center;

    .profession-item-icon-svg {
      font-size: 40px;
    }

    .profession-item-title {
      color: #242933;
      text-align: center;
      font-family: "PingFang SC";
      font-size: 20px;
      font-style: normal;
      font-weight: 600;
      line-height: 28px;
      margin-left: 12px;
    }
  }

  .profession-item-middle {
    flex-shrink: 0;
    border-radius: 1px;
    border-top: 1px solid var(--01-, #E5EFFF);
    width: 100%;
    margin: 24px 0 16px 0;
  }

  .profession-item-bottom {
    width: 234px;
    color: #595959;
    font-family: "PingFang SC";
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
}

.ai-content {
  width: 1200px;
}

.ai-model-box {
  margin-top: 20px;
  width: 100%;
  display: flex;
  justify-content: center;
  height: 640px;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-image: url('@/assets/img/login/ai-model-bg.png');

  .ai-model-title {
    padding: 60px 0 50px;
    color: #262626;
    text-align: center;
    font-family: "PingFang SC";
    font-size: 28px;
    font-style: normal;
    font-weight: 600;
    line-height: 36px;
    width: 1200px;

    .ai-model-title-stress {
      color: #2475fc;
      text-align: center;
      font-family: "PingFang SC";
      font-size: 28px;
      font-style: normal;
      font-weight: 600;
      line-height: 36px;
    }
  }

  .ai-model-item-box {
    display: flex;
    flex-wrap: wrap;
    align-items: flex-start;
    justify-content: space-between;
    padding: 20px;
    height: 450px;
    width: 1200px;

    .ai-model-item {
      display: flex;
      align-items: center;
      justify-content: center;
      flex-basis: calc((100% - 60px)/ 4);
      height: 82px;
      background: #fff;
      border-radius: 8px;
      padding: 24px;

      .ai-model-item-svg {
        font-size: 34px;
      }

      .ai-model-item-svg-gemini {
        font-size: 94px;
      }

      .ai-model-item-svg-anthrop {
        font-size: 155px;
      }

      .ai-model-item-title {
        margin-left: 8px;
        color: #262626;
        text-align: center;
        font-family: "PingFang SC";
        font-size: 20px;
        font-style: normal;
        font-weight: 600;
        line-height: 32px;
      }
    }
  }
}

.open-source-box{
  // width: 100%;
  width: 1200px;

  .open-source-title {
    padding: 60px 0 2px;
    color: #262626;
    text-align: center;
    font-family: "PingFang SC";
    font-size: 28px;
    font-style: normal;
    font-weight: 600;
    line-height: 36px;
  }

  .open-source-item-box {
    display: flex;
    background-color: #fff;
    padding: 30px;

    .open-source-item:first-child {
      border-top-left-radius: 16px;
      border-bottom-left-radius: 16px;
    }

    .open-source-item:last-child {
      border-top-right-radius: 16px;
      border-bottom-right-radius: 16px;
      border-right: 1px solid #EDEFF2;
    }

    .open-source-item {
      position: relative;
      flex: 1;
      min-height: 488px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: flex-start;
      background: #fff;
      padding: 24px 24px 138px 24px;
      border-top: 1px solid #EDEFF2;
      border-left: 1px solid #EDEFF2;
      border-bottom: 1px solid #EDEFF2;

      .open-source-item-svg {
        position: absolute;
        bottom: 0px;
        right: 0px;
        font-size: 80px;
      }

      .open-source-item-title {
        color: #262626;
        font-family: "PingFang SC";
        font-size: 20px;
        font-style: normal;
        font-weight: 600;
        width: 100%;
      }

      .open-source-item-info {
        color: #262626;
        font-family: "PingFang SC";
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        padding: 38px 0 46px;
        width: 100%;
      }

      .open-source-item-list {
        display: flex;
        justify-content: flex-start;
        align-items: center;
        width: 100%;
        line-height: 32px;

        .open-source-item-list-svg {
          font-size: 16px;
        }

        .open-source-item-list-title {
          color: #595959;
          font-family: "PingFang SC";
          font-size: 14px;
          font-style: normal;
          font-weight: 400;
          padding-left: 4px;
        }
      }
    }
  }
}

.layout-footer {
  height: 69px;
  padding: 16px 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-top: 1px solid #F0F0F0;

  .copyright-text-box {
    display: flex;
    align-items: center;
    cursor: pointer;
    line-height: 20px;
    font-size: 12px;
    color: #8c8c8c;
    text-align: center;

    .copyright-text-info {
      cursor: auto;
      margin-right: 8px;
    }
  }

  .copyright-text {
    cursor: pointer;
    line-height: 20px;
    font-size: 12px;
    color: #8c8c8c;
    text-align: center;

    &:hover {
      color: #2475fc;
    }
  }
  .user-agreement-box {
    display: flex;
    justify-content: center;
    line-height: 20px;
    margin-bottom: 4px;
    font-size: 12px;

    .link-item {
      margin: 0 8px;
      color: #595959;
    }
  }

  .footer-line {
    border-right: 2px solid #DCDCDC;
    height: 15px;
    margin: 0 10px;
  }
}

.login-box.lang-en-US {
  .open-source-item-box{
    flex-flow: column;
    padding: 0;
    margin: 30px 0 60px 0;
    border: 1px solid #EDEFF2;
    border-radius: 16px;
    overflow: hidden;

    .open-source-item{
      min-height: 200px;
      padding-bottom: 50px;
      border: none;
      border-radius: 0;
      border-bottom: 1px solid #EDEFF2;
    }

    .open-source-item:last-child {
      border: none;
    }
  }
}
</style>