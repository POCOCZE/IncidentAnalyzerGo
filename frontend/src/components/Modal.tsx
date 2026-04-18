import { useRef } from 'react'

interface ModalProps {
    buttonText: React.ReactNode
    title: React.ReactNode
    description: React.ReactNode
    closeButtonText: React.ReactNode
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