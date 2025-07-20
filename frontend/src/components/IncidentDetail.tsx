export default function IncidentDetail({ incident, onClose }: { incident: any; onClose: () => void }) {
  return (
    <div className="p-4 mt-4 border bg-white rounded shadow">
      <h3 className="font-bold text-lg mb-2">Incident Details</h3>
      <p><strong>Title:</strong> {incident.title}</p>
      <p><strong>Description:</strong> {incident.description}</p>
      <p><strong>Status:</strong> {incident.status}</p>
      <p><strong>Priority:</strong> {incident.priority}</p>
      <p><strong>AI Severity:</strong> {incident.ai_severity}</p>
      <p><strong>AI Category:</strong> {incident.ai_category}</p>
      <button onClick={onClose} className="mt-4 text-red-500 underline">
        Close
      </button>
    </div>
  );
}
