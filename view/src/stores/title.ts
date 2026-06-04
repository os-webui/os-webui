import { DefaultHttpClient } from '@/internal/api'
import { defineStore } from 'pinia'
import { computed, ref, shallowRef, type ShallowRef } from 'vue'

export const useTitle = defineStore('title', () => {
  const value :ShallowRef<string[]>= shallowRef([])
  const get = computed(() => value.value)
  const set = (val: string | string[]) => {
    value.value = typeof val === 'string' ? [val] : val
  }
  const title = ref('OS WebUI')
  let fetch = false
  const updateTitle = async () => {
    if (fetch) {
      return
    }
    fetch = true
    try {
      const resp = await DefaultHttpClient.get<string>('/api/title')
      if (typeof resp.data === "string") {
        title.value = resp.data
        return
      }
    } catch (e) {
      console.log(e)
    }
    fetch = false
  }
  return {
    get,
    set,
    title,
    updateTitle,
  }
})

