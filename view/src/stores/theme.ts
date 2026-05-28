import { shallowRef, computed } from 'vue'
import { defineStore } from 'pinia'
import { darkTheme, lightTheme, useOsTheme } from 'naive-ui'
import { saveValue } from '@/internal/utils'

const Key = 'theme'

/**
 * 這個 store 用於存儲用戶選擇的頁面主題風格
 */
export const useThemeStore = defineStore('theme', () => {
  /**
   * 跟蹤操作系統的配色風格
   */
  const os = shallowRef(useOsTheme())

  /**
   * 用戶設置
   */
  const userdata = shallowRef(loadValue())

  const name = computed(() => {
    switch (userdata.value) {
      case 'dark':
        return 'dark'
      case 'light':
        return 'light'
    }
    return os.value == 'dark' ? 'dark' : 'light'
  })

  const theme = computed(() => {
    switch (userdata.value) {
      case 'dark':
        return darkTheme
      case 'light':
        return lightTheme
    }
    return os.value == 'dark' ? darkTheme : lightTheme
  })

  const choose = computed({
    get() {
      return userdata.value as string
    },
    set(val) {
      switch (val) {
        case 'dark':
        case 'light':
        case 'auto':
          break
        default:
          console.warn('save auto theme, because input unknown theme:', val)
          val = 'auto'
          break
      }
      if (val != userdata.value) {
        saveValue(Key, val, 'auto')
        userdata.value = val
      }
    },
  })
  return {
    /**
     * 返回當前 正在使用的 主題名稱
     */
    name,
    /**
     * 返回當前 應該採用的 主題環境
     */
    theme,
    /**
     * 用戶選擇的主題
     */
    choose,
  }
})

function loadValue(): string {
  try {
    const s = localStorage.getItem(Key) || 'auto'
    switch (s) {
      case 'dark':
      case 'light':
      case 'auto':
        console.info(`load "${Key}":`, s)
        return s
      default:
        console.warn(`use "auto" theme, because load "${Key}" unknown:`, s)
        return 'auto'
    }
  } catch (e) {
    console.warn(`load "${Key}" fail:`, e)
    return 'auto'
  }
}
