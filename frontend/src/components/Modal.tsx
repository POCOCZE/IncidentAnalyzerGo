import { useState, useEffect, useRef } from 'react'
import type { Incident } from './types'
import ToastNotification from './ToastNotification'

interface ModalProps {
    buttonText: React.ReactNode
    title: React.ReactNode
    description: React.ReactNode
    closeButtonText: React.ReactNode
}

interface IncidentEditProps {
    editIcon: React.ReactNode
    incidentID: string
}

export const OneButtonModal = ({buttonText, title, description, closeButtonText}: ModalProps) => {
    const modalRef = useRef<HTMLDialogElement>(null)

    return (
        <>
        {/* Open modal */}
        <button className='btn btn-xs hover:btn-error btn-ghost' onClick={() => modalRef.current?.showModal()}>{buttonText}</button>
        <dialog ref={modalRef} className="modal modal-bottom sm:modal-middle">
            <div className="modal-box">
                <span className="font-bold text-xl">{title}</span>
                <p className="pt-3">{description}</p>
                <div className="modal-action">
                    <form method="dialog">
                        {/* if there is a button in form, it will close the modal */}
                        <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
                        <span>{closeButtonText}</span>
                    </form>
                </div>
            </div>
        </dialog>
        </>
    )
}

export const IncidentEditModal = ({editIcon, incidentID}: IncidentEditProps) => {
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string | null>(null)
    const [isSuccess, setIsSuccess] = useState<boolean>(false)
    const [submitBtnBg, setSubmitBtnBg] = useState<string>('btn-neutral')

    const [id, setID] = useState<string>('')
    const [title, setTitle] = useState<string>('')
    const [severity, setSeverity] = useState<string>('')
    const [serviceName, setServiceName] = useState<string>('')
    const [startedAt, setStartedAt] = useState<string>('')
    const [resolvedAt, setResolvedAt] = useState<string | undefined>(undefined)
    const [isResolved, setIsResolved] = useState<boolean | undefined>(undefined)
    const [originalResolvedAt, setOriginalResolvedAt] = useState<string | undefined>(undefined)
    const [enableResolvedAtOption, setEnableResolvedAtOption] = useState<boolean>(false)

    const fetchIncidentByID = async () => {
        try {
            const response = await fetch(`http://localhost:8080/incidents/${incidentID}`)
            if (!response?.ok) {
                const errorData = await response?.json()
                console.log('Error response', errorData)
                throw new Error(`- ${errorData.error}`)
            }
            const data: Incident = await response.json()
            setLoading(false)
            setID(data.id)
            setTitle(data.title)
            setSeverity(data.severity)
            setServiceName(data.service_name)
            setStartedAt(data.started_at)
            if (data.resolved_at === null) {
                console.log("Fetched: ResolvedAt is null. IsResolved:", data.is_resolved)
                setResolvedAt('')
                setIsResolved(data.is_resolved)
            } else {
                console.log("Fetched: IsResolved:", data.is_resolved)
                setResolvedAt(data.resolved_at)
                console.log("ResolvedAt is", data.resolved_at)
                setOriginalResolvedAt(data.resolved_at)
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Unknown error')
            setLoading(false)
        }
    }

    const RenderLoading = () => {
        if (loading) {
            return (
                <div className="flex w-4 justify-center items-center flex-col gap-4">
                    <div className="skeleton h-4 w-full"></div>
                </div>
            )   
        }
    }

    const ResolvedAtToUTC = async ():Promise<string | undefined> => {
        const now = new Date()
        const timezoneOffsetMinutes = now.getTimezoneOffset()

        // Convert to UTC
        if (resolvedAt === undefined) {
            throw new Error("ResolvedAt is undefined.")
        }
        const userDate = new Date(resolvedAt)
        const utcDate = new Date(userDate.getTime() + timezoneOffsetMinutes * 60000)
        return utcDate.toISOString()
    }

    const patchIncident = async () => {
        try {
            // Correct time with timezone; applies for startedAt and also for resolvedAt (if defined)
            // // ! Error: I should probsbly boolean check whether time are within the same format - if not convert them to one unique one - the must respect timezones when adding or editing incdients and then chanding them to the propriate times. It's for some reason causing trouble there with error: `ERR: parsing time "2026-03-17T01:00+02:00" as "2006-01-02T15:04:05Z07:00": cannot parse "+02:00" as ":"` - the time is obviously correct, but still the backend cannot accept it after adding the timezone browser correction "+02:00" part. IncidentAdd have no problems with this - this problem is. very strange. i am not aware of any restrictions that would check two times against each other - so when one time would be with the "Z07:00" format and other with the other format - they dont match - and this could be the problem - but i am not comparing this. Wierd.
            
            // Todo: add a check when resolvedAt already contains the timezone format - do not add one extra - this could happen in cases when editing ResolvedAt time multiple times which would mean adding this offset over and over again.
            // Incident is resolved and not changed in any way. Then adding timezone does not make sense

            console.log(`Trying so send this data:, ${id}, ${title}, ${severity}, ${serviceName}, ${startedAt}, ${ await ResolvedAtToUTC()}`)

            const response = await fetch(`http://localhost:8080/incident/${incidentID}`, {
                method: 'PATCH',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({id: id, title: title, severity: severity, service_name: serviceName, started_at: startedAt, resolved_at: await ResolvedAtToUTC()})
            })
            if (!response?.ok) {
                const errorData = await response.json()
                console.log("Error response", errorData)
                throw new Error(`- ${errorData.error}`)
            }
            setIsSuccess(true)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Unknown error')
        }
    }

    const modalRef = useRef<HTMLDialogElement>(null)
    const RenderToast = () => {
        if (error) {
            return (
                <ToastNotification duration={10000} message={`Error: ${error}`} toastLevel='alert-error' toastPos='toast-top toast-right' />
            )
        }
        if (isSuccess) {
            return (
                <ToastNotification duration={5000} message='Successfully edited new incident' toastLevel='alert-success' toastPos='toast-top toast-right' />
            )
        }
    }

    const RenderResolvedAtField = () => {
        // When incident is already marked as resolved - the render button that enable this field. Not resolved incidents will get field straight.
        const ResolvedAtField = () => {
            return (
                <>
                <label className="flex items-center label">Resolved at<span className="badge badge-soft badge-xs ml-0.5">Optional</span></label>
                <label className="input mb-4 lg:min-w-[20vw]">
                    <input value={resolvedAt} onChange={(e) => {
                        setResolvedAt(e.target.value)
                    }} type="datetime-local" />
                </label>
                </>
            )
        }

        if (isResolved !== false || originalResolvedAt !== undefined) {
            return (
                <>
                <div className='flex justify-between w-55 items-center'>
                    <div className='flex flex-col'>
                        <span>Change ResolvedAt ?</span>
                        <span className='text-gray-500'>({originalResolvedAt})</span>
                    </div>
                    <input type='checkbox' className='toggle' onChange={() => setEnableResolvedAtOption(!enableResolvedAtOption)} checked={enableResolvedAtOption}></input>
                </div>
                <div className='mt-2'>
                    {enableResolvedAtOption && <ResolvedAtField />}
                </div>
                </>
            )
        } else {
            return <ResolvedAtField />
        }
    }

    return (
        <>
        {RenderToast()}
        <button className='btn btn-xs hover:btn-info btn-ghost' onClick={() => {
            modalRef.current?.showModal()
            fetchIncidentByID()
            {RenderLoading()}
        }}>{editIcon}</button>
        <dialog ref={modalRef} className="modal modal-bottom sm:modal-middle">
            <div className="modal-box w-fit bg-base-300">
                <div className='m-1'>
                    <span className='text-lg font-semibold'>Editing incident </span>
                    <span className='text-lg font-bold bg-linear-to-r from-violet-600 to-blue-500 text-transparent bg-clip-text'>{incidentID}</span>
                </div>
                <fieldset className="fieldset rounded-box shadow p-4 bg-base-200">
                    <div className="flex flex-col">
                        {/* <label className="label">Incident name</label>
                        <input value={id} onChange={(e) => setID(e.target.value)} type="text" className="input mb-4 lg:min-w-[20vw] bg-base-100" disabled/> */}

                        <label className="label">Title</label>
                        <input value={title} onChange={(e) => setTitle(e.target.value)} type="text" className="input mb-4 lg:min-w-[20vw]" />

                        <label className="label">Severity</label>
                        <label className="select mb-4 lg:min-w-[20vw]">
                            <select value={severity} onChange={(e) => setSeverity(e.target.value)}>
                                <option value='critical'>Critical</option>
                                <option value='high'>High</option>
                                <option value='medium'>Medium</option>
                                <option value='low'>Low</option>
                            </select>
                        </label>

                        <label className="label">Service Name</label>
                        <input value={serviceName} onChange={(e) => setServiceName(e.target.value)} type="text" className="input mb-4 lg:min-w-[20vw]" />

                        {/* It's better to disable this so users cant change when the incident started to remain objective. */}
                        {/* <label className="label">Started at</label>
                        <label className="input mb-4 lg:min-w-[20vw] bg-base-100">
                            <input value={startedAt} onChange={(e) => setStartedAt(e.target.value)} type="datetime-local" disabled />
                        </label> */}
                        <RenderResolvedAtField />
                    </div>
                    <div className="flex flex-col items-center">
                        {/* Currently this button does nothing, because onClick event handler that calls function is not ready yet. Yes, i could do "editing incidents" in a way that i would delete old one and recreate new one - i can do this. But i would like to learn HTTP PATCH or HTTP PUT method. This should be the correct way - but this would mean chanding backend in some way - dont have knowledge for that. */}
                        <button onClick={() => {
                            patchIncident()
                            modalRef.current?.close()
                        }} className={`btn ${submitBtnBg} border-base-content mt-2 px-6`} onMouseEnter={() => setSubmitBtnBg("btn-info")} onMouseLeave={() => setSubmitBtnBg("btn-neutral")}>Edit</button>
                    </div>
                </fieldset>
                <div className="modal-action">
                    <form method="dialog">
                        {/* if there is a button in form, it will close the modal */}
                        <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
                        {/* <span>Close</span> */}
                    </form>
                </div>
            </div>
        </dialog>
        </>
    )
}