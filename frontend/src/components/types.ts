export interface Incident {
    id: string
    title: string
    severity: string
    service_name: string
    started_at: string
    resolved_at: string | null
    message: string
    is_resolved: boolean
}