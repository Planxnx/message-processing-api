import axios from "axios";

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_URL || "http://localhost:8080",
  timeout: 1000 * 15,
});

export default axiosInstance;
