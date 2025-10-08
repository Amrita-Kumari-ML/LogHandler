import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Button, Spinner, Alert } from 'react-bootstrap';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  PointElement,
  LineElement,
} from 'chart.js';
import { Bar, Pie, Line } from 'react-chartjs-2';
import axios from 'axios';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  PointElement,
  LineElement
);

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:7002';

function App() {
  const [dashboardData, setDashboardData] = useState(null);
  const [statusStats, setStatusStats] = useState(null);
  const [ipStats, setIpStats] = useState(null);
  const [timeStats, setTimeStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const [dashboardRes, statusRes, ipRes, timeRes] = await Promise.all([
        axios.get(`${API_BASE_URL}/stats/dashboard`),
        axios.get(`${API_BASE_URL}/stats/status`),
        axios.get(`${API_BASE_URL}/stats/ip`),
        axios.get(`${API_BASE_URL}/stats/time`)
      ]);

      setDashboardData(dashboardRes.data.data);
      setStatusStats(statusRes.data.data);
      setIpStats(ipRes.data.data);
      setTimeStats(timeRes.data.data);
    } catch (err) {
      setError('Failed to fetch data from the API');
      console.error('Error fetching data:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, []);

  const statusChartData = {
    labels: statusStats?.map(stat => `Status ${stat.status}`) || [],
    datasets: [
      {
        label: 'Request Count',
        data: statusStats?.map(stat => stat.count) || [],
        backgroundColor: [
          'rgba(75, 192, 192, 0.6)',
          'rgba(255, 99, 132, 0.6)',
          'rgba(255, 205, 86, 0.6)',
          'rgba(54, 162, 235, 0.6)',
        ],
        borderColor: [
          'rgba(75, 192, 192, 1)',
          'rgba(255, 99, 132, 1)',
          'rgba(255, 205, 86, 1)',
          'rgba(54, 162, 235, 1)',
        ],
        borderWidth: 1,
      },
    ],
  };

  const ipChartData = {
    labels: ipStats?.map(stat => stat.ip_address) || [],
    datasets: [
      {
        label: 'Requests per IP',
        data: ipStats?.map(stat => stat.request_count) || [],
        backgroundColor: 'rgba(153, 102, 255, 0.6)',
        borderColor: 'rgba(153, 102, 255, 1)',
        borderWidth: 1,
      },
    ],
  };

  const avgBytesChartData = {
    labels: statusStats?.map(stat => `Status ${stat.status}`) || [],
    datasets: [
      {
        label: 'Average Bytes',
        data: statusStats?.map(stat => stat.avg_bytes) || [],
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        tension: 0.1,
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top',
      },
    },
  };

  if (loading && !dashboardData) {
    return (
      <Container className="mt-5">
        <div className="loading-spinner">
          <Spinner animation="border" role="status">
            <span className="visually-hidden">Loading...</span>
          </Spinner>
        </div>
      </Container>
    );
  }

  return (
    <div className="App">
      <div className="dashboard-header">
        <Container>
          <Row>
            <Col>
              <h1 className="text-center mb-0">Log Analytics Dashboard</h1>
              <p className="text-center mb-0">Real-time monitoring and statistics</p>
            </Col>
          </Row>
        </Container>
      </div>

      <Container>
        {error && (
          <Alert variant="danger" className="mb-4">
            {error}
          </Alert>
        )}

        <Row className="mb-4">
          <Col className="text-center">
            <Button 
              className="refresh-btn" 
              onClick={fetchData} 
              disabled={loading}
            >
              {loading ? (
                <>
                  <Spinner
                    as="span"
                    animation="border"
                    size="sm"
                    role="status"
                    aria-hidden="true"
                    className="me-2"
                  />
                  Refreshing...
                </>
              ) : (
                'Refresh Data'
              )}
            </Button>
          </Col>
        </Row>

        {/* Summary Cards */}
        <Row>
          <Col md={3} sm={6}>
            <Card className="stat-card text-center">
              <Card.Body>
                <h3 className="text-primary">{dashboardData?.total_logs || 0}</h3>
                <p className="mb-0">Total Logs</p>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3} sm={6}>
            <Card className="stat-card text-center">
              <Card.Body>
                <h3 className="text-success">{dashboardData?.unique_ips || 0}</h3>
                <p className="mb-0">Unique IPs</p>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3} sm={6}>
            <Card className="stat-card text-center">
              <Card.Body>
                <h3 className="text-warning">
                  {dashboardData?.avg_response_size ? 
                    Math.round(dashboardData.avg_response_size) : 0}
                </h3>
                <p className="mb-0">Avg Response Size</p>
              </Card.Body>
            </Card>
          </Col>
          <Col md={3} sm={6}>
            <Card className="stat-card text-center">
              <Card.Body>
                <h3 className="text-info">
                  {dashboardData?.last_log_time ? 
                    new Date(dashboardData.last_log_time).toLocaleTimeString() : 'N/A'}
                </h3>
                <p className="mb-0">Last Log</p>
              </Card.Body>
            </Card>
          </Col>
        </Row>

        {/* Charts */}
        <Row>
          <Col lg={6}>
            <div className="chart-container">
              <h5 className="mb-3">HTTP Status Codes Distribution</h5>
              <div style={{ height: '300px' }}>
                <Pie data={statusChartData} options={chartOptions} />
              </div>
            </div>
          </Col>
          <Col lg={6}>
            <div className="chart-container">
              <h5 className="mb-3">Requests by IP Address</h5>
              <div style={{ height: '300px' }}>
                <Bar data={ipChartData} options={chartOptions} />
              </div>
            </div>
          </Col>
        </Row>

        <Row>
          <Col lg={12}>
            <div className="chart-container">
              <h5 className="mb-3">Average Response Size by Status Code</h5>
              <div style={{ height: '300px' }}>
                <Line data={avgBytesChartData} options={chartOptions} />
              </div>
            </div>
          </Col>
        </Row>

        {/* Top IPs and Status Codes Tables */}
        <Row>
          <Col lg={6}>
            <div className="chart-container">
              <h5 className="mb-3">Top IP Addresses</h5>
              <div className="table-responsive">
                <table className="table table-striped">
                  <thead>
                    <tr>
                      <th>IP Address</th>
                      <th>Requests</th>
                    </tr>
                  </thead>
                  <tbody>
                    {dashboardData?.top_ips?.map((ip, index) => (
                      <tr key={index}>
                        <td>{ip.ip}</td>
                        <td>{ip.count}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </Col>
          <Col lg={6}>
            <div className="chart-container">
              <h5 className="mb-3">Top Status Codes</h5>
              <div className="table-responsive">
                <table className="table table-striped">
                  <thead>
                    <tr>
                      <th>Status Code</th>
                      <th>Count</th>
                    </tr>
                  </thead>
                  <tbody>
                    {dashboardData?.top_status_codes?.map((status, index) => (
                      <tr key={index}>
                        <td>{status.status}</td>
                        <td>{status.count}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </Col>
        </Row>
      </Container>
    </div>
  );
}

export default App;
