import { useEffect, useState, type SetStateAction } from 'react'
import { OneButtonModal } from './Modal'
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
    resolvedFilter: "all" | "unresolved"
    setCurrentIncCount: React.Dispatch<SetStateAction<number | null>>
    searchKeyword: string
}

const SortableTable = ({columns, data, onDelete, onError, filter, resolvedFilter, setCurrentIncCount, searchKeyword}: SortableTableProps) => {
    const [sortKey, setSortKey] = useState<string>('id')
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc')
    // Todo: Variable to set table background when cursor hovers over the row. Did not work using `onMouseEnter` event handler - targets whole table, not row.
    // const [hoverRowBg, setHoverRowBg] = useState<string>('')
    // create a copy of data that gets then filtered after incident gets deleted: so after each deletion no /incidents call is needed.

    const filterSearchKeyword = () => {
        if (searchKeyword === '') {
            return data
        } else {
            return data.filter(inc => inc.id.toLowerCase().includes(searchKeyword.toLowerCase()))
        }
    }

    const filterResolvedIncidents = () => {
        if (resolvedFilter === "all") {
            return filterSearchKeyword()
        } else {
            return filterSearchKeyword().filter(inc => inc.is_resolved === false)
        }
    }

    // Array of filtered incidents
    const filteredIncidents = () => {
        if (filter === null) {
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

    return (
        <table className='table table-md table-pin-rows table-pin-cols overflow-hidden'>
            <thead>
                <tr>
                { columns.map(col => (
                    <th key={col.key} onClick={() => handleSort(col.key)} className='cursor-pointer select-none'>
                        {col.label}
                        {sortKey === col.key && (sortDirection === 'asc' ? ' ↑' : ' ↓')}
                    </th>
                ))}
                </tr>
            </thead>
            <tbody>
                {sortedData.map((row, index) => (
                <tr key={row.id || index}>
                    {columns.map(col => (
                        <td key={col.key}>
                            {col.render ? col.render(row[col.key]) : row[col.key]}
                            {col.key === 'resolved_at' && row[col.key] === null && 'Unresolved'}
                            {col.key === 'is_resolved' && row[col.key] && '✓'}
                            {col.key === 'is_resolved' && row[col.key] === false && '✗'}
                            {col.key === 'delete' && renderDeleteButton(row.id) }
                        </td>
                    ))}
                </tr>
                ))}
            </tbody>
        </table>
    )
}

export default SortableTable