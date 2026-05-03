import { useState, useEffect, createContext, useContext } from "react"
import ToastNotification from "./ToastNotification"
import type { Incident } from "./types"

type IncidentsByID = Record<string, Incident>
type IncidentsGroupped = Record<string, Record<string, Incident>>

interface IncidentReport {
    incidents_count: number
    unresolved_ids: string[]
    mttr: string
    by_services: IncidentsGroupped
    by_severity: IncidentsGroupped
    by_id: IncidentsByID
}

interface IncidentReportContextType {
    report: IncidentReport | undefined
    loading: boolean
    error: string | null
}

const IncidentReportContext = createContext<IncidentReportContextType | undefined>(undefined)

export const IncidentReportProvider = ({ children }: { children: React.ReactNode }) => {
    const [error, setError] = useState<string | null>(null)
    const [loading, setLoading] = useState<boolean>(true)
    const [report, setReport] = useState<IncidentReport | undefined>(undefined)

    useEffect(() => {
        const fetchReport = async () => {
            try {
                const response = await fetch("http://localhost:8080/report")
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}`)
                }
                const data: IncidentReport = await response.json()
                setReport(data)
                setLoading(false)
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Unknown error')
                setLoading(false)
            }
        }
        fetchReport()
    }, [])

    return (
        <IncidentReportContext.Provider value={{ report, loading, error }}>
            {children}
        </IncidentReportContext.Provider>
    )
}

const useIncidentReport = () => {
    const context = useContext(IncidentReportContext)
    if (!context) {
        throw new Error("useIncidentReport must be within IncidentReportProvider")
    }
    return context
}

export const IncidentReportCenter = () => {
    // const [group, setGroup] = useState<string>('id')
    const { report, error, loading } = useIncidentReport()
    
    if (error) {
        return <ToastNotification duration={10000} message={`Error ${error}`} toastLevel='alert-error' toastPos='toast-top toast-right'/>
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
    if (!report) {
        return (
            <span>Report is empty</span>
        )
    }
    
    const prepareExport = (): string => {
        const blob = new Blob([JSON.stringify(report, null, 2)], { type: 'application/json' })
        return URL.createObjectURL(blob)
    }

    // const filterReport = (group: string) => {
    //     switch (group) {
    //         case 'id': return report?.by_id ?? {}
    //         case 'severity': return report?.by_severity ?? {}
    //         case 'services': return report?.by_services ?? {}
    //         default : return report?.by_id
    //     }
    // }w

    return (
        <div className="flex flex-col items-center grow rounded-xl h-[97vh] m-4 bg-base-300">
            <span className="flex bg-linear-to-r text-3xl font-bold p-2 from-[#3388ff] to-[#2f64b9] text-transparent bg-clip-text">Report</span>
            <div className="flex flex-col justify-start overflow-auto">
                <div className="flex flex-col overflow-auto">
                    <span className="text-xl font-semibold text-right mx-2">Preview</span>
                    <pre className="text-xs bg-black text-base-100 p-4 rounded-xl overflow-auto">{JSON.stringify(report, null, 2)}</pre>
                </div>
                <div className="flex flex-col items-center m-2">
                    <span className="text font-semibold">Download report</span>
                    <button className="btn btn-sm border-base-100" onClick={() => {
                        const url = prepareExport()
                        const link = document.createElement('a')
                        const time = new Date()
                        link.href = url
                        link.download = `report-${time.toLocaleDateString()}-${time.toLocaleTimeString()}.json`
                        link.click()
                        URL.revokeObjectURL(url)
                    }}>
                        Download
                    </button>
                </div>
            </div>
        </div>
    )
}

export const IncidentReportSidebar = () => {
    const { report, loading } = useIncidentReport()

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

    return (
        <>
            <div className="stat">
                <div className="stat-title">MTTR</div>
                <div className="stat-value">{report?.mttr}</div>
                <div className="stat-desc">Mean time to average</div>
            </div>
            <div className="stat">
                <div className="stat-title">Total Unresolved</div>
                <div className="stat-value">{report?.unresolved_ids.length}</div>
            </div>
        </>
    )
}