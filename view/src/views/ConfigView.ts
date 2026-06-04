import { computed, ref } from "vue"
import type { PluginInfo } from "./HomeView"
import { DefaultHttpClient } from "@/internal/api"
import { useMessage } from "naive-ui"
import { useI18n } from "vue-i18n"
import { ResponseType } from "@/internal/core/http_client"
import { errorString } from "@/internal/core/strings"
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
  const originValue = ref(props.data.data)
  const textValue = ref(props.data.data)
  const message = useMessage()
  const i18n = useI18n()

  const disabled = ref(false)
  const disabledClear = computed(() => disabled.value || textValue.value === '')
  const disabledReset = computed(() => disabled.value || textValue.value === originValue.value)
  return {
    textValue,
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
        textValue.value = originValue.value
      },
      async save() {
        if (disabled.value) {
          return
        }
        disabled.value = true
        const value = textValue.value
        try {
          await DefaultHttpClient.post(`/api/v1/plugins/${encodeURIComponent(props.plugin)}/config`, {
            body: JSON.stringify({
              data: value,
            }),
            signal: props.signal,
          }, {
            type: ResponseType.text,
          })
          originValue.value = value
        } catch (e) {
          console.log(e)
          message.error(errorString(e))
          return
        } finally {
          disabled.value = false
        }
        message.success(i18n.t('ui.saveSuccess'))
      },
    },
  }
}
