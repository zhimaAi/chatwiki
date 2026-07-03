import { ref } from 'vue'
import { defineStore } from 'pinia'
import { getMiniCardList, addMiniCard, updateMiniCard, deleteMiniCard } from '@/api/mini-card/index'

export const useMiniCardStore = defineStore('mini-card', () => {
  const cardList = ref([])
  const loading = ref(false)
  let lastFetchTime = 0
  const CACHE_DURATION = 5 * 60 * 1000 // 5分钟缓存

  const fetchCardList = async (force = false) => {
    // stale-while-revalidate：缓存未过期时直接返回，过期则后台刷新
    if (!force && cardList.value.length > 0 && Date.now() - lastFetchTime < CACHE_DURATION) {
      return cardList.value
    }

    loading.value = true
    try {
      const res = await getMiniCardList()
      cardList.value = res.data?.list || []
      lastFetchTime = Date.now()
      return cardList.value
    } finally {
      loading.value = false
    }
  }

  const addCard = async (data) => {
    const res = await addMiniCard(data)
    await fetchCardList(true)
    return res
  }

  const editCard = async (data) => {
    const res = await updateMiniCard(data)
    await fetchCardList(true)
    return res
  }

  const removeCard = async (data) => {
    const res = await deleteMiniCard(data)
    await fetchCardList(true)
    return res
  }

  return {
    cardList,
    loading,
    fetchCardList,
    addCard,
    editCard,
    removeCard
  }
})
