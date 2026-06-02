import { useEffect, useState, type SetStateAction } from 'react'
import { OneButtonModal } from './Modal'
import { IncidentEdit } from './IncidentEdit'
import type { Incident } from './types'

interface Column {
    key: string
    label: string
    render?: (value: any) => React.ReactNode
}

interface SortableTableProps {
    columns: Column[]
    data: any[]
    onDelete: (value: Incident[]) => void
    onError: (value: string | null) => void
    filter: string | null
    resolvedFilter: string | null
    setCurrentIncCount: React.Dispatch<SetStateAction<number | null>>
    searchKeyword: string
}

const SortableTable = ({columns, data, onDelete, onError, filter, resolvedFilter, setCurrentIncCount, searchKeyword}: SortableTableProps) => {
    const [sortKey, setSortKey] = useState<string>('id')
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc')
    const [currentHoveredRowID, setcurrentHoveredRowID] = useState<string | null>(null)

    const filterSearchKeyword = () => {
        if (searchKeyword === '') {
            return data
        } else {
            return data.filter(inc => inc.id.toLowerCase().includes(searchKeyword.toLowerCase()))
        }
    }

    const filterResolvedIncidents = () => {
        if (resolvedFilter === 'unresolved') {
            return filterSearchKeyword().filter(inc => inc.is_resolved === false)
        }
        if (resolvedFilter === 'resolved') {
            return filterSearchKeyword().filter(inc => inc.is_resolved === true)
        } else {
            return filterSearchKeyword()
        }
    }

    // Array of filtered incidents
    const filteredIncidents = () => {
        if (filter === null || filter === 'all') {
            return filterResolvedIncidents()
        } else {
            return filterResolvedIncidents().filter(inc => inc.severity === filter)
        }
    }

    // Only recalculate Incident count when the values in the square brackets change. This is the correct pattern.
    useEffect(() => {
        setCurrentIncCount(filteredIncidents().length)
    }, [filter, resolvedFilter, searchKeyword, data])

    const handleSort = async (colKey: string) => {
        // Do not sort "message" and "other" columns
        if (colKey === 'new_message' || colKey === 'other') {
            return
        }
        if (colKey === sortKey) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
        } else {
            setSortKey(colKey)
            setSortDirection('asc')
        }
    }

    const severityOrder: Record<string, number> = {
        critical: 0,
        high: 1,
        medium: 2,
        low: 3,
    }

    const sortedData = [...filteredIncidents()].sort((a, b) => {
        if (sortKey === 'severity') {
            if (severityOrder[a.severity] < severityOrder[b.severity]) return sortDirection === 'asc' ? -1 : 1
            if (severityOrder[a.severity] > severityOrder[b.severity]) return sortDirection === 'asc' ? 1 : -1
            return 0
        }
        if (a[sortKey] < b[sortKey]) return sortDirection === 'asc' ? -1 : 1
        if (a[sortKey] > b[sortKey]) return sortDirection === 'asc' ? 1 : -1
        return 0
    })
    
    if (sortedData.length === 0) {
        return (
            <div className='flex justify-center bg-base-100'>
                <span className='text-sm font-bold text-error/80 m-3'>Table is empty</span>
            </div>
        )
    }

    const deleteIncident = async (incidentID: string) => {
        try {
            const response = await fetch(`http://localhost:8080/incidents/${incidentID}`, {method: "DELETE"})
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`)
            }
            onDelete(filteredIncidents().filter(inc => inc.id !== incidentID))
        } catch (err) {
            onError(err instanceof Error ? err.message : "Unknown error")
        }
    }

    const renderDeleteButton = (incidentID: string) => {
        return (
                <OneButtonModal buttonText={
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-4">
                        <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                    </svg>
                } title={
                    <>
                        <span>Are you sure you want to delete </span>
                        <span className='bg-linear-to-r from-violet-600 to-blue-500 text-transparent bg-clip-text'>{incidentID}</span>
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
                        <button className='btn btn-sm btn-error absolute bottom-4 right-4' onClick={() => deleteIncident(incidentID)}>Delete</button>
                    </>
                } />
        )
    }

    const subtractDates = (date1: string, date2: string): number => {
        const time1 = new Date(date1).getTime()
        const time2 = new Date(date2).getTime()
        let timeDiff: number
        if (time1 > time2) {
            timeDiff = time1 - time2
        } else {
            timeDiff = time2 - time1
        }
        return Math.round(timeDiff / 1000)
    }

    const secondsToTime = (seconds: number): string => {
        const hours = Math.floor(seconds / 3600)
        seconds %= 3600
        const minutes = Math.floor(seconds / 60)
        if (hours === 0) {
            return `${minutes}m`
        } else {
            return `${hours}h${minutes}m`
        }
    }

    const UTCToBrowserTime = (time: string):string => {
        const d = new Date(time).toLocaleString('default', {
            month: 'short',
            day: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            hour12: false,
            minute: '2-digit'
        })
        return d
    }

    const setupMessage = (row: any) => {
        if (row.is_resolved) {
            const timeDiff = subtractDates(row.started_at, row.resolved_at)
            const time = secondsToTime(timeDiff)
            const browserTime = UTCToBrowserTime(row.resolved_at)
            return (
                <div className='flex flex-col justify-center items-end h-8'>
                    <span>Resolved in {time}</span>
                    <span className='text-xs text-base-content/70'>Ended: {browserTime}</span>
                </div>
            )
        } else {
            const browserTime = UTCToBrowserTime(row.started_at)
            return (
                <div className='flex flex-col justify-center items-end h-8'>
                    <span>Pending ...</span>
                    <span className='text-xs text-base-content/70'>Started: {browserTime}</span>
                </div>
            )
        }
    }

    const renderTableButtons = (row: any) => {
        if (row.id === currentHoveredRowID) {
            return (
                <div className='flex justify-end items-center h-8'>
                    <IncidentEdit incidentID={row.id} setError={onError} />
                    {renderDeleteButton(row.id)}
                </div>
            )
        } else {
            return setupMessage(row)
        }
    }

    const renderIsResolved = (colVal: any) => {
        if (colVal) {
            return '✓'
        } else {
            return '✗'
        }
    }

    return (
        <table className='table table-md table-pin-rows table-pin-cols overflow-hidden'>
            <thead>
                <tr>
                { columns.map(col => (
                    <th key={col.key} onClick={() => handleSort(col.key)} className={`cursor-pointer select-none ${col.key === 'new_message' && 'w-50'} ${col.key === 'severity' && 'text-center'}`}>
                        {col.label}
                        {sortKey === col.key && (sortDirection === 'asc' ? ' ↑' : ' ↓')}
                    </th>
                ))}
                </tr>
            </thead>
            <tbody>
                {sortedData.map((row, index) => (
                <tr className={`${row.id === currentHoveredRowID && 'bg-base-200'}`} key={row.id || index}>
                    {columns.map(col => (
                        <td key={col.key} onMouseEnter={() => setcurrentHoveredRowID(row.id)} onMouseLeave={() => setcurrentHoveredRowID(null)}>
                            {col.render ? col.render(row[col.key]) : row[col.key]}
                            {col.key === 'resolved_at' && row[col.key] === null && 'Unresolved'}
                            {col.key === 'new_message' && renderTableButtons(row)}
                            {/* {col.key === 'new_message' && setupMessage(row)} */}
                            {col.key === 'is_resolved' && renderIsResolved(row[col.key])}
                        </td>
                    ))}
                </tr>
                ))}
            </tbody>
        </table>
    )
}

export default SortableTable