import { createRouter, createWebHistory } from 'vue-router'
import OpencLayout from '@/layouts/open-layout/index.vue'
import { useOpenDocStore } from '@/stores/open-doc'
import NotFound from '@/views/404/404.vue'

const docBeforeEnter = async (to, from, next) => {
  let preview = to.query.preview
  let token = to.query.token

  const store = useOpenDocStore()

  if (preview) {
    store.setPreview(preview)
  }

  if (token) {
    store.setToken(token)
  }

  next()
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'index',
      component: OpencLayout,
      // redirect: '/404',
      beforeEnter: async (to, from, next) => {
        const store = useOpenDocStore()
        let libraryList = store.libraryList
        let defaultLibraryKey = ''
        let lastOpenDocLibraryKey = localStorage.getItem('last_open_doc_library_key')

        if (libraryList.length === 0) {
          let params = {}

          if (to.path.indexOf('/home') > -1 && !store.libraryInfo.library_key) {
            params.library_key = to.params.id
          }

          libraryList = await store.getLibraryList(params)
        }

        if (libraryList.length === 0) {
          next('/404')
          return
        }

        if (to.path != '/') {
          next()
          return
        }

        if (lastOpenDocLibraryKey) {
          let lastOpenLibrary = libraryList.filter(
            (item) => item.library_key == lastOpenDocLibraryKey,
          )

          if (lastOpenLibrary.length > 0) {
            defaultLibraryKey = lastOpenDocLibraryKey
          }
        } else {
          defaultLibraryKey = libraryList[0].library_key
        }

        localStorage.setItem('last_open_doc_library_key', defaultLibraryKey)

        next('/home/' + defaultLibraryKey)
      },
      children: [
        {
          path: '/home/:id',
          name: 'open-home',
          component: () => import('@/views/open/home/index.vue'),
          beforeEnter: docBeforeEnter,
        },
        {
          path: '/doc/:id',
          name: 'open-doc',
          component: () => import('@/views/open/doc/index.vue'),
          beforeEnter: docBeforeEnter,
        },
        {
          path: '/search/:id',
          name: 'open-search',
          component: () => import('@/views/open/search/index.vue'),
          meta: {
            hideHeader: true,
          },
          beforeEnter: docBeforeEnter,
        },
      ],
    },
    {
      path: '/404',
      name: 'not-found',
      component: NotFound,
    },
    // 通配符路由，捕获所有未匹配的路径
    {
      path: '/:pathMatch(.*)*',
      redirect: '/404',
    },
  ],
})

export default router
