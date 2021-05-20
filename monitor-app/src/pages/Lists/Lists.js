import { useEffect, useState } from "react";
import { Navbar, Table, Badge } from "react-bootstrap";
import axiosInstance from "../../utils/axios";

const getAPIsHealth = async () => {
  const { data } = await axiosInstance.get("/health");
  return data.data;
};

const ListsPage = () => {
  const [apisHealth, setAPIsHealth] = useState([]);
  useEffect(() => {
    (async () => {
      const lists = await getAPIsHealth();
      if (lists) {
        setAPIsHealth(lists);
      }
    })();
  }, []);

  return (
    <div>
      <Navbar bg="dark" variant="dark">
        <Navbar.Brand>API Monitor</Navbar.Brand>
      </Navbar>
      <Table hover responsive>
        <thead>
          <tr>
            <th>Feature</th>
            <th> </th>
            <th>Mode</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {apisHealth.map((data) => (
            <tr>
              <td colSpan="2">{data.feature}</td>
              <td>{data.executeMode.join(", ")}</td>
              <td>
                {data.status ? (
                  <Badge variant="primary">Online</Badge>
                ) : (
                  <Badge variant="danger">Offline</Badge>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
  );
};

export default ListsPage;
