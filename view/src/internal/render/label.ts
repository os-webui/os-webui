import type { Component, VNode } from 'vue'
import { RouterLink } from 'vue-router';
import { h } from 'vue'
import {
  NIcon,
} from 'naive-ui'

const RENDER_LABEL_CSS = 'flex flex-row align-items-center align-content-center gap-2'

export interface RenderLabelOptions {
  /**
   * 要渲染的文本內容
   */
  label?: string | (() => string)
  /**
   * 一個可選的圖標，會顯示在 label 之前
   */
  prefix?: Component
  /**
   * 一個可選的圖標，會顯示在 label 之後
   */
  suffix?: Component
  /**
   * 可選的動態渲染函數，會顯示在 label 之前，如果設置則忽略 prefix 屬性
   */
  prefixRender?: () => Component
  /**
   * 可選的動態渲染函數，會顯示在 label 之後，如果設置則忽略 suffix 屬性
   */
  suffixRender?: () => Component
  /**
   * 可選的 url，當最終結果爲字符串時，如果以 'http://' 或 'https://' 爲前綴將創建 a 標籤，否則創建 RouterLink
   */
  url?: string | (() => null | undefined | string)
  /**
   * 可選的自定義 CSS 類，附加到容器元素
   */
  class?: string | string[]
  /**
   * 可選的 target 屬性，僅對 <a> 標籤有效
   */
  target?: '_blank' | '_self' | '_parent' | '_top'
  /**
   * 可選的 aria-label 屬性，僅對 <a> 標籤 和 RouterLink 有效
   */
  ariaLabel?: string | (() => null | undefined | string)

  /**
   * 可選的點擊事件處理函數，應用於 <a>、RouterLink 或 div
   */
  onClick?: (event: MouseEvent) => void
  /**
      * 可選的交互事件處理函數，處理鼠標點擊、鍵盤 Enter 鍵和電視遙控器確定按鈕，應用於 <a>、RouterLink 或 div
      */
  onTab?: (event: Event) => void;
}

function renderDivStatic() {
  return h('div', { class: 'flex-grow-1' });
}
function renderIconStatic() {
  return h('div', { class: 'flex-grow-0', 'aria-hidden': 'true' });
}
function renderUrlStatic(url: null | undefined | string,
  children: Array<() => VNode>,
  css: string,
  ariaLabel?: string | (() => null | undefined | string),
  target?: '_blank' | '_self' | '_parent' | '_top',
  onClick?: (event: MouseEvent) => void,
  onTab?: (event: Event) => void,
) {
  const handleTab = onTab
    ? (event: KeyboardEvent) => {
      if (event.type === 'keyup' && (event as KeyboardEvent).key === 'Enter') {
        onTab(event)
      }
    }
    : undefined;
  if (onTab) {
    onClick = onTab
  }

  if (url === null || url === undefined) {
    return h('div',
      {
        class: css,
        onClick: onClick,
        onKeyup: handleTab,
        tabindex: onTab ? '0' : undefined,
      },
      children.map((f) => f()))
  }

  const finalLabel = typeof ariaLabel === 'function' ? ariaLabel() || `go to ${url}` : ariaLabel || `go to ${url}`;

  if (url.startsWith('http://') || url.startsWith('https://')) {
    target = target || '_blank'
    return h('a',
      {
        href: url,
        class: css,
        target: target,
        rel: target === '_blank' ? 'noopener noreferrer' : undefined,
        'aria-label': finalLabel,
        onClick: onClick,
        onKeyup: handleTab,
        tabindex: onTab ? '0' : undefined,
      },
      children.map((f) => f()))
  }
  return h(RouterLink,
    {
      to: url,
      class: css,
      'aria-label': finalLabel,
      onClick: onClick,
      onKeyup: handleTab,
      tabindex: onTab ? '0' : undefined,
    },
    () => children.map((f) => f()))
}
/**
 * 返回一個 render 函數，用於渲染一個 label
 * - 它可以有個可選的前綴圖標
 * - 它可以有個可選的後綴圖標
 * - 如果 url 被設置，它將被渲染在 a 標籤 或 RouterLink 中
 * @param opts 渲染選項
 */
export function renderLabel(options: RenderLabelOptions) {
  let renderPrefix: () => VNode
  if (options.prefixRender) {
    const render = options.prefixRender
    renderPrefix = () => h('div',
      { 'class': 'flex align-items-center justify-content-center' },
      h(NIcon,
        { 'aria-hidden': 'true' },
        {
          default: () => h(render())
        },
      ),
    )
  } else if (options.prefix) {
    const tag = options.prefix
    renderPrefix = () => h('div',
      { 'class': 'flex align-items-center justify-content-center' },
      h(NIcon,
        { 'aria-hidden': 'true' },
        {
          default: () => h(tag)
        },
      ),
    )
  } else {
    renderPrefix = renderIconStatic
  }
  let renderLabel: () => VNode
  const label = options.label
  switch (typeof label) {
    case 'string':
      renderLabel = () => h('div', { class: 'flex align-items-center justify-content-start flex-grow-1' }, label)
      break
    case 'function':
      renderLabel = () => h('div', { class: 'flex align-items-center justify-content-start flex-grow-1' }, label())
      break
    default:
      renderLabel = renderDivStatic
      break
  }
  let renderSuffix: () => VNode
  if (options.suffixRender) {
    const render = options.suffixRender
    renderSuffix = () => h('div',
      { 'class': 'flex align-items-center justify-content-center' },
      h(NIcon,
        { 'aria-hidden': 'true' },
        {
          default: () => h(render())
        },
      ),
    )
  } else if (options.suffix) {
    const tag = options.suffix
    renderSuffix = () => h('div',
      { 'class': 'flex align-items-center justify-content-center' },
      h(NIcon,
        { 'aria-hidden': 'true' },
        {
          default: () => h(tag)
        },
      ),
    )
  } else {
    renderSuffix = renderIconStatic
  }
  const css = [RENDER_LABEL_CSS, ...(Array.isArray(options.class) ? options.class : [options.class])]
    .filter(Boolean)
    .join(' ');
  const url = options.url
  switch (typeof url) {
    case 'string':
      return () => renderUrlStatic(url,
        [
          renderPrefix,
          renderLabel,
          renderSuffix,
        ],
        css,
        options.ariaLabel, options.target,
        options.onClick, options.onTab,
      )
    case 'function':
      return () => renderUrlStatic(url(),
        [
          renderPrefix,
          renderLabel,
          renderSuffix,
        ],
        css,
        options.ariaLabel, options.target,
        options.onClick, options.onTab,
      )
    default:
      return () => renderUrlStatic(null,
        [
          renderPrefix,
          renderLabel,
          renderSuffix,
        ],
        css,
        undefined, undefined,
        options.onClick, options.onTab,
      )
  }
}
