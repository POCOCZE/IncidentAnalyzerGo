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
    resolvedFilter: "all" | "unresolved"
    setResolvedFilter: React.Dispatch<SetStateAction<"all" | "unresolved">>
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
    const [resolvedFilter, setResolvedFilter] = useState<"all" | "unresolved">('all')
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

export const IncidentListCenter = () => {
    const { incidents, setIncidents, setCurrentIncCount, filter, resolvedFilter, loading, error, setError } = useIncidentList()
    const [searchBg, setSearchBg] = useState<string>('')
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
        critical: 'badge border-base-content bg-red-300 text-[#ba1c09] font-semibold',
        high: 'badge border-base-content bg-[#ffb954] text-[#966825] font-semibold',
        medium: 'badge border-base-content bg-[#FFCE47] text-[#966825] font-semibold',
        low: 'badge border-base-content bg-[#E8DB27] text-[#988F12] font-semibold',
    };

    const columns = [
        { key: 'id', label: 'ID'},
        { key: 'title', label: 'Title'},
        { key: 'severity', label: 'Severity', 
            render: (value: string) => <span className={severityColor[value]}>{value}</span>},
        { key: 'service_name', label: 'Service name'},
        // { key: 'started_at', label: 'Started at'},
        // { key: 'resolved_at', label: 'Resolved at'},
        { key: 'message', label: 'Message'},
        { key: 'is_resolved', label: 'Is resolved'},
        // delete button...
        { key: 'delete', label: ''}
    ]

    return (
        <div className='flex flex-col px-3 pt-3 grow rounded-xl h-[97vh] m-4 bg-base-300 shadow'>
            <span className='flex text-3xl font-bold justify-center mb-4 bg-linear-to-r from bg-orange-500 to-yellow-500 bg-clip-text text-transparent'>Incident list</span>
            <div className='flex justify-end my-2'>
                <label className={`input rounded-2xl w-60 border-base-content/50 ${searchBg}`} onMouseEnter={() => setSearchBg('bg-base-200')} onMouseLeave={() => setSearchBg('')}>
                    <svg className="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g strokeLinejoin="round" strokeLinecap="round" strokeWidth="2.5" fill="none" stroke="currentColor"><circle cx="11" cy="11" r="8"></circle><path d="m21 21-4.3-4.3"></path></g></svg>
                    <input type="search" placeholder="Search incidents" value={searchKeyword} onChange={(e) => setSearchKeyword(e.currentTarget.value)} />
                </label>
            </div>
            <div className='rounded-box border border-base-content/50 m-1 mb-4 overflow-auto'>
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

    const severityFilterClasses = "btn btn-xs btn-neutral border-base-content m-0.5 px-5.5"

    const handleIncidentCount = () => {
        if (currentIncCount === null) {
            return (
                <div className="stat-value">NaN</div>
            )
        }
        return (
            <div className="stat-value">{currentIncCount}/{incidents?.length}</div>
        )
    }

    return (
        <div className='flex flex-col bg-base-300 rounded-xl my-4 ml-2 mr-4 h-[97vh] w-60 items-center shadow'>
            <div className='flex flex-col'>
                <div className="flex flex-col stats stats-vertical bg-base-200 m-4 shadow">
                    <div className="stat">
                        <div className="stat-title">Incident Count</div>
                        {handleIncidentCount()}
                        <div className="stat-desc">Current count out of all</div>
                    </div>
                    <IncidentReportProvider>
                        <IncidentReportSidebar />
                    </IncidentReportProvider>
                </div>
                <div className='flex flex-col rounded-xl shadow mx-4 bg-base-200'>
                    <span className='text-xl font-bold text-center mt-2'>Filters</span>
                    <div className='flex flex-col mx-4 my-2 items-center'>
                        <span>Severity</span>
                        <div className='flex flex-col'>
                            <input className={`checked:bg-green-200 checked:text-[#41ba09] ${severityFilterClasses}`} type="radio" name="frameworks" aria-label="All" value="all" onClick={() => setFilter(null)} defaultChecked/>
                            <div className='flex'>
                                <input className={`checked:bg-red-300 checked:text-[#ba1c09] ${severityFilterClasses}`} type="radio" name="frameworks" aria-label="Critical" value='critical' onClick={(e) => setFilter(e.currentTarget.value)}/>
                                <input className={`checked:bg-[#ffb954] checked:text-[#966825] ${severityFilterClasses}`} type="radio" name="frameworks" aria-label="High" value='high' onClick={(e) => setFilter(e.currentTarget.value)}/>
                            </div>
                            <div className='flex'>
                                <input className={`checked:bg-[#FFCE47] checked:text-[#966825] ${severityFilterClasses}`} type="radio" name="frameworks" aria-label="Medium" value='medium' onClick={(e) => setFilter(e.currentTarget.value)}/>
                                <input className={`checked:bg-[#E8DB27] checked:text-[#988F12] ${severityFilterClasses}`} type="radio" name="frameworks" aria-label="Low" value='low' onClick={(e) => setFilter(e.currentTarget.value)}/>
                            </div>
                        </div>
                    </div>
                    <div className='text-center mb-4'>
                        <span>Show incidents</span>
                        <div className='items-center'>
                            <label className="swap">
                                <input type="checkbox"/>
                                <span className="bg-green-200 text-[#41ba09] rounded-xl text-sm px-2 swap-off" onClick={() => setResolvedFilter("all")} defaultChecked>All</span>
                                <span className="bg-[#FFCE47] text-[#966825] rounded-xl text-sm px-2 swap-on" onClick={() => setResolvedFilter("unresolved")}>Unresolved</span>
                            </label>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}