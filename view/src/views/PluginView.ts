import { type PluginInfo } from './HomeView'
export interface FeatureInfo {
  id: string
  name: string
  description: string
}
export interface PluginFeatures {
  info: PluginInfo
  features: FeatureInfo[]
}
export interface Props {
  plugin: PluginFeatures
  abort: AbortSignal
}
