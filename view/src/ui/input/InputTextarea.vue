<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import { NInput } from 'naive-ui';
import { onMounted, ref, computed, onBeforeUnmount } from 'vue';

interface Props {
  readonly?: boolean;
  value?: string;
  autosize?: any;
  disabled?: boolean;
}

const props = defineProps<Props>();

// Calculate total lines based on newline characters to drive line numbers
// 透過換行符號計算總行數以驅動行號顯示
const totalLines = computed<number>(() => {
  if (!props.value) {
    return 1;
  }
  return props.value.split('\n').length || 1;
});

const lineNumbersRef = ref<HTMLDivElement | null>(null);
const isFocused = ref<boolean>(false); // Controls the Naive UI outer focus ring state / 控制外框呼吸燈的狀態變數

// Synchronize the vertical scroll position of line numbers with the textarea
// 讓左側行號的垂直捲動位置與右側文字框完全同步
const handleScroll = (event: Event): void => {
  const target = event.target as HTMLTextAreaElement;
  if (lineNumbersRef.value) {
    lineNumbersRef.value.scrollTop = target.scrollTop;
  }
};

const editorWrapperRef = ref<HTMLDivElement | null>(null);
let scrollTarget: HTMLElement | null = null;
let quit = false;

onMounted(() => {
  // Use a slight delay to ensure Naive UI has fully rendered the internal textarea element
  // 使用微小的延遲，確保 Naive UI 已經完全渲染出內部的 textarea 元素
  setTimeout(() => {
    if (editorWrapperRef.value && !quit) {
      scrollTarget = editorWrapperRef.value.querySelector('textarea');
      if (scrollTarget) {
        scrollTarget.addEventListener('scroll', handleScroll);
      }
    }
  }, 100);
});

onBeforeUnmount(() => {
  quit = true;
  if (scrollTarget) {
    scrollTarget.removeEventListener('scroll', handleScroll);
  }
});

const emit = defineEmits<{
  'update:value': [value: string];
}>();

function handleUpdateValue(value: string) {
  // block input if it is read-only or disabled pseudo-readonly
  // 如果是唯讀或禁用的偽唯讀狀態，則攔截更新
  if (!props.readonly && !props.disabled) {
    emit('update:value', value);
  }
}

const handleKeyDown = (event: KeyboardEvent): void => {
  // If disabled, block all shortcut actions immediately
  // 如果是禁用狀態，立刻攔截所有快捷鍵動作
  if (props.disabled) return;

  const textarea = event.target as HTMLTextAreaElement;
  const { selectionStart, selectionEnd } = textarea;
  const currentText = props.value ?? '';

  const insertTextWithHistory = (startPos: number, endPos: number, text: string, newCursorPos: number) => {
    textarea.focus();
    textarea.setSelectionRange(startPos, endPos);

    const success = document.execCommand('insertText', false, text);

    if (!success) {
      handleUpdateValue(
        currentText.substring(0, startPos) +
        text +
        currentText.substring(endPos)
      );
    }

    setTimeout(() => {
      textarea.selectionStart = textarea.selectionEnd = newCursorPos;
    }, 0);
  };

  // --- 1. Shift + Tab: Insert Indentation ---
  // --- 1. Shift + Tab: 插入縮排 ---
  if (event.key === 'Tab' && event.shiftKey) {
    event.preventDefault();
    const tabCharacter = '\t';
    insertTextWithHistory(selectionStart, selectionEnd, tabCharacter, selectionStart + tabCharacter.length);
    return;
  }

  // --- 2. Ctrl + Enter: Insert New Line BELOW Current Line ---
  // --- 2. Ctrl + Enter: 在當前行下方插入新行 ---
  if (event.key === 'Enter' && event.ctrlKey && !event.shiftKey) {
    event.preventDefault();
    const nextLineBreak = currentText.indexOf('\n', selectionStart);
    const insertPos = nextLineBreak === -1 ? currentText.length : nextLineBreak;
    insertTextWithHistory(insertPos, insertPos, '\n', insertPos + 1);
    return;
  }

  // --- 3. Shift + Enter: Insert New Line ABOVE Current Line ---
  // --- 3. Shift + Enter: 在當前行上方插入新行 ---
  if (event.key === 'Enter' && event.shiftKey && !event.ctrlKey) {
    event.preventDefault();
    const prevLineBreak = currentText.lastIndexOf('\n', selectionStart - 1);
    const insertPos = prevLineBreak === -1 ? 0 : prevLineBreak + 1;
    insertTextWithHistory(insertPos, insertPos, '\n', insertPos);
    return;
  }
};
</script>

