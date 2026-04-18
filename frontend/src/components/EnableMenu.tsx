interface MenuProps {
    selected: string
    onSelect: (value: string) => void
}

const EnableMenu = ({selected, onSelect}: MenuProps) => {
    const selectedLogic = (buttonName: string) => {
        if (selected !== buttonName) {
            onSelect(buttonName)
        }
        // This means when the button is clicked for the second time
        // Not used since its not really clever and best practice.
        // if (selected === buttonName) {
        //     onSelect(null)
        // }
    }

    return(
        <ul className="menu bg-base-200 rounded-box">
            <li className={selected === 'Home' ? 'menu-active rounded-2xl': ''} onClick={() => selectedLogic('Home')}>
                <a className="tooltip tooltip-right" data-tip="Home">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
                </svg>
                </a>
            </li>
            <li className={selected === 'Add' ? 'menu-active rounded-2xl': ''} onClick={() => {selectedLogic('Add')}}>
                <a className="tooltip tooltip-right" data-tip="Add">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                </svg>
                </a>
            </li>
            <li className={selected === 'Incidents' ? 'menu-active rounded-2xl': ''} onClick={() => {selectedLogic('Incidents')}}>
                <a className="tooltip tooltip-right" data-tip="Incidents">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 0 1 0 3.75H5.625a1.875 1.875 0 0 1 0-3.75Z" />
                </svg>
                </a>
            </li>
            <li className={selected === 'Report' ? 'menu-active rounded-2xl': ''} onClick={() => {selectedLogic('Report')}}>
                <a className="tooltip tooltip-right" data-tip="Report">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z" />
                </svg>
                </a>
            </li>
        </ul>
    )
}

export default EnableMenu