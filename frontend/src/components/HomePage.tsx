const HomePage = () => {
    return (
        <div className="flex flex-col grow rounded-xl justify-start items-center h-[97vh] m-4 bg-base-300">
            <span className='text-center text-3xl mt-4 font-bold bg-linear-to-r from-gray-400 to-gray-500 text-transparent bg-clip-text'>Coming Soon!</span>
            <div className="flex w-52 justify-center items-center flex-col gap-4 p-4 mt-8">
                <div className="skeleton h-45 w-full bg-base-200"></div>
                <div className="skeleton h-4 w-28 bg-base-200"></div>
                <div className="skeleton h-4 w-40 bg-base-200"></div>
                <div className="skeleton h-4 w-34 bg-base-200"></div>
            </div>
        </div>
    )
}

export default HomePage