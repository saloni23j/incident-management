import axios from "axios";

const BASE_URL = "http://localhost:8080"; // Change if needed

export const fetchIncidents = () => axios.get(`${BASE_URL}/incidents`);

export const createIncident = (data: {
  title: string;
  description: string;
  status?: string;
  priority?: string;
}) => axios.post(`${BASE_URL}/incidents`, data);
