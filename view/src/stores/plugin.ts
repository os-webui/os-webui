import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useNav = defineStore('nav', () => {
  const pluginId = ref('')
  const pluginName = ref('')
  const setPlugin = (val: string, name?: string) => {
    pluginId.value = val
    pluginName.value = name ?? val
  }
  return {
    pluginId,
    pluginName,
    setPlugin,
  }
})

