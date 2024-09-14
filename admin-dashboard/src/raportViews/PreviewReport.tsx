import React from "react";
import { ViewModel, ResponseMessage } from "../controllers/Types";

import { useEffect, useState } from "react";
import { ViewsController } from "../controllers/ViewsController";

function PreviewReport({ viewName }: { viewName: string }) {
  const [viewResult, setViewResult] = useState<Map<string, any>[]>([]);

  useEffect(() => {
    downloadReportView();
    console.log("loading");
  }, []);

  const downloadReportView = async () => {
    let model: ViewModel = {
      viewName: viewName,
    };
    const response: ResponseMessage = await ViewsController.GetView(model);
    setViewResult(response.data);
  };

  return (
    <div>
      {/* Table */}
      <div className="text-black font-semibold">View raport / {viewName}</div>
      {viewResult.length > 0 && (
        <div className="relative overflow-x-auto text-sm">
          <table className="w-full bg-white">
            <thead className="bg-white">
              <tr>
                {Object.keys(viewResult[0]).map((name, index) => (
                  <th
                    key={index + name}
                    className="py-2 px-4 text-black font-semibold border-b border-gray-200 text-left min-w-[13rem]"
                  >
                    {name}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {Object.values(viewResult).map((item, index) => (
                <tr
                  key={index}
                  className="hover:bg-slate-200 transition ease-in-out duration-300 "
                >
                  {Object.keys(viewResult[0]).map((key, i) => (
                    <td
                      key={i}
                      className="py-2 px-4 text-black font-semibold border-b border-gray-200 text-left max-w-[12rem] overflow-x-hidden"
                    >
                      {item[key] as string}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default PreviewReport;
