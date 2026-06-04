<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import {
  NInput,
} from 'naive-ui';
import { onMounted, ref, computed, onBeforeUnmount } from 'vue';
interface Props {
  readonly?: boolean
  value?: string
  autosize?: any
  disabled?: boolean
}
const props = defineProps<Props>()
const totalLines = computed<number>(() => {
  if (!props.value) {
    return 1
  }
  return props.value.split('\n').length || 1;
})
const lineNumbersRef = ref<HTMLDivElement | null>(null)
const handleScroll = (event: Event): void => {
  const target = event.target as HTMLTextAreaElement
  if (lineNumbersRef.value) {
    lineNumbersRef.value.scrollTop = target.scrollTop
  }
}
const editorWrapperRef = ref<HTMLDivElement | null>(null)
let scrollTarget: HTMLElement | null = null
let quit = false
onMounted(() => {
  setTimeout(() => {
    if (editorWrapperRef.value && !quit) {
      scrollTarget = editorWrapperRef.value.querySelector('textarea')
      if (scrollTarget) {
        scrollTarget.addEventListener('scroll', handleScroll)
      }
    }
  }, 100)
})
onBeforeUnmount(() => {
  quit = true
  if (scrollTarget) {
    scrollTarget.removeEventListener('scroll', handleScroll)
  }
})
const handleKeyDown = (event: KeyboardEvent): void => {
  const textarea = event.target as HTMLTextAreaElement;
  const { selectionStart, selectionEnd } = textarea;
  const currentText = props.value ?? '';

  // Helper function to safely execute text insertion while keeping Ctrl+Z history
  const insertTextWithHistory = (startPos: number, endPos: number, text: string, newCursorPos: number) => {
    textarea.focus();
    // 1. Highlight the target text area to be replaced
    textarea.setSelectionRange(startPos, endPos);

    // 2. Use browser command to insert text (this preserves Ctrl+Z)
    const success = document.execCommand('insertText', false, text);

    // 3. Fallback to manual string manipulation if execCommand fails
    if (!success) {
      handleUpdateValue(
        currentText.substring(0, startPos) +
        text +
        currentText.substring(endPos));
    }

    // 4. Update cursor position precisely
    setTimeout(() => {
      textarea.selectionStart = textarea.selectionEnd = newCursorPos;
    }, 0);
  };

  // --- 1. Shift + Tab: Insert Indentation ---
  if (event.key === 'Tab' && event.shiftKey) {
    event.preventDefault();
    const tabCharacter = '\t';
    insertTextWithHistory(selectionStart, selectionEnd, tabCharacter, selectionStart + tabCharacter.length);
    return;
  }

  // --- 2. Ctrl + Enter: Insert New Line BELOW Current Line ---
  if (event.key === 'Enter' && event.ctrlKey && !event.shiftKey) {
    event.preventDefault();

    const nextLineBreak = currentText.indexOf('\n', selectionStart);
    const insertPos = nextLineBreak === -1 ? currentText.length : nextLineBreak;

    // We insert "\n" at the end of current line, and move cursor to next line (insertPos + 1)
    insertTextWithHistory(insertPos, insertPos, '\n', insertPos + 1);
    return;
  }

  // --- 3. Shift + Enter: Insert New Line ABOVE Current Line ---
  if (event.key === 'Enter' && event.shiftKey && !event.ctrlKey) {
    event.preventDefault();

    const prevLineBreak = currentText.lastIndexOf('\n', selectionStart - 1);
    const insertPos = prevLineBreak === -1 ? 0 : prevLineBreak + 1;

    // We insert "\n" at the start of current line, and move cursor to the newly created line (insertPos)
    insertTextWithHistory(insertPos, insertPos, '\n', insertPos);
    return;
  }
}
const emit = defineEmits<{
  'update:value': [value: string]
}>()
function handleUpdateValue(value: string) {
  if (!props.readonly) {
    emit('update:value', value)
  }
}
</script>
<template>
  <div ref="editorWrapperRef" class="code-editor-container">
    <div ref="lineNumbersRef" class="line-numbers">
      <div v-for="line in totalLines" :key="line" class="line-number-item">
        {{ line }}
      </div>
    </div>
    <n-input type="textarea" :value="value" @update:value="handleUpdateValue" @keydown="handleKeyDown"
      :placeholder="$t('ui.configPlaceholder')" :disabled="disabled" :autosize="autosize"
      :input-props="{ style: { whiteSpace: 'pre', overflowX: 'auto' } }" autofocus show-count autocapitalize="off"
      :readonly="readonly" autocomplete="off" autocorrect="off" spellcheck="false" />
  </div>


