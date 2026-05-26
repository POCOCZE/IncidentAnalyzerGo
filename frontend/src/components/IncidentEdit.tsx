import { useEffect, useState } from 'react'
// import type { Incident } from './types'
import { IncidentEditModal } from './Modal'

interface IncidentEditProps {
    incidentID: string
    setError: (value: string | null) => void
}

export const IncidentEdit = ({incidentID, setError}: IncidentEditProps) => {
    // const [loading, setLoading] = useState<boolean>(true)
    // const [incident, setIncident] = useState<Incident>()

    // useEffect(() => {
    //     const fetchIncidentByID = async () => {
    //         try {
    //             const response = await fetch(`http://localhost:8080/incidents/${incidentID}`)
    //             if (!response?.ok) {
    //                 const errorData = await response?.json()
    //                 console.log('Error response', errorData)
    //                 throw new Error(`- ${errorData.error}`)
    //             }
    //             const data: Incident = await response.json()
    //             setIncident(data)
    //             setLoading(false)
    //         } catch (err) {
    //             setError(err instanceof Error ? err.message : 'Unknown error')
    //             setLoading(false)
    //         }
    //     }
    //     fetchIncidentByID()
    // }, [])

    // if (loading) {
    //     return (
    //         <div className="flex w-4 justify-center items-center flex-col gap-4">
    //             <div className="skeleton h-4 w-full"></div>
    //         </div>
    //     )
    // }

    return (
        <>
        <IncidentEditModal editIcon={
            <>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="size-4">
                <path d="M21.731 2.269a2.625 2.625 0 0 0-3.712 0l-1.157 1.157 3.712 3.712 1.157-1.157a2.625 2.625 0 0 0 0-3.712ZM19.513 8.199l-3.712-3.712-12.15 12.15a5.25 5.25 0 0 0-1.32 2.214l-.8 2.685a.75.75 0 0 0 .933.933l2.685-.8a5.25 5.25 0 0 0 2.214-1.32L19.513 8.2Z" />
                </svg>
            </>
        } incidentID={incidentID} />
        </>
    )
}