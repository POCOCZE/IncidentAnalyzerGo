import { useState } from 'react'

const HomePage = () => {
    // const [loading, setLoading] = useState<boolean>(true)

    // // Setting loading to false for now. But there will be components here that requires connecting to backend.
    // setLoading(false)
    // if (loading) {
    //     return (
    //         <div className="flex w-52 justify-center items-center flex-col gap-4 p-4 mt-8">
    //             <div className="skeleton h-45 w-full bg-base-200"></div>
    //             <div className="skeleton h-4 w-28 bg-base-200"></div>
    //             <div className="skeleton h-4 w-40 bg-base-200"></div>
    //             <div className="skeleton h-4 w-34 bg-base-200"></div>
    //         </div>
    //     )
    // }

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
                        <div className="flex w-45 justify-center items-center flex-col gap-4 p-2 m-4 bg-base-200 rounded-xl shadow">
                            <span className='text-lg font-semibold'>Download report</span>
                            <div className="skeleton h-10 w-full bg-base-200"></div>
                            <div className="skeleton h-4 w-28 bg-base-200"></div>
                            <div className="skeleton h-4 w-40 bg-base-200"></div>
                            <div className="skeleton h-4 w-34 bg-base-200"></div>
                        </div>
                        <div className="flex w-45 justify-center items-center flex-col gap-4 p-2 m-4 bg-base-200 rounded-xl shadow">
                            <span className='text-lg font-semibold'>Export incidents</span>
                            <div className="skeleton h-10 w-full bg-base-200"></div>
                            <div className="skeleton h-4 w-28 bg-base-200"></div>
                            <div className="skeleton h-4 w-40 bg-base-200"></div>
                            <div className="skeleton h-4 w-34 bg-base-200"></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default HomePage