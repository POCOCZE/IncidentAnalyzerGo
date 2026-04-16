import {useState, useEffect } from 'react'

interface Incident {
    id: string
    title: string
    severity: string
    service_name: string
    started_at: string
    resolved_at: string | null
}

const IncidentList = () => {
    const [incidents, setIncidents] = useState<Incident[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string | null>(null)
    // Variable for severity dropdown menu
    const [filter, setFilter] = useState<string>("all")

    useEffect(() => {
        const fetchIncidents = async () => {
            try {
                const response = await fetch("http://localhost:8080/incidents")
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}`)
                }
                const data: Incident[] = await response.json()
                setIncidents(data)
                setLoading(false)
            } catch (err) {
                setError(err instanceof Error ? err.message : "Unknown error")
                setLoading(false)
            }
        }
        fetchIncidents()
    }, [])

    if (loading) {
        return <p className='text-orange-300'>Loading...</p>
    }

    if (error) {
        return <p className='text-red-600'>Error: {error}</p>
    }

    const severityColor: Record<string, string> = {
        critical: 'text-red-600/70',
        high: 'text-amber-600',
        medium: 'text-amber-400/80',
    };

    // Array of filtered incidents
    const filteredIncidents = filter === "all"
        ? incidents
        : incidents.filter(inc => inc.severity === filter)

    return (
        <div className='rounded-xl backdrop-blur-md px-3'>
            <h2 className='text-2xl text-center font-semibold p-1'>Incident list</h2>
            {/* key= forces to show the current value, onChange= when user picks something, state needs to be updated (setSevFilter), this causes rerender, which updates the dropdown to match */}
            <div className='flex justify-end items-center'>
                <p className='text-xs text-base-content/50 mr-1 w-12 mb-1'>Filter by severity</p>
                <select value={filter} onChange={(e) => setFilter(e.target.value)} className='select select-sm w-23 text-base-content rounded-md mb-1'>
                    <option value="all">All</option>
                    <option value="critical">Critical</option>
                    <option value="high">High</option>
                    <option value="medium">Medium</option>
                    <option value="low">Low</option>
                </select>
            </div>
            <div className='overflow-x-auto rounded-md bg-base-200 h-[calc(100vh-14vh)]'>
                <table className='table table-md table-pin-rows table-pin-cols'>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Title</th>
                            <th>Severity</th>
                            <th>Service</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        {filteredIncidents.map(incident => (
                            <tr key={incident.id}>
                                <td>{incident.id}</td>
                                <td>{incident.title}</td>
                                <td className={severityColor[incident.severity]}>{incident.severity}</td>
                                <td>{incident.service_name}</td>
                                <td>{incident.resolved_at ? 'Resolved' : 'Pending'}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            <div className='flex text-xs text-base-content/60 justify-center items-center m-1'>
                <p>Showing {filteredIncidents.length} of {incidents.length}</p>
            </div>
        </div> 
    )
}

export default IncidentList