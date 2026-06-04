<script setup lang="ts">
/**
 * 一個用於創建透明圖標佔位符的組件
 * 它渲染一個指定尺寸的透明空間，可用於在佈局中佔位
 */
const props = defineProps<{
    /**
     * 盒子寬度
     */
    width?: string
    /**
     * 盒子高度
     */
    height?: string
}>()
</script>
<template>
    <div :style="{
        width: props.width ?? '14px',
        height: props.height ?? '14px',
        display: 'inline-block',
    }"></div>
</template>
