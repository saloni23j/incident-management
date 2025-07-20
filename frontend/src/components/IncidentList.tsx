import { useEffect, useState } from "react";
import { fetchIncidents } from "../api";
import IncidentDetail from "./IncidentDetail";

export default function IncidentList() {
  const [incidents, setIncidents] = useState([]);
  const [selected, setSelected] = useState(null);

  const loadIncidents = async () => {
    const res = await fetchIncidents();
    setIncidents(res.data);
  };

  useEffect(() => {
    loadIncidents();
  }, []);

  return (
    <div className="p-4">
      <h2 className="text-lg font-bold">All Incidents</h2>
      <ul className="space-y-2">
        {incidents.map((incident: any) => (
          <li key={incident.id} className="p-4 border rounded bg-gray-100">
            <div className="flex justify-between">
              <div>
                <p className="font-semibold">{incident.title}</p>
                <p className="text-sm">{incident.ai_severity} | {incident.ai_category}</p>
              </div>
              <button
                className="text-blue-600 underline"
                onClick={() => setSelected(incident)}
              >
                View Details
              </button>
            </div>
          </li>
        ))}
      </ul>

      {selected && (
        <IncidentDetail incident={selected} onClose={() => setSelected(null)} />
      )}
    </div>
  );
}
