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
        // return <p className="text-orange-300">Loading...</p>
        return (
            <span className="loading loading-spinner loading-xs">Loading</span>
        )
    }

    if (error) {
        return <p className="text-red-600">Error: {error}</p>
    }

    return <p className="text-green-500">Status: {status}</p>
}

export default HealthStatus