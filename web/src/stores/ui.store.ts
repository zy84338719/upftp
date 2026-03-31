import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUIStore = defineStore('ui', () => {
  const searchQuery = ref('')
  const isLoading = ref(false)

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  function setLoading(loading: boolean) {
    isLoading.value = loading
  }

  return {
    searchQuery,
    isLoading,
    setSearchQuery,
    setLoading
  }
})
