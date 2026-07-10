import { apiClient } from './client'

export interface ActivityCenterConfig {
  enabled: boolean
}

export interface ActivityCenterItem {
  id: number
  slug: string
  title: string
  subtitle: string
  description: string
  icon: string
  cover_url: string
  route_path: string
  external_url: string
  action_label: string
  activity_type: string
  enabled: boolean
  sort_order: number
  start_at?: string
  end_at?: string
  metadata: Record<string, unknown>
  created_by?: number
  created_at: string
  updated_at: string
}

export interface ActivityCenterItemInput {
  slug: string
  title: string
  subtitle: string
  description: string
  icon: string
  cover_url: string
  route_path: string
  external_url: string
  action_label: string
  activity_type: string
  enabled: boolean
  sort_order: number
  start_at?: string | null
  end_at?: string | null
  metadata?: Record<string, unknown>
}

export interface ActivityCenterListResponse {
  enabled: boolean
  items: ActivityCenterItem[]
}

export function listActivityCenter() {
  return apiClient.get<ActivityCenterListResponse>('/activity-center').then(res => res.data)
}

export function getActivityCenterConfig() {
  return apiClient.get<ActivityCenterConfig>('/admin/activity-center/config').then(res => res.data)
}

export function updateActivityCenterConfig(payload: ActivityCenterConfig) {
  return apiClient.put<ActivityCenterConfig>('/admin/activity-center/config', payload).then(res => res.data)
}

export function listAdminActivityCenterItems() {
  return apiClient.get<ActivityCenterItem[]>('/admin/activity-center/items').then(res => res.data)
}

export function createActivityCenterItem(payload: ActivityCenterItemInput) {
  return apiClient.post<ActivityCenterItem>('/admin/activity-center/items', payload).then(res => res.data)
}

export function updateActivityCenterItem(id: number, payload: ActivityCenterItemInput) {
  return apiClient.put<ActivityCenterItem>(`/admin/activity-center/items/${id}`, payload).then(res => res.data)
}

export function deleteActivityCenterItem(id: number) {
  return apiClient.delete<{ deleted: boolean }>(`/admin/activity-center/items/${id}`).then(res => res.data)
}
