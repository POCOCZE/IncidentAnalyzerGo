import {useState, useEffect } from 'react'
import { OneButtonModal } from './Modal'

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
    
    const deleteIncident = async (incidentID: string) => {
        try {
            const response = await fetch(`http://localhost:8080/incidents/${incidentID}`, {method: "DELETE"})
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`)
            }
            setIncidents(incidents.filter(inc => inc.id !== incidentID))
        } catch (err) {
            setError(err instanceof Error ? err.message : "Unknown error")
        }
    }

    if (loading) {
        return (
            <div className="flex w-52 justify-center items-center flex-col gap-4">
                <div className="skeleton h-4 w-full"></div>
                <div className="skeleton h-4 w-28"></div>
                <div className="skeleton h-4 w-40"></div>
                <div className="skeleton h-4 w-34"></div>
            </div>
        )
    }

    if (error) {
        return <p className='text-red-600'>Error: {error}</p>
    }

    const severityColor: Record<string, string> = {
        critical: 'badge bg-red-600/70 text-base-content',
        high: 'badge bg-amber-600/60 text-base-content',
        medium: 'badge bg-amber-400/60 text-base-content',
        low: 'badge bg-amber-400/30 text-base-content',
    };

    // Array of filtered incidents
    const filteredIncidents = filter === "all"
        ? incidents
        : incidents.filter(inc => inc.severity === filter)

    return (
        <div className='flex flex-col px-3 pt-3 justify-center grow rounded-xl h-[97vh] m-4 bg-base-300'>
            <span className='text-2xl text-center font-semibold'>Incident list</span>
            <div className='flex justify-end items-center'>
                <p className='text-xs text-base-content/50 mr-1 w-12 mb-2'>Filter by severity</p>
                <select value={filter} onChange={(e) => setFilter(e.target.value)} className='select select-sm w-23 text-base-content rounded-xl mb-1'>
                    <option value="all">All</option>
                    <option value="critical">Critical</option>
                    <option value="high">High</option>
                    <option value="medium">Medium</option>
                    <option value="low">Low</option>
                </select>
            </div>
            <div className='rounded-xl bg-base-200 overflow-auto lg:h-[89.5vh] h-[84.5vh]'>
                <table className='table table-md table-pin-rows table-pin-cols min-w-full'>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Title</th>
                            <th className='lg:flex lg:justify-center-safe'>Severity</th>
                            <th>Service</th>
                            <th>Status</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>
                        {filteredIncidents.map(incident => (
                            <tr key={incident.id}>
                                <td>{incident.id}</td>
                                <td>{incident.title}</td>
                                <td className='lg:flex lg:justify-center-safe'>
                                    <div className={severityColor[incident.severity]}>
                                        <span>{incident.severity}</span>
                                    </div>
                                </td>
                                <td>{incident.service_name}</td>
                                <td>{incident.resolved_at ? 'Resolved' : 'Pending'}</td>
                                <td>
                                    <OneButtonModal buttonText={
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-4">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                                        </svg>
                                    } title={
                                        <>
                                            <span>Are you sure you want to delete </span>
                                            <span className='bg-linear-to-r from-violet-600 to-blue-500 text-transparent bg-clip-text'>{incident.id}</span>
                                            <span> ?</span>
                                        </>
                                    } description={
                                        <>
                                            This action is irreversible.
                                            <br/>
                                            Incident will be removed forever after its deleted!
                                        </>
                                        } closeButtonText={
                                            <>
                                                <button className='btn btn-sm btn-error absolute bottom-4 right-4' onClick={() => deleteIncident(incident.id)}>Delete</button>
                                            </>
                                        } />
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            <span className='text-xs text-base-content/60 text-center'>Showing {filteredIncidents.length} of {incidents.length}</span>
        </div> 
    )
}

export default IncidentList