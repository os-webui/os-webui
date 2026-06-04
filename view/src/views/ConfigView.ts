import { computed, ref } from "vue"
import type { PluginInfo } from "./HomeView"
import type { Save20Regular } from "@vicons/fluent"
export interface Data {
  info: PluginInfo
  data: string
}
export interface Props {
  signal: AbortSignal
  plugin: string
  data: Data
}
export function createConfigView(props: Props) {
  const textValue = ref(props.data.data)
  const handleKeyDown = (event: KeyboardEvent): void => {
    const textarea = event.target as HTMLTextAreaElement;
    const { selectionStart, selectionEnd } = textarea;
    const currentText = textValue.value;

    // Helper function to safely execute text insertion while keeping Ctrl+Z history
    const insertTextWithHistory = (startPos: number, endPos: number, text: string, newCursorPos: number) => {
      textarea.focus();
      // 1. Highlight the target text area to be replaced
      textarea.setSelectionRange(startPos, endPos);

      // 2. Use browser command to insert text (this preserves Ctrl+Z)
      const success = document.execCommand('insertText', false, text);

      // 3. Fallback to manual string manipulation if execCommand fails
      if (!success) {
        textValue.value =
          currentText.substring(0, startPos) +
          text +
          currentText.substring(endPos);
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
  const disabled = ref(false)
  const disabledClear = computed(() => disabled.value || textValue.value === '')
  const disabledReset = computed(() => disabled.value || textValue.value === props.data.data)
  return {
    textValue,
    handleKeyDown,
    disabled, disabledClear, disabledReset,
    actions: {
      clear() {
        if (disabled.value) {
          return
        }
        textValue.value = ''
      },
      reset() {
        if (disabled.value) {
          return
        }
        textValue.value = props.data.data
      },
      save() {
        if (disabled.value) {
          return
        }
        console.log('save')
      },
    },
  }
}
