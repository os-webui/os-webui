export interface PluginInfo {
  id: string
  version: string
  author: string
  name: string
  platform: string[]
  description: string
}
export interface Props {
  items?: PluginInfo[]
  signal: AbortSignal
}
