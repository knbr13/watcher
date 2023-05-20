const OutputWindow = ({ outputDetails }) => {
    const getOutput = () => {
      let statusId = outputDetails?.status?.id;
  
      if (statusId === 6) {
        return (
          <pre className="px-2 py-1 font-normal text-xs text-red-500">
            {window.atob(outputDetails?.compile_output)}
          </pre>
        );
      } else if (statusId === 3) {
        return (
          <pre className="px-2 py-1 font-normal text-xs text-green-500">
            {window.atob(outputDetails.stdout) !== null
              ? `${window.atob(outputDetails.stdout)}`
              : null}
          </pre>
        );
      } else if (statusId === 5) {
        return (
          <pre className="px-2 py-1 font-normal text-xs text-red-500">
            {`Time Limit Exceeded`}
          </pre>
        );
      } else {
        return (
          <pre className="px-2 py-1 font-normal text-xs text-red-500">
            {window.atob(outputDetails?.stderr)}
          </pre>
        );
      }
    };
    return (
      <>
        <div className="w-full h-56 bg-[#1e293b] rounded-md text-white font-normal text-sm overflow-y-auto">
          {outputDetails ? <>{getOutput()}</> : "Output: ..."}
        </div>
      </>
    );
  };
  
  export default OutputWindow;