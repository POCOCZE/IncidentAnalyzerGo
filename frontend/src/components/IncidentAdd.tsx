import { useState } from "react";
import ToastNotification from "./ToastNotification";

// Maybe switch const to interface in the future to make it cleaner
// interface Incident {
//     id: string
//     title: string
//     severity: string
//     service_name: string
//     started_at: string
//     resolved_at: string | undefined
// }

const IncidentAdd = () => {
    const [isSuccess, setIsSuccess] = useState<boolean>(false)
    const [isDataSubmit, setIsDataSubmit] = useState<boolean>(false)
    const [infoMessage, setInfoMessage] = useState<string>('')
    const [error, setError] = useState<string | null>(null)

    const [id, setID] = useState<string>('')
    const [title, setTitle] = useState<string>('')
    const [severity, setSeverity] = useState<string>('critical')
    const [serviceName, setServiceName] = useState<string>('')
    const [startedAt, setStartedAt] = useState<string>('')
    const [resolvedAt, setResolvedAt] = useState<string>('')
    const [incidentsFile, setIncidentsFile] = useState<File | null>(null)

    const calcTimeZoneOffset = async ():Promise<string> => {
        const formatter = new Intl.DateTimeFormat('en-US', {
            timeZoneName: 'longOffset'
        })
        const formatted = formatter.format(new Date())
        const offset = formatted.split('GMT')[1] // 
        const result = ':00' + offset
        return result
    }

    const postIncident = async () => {
        try {
            // Correct time with timezone; applies for startedAt and also for resolvedAt (if defined)
            const result = await calcTimeZoneOffset()
            const finalStartedAt = startedAt + result
            const finalResolvedAt = resolvedAt !== '' ? resolvedAt + result : null

            // Reset variables before trying again
            setIsSuccess(false)
            setError(null)

            // If user did not imported list of incidents, add one incident
            if (incidentsFile === null) {
                const response = await fetch(`http://localhost:8080/incident`, {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({id: id, title: title, severity: severity, service_name: serviceName, started_at: finalStartedAt, resolved_at: finalResolvedAt})
                })
                if (!response?.ok) {
                    const errorData = await response?.json()
                    console.log('Error response', errorData)
                    throw new Error(`- ${errorData.error}`)
                } else {
                    const data = await response?.json()
                    console.log(data)
                    setInfoMessage(data.info)
                }
            } else {
                const contents = await incidentsFile.text()
                const response = await fetch(`http://localhost:8080/incidents`, {
                    method: 'POST',
                    headers: {"Content-Type": "application/json"},
                    body: contents
                })
                if (!response?.ok) {
                    const errorData = await response?.json()
                    console.log('Error response', errorData)
                    throw new Error(`- ${errorData.error}`)
                } else {
                    const data = await response?.json()
                    console.log(data)
                    setInfoMessage(data.info)
                }
            }
            setIsSuccess(true)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Unknown error')
        }
    }

    const RenderToast = () => {
        if (error) {
            return <ToastNotification duration={10000} message={`Error ${error}`} toastLevel='alert-error' toastPos='toast-top toast-right'/>
        }
        if (infoMessage !== '') {
            return <ToastNotification duration={5000} message={`Info - ${infoMessage}`} toastLevel="alert-info" toastPos="toast-top toast-end" />
        }
        // After rendering a notification - reset the boolean
        setIsDataSubmit(false)
    }

    const minWidthFields = 'flex flex-col'
    const inputWidthFields = 'mb-4 lg:min-w-[50vw] md:min-w-[50vw] sm:min-w-[50vw] min-w-[70vw]'

    return (
        <div className="flex flex-col grow rounded-lg h-[97vh] m-4 bg-base-200 border border-base-content/15">
            <span className="flex text-3xl font-bold justify-center mt-4 bg-linear-to-r from-amber-600 to-yellow-400 text-transparent bg-clip-text">Add Incident</span>
            {isDataSubmit === true && RenderToast()}
            <fieldset className="fieldset border-base-200 rounded-box border p-4">
                {/* <legend className="fieldset-legend">Page details</legend> */}
                <div className="flex flex-col items-center">
                    <div className={minWidthFields}>
                        <label className="label">Incident name / ID</label>
                        <input value={id} onChange={(e) => setID(e.target.value)} type="text" className={`input ${inputWidthFields}`} placeholder="INC-001" />
                    </div>
                    <div className={minWidthFields}>
                        <label className="label">Title</label>
                        <input value={title} onChange={(e) => setTitle(e.target.value)} type="text" className={`input ${inputWidthFields}`} placeholder="PostgreSQL corruption" />
                    </div>
                    <div className={minWidthFields}>
                        <label className="label">Severity</label>
                        <label className={`select ${inputWidthFields}`}>
                            <select value={severity} onChange={(e) => setSeverity(e.target.value)}>
                                <option value='critical'>Critical</option>
                                <option value='high'>High</option>
                                <option value='medium'>Medium</option>
                                <option value='low'>Low</option>
                            </select>
                        </label>
                    </div>
                    <div className={minWidthFields}>
                        <label className="label">Service Name</label>
                        <input value={serviceName} onChange={(e) => setServiceName(e.target.value)} type="text" className={`input ${inputWidthFields}`} placeholder="postgresql-replica-1" />
                    </div>
                    <div className={minWidthFields}>
                        <label className="label">Started at</label>
                        <label className={`input ${inputWidthFields}`}>
                            <input value={startedAt} onChange={(e) => setStartedAt(e.target.value)} type="datetime-local" />
                        </label>
                    </div>
                    <div className={minWidthFields}>
                        <label className="flex items-center label">Resolved at<span className="badge badge-soft badge-xs ml-0.5">Optional</span></label>
                        <label className={`input ${inputWidthFields}`}>
                            <input value={resolvedAt} onChange={(e) => setResolvedAt(e.target.value)} type="datetime-local" />
                        </label>
                    </div>
                </div>
                <div className="flex flex-col items-center">
                    {/* Todo: Set maximum file size for the file */} 
                    <span className="text-gray-600">Import incidents from file:</span>
                    <input type="file" placeholder="You can't touch this" className="file-input" onChange={(e) => {
                        const file = e.currentTarget.files?.[0] ?? null
                        setIncidentsFile(file)
                    }} />
                    <button onClick={() => {
                        postIncident()
                        setIsDataSubmit(true)
                    }} className="btn btn-neutral bg-base-100 border-base-content/15 active:text-base-content mt-6 lg:w-40">Submit</button>
                </div>
            </fieldset>
        </div>
    )
}

export default IncidentAdd