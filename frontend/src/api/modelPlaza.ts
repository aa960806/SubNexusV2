/**
 * Model Plaza API (display-only pricing showcase).
 *
 * 纯展示功能：按管理员手动配置的「分组 -> 模型 -> 价格」渲染。
 * 后端开关关闭时该接口返回 404（功能等同不存在）。
 */

import { apiClient } from './client'

/** 单个模型的四类价格（每百万 Token；均为管理员自填字符串）。 */
export interface ModelPlazaPrice {
  input: string
  output: string
  cache_read: string
  cache_write: string
}

/** 模型广场中的一行模型。 */
export interface ModelPlazaModel {
  name: string
  price: ModelPlazaPrice
  note?: string
}

/** 展示分组（与项目真实分组解耦，由管理员手动配置）。 */
export interface ModelPlazaGroup {
  name: string
  description?: string
  /** 组级默认货币符号，如 "$" / "¥"。 */
  currency?: string
  models: ModelPlazaModel[]
}

/** 固定的一级模型类型 key（展示名由前端 i18n 决定）。 */
export type ModelPlazaCategoryKey = 'claude' | 'openai' | 'gemini' | 'domestic' | 'image'

/** 固定类型 + 其下配置的分组。 */
export interface ModelPlazaCategory {
  key: ModelPlazaCategoryKey
  groups: ModelPlazaGroup[]
}

/** 模型广场展示配置（不含开关；开关来自 public settings 的 model_plaza_enabled）。 */
export interface ModelPlazaConfig {
  categories: ModelPlazaCategory[]
}

export interface ModelPlazaAdminConfig {
  enabled: boolean
  config: ModelPlazaConfig
}

/** 固定一级类型顺序（与后端 ModelPlazaCategoryKeys 对齐）。 */
export const MODEL_PLAZA_CATEGORY_KEYS: ModelPlazaCategoryKey[] = [
  'claude',
  'openai',
  'gemini',
  'domestic',
  'image',
]

/** 获取模型广场展示数据（用户端，需登录）。开关关闭时后端 404。 */
export async function getModelPlaza(options?: { signal?: AbortSignal }): Promise<ModelPlazaConfig> {
  const { data } = await apiClient.get<ModelPlazaConfig>('/settings/model-plaza', {
    signal: options?.signal,
  })
  return data
}

export async function getModelPlazaAdminConfig(): Promise<ModelPlazaAdminConfig> {
  const { data } = await apiClient.get<ModelPlazaAdminConfig>('/admin/settings/model-plaza')
  return data
}

export async function updateModelPlazaAdminConfig(
  payload: ModelPlazaAdminConfig,
): Promise<ModelPlazaAdminConfig> {
  const { data } = await apiClient.put<ModelPlazaAdminConfig>('/admin/settings/model-plaza', payload)
  return data
}

export const modelPlazaAPI = {
  getModelPlaza,
  getModelPlazaAdminConfig,
  updateModelPlazaAdminConfig,
}

export default modelPlazaAPI
