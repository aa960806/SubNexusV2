import { apiClient } from './client'

export interface BroadcastConfig {
  enabled: boolean
}

export interface ActivityBroadcast {
  id: number
  title: string
  content: string
  source: string
  enabled: boolean
  priority: number
  start_at?: string | null
  end_at?: string | null
  created_at: string
  updated_at: string
}

export interface ActivityBroadcastInput {
  title: string
  content: string
  enabled: boolean
  priority: number
  start_at?: string | null
  end_at?: string | null
}

export async function listActiveBroadcasts(limit = 12): Promise<ActivityBroadcast[]> {
  const { data } = await apiClient.get<ActivityBroadcast[]>('/activity/broadcasts', { params: { limit } })
  return data
}

export async function getBroadcastConfig(): Promise<BroadcastConfig> {
  const { data } = await apiClient.get<BroadcastConfig>('/admin/broadcasts/config')
  return data
}

export async function updateBroadcastConfig(config: BroadcastConfig): Promise<BroadcastConfig> {
  const { data } = await apiClient.put<BroadcastConfig>('/admin/broadcasts/config', config)
  return data
}

export async function listBroadcasts(): Promise<ActivityBroadcast[]> {
  const { data } = await apiClient.get<ActivityBroadcast[]>('/admin/broadcasts')
  return data
}

export async function createBroadcast(input: ActivityBroadcastInput): Promise<ActivityBroadcast> {
  const { data } = await apiClient.post<ActivityBroadcast>('/admin/broadcasts', input)
  return data
}

export async function updateBroadcast(id: number, input: ActivityBroadcastInput): Promise<ActivityBroadcast> {
  const { data } = await apiClient.put<ActivityBroadcast>(`/admin/broadcasts/${id}`, input)
  return data
}

export async function deleteBroadcast(id: number): Promise<void> {
  await apiClient.delete(`/admin/broadcasts/${id}`)
}

export async function cleanupExpiredBroadcasts(retentionDays = 7): Promise<number> {
  const { data } = await apiClient.post<{ deleted: number }>('/admin/broadcasts/cleanup', {
    retention_days: retentionDays,
  })
  return data.deleted
}
