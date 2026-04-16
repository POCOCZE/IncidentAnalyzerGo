import { useState } from 'react';
import HealthStatus from './components/HealthResponse'
import IncidentList from './components/IncidentList';

const App = () => {
    const [renderIncidents, setRenderIncidents] = useState<boolean>(false)
    const toggleIncidentsVisibility = () => {
        setRenderIncidents(!renderIncidents);
    }

    return (
        <div className='dark:bg-linear-to-r from-[#252525] to-[#363636]'>
            <div className='flex justify-between items-center'>
                <h1 className='text-left text-md font-semibold ml-2'>Incident Dashboard</h1>
                <div className='flex text-xs text-right mr-2'>
                    <HealthStatus />
                </div>
            </div>
            <div className='flex'>
                {/* Left side */}
                <div className='min-w-45 bg-base-200 rounded-xl m-2 h-[calc(100vh-4vh)]'>
                    <div className='flex mt-1 items-center justify-center'>
                        <span className='text-sm mr-2'>
                            <p>Show incidents</p>
                        </span>
                        <span>
                            <input type="checkbox" className="toggle" onClick={toggleIncidentsVisibility} />
                        </span>
                    </div>
                    <div>
                        <button className='btn'>Nothing</button>
                    </div>
                </div>
                {/* Center */}
                <div className='flex justify-center items-top bg-base-200 rounded-xl mx-2 m-2 h-[calc(100vh-4vh)]'>
                    {renderIncidents && <IncidentList />
                    }
                </div>
            </div>
        </div>
    )
}

export default App;