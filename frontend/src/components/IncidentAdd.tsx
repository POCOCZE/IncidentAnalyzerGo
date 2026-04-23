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
    const [error, setError] = useState<string | null>(null)

    const [id, setID] = useState<string>('')
    const [title, setTitle] = useState<string>('')
    const [severity, setSeverity] = useState<string>('critical')
    const [serviceName, setServiceName] = useState<string>('')
    const [startedAt, setStartedAt] = useState<string>('')
    const [resolvedAt, setResolvedAt] = useState<string>('')

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
            const response = await fetch(`http://localhost:8080/incidents`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({id: id, title: title, severity: severity, service_name: serviceName, started_at: finalStartedAt, resolved_at: finalResolvedAt})
            })
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`)
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
        if (isSuccess) {
            return <ToastNotification duration={5000} message="Successfully added new incident" toastLevel="alert-success" toastPos="toast-top toast-end" />
        }
    }

    return (
        <div className="flex flex-col grow rounded-xl h-[97vh] m-4 bg-base-300">
            <span className="flex text-3xl font-bold justify-center mt-4 bg-linear-to-r from-amber-600 to-yellow-400 text-transparent bg-clip-text">Add Incident</span>
            <RenderToast />
            <fieldset className="fieldset border-base-300 rounded-box border p-4">
                {/* <legend className="fieldset-legend">Page details</legend> */}
                <div className="lg:flex md:flex justify-center">
                    <div className="flex flex-col md:px-2 lg:px-[8vw]">
                        <label className="label">Incident name</label>
                        <input value={id} onChange={(e) => setID(e.target.value)} type="text" className="input mb-4 lg:min-w-[20vw]" placeholder="INC-001" />

                        <label className="label">Title</label>
                        <input value={title} onChange={(e) => setTitle(e.target.value)} type="text" className="input mb-4 lg:min-w-[20vw]" placeholder="PostgreSQL corruption" />

                        <label className="label">Severity</label>
                        <label className="select mb-4 lg:min-w-[20vw]">
                            <select value={severity} onChange={(e) => setSeverity(e.target.value)}>
                                <option value='critical'>Critical</option>
                                <option value='high'>High</option>
                                <option value='medium'>Medium</option>
                                <option value='low'>Low</option>
                            </select>
                        </label>
                    </div>
                    <div className='flex flex-col md:px-2 lg:px-[8vw]'>
                        <label className="label">Service Name</label>
                        <input value={serviceName} onChange={(e) => setServiceName(e.target.value)} type="text" className="input mb-4 lg:min-w-[20vw]" placeholder="postgresql-replica-1" />

                        <label className="label">Started at</label>
                        <label className="input mb-4 lg:min-w-[20vw]">
                            <input value={startedAt} onChange={(e) => setStartedAt(e.target.value)} type="datetime-local" />
                        </label>

                        <label className="flex items-center label">Resolved at<span className="badge badge-soft badge-xs ml-0.5">Optional</span></label>
                        <label className="input mb-4 lg:min-w-[20vw]">
                            <input value={resolvedAt} onChange={(e) => setResolvedAt(e.target.value)} type="datetime-local" />
                        </label>
                    </div>
                </div>
                <div className="flex flex-col items-center">
                    <span className="text-gray-600"><span className="badge badge-soft badge-xs mr-0.5">Coming Soon!</span>Import multiple incidents from existing file:</span>
                    <input type="file" placeholder="You can't touch this" className="file-input" disabled />
                    <button onClick={postIncident} className="btn btn-neutral bg-base-100 border-base-content/15 active:bg-success active:text-base-content mt-6 lg:w-40">Submit</button>
                </div>
            </fieldset>
        </div>
    )
}

export default IncidentAdd