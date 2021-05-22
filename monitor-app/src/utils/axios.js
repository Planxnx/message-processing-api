import axios from "axios";

const axiosInstance = axios.create({
  baseURL:
    process.env.REACT_APP_API_URL ||
    "https://26e6ccc2a0a8.ngrok.io" ||
    "http://localhost:8080",
  timeout: 1000 * 15,
});

export default axiosInstance;
