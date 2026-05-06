import { DownloadReport, IncidentReportProvider } from './IncidentReport'
import { ExportIncidents, IncidentListProvider } from './IncidentList'

const HomePage = () => {
    return (
        <div className="flex flex-col grow rounded-xl justify-start h-[97vh] m-4 bg-base-300">
            <span className='text-center text-4xl mt-4 font-bold bg-linear-to-r from-blue-500 to-cyan-400 text-transparent bg-clip-text'>Incident Analyzer</span>
            <div className='flex flex-col ml-20 m-2'>
                <div className='flex flex-col mt-6'>
                    <span className='text-lg font-semibold'>This tool allows to quickly add or remove incidents, download incident report or see all incidents in a sortable table with filters.</span>
                    <span className='font-light'>Quick help: Navigate through various features on the right side to do actions.</span>
                </div>
                <div className='flex flex-col mt-12'>
                    <span className='text-xl font-bold'>Quick actions</span>
                    <div className='flex mt-2'>
                        <div className="flex flex-col p-4 m-4 bg-base-200 rounded-xl shadow h-fit">
                            <IncidentReportProvider>
                                <DownloadReport />
                            </IncidentReportProvider>
                        </div>
                        <div className="flex flex-col p-4 m-4 bg-base-200 rounded-xl shadow h-fit">
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

export default HomePage