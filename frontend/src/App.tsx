import { useState } from 'react';
import HealthStatus from './components/HealthResponse'
import IncidentList from './components/IncidentList';
import IncidentAdd from './components/IncidentAdd';
import EnableMenu from './components/EnableMenu';

const App = () => {
    const [selectedMenu, setSelectedMenu] = useState<string>('Home')

    return (
        <div className='dark:bg-linear-to-r dark:from-[#252525] dark:to-[#363636]'>
            {/* Top header */}
            {/* <div className='flex justify-between items-center h-[3vh] bg-amber-500'>
                <span className='text-left text-md font-semibold ml-2'>Incident Dashboard</span>
                <div className='flex text-xs text-right mr-2'>
                    <HealthStatus />
                </div>
            </div> */}
            <div className='flex justify-between'>
                {/* Left side */}
                <div className='flex flex-col bg-base-300 rounded-xl m-2 h-[98.5vh]'>
                    <div className='flex justify-center items center m-2'>
                        <EnableMenu selected={selectedMenu} onSelect={setSelectedMenu} />
                    </div>
                    <div className='text-xs absolute bottom-3 left-4'>
                        <HealthStatus />
                    </div>
                </div>
                {/* Center */}
                <div className='grow rounded-xl h-[98.5vh] m-2 bg-base-300'>
                    { selectedMenu === 'Home' &&
                    <div className='flex flex-col justify-center items-center m-2 mt-8'>
                        <span className='text-center text-4xl font-bold bg-linear-to-r from-gray-400 to-gray-500 text-transparent bg-clip-text'>Coming Soon!</span>
                        <div className="flex w-52 justify-center items-center flex-col gap-4 p-4 mt-8">
                            <div className="skeleton h-45 w-full bg-base-200"></div>
                            <div className="skeleton h-4 w-28 bg-base-200"></div>
                            <div className="skeleton h-4 w-40 bg-base-200"></div>
                            <div className="skeleton h-4 w-34 bg-base-200"></div>
                        </div>
                    </div>}
                    { selectedMenu === 'Incidents' && <IncidentList />}
                    { selectedMenu === 'Add' && <IncidentAdd />}
                    { selectedMenu === 'Report' && 
                    <div className='flex flex-col justify-center items-center m-2 mt-8'>
                        <div className='flex'>
                            <span className='text-base-content m-1'>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-8">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M10.34 15.84c-.688-.06-1.386-.09-2.09-.09H7.5a4.5 4.5 0 1 1 0-9h.75c.704 0 1.402-.03 2.09-.09m0 9.18c.253.962.584 1.892.985 2.783.247.55.06 1.21-.463 1.511l-.657.38c-.551.318-1.26.117-1.527-.461a20.845 20.845 0 0 1-1.44-4.282m3.102.069a18.03 18.03 0 0 1-.59-4.59c0-1.586.205-3.124.59-4.59m0 9.18a23.848 23.848 0 0 1 8.835 2.535M10.34 6.66a23.847 23.847 0 0 0 8.835-2.535m0 0A23.74 23.74 0 0 0 18.795 3m.38 1.125a23.91 23.91 0 0 1 1.014 5.395m-1.014 8.855c-.118.38-.245.754-.38 1.125m.38-1.125a23.91 23.91 0 0 0 1.014-5.395m0-3.46c.495.413.811 1.035.811 1.73 0 .695-.316 1.317-.811 1.73m0-3.46a24.347 24.347 0 0 1 0 3.46" />
                                </svg>
                            </span>
                            <span className='text-4xl font-bold bg-linear-to-r from-lime-400 to-green-500 text-transparent bg-clip-text'>
                            Please hang up and try again.</span>
                        </div>
                        <div className="flex w-52 justify-center items-center flex-col gap-4 p-4 mt-8">
                            <div className="skeleton h-45 w-full mb-2 bg-base-200"></div>
                            <div className="skeleton h-4 w-28 bg-base-200"></div>
                            <div className="skeleton h-4 w-full bg-base-200"></div>
                            <div className="skeleton h-4 w-34 bg-base-200"></div>
                            <div className="skeleton h-4 w-40 bg-base-200"></div>
                            <div className="skeleton h-4 w-full bg-base-200"></div>
                        </div>
                    </div>}
                </div>
                {/* Right side */}
                <div className='flex-col bg-base-300 rounded-xl m-2 h-[98.5vh] w-50 justify-end'>
                    <span></span>
                </div>
            </div>
        </div>
    )
}

export default App;