</template>

<style scoped>
.n-input {
  height: 100%;
}

.code-editor-container {
  width: 100%;

  --editor-font-family: 'Fira Code', Consolas, Monaco, 'Courier New', Courier, monospace;
  --editor-font-size: 14px;
  --editor-line-height: 1.5;
  --editor-padding-top: 8px;
  --editor-padding-bottom: 8px;

  /* =============================================================
     核心修正：直接劫持並讀取 Naive UI 當前主題注入的 CSS 變數
     不論你在父層綁定的是亮色還是深色物件，這些變數的值都會被 Naive UI 自動更新
     ============================================================= */
  background-color: var(--n-color);
  /* 讀取當前主題的輸入框背景色 */
  border: 1px solid var(--n-border-color);
  /* 讀取當前主題的邊框顏色 */

  /* 行號區域的配色自動化調校 */
  /* 利用與主背景色微弱的混合，自動生成不論亮暗色都完美的行號背景 */
  --line-number-bg: var(--n-code-bg-color, rgba(0, 0, 0, 0.03));
  --line-number-color: var(--n-placeholder-color);
  /* 讀取當前主題的提示文字顏色作為行號色 */

  display: flex !important;
  border-radius: 3px;
  overflow: hidden;
  box-sizing: border-box;

  /* 呼吸燈動畫過渡，時間比照 Naive UI 官方規範 */
  transition: border-color 0.2s var(--n-cubic-bezier-ease-in-out),
    box-shadow 0.2s var(--n-cubic-bezier-ease-in-out),
    background-color 0.3s var(--n-cubic-bezier-ease-in-out);
}

/* =============================================================
   狀態控制：完全同步 Naive UI 官方的 Hover 與 Focus 呼吸燈效果
   ============================================================= */

/* Hover 狀態 */
.code-editor-container:hover {
  border-color: var(--n-border-color-hover);
}

/* Focus 狀態（當你在外層 div 加上 :class="{ 'is-focused': isFocused }" 時觸發） */
.code-editor-container.is-focused {
  border-color: var(--n-border-color-hover);
  box-shadow: var(--n-box-shadow-focus);
}

/* 移除 Naive UI 內建的二次重複邊框 */
:deep(.n-input) {
  border: none !important;
  background-color: transparent !important;
  --n-border: none !important;
  --n-border-hover: none !important;
  --n-border-focus: none !important;
  --n-box-shadow-focus: none !important;
}

/* 1. Left Side: Line Numbers Styling */
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
  /* Hide scrollbar */
  border-right: 1px solid var(--n-border-color);
  /* 共用 Naive UI 邊框變數 */
  opacity: 0.75;
  box-sizing: border-box;
}

.line-number-item {
  padding-right: 8px;
  height: calc(var(--editor-font-size) * var(--editor-line-height));
}

/* 2. Right Side: Force Naive UI Internal Textarea to Match Exactly */
:deep(.n-input .n-input__textarea-el) {
  /* 這裡不用手動寫死顏色，因為 Naive UI 內部的原生 textarea 已經由物件注入了正確的顏色 */
  font-family: var(--editor-font-family) !important;
  font-size: var(--editor-font-size) !important;
  line-height: var(--editor-line-height) !important;
  padding-top: var(--editor-padding-top) !important;
  padding-bottom: var(--editor-padding-bottom) !important;
  padding-left: 10px !important;
  padding-right: 10px !important;
}

/* 現代自訂滾動條：滑塊顏色同樣讀取自 Naive UI 的提示色，亮暗模式皆透明適中 */
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

/* Mirror padding for wrapper if Naive UI introduces extra offsets */
:deep(.n-input-wrapper) {
  padding: 0 !important;
}

/* 精準對齊 placeholder 的水平與垂直位置 */
:deep(.n-input .n-input__placeholder) {
  left: 10px !important;
  /* 與文字的 padding-left 一致 */
  top: 0px !important;
  /* 與文字的 padding-top 一致，自動跟隨變數 */
  transform: none !important;
  /* 拔掉 Naive UI 預設的垂直置中位移 */
}
</style>
