import { DownloadReport, IncidentReportProvider, useIncidentReport } from './IncidentReport'
import { ExportIncidents, IncidentListProvider } from './IncidentList'

const HomePage = () => {
    return (
        <div className="flex flex-col grow rounded-lg justify-start h-[97vh] m-4 bg-base-200 border border-base-content/15">
            <span className='text-center text-4xl mt-4 font-bold bg-linear-to-r from-blue-500 to-cyan-400 text-transparent bg-clip-text'>Incident Analyzer</span>
            <div className='flex flex-col ml-20 m-2'>
                <div className='flex flex-col mt-4'>
                    <div role="alert" className="alert alert-soft w-fit mb-1 bg-linear-to-r from-emerald-300 to-green-400 dark:from-emerald-800 dark:to-green-900">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-info h-6 w-6 shrink-0">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                        <span className='text-lg'>Quickly add, edit or remove incidents. Download report or see all incidents in a sortable table with filters.</span>
                    </div>
                    <span className='font-light'>Quick help: Navigate through various features on the right side to do actions.</span>
                </div>
                <div className='flex flex-col mt-12'>
                    <span className='text-xl font-bold'>Quick actions</span>
                    <div className='flex mt-2'>
                        <IncidentReportProvider>
                            <IncidentsUnresolved />
                        </IncidentReportProvider>
                        <div className="flex flex-col shadow p-4 m-4 bg-base-300 rounded-lg h-fit border border-base-content/15">
                            <IncidentReportProvider>
                                <DownloadReport />
                            </IncidentReportProvider>
                        </div>
                        <div className="flex flex-col shadow p-4 m-4 bg-base-300 rounded-lg h-fit border border-base-content/15">
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
        <div className="flex flex-col shadow p-4 m-4 bg-base-300 rounded-lg h-fit border border-base-content/15">
            {displayIfUnresolved()}
        </div>
    )
    
}

export default HomePage