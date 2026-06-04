import {useState, useEffect, createContext, useContext, type SetStateAction } from 'react'
import SortableTable from './SortableTable'
import { IncidentReportProvider, IncidentReportSidebar } from './IncidentReport'
import type { Incident } from './types'

interface IncidentContextType {
    incidents: Incident[]
    setIncidents: React.Dispatch<SetStateAction<Incident[]>>
    currentIncCount: number | null
    setCurrentIncCount: React.Dispatch<SetStateAction<number | null>>
    filter: string | null
    setFilter: React.Dispatch<SetStateAction<string | null>>
    resolvedFilter: string | null
    setResolvedFilter: React.Dispatch<SetStateAction<string | null>>
    loading: boolean
    error: string | null
    setError: React.Dispatch<SetStateAction<string | null>>
}

const IncidentListContext = createContext<IncidentContextType | undefined>(undefined)

export const IncidentListProvider = ({ children }: { children: React.ReactNode }) => {
    const [incidents, setIncidents] = useState<Incident[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string | null>(null)
    // Severity dropdown
    const [filter, setFilter] = useState<string | null>(null)
    const [resolvedFilter, setResolvedFilter] = useState<string | null>(null)
    const [currentIncCount, setCurrentIncCount] = useState<number | null>(null)

    
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

    return (
        <IncidentListContext.Provider value={{ incidents, currentIncCount, setCurrentIncCount, setIncidents, filter, setFilter, resolvedFilter, setResolvedFilter, loading, error, setError }}>
            {children}
        </IncidentListContext.Provider>
    )
}

const useIncidentList = () => {
    const context = useContext(IncidentListContext)
    if (!context) {
        throw new Error("useIncidentList must be within IncidentListProvider")
    }
    return context
}

export const ExportIncidents = () => {
    const incidents = useIncidentList()
    const wrapped = { incidents: incidents.incidents}

    const prepareExport = () => {
        const blob = new Blob([JSON.stringify(wrapped, null, 2)], { type: "application/json" } )
        return URL.createObjectURL(blob)
    }

    return (
        <div className="flex flex-col items-center m-2">
            <span className="text font-semibold">Export incidents</span>
            <button className="btn btn-sm btn-neutral border-base-content mt-4 m-0.5 px-5.5" onClick={() => {
                const url = prepareExport()
                const link = document.createElement('a')
                const time = new Date()
                link.href = url
                link.download = `incidents-${time.toLocaleDateString()}-${time.toLocaleTimeString()}.json`
                link.click()
                URL.revokeObjectURL(url)
            }}>
                Export
            </button>
        </div>
    )
}

export const IncidentListCenter = () => {
    const { incidents, setIncidents, setCurrentIncCount, filter, resolvedFilter, loading, error, setError } = useIncidentList()
    const [searchBg, setSearchBg] = useState<string>('bg-base-300')
    const [searchKeyword, setSearchKeyword] = useState<string>('')
    
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

    const emptyIncidentsGuard = () => {
        if (!incidents) {
            return (
                <div className='flex justify-center bg-base-100'>
                    <span className='text-sm font-bold text-error/80 m-3'>Table is empty</span>
                </div>
            )
        }
    }

    const severityColor: Record<string, string> = {
        critical: 'badge badge-sm border-base-content dark:border-black bg-error/80 text-base-content dark:text-black font-semibold',
        high: 'badge badge-sm border-base-content dark:border-black bg-warning text-base-content dark:text-black font-semibold',
        medium: 'badge badge-sm border-base-content dark:border-black bg-orange-400 text-base-content dark:text-black font-semibold',
        low: 'badge badge-sm border-base-content dark:border-black bg-yellow-500 text-base-content dark:text-black font-semibold',
    };

    const severityLabel: Record<string, string> = {
        critical: 'Critical',
        high: 'High',
        medium: 'Medium',
        low: 'Low',
    }

    const columns = [
        { key: 'is_resolved', label: ''},
        { key: 'id', label: 'ID',
            render: (value: string) => <span className='font-semibold'>{value}</span>},
        { key: 'title', label: 'Title'},
        { key: 'severity', label: 'Severity',
            render: (value: string) => <div className='flex justify-center'><span className={`text-xs ${severityColor[value]}`}>{severityLabel[value]}</span></div>},
        { key: 'service_name', label: 'Service name'},
        { key: 'new_message', label: ''}
        // { key: 'other', label: ''}
    ]

    const searchIncidentsField = () => {
        return (
            <label className={`input rounded-lg shadow-2xl border border-base-content/15 w-60 ${searchBg}`} onMouseEnter={() => setSearchBg('bg-base-100')} onMouseLeave={() => setSearchBg('bg-base-300')}>
                <svg className="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g strokeLinejoin="round" strokeLinecap="round" strokeWidth="2.5" fill="none" stroke="currentColor"><circle cx="11" cy="11" r="8"></circle><path d="m21 21-4.3-4.3"></path></g></svg>
                <input type="search" placeholder="Search incident IDs" value={searchKeyword} onChange={(e) => setSearchKeyword(e.currentTarget.value)} />
            </label>
        )
    }

    return (
        <div className='flex flex-col px-3 pt-3 grow rounded-lg h-[97vh] m-4 bg-base-200 shadow border border-base-content/15'>
            <span className='flex text-3xl font-bold justify-center bg-linear-to-r from bg-orange-500 to-yellow-500 bg-clip-text text-transparent'>Incident list</span>
            <div className='flex justify-end m-1'>
                {searchIncidentsField()}
            </div>
            <div className='rounded-box shadow-md bg-base-300 m-1 mb-4 border border-base-content/15 overflow-auto max-w-[80vw]'>
                {incidents ?
                <SortableTable columns={columns} data={incidents} onDelete={setIncidents} onError={setError} filter={filter} resolvedFilter={resolvedFilter} setCurrentIncCount={setCurrentIncCount} searchKeyword={searchKeyword}/>
                : emptyIncidentsGuard()}
            </div>
        </div> 
    )
}

export const IncidentListSidebar = () => {
    const { incidents, currentIncCount, setFilter, setResolvedFilter, loading } = useIncidentList()

    if (loading) {
        return (
            <div className="flex w-30 justify-center items-center flex-col gap-4">
                <div className="skeleton h-4 w-full"></div>
                <div className="skeleton h-4 w-20"></div>
                <div className="skeleton h-4 w-full"></div>
                <div className="skeleton h-4 w-14"></div>
            </div>
        )
    }

    const severityFilterBtnClasses = "btn btn-xs btn-neutral checked:text-black dark:checked:text-white border-base-content m-0.5 px-6"

    const incFilterBtnClasses = "btn btn-xs btn-neutral checked:text-black dark:checked:text-white border-base-content m-0.5 px-3"

    const statClassName = 'stat-value text-xl'

    const handleIncidentCount = () => {
        if (currentIncCount === null) {
            return (
                <div className={statClassName}>NaN</div>
            )
        }
        return (
            <div className={statClassName}>{currentIncCount}/{incidents?.length}</div>
        )
    }

    const severityFilter = () => {
        return (
            <div className='flex flex-col mx-4 my-2 items-center'>
                <span className='text-sm font-semibold m-1'>Severity</span>
                <div className='flex flex-col w-full'>
                    <input className={`checked:bg-success ${severityFilterBtnClasses}`} type="radio" name="frameworks" aria-label="All" value="all" onClick={() => setFilter(null)} defaultChecked/>
                    <div className='flex justify-center'>
                        <input className={`checked:bg-red-400 ${severityFilterBtnClasses}`} type="radio" name="frameworks" aria-label="Critical" value='critical' onClick={(e) => setFilter(e.currentTarget.value)}/>
                        <input className={`checked:bg-orange-400 ${severityFilterBtnClasses}`} type="radio" name="frameworks" aria-label="High" value='high' onClick={(e) => setFilter(e.currentTarget.value)}/>
                    </div>
                    <div className='flex justify-center'>
                        <input className={`checked:bg-orange-300 ${severityFilterBtnClasses}`} type="radio" name="frameworks" aria-label="Medium" value='medium' onClick={(e) => setFilter(e.currentTarget.value)}/>
                        <input className={`checked:bg-yellow-500 ${severityFilterBtnClasses}`} type="radio" name="frameworks" aria-label="Low" value='low' onClick={(e) => setFilter(e.currentTarget.value)}/>
                    </div>
                </div>
            </div>
        )
    }

    const showIncidentsFilter = () => {
        return (
            <div className='flex flex-col items-center mx-4 my-2'>
                <span className='text-sm font-semibold m-1'>Show incidents</span>
                <div className='flex flex-col w-full'>
                    <input type='radio' name='frameworks2' aria-label='All' value='all' className={`checked:bg-success checked:text-black ${incFilterBtnClasses}`} onClick={(e) => setResolvedFilter(e.currentTarget.value)} defaultChecked/>
                    <div className='flex'>
                        <input type='radio' name='frameworks2' aria-label='Unresolved' value='unresolved' className={`checked:bg-warning checked:text-black ${incFilterBtnClasses}`} onClick={(e) => setResolvedFilter(e.currentTarget.value)}/>
                        <input type='radio' name='frameworks2' aria-label='Resolved' value='resolved' className={`checked:bg-success checked:text-black ${incFilterBtnClasses}`} onClick={(e) => setResolvedFilter(e.currentTarget.value)}/>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div className='flex flex-col bg-base-200 rounded-lg my-4 ml-2 mr-4 h-[97vh] w-60 items-center shadow border border-base-content/15'>
            <div className='flex flex-col'>
                <div className="flex flex-col stats stats-vertical bg-base-300 m-4 shadow border border-base-content/15">
                    <div className="stat">
                        <div className="stat-title">Incident Count</div>
                        {handleIncidentCount()}
                        <div className="stat-desc">Current count out of all</div>
                    </div>
                    <IncidentReportProvider>
                        <IncidentReportSidebar />
                    </IncidentReportProvider>
                </div>
                <div className='flex flex-col rounded-lg shadow mx-4 bg-base-300 border border-base-content/15'>
                    <span className='text-xl font-bold text-center mt-2'>Filters</span>
                    {severityFilter()}
                    {showIncidentsFilter()}
                </div>
            </div>
        </div>
    )
}