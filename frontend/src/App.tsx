import { useEffect, useState } from "react";

interface Incident {
  id: number;
  title: string;
  description: string;
  service: string;
  ai_severity: string;
  ai_category: string;
}

function App() {
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [service, setService] = useState("");

  const fetchIncidents = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/v1/incidents");
      const data = await res.json();
      setIncidents(data);
    } catch (err) {
      console.error("Failed to fetch incidents", err);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/api/v1/incidents", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title, description, service }),
    });

    if (res.ok) {
      setTitle("");
      setDescription("");
      setService("");
      fetchIncidents();
    } else {
      alert("Failed to submit incident");
    }
  };

  useEffect(() => {
    fetchIncidents();
  }, []);

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-4">Incident Reporting</h2>
      <form onSubmit={handleSubmit} className="space-y-4 mb-6">
        <input
          type="text"
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          className="w-full border p-2 rounded"
        />
        <textarea
          placeholder="Description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          required
          className="w-full border p-2 rounded"
        />
        <input
          type="text"
          placeholder="Affected Service"
          value={service}
          onChange={(e) => setService(e.target.value)}
          required
          className="w-full border p-2 rounded"
        />
        <button
          type="submit"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Submit Incident
        </button>
      </form>

      <h2 className="text-xl font-semibold mb-2">All Incidents</h2>
      {incidents.length === 0 ? (
        <p>No incidents found.</p>
      ) : (
        <div className="space-y-4">
          {incidents.map((incident) => (
            <div
              key={incident.id}
              className="border rounded p-4 shadow-sm bg-white"
            >
              <h3 className="text-lg font-semibold">{incident.title}</h3>
              <p>{incident.description}</p>
              <p>
                <strong>Service:</strong> {incident.service}
              </p>
              <p>
                <strong>Severity:</strong>{" "}
                {incident.ai_severity || "N/A"}
              </p>
              <p>
                <strong>Category:</strong>{" "}
                {incident.ai_category || "N/A"}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default App;
