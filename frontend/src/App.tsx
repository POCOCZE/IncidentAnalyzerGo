import { useState } from 'react';
import HealthStatus from './components/HealthResponse'
import { IncidentListProvider, IncidentListCenter, IncidentListSidebar} from './components/IncidentList';
import IncidentAdd from './components/IncidentAdd';
import { IncidentReportProvider, IncidentReportCenter } from './components/IncidentReport';
import EnableMenu from './components/EnableMenu';
import HomePage from './components/HomePage';

const App = () => {
    const [selectedMenu, setSelectedMenu] = useState<string>('Home')

    return (
        <div className='flex bg-base-100'>
            {/* Left side */}
            <div className='flex flex-col justify-between items-center rounded-xl ml-4 my-4 h-[97vh]'>
                <EnableMenu selected={selectedMenu} onSelect={setSelectedMenu} />
                <HealthStatus />
            </div>
            {/* Center */}
            { selectedMenu === 'Home' && <HomePage />}
            { selectedMenu === 'Incidents' && 
                <IncidentListProvider>
                    <IncidentListCenter />
                    <IncidentListSidebar />
                </IncidentListProvider>}
            { selectedMenu === 'Add' && <IncidentAdd />}
            { selectedMenu === 'Report' &&
                <IncidentReportProvider>
                    <IncidentReportCenter />
                    {/* <IncidentReportSidebar /> */}
                </IncidentReportProvider>}
        </div>
    )
}

export default App;