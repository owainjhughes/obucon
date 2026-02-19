import axios from "axios";

const apiBaseUrl = process.env.REACT_APP_API_URL || "http://localhost:8080";

export const apiClient = axios.create({
  baseURL: apiBaseUrl,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});
