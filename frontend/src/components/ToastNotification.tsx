import { useEffect, useState } from 'react'

interface Notification {
    duration: number
    message: string
    toastLevel: string
    toastPos: string
}

const ToastNotification = ({duration, message, toastLevel, toastPos}:Notification) => {
    const [isToastVisible, setIsToastVisible] = useState<boolean>(true)

    useEffect(() => {
        const timer = setTimeout(() => {
            setIsToastVisible(false)
        }, duration)

        return () => clearTimeout(timer)
    })

    if (!isToastVisible) return null

    return (
        <div className={`toast ${toastPos}`}>
            {/* <div className={`alert ${toastLevel}`}> */}
            {/* {toastLevel ? 'alert-success' : } */}
            <div className={`alert ${toastLevel ? toastLevel : 'alert-success'}`}>
                <span>{message}</span>
                <button onClick={() => {setIsToastVisible(false)}}>Dismiss</button>
            </div>
        </div>
    )
}

export default ToastNotification