<template>
  <div ref="editorWrapperRef" class="code-editor-container"
    :class="{ 'is-focused': isFocused, 'is-disabled': disabled }">
    <div ref="lineNumbersRef" class="line-numbers">
      <div v-for="line in totalLines" :key="line" class="line-number-item">
        {{ line }}
      </div>
    </div>

    <n-input type="textarea" :value="value" @update:value="handleUpdateValue" @keydown="handleKeyDown"
      @focus="isFocused = true" @blur="isFocused = false" :placeholder="$t('ui.configPlaceholder')" :disabled="false"
      :readonly="readonly || disabled" :autosize="autosize"
      :input-props="{ style: { whiteSpace: 'pre', overflowX: 'auto' } }" autofocus show-count autocapitalize="off"
      autocomplete="off" autocorrect="off" spellcheck="false" />
  </div>
</template>

<style scoped>
.n-input {
  height: 100%;
}

.code-editor-container {
  width: 100%;

  /* Variables for precise alignment between text and line numbers */
  /* 文字與行號精密對齊專用的尺寸變數 */
  --editor-font-family: 'Fira Code', Consolas, Monaco, 'Courier New', Courier, monospace;
  --editor-font-size: 14px;
  --editor-line-height: 1.5;
  --editor-padding-top: 8px;
  --editor-padding-bottom: 8px;

  /* Automatically inherit CSS variables injected by Naive UI's current theme */
  /* 自動繼承 Naive UI 當前主題注入的 CSS 變數環境 */
  background-color: var(--n-color);
  border: 1px solid var(--n-border-color);

  /* Compute dynamic styling for line numbers based on theme context */
  /* 根據主題內文自動調配行號區域的背景與文字色 */
  --line-number-bg: var(--n-code-bg-color, rgba(0, 0, 0, 0.03));
  --line-number-color: var(--n-placeholder-color);

  display: flex !important;
  border-radius: 3px;
  overflow: hidden;
  box-sizing: border-box;

  /* Replicate standard Naive UI transition curves */
  /* 複製標準的 Naive UI 漸變動畫曲線 */
  transition: border-color 0.2s var(--n-cubic-bezier-ease-in-out),
    box-shadow 0.2s var(--n-cubic-bezier-ease-in-out),
    background-color 0.3s var(--n-cubic-bezier-ease-in-out);
}

/* Hover effect (Only apply when not pseudo-disabled) */
/* 懸停狀態 (僅在非禁用狀態下啟用) */
.code-editor-container:not(.is-disabled):hover {
  border-color: var(--n-border-color-hover);
}

/* Focus glow effect (Only apply when not pseudo-disabled) */
/* 聚焦發光狀態 (僅在非禁用狀態下啟用) */
.code-editor-container.is-focused:not(.is-disabled) {
  border-color: var(--n-border-color-hover);
  box-shadow: var(--n-box-shadow-focus);
}

/* 🎨 SPECIAL CSS FOR PSEUDO-DISABLED STATE */
/* 🎨 核心修改：針對偽裝禁用狀態（is-disabled）的專屬特殊 CSS */
.code-editor-container.is-disabled {
  background-color: var(--n-color-disabled) !important;
  border-color: var(--n-border-color) !important;
  /* Lock the border without glow / 鎖定原始邊框顏色不發光 */
  box-shadow: none !important;
}

