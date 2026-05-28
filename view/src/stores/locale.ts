/* eslint-disable @typescript-eslint/no-explicit-any */
import { ref, computed, shallowRef } from 'vue'
import { defineStore } from 'pinia'

import {
  zhTW, dateZhTW,
  zhCN, dateZhCN,
  enUS, dateEnUS,
} from 'naive-ui'
import type { NDateLocale } from 'naive-ui'
import { equal, once, saveValue } from '@/internal/utils'

const Key = 'locale'

export interface Locale {
  id: string
  name: string
  message: any
  date: NDateLocale
}

/**
 * 支持的語言列表
 */
export const Locales: Locale[] = [
  {
    id: 'en-us',
    name: '🇺🇸 English',
    message: enUS,
    date: dateEnUS,
  },
  {
    id: 'zh-tw',
    name: '🇹🇼 繁體中文',
    message: zhTW,
    date: dateZhTW,
  },
  {
    id: 'zh-cn',
    name: '🇨🇳 简体中文',
    message: zhCN,
    date: dateZhCN,
  },
]
const indexEN_US = 0
const indexZH_TW = 1
const indexZH_CN = 2
/**
 * 計算當前應該顯示的語言
 */
function calculate(id: string, languages: readonly string[]): Locale {
  if (id != 'auto') {
    for (const locale of Locales) {
      if (locale.id === id) {
        return locale
      }
    }
  }
  for (const lang of languages) {
    switch (lang) {
      case 'zh':
        return Locales[indexZH_TW]!
      case 'en':
        return Locales[indexEN_US]!
    }
    if (lang.length > 3 && lang[2] == '-') {
      if (lang.startsWith('zh-')) {
        id = lang.substring(3).toLocaleLowerCase()
        if (id.indexOf('hans') != -1 || id.indexOf('cn') != -1) {
          return Locales[indexZH_CN]!
        } else {
          return Locales[indexZH_TW]!
        }
      } else if (lang.startsWith('en-')) {
        return Locales[indexEN_US]!
      }
    }
  }
  return Locales[indexEN_US]!
}

/**
 * 這個 store 用於存儲用戶選擇的頁面語言
 */
export const useLocaleStore = defineStore('locale', () => {
  /**
   * 系統設置
   */
  const languages = shallowRef(Array.from(navigator.languages || []))
  const handler = () => {
    const vals = navigator.languages || []
    if (!equal(languages.value, vals)) {
      languages.value = Array.from(vals)
    }
  }

  /**
   * 用戶設置
   */
  const userdata = ref(loadValue())

  const locale = computed(() => calculate(userdata.value, languages.value))
  const choose = computed({
    get() {
      return userdata.value
    },
    set(val) {
      if (val !== 'auto') {
        for (const locale of Locales) {
          if (locale.id === val) {
            val = locale.id
            break
          }
        }
      }
      if (val != userdata.value) {
        saveValue(Key, val, 'auto')
        userdata.value = val
      }
    },
  })

  let counter = 0
  function start() {
    if (counter == 0) {
      window.addEventListener('languagechange', handler)
    }
    counter++

    return once(() => {
      counter--
      if (counter === 0) {
        window.removeEventListener('languagechange', handler)
      }
    })
  }
  return {
    /**
     * 返回 當前應該顯示的語言環境
     */
    locale,
    /**
     * 用戶選擇的語言
     */
    choose,
    /**
     * 開始監聽系統語言環境變化
     * @returns {Function} 清理函數，必須在組件卸載時調用（例如，onUnmounted）
     * @example
     * onMounted(() => { cleanup = start(); });
     * onUnmounted(() => { cleanup(); });
     */
    start,
  }
})


function loadValue() {
  try {
    const s = localStorage.getItem(Key) || 'auto'
    if (s === 'auto') {
      console.info(`load "${Key}":`, s)
      return s
    }
    for (const locale of Locales) {
      if (locale.id == s) {
        console.info(`load "${Key}":`, s)
        return s
      }
    }
    console.warn(`use "auto" locale, because load "${Key}" unknown:`, s)
  } catch (e) {
    console.warn(`load "${Key}" fail:`, e)
  }
  return 'auto'
}
