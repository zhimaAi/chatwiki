import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat', () => {
    const baseApiUrl= ref(import.meta.env.VITE_BASE_API_URL)

    return {
        baseApiUrl
    }
})