import { apiClient } from '../client'

export interface BackupS3Config {
  endpoint: string
  region: string
  bucket: string
  access_key_id: string
  secret_access_key?: string
  prefix: string
  force_path_style: boolean
}

export interface BackupScheduleConfig {
  enabled: boolean
  cron_expr: string
  retain_days: number
  retain_count: number
}

export interface BackupWebDAVConfig {
  url: string
  username: string
  password?: string
  prefix: string
}

export interface BackupRecord {
  id: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  backup_type: string
  file_name: string
  storage_type?: string
  s3_key: string
  size_bytes: number
  triggered_by: string
  error_message?: string
  started_at: string
  finished_at?: string
  expires_at?: string
  progress?: string
  restore_status?: string
  restore_error?: string
  restored_at?: string
}

export interface CreateBackupRequest {
  expire_days?: number
}

export interface TestS3Response {
  ok: boolean
  message: string
}

// S3 Config
export async function getS3Config(): Promise<BackupS3Config> {
  const { data } = await apiClient.get<BackupS3Config>('/admin/backups/s3-config')
  return data
}

export async function updateS3Config(config: BackupS3Config): Promise<BackupS3Config> {
  const { data } = await apiClient.put<BackupS3Config>('/admin/backups/s3-config', config)
  return data
}

export async function testS3Connection(config: BackupS3Config): Promise<TestS3Response> {
  const { data } = await apiClient.post<TestS3Response>('/admin/backups/s3-config/test', config)
  return data
}

// WebDAV Config
export async function getWebDAVConfig(): Promise<BackupWebDAVConfig> {
  const { data } = await apiClient.get<BackupWebDAVConfig>('/admin/backups/webdav-config')
  return data
}

export async function updateWebDAVConfig(config: BackupWebDAVConfig): Promise<BackupWebDAVConfig> {
  const { data } = await apiClient.put<BackupWebDAVConfig>('/admin/backups/webdav-config', config)
  return data
}

export async function testWebDAVConnection(config: BackupWebDAVConfig): Promise<TestS3Response> {
  const { data } = await apiClient.post<TestS3Response>('/admin/backups/webdav-config/test', config)
  return data
}

// Storage Type
export async function getStorageType(): Promise<{ storage_type: string }> {
  const { data } = await apiClient.get<{ storage_type: string }>('/admin/backups/storage-type')
  return data
}

export async function updateStorageType(storageType: string): Promise<{ storage_type: string }> {
  const { data } = await apiClient.put<{ storage_type: string }>('/admin/backups/storage-type', { storage_type: storageType })
  return data
}

// Schedule
export async function getSchedule(): Promise<BackupScheduleConfig> {
  const { data } = await apiClient.get<BackupScheduleConfig>('/admin/backups/schedule')
  return data
}

export async function updateSchedule(config: BackupScheduleConfig): Promise<BackupScheduleConfig> {
  const { data } = await apiClient.put<BackupScheduleConfig>('/admin/backups/schedule', config)
  return data
}

// Backup operations
export async function createBackup(req?: CreateBackupRequest): Promise<BackupRecord> {
  const { data } = await apiClient.post<BackupRecord>('/admin/backups', req || {})
  return data
}

export async function listBackups(): Promise<{ items: BackupRecord[] }> {
  const { data } = await apiClient.get<{ items: BackupRecord[] }>('/admin/backups')
  return data
}

export async function getBackup(id: string): Promise<BackupRecord> {
  const { data } = await apiClient.get<BackupRecord>(`/admin/backups/${id}`)
  return data
}

export async function deleteBackup(id: string): Promise<void> {
  await apiClient.delete(`/admin/backups/${id}`)
}

export async function getDownloadURL(id: string): Promise<{ url: string }> {
  const { data } = await apiClient.get<{ url: string }>(`/admin/backups/${id}/download-url`)
  return data
}

// Restore
export async function restoreBackup(id: string, password: string): Promise<BackupRecord> {
  const { data } = await apiClient.post<BackupRecord>(`/admin/backups/${id}/restore`, { password })
  return data
}

export const backupAPI = {
  getS3Config,
  updateS3Config,
  testS3Connection,
  getWebDAVConfig,
  updateWebDAVConfig,
  testWebDAVConnection,
  getStorageType,
  updateStorageType,
  getSchedule,
  updateSchedule,
  createBackup,
  listBackups,
  getBackup,
  deleteBackup,
  getDownloadURL,
  restoreBackup,
}

export default backupAPI
