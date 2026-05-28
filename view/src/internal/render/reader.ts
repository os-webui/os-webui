import type { Component, VNode } from 'vue'
import { h } from 'vue'
import {
    NIcon,
} from 'naive-ui'
import type { IconProps } from 'naive-ui'

/**
 * 返回一個 render 函數，用於渲染一個靜態圖標
 * @param icon 要渲染的圖標組件
 */
export function renderIconStatic(icon: Component, props?: IconProps) {
    return () => h(NIcon, props, {
        default: () => h(icon)
    })
}