/* Change cursor behavior for the entire container and inner elements */
/* 變更整個容器及其內部所有元素的滑鼠游標為不允許狀態 */
.code-editor-container.is-disabled,
.code-editor-container.is-disabled .line-numbers,
.code-editor-container.is-disabled :deep(.n-input),
.code-editor-container.is-disabled :deep(.n-input .n-input__textarea-el) {
  cursor: not-allowed !important;
}

/* Force dim the text and line numbers colors using Naive UI's native disabled token */
/* 使用 Naive UI 原生的禁用文字變數，強行將行號與內文顏色調暗 */
.code-editor-container.is-disabled .line-numbers {
  background-color: var(--n-color-disabled);
  color: var(--n-text-color-disabled) !important;
  opacity: 0.6;
}

.code-editor-container.is-disabled :deep(.n-input .n-input__textarea-el) {
  color: var(--n-text-color-disabled) !important;
}

/* Strip internal Naive UI defaults to let container handle borders */
/* 移除 Naive UI 內部預設樣式，改由最外層 container 統一掌控邊框 */
:deep(.n-input) {
  border: none !important;
  background-color: transparent !important;
  --n-border: none !important;
  --n-border-hover: none !important;
  --n-border-focus: none !important;
  --n-box-shadow-focus: none !important;
}

/* 1. Left Side: Line Numbers Styling */
/* 1. 左側：行號區域樣式 */
.line-numbers {
  background-color: var(--line-number-bg);
  color: var(--line-number-color);
  font-family: var(--editor-font-family);
  font-size: var(--editor-font-size);
  line-height: var(--editor-line-height);
  padding-top: var(--editor-padding-top);
  padding-bottom: var(--editor-padding-bottom);
  text-align: right;
  width: 45px;
  user-select: none;
  overflow: hidden;
  /* Hide line numbers' native scrollbar / 隱藏行號欄自身的原生捲動條 */
  border-right: 1px solid var(--n-border-color);
  opacity: 0.75;
  box-sizing: border-box;
}

.line-number-item {
  padding-right: 8px;
  height: calc(var(--editor-font-size) * var(--editor-line-height));
}

/* 2. Right Side: Force Naive UI Internal Textarea to Match Exactly */
/* 2. 右側：強行覆寫 Naive UI 內部文字框，確保像素級對齊 */
:deep(.n-input .n-input__textarea-el) {
  font-family: var(--editor-font-family) !important;
  font-size: var(--editor-font-size) !important;
  line-height: var(--editor-line-height) !important;
  padding-top: var(--editor-padding-top) !important;
  padding-bottom: var(--editor-padding-bottom) !important;
  padding-left: 10px !important;
  padding-right: 10px !important;
}

/* Modern minimalist custom scrollbars syncing with current theme color */
/* 現代化極簡自訂捲動條，顏色隨當前主題自動變更 */
:deep(.n-input .n-input__textarea-el::-webkit-scrollbar) {
  width: 8px;
  height: 8px;
}

:deep(.n-input .n-input__textarea-el::-webkit-scrollbar-track) {
  background: transparent;
}

:deep(.n-input .n-input__textarea-el::-webkit-scrollbar-thumb) {
  background: var(--n-placeholder-color);
  opacity: 0.3;
  border-radius: 4px;
}

:deep(.n-input .n-input__textarea-el::-webkit-scrollbar-thumb:hover) {
  background: var(--n-text-color);
  opacity: 0.5;
}

/* Remove default padding offsets added by Naive UI wrappers */
/* 移除 Naive UI wrapper 造成的預設內距偏移 */
:deep(.n-input-wrapper) {
  padding: 0 !important;
}

/* Precision alignment for the floating placeholder container */
/* 精準對齊絕對定位的 placeholder 浮動容器 */
:deep(.n-input .n-input__placeholder) {
  left: 10px !important;
  top: var(--editor-padding-top) !important;
  /* Dynamically inherits the exact text top padding / 動態積累文字的頂部內距 */
  transform: none !important;
  /* Strips Naive's single-line vertical centering / 拔除 Naive UI 針對單行設計的垂直居中位移 */
}
</style>
