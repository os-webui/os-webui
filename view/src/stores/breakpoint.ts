import { computed, shallowRef } from 'vue'
import { defineStore } from 'pinia'

import { once } from '@/internal/utils'

/**
 * 定義屏幕尺寸
 */
export enum Breakpoint {
  /**
   * < 576px 的屏幕，例如：極小屏幕、穿戴式設備
   */
  mini = 0,
  /**
   * [576px, 768px) 的小屏幕，通常是手機屏幕
   */
  sm = 1,
  /**
   * [768px, 992px) 的中等屏幕，通常是大屏手機或平板
   */
  md = 2,
  /**
   * [992px, 1200px) 的大屏幕，通常是筆記本或桌面機
   */
  lg = 3,
  /**
   * >= 1200px 的超大屏幕，通常是桌面機或電視
   */
  xl = 4,
}

/**
 * 定義斷點 臨界值
 */
export const BREAKPOINTS = {
  sm: 576,
  md: 768,
  lg: 992,
  xl: 1200,
};
/**
 * 這個 store 用於窗口查詢結果以便與爲不同窗口實現響應式佈局
 */
export const useBreakpointStore = defineStore('breakpoint', () => {
  const isBrowser = typeof window !== 'undefined' && typeof window.matchMedia === 'function';
  const matchs = isBrowser ? [
    window.matchMedia(`(min-width: ${BREAKPOINTS.sm}px)`),
    window.matchMedia(`(min-width: ${BREAKPOINTS.md}px)`),
    window.matchMedia(`(min-width: ${BREAKPOINTS.lg}px)`),
    window.matchMedia(`(min-width: ${BREAKPOINTS.xl}px)`),
  ] : undefined
  const breakpoint = shallowRef(matchs ? calculate(matchs[0]!.matches, matchs[1]!.matches, matchs[2]!.matches, matchs[3]!.matches) : Breakpoint.lg)
  const handler = matchs ? (_: MediaQueryListEvent) => {
    breakpoint.value = calculate(matchs[0]!.matches, matchs[1]!.matches, matchs[2]!.matches, matchs[3]!.matches)
  } : undefined

  const miniOnly = computed(() => breakpoint.value == Breakpoint.mini)
  const smOnly = computed(() => breakpoint.value == Breakpoint.sm)
  const mdOnly = computed(() => breakpoint.value == Breakpoint.md)
  const lgOnly = computed(() => breakpoint.value == Breakpoint.lg)
  const xlOnly = computed(() => breakpoint.value == Breakpoint.xl)
  const sm = computed(() => breakpoint.value >= Breakpoint.sm)
  const md = computed(() => breakpoint.value >= Breakpoint.md)
  const lg = computed(() => breakpoint.value >= Breakpoint.lg)
  const xl = computed(() => breakpoint.value >= Breakpoint.xl)

  let counter = 0
  function start() {
    if (!matchs) {
      return nothing
    }
    if (counter == 0) {
      for (const match of matchs) {
        match.addEventListener("change", handler!)
      }
    }
    counter++

    return once(() => {
      counter--
      if (counter === 0) {
        for (const match of matchs) {
          match.removeEventListener("change", handler!)
        }
      }
    })
  }
  return {
    /**
     * 當前屏幕尺寸
     * @remarks 在 SSR 環境中，默認為 Breakpoint.lg（桌面），可通過 App.vue 或其他頂層組件修改
     * @example
     * // 在 App.vue 中設置 SSR 斷點
     * import { useBreakpointStore, Breakpoint } from '@/stores/breakpoint';
     * const { breakpoint } = useBreakpointStore();
     * if (!process.client) {
     *   breakpoint.value = Breakpoint.sm; // 模擬手機斷點
     * }
     */
    breakpoint,

    /**
     * block sm:hidden
     * - v-if="$breakpoint.miniOnly.value" 只在 < 576px 的屏幕上渲染，可能是電子手錶
     */
    miniOnly,
    /**
     * hidden sm:block md:hidden
     * - v-if="$breakpoint.smOnly.value" 只在 [576px, 768px) 的小屏幕上渲染，通常是手機屏幕
     */
    smOnly,
    /**
     * hidden md:block lg:hidden
     * - v-if="$breakpoint.mdOnly.value" 只在 [768px, 992px) 的中等屏幕渲染，通常是大屏手機或平板
     */
    mdOnly,
    /**
     * hidden lg:block xl:hidden
     * - v-if="$breakpoint.lgOnly.value" 只在 [992px, 1200px) 的大屏幕渲染，通常是筆記本或桌面機
     */
    lgOnly,
    /**
     * hidden xl:block
     * - v-if="$breakpoint.xlOnly.value" 只在 >= 1200px 的超大屏幕渲染，通常是桌面機或電視
     */
    xlOnly,
    /**
     * hidden sm:block
     * - v-if="$breakpoint.sm.value" 在屏幕尺寸 >= 576px 時渲染
     * - v-if="!$breakpoint.sm.value" 在屏幕尺寸 < 576px 時渲染
     */
    sm,
    /**
     * hidden md:block
     * - v-if="$breakpoint.md.value" 在屏幕尺寸 >= 768px 時渲染
     * - v-if="!$breakpoint.md.value" 在屏幕尺寸 < 768px 時渲染
     */
    md,
    /**
     * hidden lg:block
     * - v-if="$breakpoint.lg.value" 在屏幕尺寸 >= 992px 時渲染
     * - v-if="!$breakpoint.lg.value" 在屏幕尺寸 < 992px 時渲染
     */
    lg,
    /**
     * hidden xl:block
     * - v-if="$breakpoint.xl.value" 在屏幕尺寸 >= 1200px 時渲染
     * - v-if="!$breakpoint.xl.value" 在屏幕尺寸 < 1200px 時渲染
     */
    xl,

    /**
     * 開始監聽屏幕尺寸變化
     * @returns {Function} 清理函數，必須在組件卸載時調用（例如，onUnmounted）
     * @example
     * onMounted(() => { cleanup = start(); });
     * onUnmounted(() => { cleanup(); });
     */
    start,
  }
})
function nothing() { }
/**
 * 計算屏幕尺寸屬於哪個級別
 */
function calculate(sm: boolean, md: boolean, lg: boolean, xl: boolean): Breakpoint {
  if (xl) {
    return Breakpoint.xl
  }
  if (lg) {
    return Breakpoint.lg
  }
  if (md) {
    return Breakpoint.md
  }
  if (sm) {
    return Breakpoint.sm
  }
  return Breakpoint.mini
}
