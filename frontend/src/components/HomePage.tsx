import { DownloadReport, IncidentReportProvider, useIncidentReport } from './IncidentReport'
import { ExportIncidents, IncidentListProvider } from './IncidentList'

const HomePage = () => {
    return (
        <div className="flex flex-col grow rounded-lg justify-start h-[97vh] m-4 bg-base-200">
            <span className='text-center text-4xl mt-4 font-bold bg-linear-to-r from-blue-500 to-cyan-400 text-transparent bg-clip-text'>Incident Analyzer</span>
            <div className='flex flex-col ml-20 m-2'>
                <div className='flex flex-col mt-6'>
                    <span className='text-lg font-semibold'>This tool allows to quickly add or remove incidents, download incident report or see all incidents in a sortable table with filters.</span>
                    <span className='font-light'>Quick help: Navigate through various features on the right side to do actions.</span>
                </div>
                <div className='flex flex-col mt-12'>
                    <span className='text-xl font-bold'>Quick actions</span>
                    <div className='flex mt-2'>
                        <IncidentReportProvider>
                            <IncidentsUnresolved />
                        </IncidentReportProvider>
                        <div className="flex flex-col shadow-xl p-4 m-4 bg-base-300 rounded-lg h-fit">
                            <IncidentReportProvider>
                                <DownloadReport />
                            </IncidentReportProvider>
                        </div>
                        <div className="flex flex-col shadow-xl p-4 m-4 bg-base-300 rounded-lg h-fit">
                            <IncidentListProvider>
                                <ExportIncidents />
                            </IncidentListProvider>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

const IncidentsUnresolved = () => {
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

    const statClassName = 'stat-value text-xl'

    const displayIfUnresolved = () => {
        if (report?.unresolved_ids === undefined) {
            return <div className={statClassName}>NaN</div>
        } else {
            return (
                <div className="flex flex-col items-center m-2">
                <div className="flex">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="red" className="size-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
                </svg>
                <span className="text font-semibold ml-1">Unresolved</span>
                </div>
                <span className="text-xl text-red-500 font-bold">{report?.unresolved_ids?.length}</span>
                <button className="btn btn-sm btn-neutral border-base-content mt-4 m-0.5 px-5.5" disabled>More</button>
                </div>
            )
        }
    }
    if (report?.unresolved_ids === null) {
        return null
    }
    return (
        <div className="flex flex-col shadow-xl p-4 m-4 bg-base-300 rounded-lg h-fit">
            {displayIfUnresolved()}
        </div>
    )
    
}

export default HomePage