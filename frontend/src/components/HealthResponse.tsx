import { useState, useEffect } from "react"

interface HealthResponse {
    status: string
}

const HealthStatus = () => {
    const [status, setStatus] = useState<string | null>(null)
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string | null >(null)

    useEffect(() => {
        const fetchHealth = async () => {
            try {
                const response = await fetch("http://localhost:8080/healthz")
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}`)
                }
                const data: HealthResponse = await response.json()
                setStatus(data.status)
                setLoading(false)
            } catch (err) {
                setError(err instanceof Error ? err.message : "Unknown error")
                setLoading(false)
            }
        }
        fetchHealth()
    }, []) // Empty array - run only once

    if (loading) {
        return (
            <span className="loading loading-spinner loading-xs rounded">Loading</span>
        )
    }

    if (error) {
        // return <ToastNotification duration={10000} message={`Error ${error}`} toastLevel='alert-error' toastPos='toast-top toast-right'/>
        return (
            <div className="flex items-center rounded-xl p-1 bg-red-100">
                <div aria-label="error" className="status status-error animate-pulse"></div>
                <span className="text-sm text-base-content pl-1">Status</span>
            </div>
        )
    }

    return (
        <div className="flex items-center bg-base-200 rounded-xl p-1 shadow">
            <div aria-label="success" className="status status-success animate-none"></div>
            <span className="text-xs text-base-content pl-1">Status</span>
        </div>
    )
}

export default HealthStatus