// API Configuration
// Docker container ports
const API_BASE_URL = '/api/parser'; // proxied to logparser:8082 by nginx
const GENERATOR_API_URL = '/api/generator'; // proxied to loggenerate:8081 by nginx
const WINDOW_SIZE = 5; // Max items to show per chart/table window

// Chart instances
let statusChart, ipChart, avgBytesChart, timeChart;
let anomalyChart, trendChart, clusterChart, predictionChart;

// Global state
let autoRefreshInterval = null;
let generationInterval = null;
let isGenerating = false;
let previousData = {};
let currentPage = 1;
let logsPerPage = 50;
// Cursor-based pagination state for Log Viewer
let currentCursorToken = null; // e.g., "2025-10-07T12:34:56Z&id=123"
let cursorHistory = []; // stack of tokens to support Prev navigation
let lastPaging = null; // latest paging object from API

// Pagination helpers for Log Viewer
function resetPagination() {
    currentPage = 1;
    currentCursorToken = null;
    cursorHistory = [];
    updatePaginationUI(null);
}

function updatePaginationUI(paging) {
    const currentPageEl = document.getElementById('current-page');
    if (currentPageEl) currentPageEl.textContent = String(currentPage);

    const prevLi = document.getElementById('prev-page')?.parentElement;
    const nextLi = document.getElementById('next-page')?.parentElement;

    if (prevLi) prevLi.classList.toggle('disabled', cursorHistory.length === 0);
    const hasNext = !!(paging && paging.next_cursor);
    if (nextLi) nextLi.classList.toggle('disabled', !hasNext);
}

function onNextPageClick(e) {
    e.preventDefault();
    if (!lastPaging || !lastPaging.next_cursor) return;
    cursorHistory.push(currentCursorToken);
    currentCursorToken = lastPaging.next_cursor;
    currentPage += 1;
    fetchLogs();
}

function onPrevPageClick(e) {
    e.preventDefault();
    if (cursorHistory.length === 0) return;
    currentCursorToken = cursorHistory.pop() || null;
    currentPage = Math.max(1, currentPage - 1);
    fetchLogs();
}

// Initialize dashboard
document.addEventListener('DOMContentLoaded', function() {
    initializeCharts();
    initializeEventListeners();

    // Theme: apply saved or system preference
    const themeSelect = document.getElementById('themeSelect');
    const savedThemeMode = localStorage.getItem('themeMode') || 'system';
    if (themeSelect) {
        themeSelect.value = savedThemeMode;
        themeSelect.addEventListener('change', (e) => {
            const mode = e.target.value;
            localStorage.setItem('themeMode', mode);
            applyTheme(mode);
        });
    }
    applyTheme(savedThemeMode);

    // Default to live updates enabled and fast interval
    const autoToggle = document.getElementById('autoRefreshToggle');
    const intervalSel = document.getElementById('refreshInterval');
    if (autoToggle) autoToggle.checked = true;
    if (intervalSel) intervalSel.value = '5';

    fetchData();
    setupAutoRefresh();
    checkInitialGeneratorStatus();
    initializeMLFeatures();

    // Initial load of logs based on selected parameters
    resetPagination();
    fetchLogs();

    // Pause/resume live updates when tab visibility changes
    document.addEventListener('visibilitychange', handleVisibilityChange);

    // React to system theme changes when in system mode
    if (window.matchMedia) {
        const media = window.matchMedia('(prefers-color-scheme: dark)');
        try {
            media.addEventListener('change', () => {
                const mode = localStorage.getItem('themeMode') || 'system';
                if (mode === 'system') applyTheme('system');
            });
        } catch (e) {
            media.addListener && media.addListener(() => {
                const mode = localStorage.getItem('themeMode') || 'system';
                if (mode === 'system') applyTheme('system');
            });
        }
    }
});

// Initialize event listeners
function initializeEventListeners() {
    // Refresh controls
    document.getElementById('refresh-btn').addEventListener('click', fetchData);
    document.getElementById('autoRefreshToggle').addEventListener('change', handleAutoRefreshToggle);
    document.getElementById('refreshInterval').addEventListener('change', setupAutoRefresh);

    // Time unit change handler removed (no longer needed)

    // Log generation form
    document.getElementById('log-generation-form').addEventListener('submit', handleLogGeneration);
    document.getElementById('stop-generation').addEventListener('click', stopLogGeneration);

    // Log viewer controls
    document.getElementById('refresh-logs').addEventListener('click', (e) => { e.preventDefault(); resetPagination(); fetchLogs(); });
    document.getElementById('logLimit').addEventListener('change', () => { resetPagination(); fetchLogs(); });
    document.getElementById('statusFilter').addEventListener('change', () => { resetPagination(); fetchLogs(); });
    // Pagination buttons
    const prevBtn = document.getElementById('prev-page');
    const nextBtn = document.getElementById('next-page');
    if (prevBtn) prevBtn.addEventListener('click', onPrevPageClick);
    if (nextBtn) nextBtn.addEventListener('click', onNextPageClick);

    // Tab change handler
    document.querySelectorAll('[data-bs-toggle="tab"]').forEach(tab => {
        tab.addEventListener('shown.bs.tab', handleTabChange);
    });
}

// Initialize empty charts
function initializeCharts() {
    // Status Chart (Pie)
    const statusCtx = document.getElementById('statusChart').getContext('2d');
    statusChart = new Chart(statusCtx, {
        type: 'pie',
        data: {
            labels: [],
            datasets: [{
                data: [],
                backgroundColor: [
                    'rgba(75, 192, 192, 0.8)',
                    'rgba(255, 99, 132, 0.8)',
                    'rgba(255, 205, 86, 0.8)',
                    'rgba(54, 162, 235, 0.8)',
                    'rgba(153, 102, 255, 0.8)',
                ],
                borderColor: [
                    'rgba(75, 192, 192, 1)',
                    'rgba(255, 99, 132, 1)',
                    'rgba(255, 205, 86, 1)',
                    'rgba(54, 162, 235, 1)',
                    'rgba(153, 102, 255, 1)',
                ],
                borderWidth: 2
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'top',
                }
            }
        }
    });

    // IP Chart (Bar)
    const ipCtx = document.getElementById('ipChart').getContext('2d');
    ipChart = new Chart(ipCtx, {
        type: 'bar',
        data: {
            labels: [],
            datasets: [{
                label: 'Requests per IP',
                data: [],
                backgroundColor: 'rgba(153, 102, 255, 0.8)',
                borderColor: 'rgba(153, 102, 255, 1)',
                borderWidth: 2,
                borderRadius: 5,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'top',
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    grid: {
                        color: 'rgba(0, 0, 0, 0.1)'
                    }
                },
                x: {
                    grid: {
                        display: false
                    }
                }
            }
        }
    });

    // Average Bytes Chart (Line)
    const avgBytesCtx = document.getElementById('avgBytesChart').getContext('2d');
    avgBytesChart = new Chart(avgBytesCtx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [{
                label: 'Average Bytes',
                data: [],
                borderColor: 'rgb(255, 99, 132)',
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                tension: 0.4,
                fill: true,
                pointBackgroundColor: 'rgb(255, 99, 132)',
                pointBorderColor: '#fff',
                pointBorderWidth: 2,
                pointRadius: 5,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'top',
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    grid: {
                        color: 'rgba(0, 0, 0, 0.1)'
                    }
                },
                x: {
                    grid: {
                        display: false
                    }
                }
            }
        }
    });

    // Time Chart (Area)
    const timeCtx = document.getElementById('timeChart').getContext('2d');
    timeChart = new Chart(timeCtx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [{
                label: 'Logs per Hour',
                data: [],
                borderColor: 'rgb(54, 162, 235)',
                backgroundColor: 'rgba(54, 162, 235, 0.2)',
                tension: 0.4,
                fill: true,
                pointBackgroundColor: 'rgb(54, 162, 235)',
                pointBorderColor: '#fff',
                pointBorderWidth: 2,
                pointRadius: 4,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'top',
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    grid: {
                        color: 'rgba(0, 0, 0, 0.1)'
                    }
                },
                x: {
                    grid: {
                        display: false
                    }
                }
            }
        }
    });

    initializeMLCharts();
}

// Initialize ML Charts
function initializeMLCharts() {
    // Anomaly Chart
    const anomalyCtx = document.getElementById('anomalyChart').getContext('2d');
    anomalyChart = new Chart(anomalyCtx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [{
                label: 'Request Rate',
                data: [],
                borderColor: 'rgb(75, 192, 192)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                tension: 0.1
            }, {
                label: 'Anomalies',
                data: [],
                borderColor: 'rgb(255, 99, 132)',
                backgroundColor: 'rgba(255, 99, 132, 0.8)',
                pointRadius: 8,
                pointHoverRadius: 10,
                showLine: false
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            aspectRatio: 1,
            scales: { y: { beginAtZero: true } }
        }
    });

    // Trend Chart
    const trendCtx = document.getElementById('trendChart').getContext('2d');
    trendChart = new Chart(trendCtx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [{
                label: 'Actual',
                data: [],
                borderColor: 'rgb(54, 162, 235)',
                backgroundColor: 'rgba(54, 162, 235, 0.2)',
                tension: 0.1
            }, {
                label: 'Trend',
                data: [],
                borderColor: 'rgb(255, 205, 86)',
                backgroundColor: 'rgba(255, 205, 86, 0.2)',
                borderDash: [5, 5],
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            aspectRatio: 1,
            scales: { y: { beginAtZero: true } }
        }
    });

    // Cluster Chart
    const clusterCtx = document.getElementById('clusterChart').getContext('2d');
    clusterChart = new Chart(clusterCtx, {
        type: 'doughnut',
        data: {
            labels: ['Heavy Users', 'Regular Users', 'Light Users'],
            datasets: [{
                data: [0, 0, 0],
                backgroundColor: [
                    'rgba(255, 99, 132, 0.8)',
                    'rgba(255, 205, 86, 0.8)',
                    'rgba(75, 192, 192, 0.8)'
                ]
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            aspectRatio: 1
        }
    });

    // Prediction Chart
    const predictionCtx = document.getElementById('predictionChart').getContext('2d');
    predictionChart = new Chart(predictionCtx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [
                {
                    label: 'Predicted',
                    data: [],
                    borderColor: 'rgb(255, 99, 132)',
                    backgroundColor: 'rgba(255, 99, 132, 0.15)',
                    tension: 0.1
                },
                {
                    label: 'Upper Bound',
                    data: [],
                    borderColor: 'rgba(255, 99, 132, 0.5)',
                    borderDash: [5, 5],
                    fill: false,
                    tension: 0.1
                },
                {
                    label: 'Lower Bound',
                    data: [],
                    borderColor: 'rgba(54, 162, 235, 0.5)',
                    borderDash: [5, 5],
                    fill: false,
                    tension: 0.1
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            aspectRatio: 1,
            scales: { y: { beginAtZero: true } }
        }
    });
}

// Fetch data from API
async function fetchData() {
    setLoading(true);
    hideError();

    try {
        const [dashboardRes, statusRes, ipRes, timeRes] = await Promise.all([
            fetch(`${API_BASE_URL}/stats/dashboard`),
            fetch(`${API_BASE_URL}/stats/status`),
            fetch(`${API_BASE_URL}/stats/ip`),
            fetch(`${API_BASE_URL}/stats/time`)
        ]);

        if (!dashboardRes.ok || !statusRes.ok || !ipRes.ok || !timeRes.ok) {
            throw new Error('Failed to fetch data from API');
        }

        const dashboardData = await dashboardRes.json();
        const statusData = await statusRes.json();
        const ipData = await ipRes.json();
        const timeData = await timeRes.json();

        // Pass full time payload (groupBy + data array)
        const timePayload = timeData && timeData.data ? timeData.data : { groupBy: 'hour', data: [] };
        updateDashboard(dashboardData.data, statusData.data, ipData.data, timePayload);
        updateMLAnalytics(dashboardData.data, statusData.data, ipData.data);
    } catch (error) {
        showError('Failed to fetch data: ' + error.message);
        console.error('Error fetching data:', error);
    } finally {
        setLoading(false);
    }
}

// Update dashboard with new data
function updateDashboard(dashboardData, statusData, ipData, timeData) {
    // Calculate trends
    const trends = calculateTrends(dashboardData);

    // Update summary cards with trends
    document.getElementById('total-logs').textContent = formatNumber(dashboardData.total_logs || 0);
    document.getElementById('unique-ips').textContent = dashboardData.unique_ips || 0;
    document.getElementById('avg-response-size').textContent =
        dashboardData.avg_response_size ? Math.round(dashboardData.avg_response_size) + ' bytes' : '0 bytes';
    document.getElementById('last-log-time').textContent =
        dashboardData.last_log_time ? new Date(dashboardData.last_log_time).toLocaleTimeString() : 'N/A';

    // Update trend indicators
    document.getElementById('logs-trend').textContent = trends.logs;
    document.getElementById('ips-trend').textContent = trends.ips;
    document.getElementById('size-trend').textContent = trends.size;
    document.getElementById('time-ago').textContent = getTimeAgo(dashboardData.last_log_time);

    // Update charts with fixed window size
    // Status chart - top WINDOW_SIZE by count
    if (statusData && statusData.length > 0) {
        const statusTop = [...statusData]
            .sort((a, b) => (b.count || 0) - (a.count || 0))
            .slice(0, WINDOW_SIZE);
        statusChart.data.labels = statusTop.map(stat => `${stat.status} (${stat.count})`);
        statusChart.data.datasets[0].data = statusTop.map(stat => stat.count);
        statusChart.update();

        // Average bytes chart aligned with the same top statuses
        avgBytesChart.data.labels = statusTop.map(stat => `Status ${stat.status}`);
        avgBytesChart.data.datasets[0].data = statusTop.map(stat => stat.avg_bytes || 0);
        avgBytesChart.update();
    }

    // IP chart - top WINDOW_SIZE by requests
    if (ipData && ipData.length > 0) {
        const ipTop = [...ipData]
            .sort((a, b) => ((b.request_count || b.count || 0) - (a.request_count || a.count || 0)))
            .slice(0, WINDOW_SIZE);
        ipChart.data.labels = ipTop.map(stat => stat.ip_address || stat.ip);
        ipChart.data.datasets[0].data = ipTop.map(stat => stat.request_count || stat.count || 0);
        ipChart.update();
    }

    // Time chart - support backend payload { groupBy, data: [{ time_unit, request_count, avg_bytes }] }
    if (timeData) {
        const groupBy = timeData.groupBy || timeData.group_by || 'hour';
        let series = Array.isArray(timeData) ? timeData : (timeData.data || []);
        if (!Array.isArray(series)) series = [];
        if (series.length > 0) {
            const last = series.slice(-WINDOW_SIZE);
            const labels = last.map(item => {
                const tu = item.time_unit;
                if (groupBy === 'hour') {
                    const h = typeof tu === 'number' ? tu : parseInt(tu, 10);
                    return `${String(isNaN(h) ? tu : h).padStart(2, '0')}:00`;
                }
                // day/month or timestamp
                try { return new Date(tu).toLocaleDateString(); } catch { return String(tu); }
            });
            const values = last.map(item => item.request_count || item.count || 0);
            timeChart.data.labels = labels;
            timeChart.data.datasets[0].data = values;
            timeChart.update();
        }
    }

    // Update tables
    updateTopIPsTable(ipData || []);
    updateTopStatusTable(statusData || []);

    // Store current data for trend calculation
    previousData = { ...dashboardData };
}

// Update ML Analytics
function updateMLAnalytics(dashboardData, statusData, ipData) {
    // Anomaly Detection
    updateAnomalyDetection(dashboardData);

    // Trend Analysis
    updateTrendAnalysis(dashboardData);

    // User Clustering
    updateUserClustering(ipData);

    // Security Risk Assessment
    updateSecurityRisk(statusData);

    // Predictive Analytics
    updatePredictiveAnalytics(dashboardData);
}

// Update top IPs table
function updateTopIPsTable(topIps) {
    const tbody = document.getElementById('top-ips-table');
    tbody.innerHTML = '';

    const list = [...topIps]
        .sort((a, b) => ((b.request_count || b.count || 0) - (a.request_count || a.count || 0)))
        .slice(0, WINDOW_SIZE);

    list.forEach(ip => {
        const row = tbody.insertRow();
        row.insertCell(0).textContent = ip.ip_address || ip.ip;
        row.insertCell(1).textContent = formatNumber(ip.request_count || ip.count || 0);
        row.insertCell(2).textContent = ip.avg_bytes ? Math.round(ip.avg_bytes) + ' bytes' : 'N/A';
    });
}

// Update top status codes table
function updateTopStatusTable(topStatus) {
    const tbody = document.getElementById('top-status-table');
    tbody.innerHTML = '';

    const top = [...topStatus]
        .sort((a, b) => (b.count || 0) - (a.count || 0))
        .slice(0, WINDOW_SIZE);

    const total = top.reduce((sum, status) => sum + (status.count || 0), 0);

    top.forEach(status => {
        const row = tbody.insertRow();
        const cnt = status.count || 0;
        const percentage = total > 0 ? ((cnt / total) * 100).toFixed(1) : 0;
        row.insertCell(0).textContent = status.status;
        row.insertCell(1).textContent = formatNumber(cnt);
        row.insertCell(2).textContent = percentage + '%';
    });
}

// Helper Functions
function formatNumber(num) {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
    return num.toString();
}

function getTimeAgo(timestamp) {
    if (!timestamp) return 'N/A';
    const now = new Date();
    const time = new Date(timestamp);
    const diffMs = now - time;
    const diffMins = Math.floor(diffMs / 60000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours}h ago`;
    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays}d ago`;
}

function calculateTrends(currentData) {
    if (!previousData.total_logs) {
        return { logs: '+0%', ips: '+0%', size: '+0%' };
    }

    const logsTrend = ((currentData.total_logs - previousData.total_logs) / previousData.total_logs * 100).toFixed(1);
    const ipsTrend = ((currentData.unique_ips - previousData.unique_ips) / previousData.unique_ips * 100).toFixed(1);
    const sizeTrend = ((currentData.avg_response_size - previousData.avg_response_size) / previousData.avg_response_size * 100).toFixed(1);

    return {
        logs: `${logsTrend >= 0 ? '+' : ''}${logsTrend}%`,
        ips: `${ipsTrend >= 0 ? '+' : ''}${ipsTrend}%`,
        size: `${sizeTrend >= 0 ? '+' : ''}${sizeTrend}%`
    };
}

// ML Analytics Functions
function updateAnomalyDetection(data) {
    // Simple anomaly detection using Z-score
    const values = [data.total_logs || 0];
    const mean = values.reduce((a, b) => a + b, 0) / values.length;
    const stdDev = Math.sqrt(values.reduce((sq, n) => sq + Math.pow(n - mean, 2), 0) / values.length);

    const isAnomaly = Math.abs((data.total_logs - mean) / stdDev) > 2;

    document.getElementById('anomaly-status').textContent = isAnomaly ? 'Anomaly Detected!' : 'Normal';
    document.getElementById('anomaly-alert').className = isAnomaly ? 'alert alert-danger' : 'alert alert-success';

    // Update anomaly chart with mock data
    const labels = Array.from({length: 20}, (_, i) => `T-${20-i}`);
    const normalData = Array.from({length: 20}, () => Math.random() * 100 + 50);
    const anomalies = [null, null, null, null, null, null, null, null, null, null,
                      null, null, null, null, null, null, null, 150, null, null];

    anomalyChart.data.labels = labels;
    anomalyChart.data.datasets[0].data = normalData;
    anomalyChart.data.datasets[1].data = anomalies;
    anomalyChart.update();
}

function updateTrendAnalysis(data) {
    // Mock trend analysis
    const trend = Math.random() > 0.5 ? 'Increasing' : 'Decreasing';
    const confidence = Math.floor(Math.random() * 30 + 70);

    document.getElementById('trend-direction').textContent = trend === 'Increasing' ? '↗️ Increasing' : '↘️ Decreasing';
    document.getElementById('trend-confidence').textContent = confidence + '%';

    // Update trend chart
    const labels = Array.from({length: 10}, (_, i) => `Day ${i+1}`);
    const actualData = Array.from({length: 10}, () => Math.random() * 100 + 50);
    const trendData = actualData.map((val, i) => val + (trend === 'Increasing' ? i * 5 : -i * 3));

    trendChart.data.labels = labels;
    trendChart.data.datasets[0].data = actualData;
    trendChart.data.datasets[1].data = trendData;
    trendChart.update();
}

function updateUserClustering(ipData) {
    if (!ipData || ipData.length === 0) return;

    // Simple clustering based on request count
    const heavy = ipData.filter(ip => (ip.request_count || ip.count) > 1000).length;
    const regular = ipData.filter(ip => (ip.request_count || ip.count) > 100 && (ip.request_count || ip.count) <= 1000).length;
    const light = ipData.filter(ip => (ip.request_count || ip.count) <= 100).length;

    document.getElementById('heavy-users').textContent = heavy;
    document.getElementById('regular-users').textContent = regular;
    document.getElementById('light-users').textContent = light;

    clusterChart.data.datasets[0].data = [heavy, regular, light];
    clusterChart.update();
}

function updateSecurityRisk(statusData) {
    if (!statusData || statusData.length === 0) return;

    const errorCodes = statusData.filter(s => s.status >= 400);
    const totalErrors = errorCodes.reduce((sum, s) => sum + s.count, 0);
    const totalRequests = statusData.reduce((sum, s) => sum + s.count, 0);
    const errorRate = totalRequests > 0 ? (totalErrors / totalRequests) * 100 : 0;

    let highRisk = 0, mediumRisk = 0, lowRisk = 0;

    if (errorRate > 10) highRisk = Math.floor(errorRate / 10);
    else if (errorRate > 5) mediumRisk = Math.floor(errorRate / 5);
    else lowRisk = Math.max(1, Math.floor(100 - errorRate));

    document.getElementById('high-risk').textContent = highRisk;
    document.getElementById('medium-risk').textContent = mediumRisk;
    document.getElementById('low-risk').textContent = lowRisk;
}

function updatePredictiveAnalytics(data) {
    // Mock prediction
    const predicted = Math.floor((data.total_logs || 0) * 1.15);
    const confidence = Math.floor(Math.random() * 20 + 70);

    document.getElementById('predicted-logs').textContent = `~${formatNumber(predicted)}`;
    document.getElementById('prediction-confidence').style.width = confidence + '%';
    document.getElementById('prediction-confidence').textContent = confidence + '% Confidence';

    // Update prediction chart
    const labels = Array.from({length: 10}, (_, i) => `H-${i+1}`);
    const historical = Array.from({length: 7}, () => Math.random() * 100 + 50);
    const predicted_data = Array.from({length: 3}, () => Math.random() * 120 + 60);

    predictionChart.data.labels = labels;
    predictionChart.data.datasets[0].data = [...historical, null, null, null];
    predictionChart.data.datasets[1].data = [null, null, null, null, null, null, null, ...predicted_data];
    predictionChart.update();
}

// Auto Refresh Functions
function setupAutoRefresh() {
    if (autoRefreshInterval) {
        clearInterval(autoRefreshInterval);
        autoRefreshInterval = null;

    }

    const isEnabled = document.getElementById('autoRefreshToggle').checked;
    const interval = parseInt(document.getElementById('refreshInterval').value) * 1000;

    const liveBadge = document.getElementById('live-badge');
    if (liveBadge) liveBadge.classList.toggle('d-none', !isEnabled);

    if (isEnabled) {
        autoRefreshInterval = setInterval(fetchData, interval);
    }
}

function handleVisibilityChange() {
    const isHidden = document.hidden;
    const autoToggle = document.getElementById('autoRefreshToggle');
    if (isHidden) {
        if (autoRefreshInterval) {
            clearInterval(autoRefreshInterval);
            autoRefreshInterval = null;
        }
    } else {
        if (autoToggle && autoToggle.checked && !autoRefreshInterval) {

            setupAutoRefresh();
            // Immediately fetch on resume
            fetchData();
        }
    }
}

// Theme application helper
function applyTheme(mode) {
    const root = document.documentElement;
    if (mode === 'system') {
        const isDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
        root.setAttribute('data-theme', isDark ? 'dark' : 'light');
    } else {
        root.setAttribute('data-theme', mode);
    }
}



function handleAutoRefreshToggle() {
    setupAutoRefresh();
}

// Show/hide loading state
function setLoading(loading) {
    const refreshBtn = document.getElementById('refresh-btn');
    const refreshText = document.getElementById('refresh-text');
    const refreshSpinner = document.getElementById('refresh-spinner');

    if (loading) {
        refreshBtn.disabled = true;
        refreshText.textContent = 'Refreshing...';
        refreshSpinner.classList.remove('d-none');
    } else {
        refreshBtn.disabled = false;
        refreshText.textContent = 'Refresh Now';
        refreshSpinner.classList.add('d-none');
    }
}

// Show/hide messages
function showError(message) {
    const errorAlert = document.getElementById('error-alert');
    errorAlert.textContent = message;
    errorAlert.classList.remove('d-none');
    setTimeout(() => errorAlert.classList.add('d-none'), 5000);
}

function showSuccess(message) {
    const successAlert = document.getElementById('success-alert');
    successAlert.textContent = message;
    successAlert.classList.remove('d-none');
    setTimeout(() => successAlert.classList.add('d-none'), 3000);
}

function showInfo(message) {
    // Use success alert for info messages with different styling
    const successAlert = document.getElementById('success-alert');
    successAlert.textContent = message;
    successAlert.classList.remove('d-none');
    setTimeout(() => successAlert.classList.add('d-none'), 4000);
}

function showWarning(message) {
    // Use error alert for warnings
    const errorAlert = document.getElementById('error-alert');
    errorAlert.textContent = message;
    errorAlert.classList.remove('d-none');
    setTimeout(() => errorAlert.classList.add('d-none'), 4000);
}

function hideError() {
    const errorAlert = document.getElementById('error-alert');
    errorAlert.classList.add('d-none');
}

// Log Generation Functions

async function handleLogGeneration(event) {
    event.preventDefault();

    if (isGenerating) {
        showError('Log generation is already in progress');
        return;
    }

    const numLogs = parseInt(document.getElementById('numLogs').value);
    const timeUnit = document.getElementById('timeUnit').value;

    if (numLogs <= 0 || numLogs > 1000) {
        showError('Number of logs must be between 1 and 1,000');
        return;
    }

    try {
        // Start log generation
        const response = await fetch(`${GENERATOR_API_URL}/logs`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                num_logs: numLogs,
                time: timeUnit
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: Failed to start log generation`);
        }

        const result = await response.json();

        if (result.status) {
            const timeUnitText = timeUnit === 's' ? 'second' : timeUnit === 'm' ? 'minute' : 'hour';
            showSuccess(`Log generation started! Generating ${numLogs} logs per ${timeUnitText}`);
            startRealTimeMonitoring(numLogs, timeUnit);
        } else {
            throw new Error(result.message || 'Failed to start log generation');
        }

    } catch (error) {
        console.error('Log generation error:', error);
        showError('Error starting log generation: ' + error.message);
    }
}

function startRealTimeMonitoring(logsPerInterval, timeUnit) {
    isGenerating = true;
    updateGeneratorStatus();

    document.getElementById('start-generation').disabled = true;
    document.getElementById('stop-generation').disabled = false;

    const startTime = Date.now();

    // Update display immediately
    const timeUnitDisplay = timeUnit === 's' ? 'second' : timeUnit === 'm' ? 'minute' : 'hour';
    document.getElementById('generation-rate-display').textContent = `${logsPerInterval}/${timeUnit}`;
    document.getElementById('generation-status').textContent =
        `Generating ${logsPerInterval} logs per ${timeUnitDisplay}`;

    // Update status alert styling
    const statusAlert = document.getElementById('generation-status-alert');
    statusAlert.className = 'alert alert-success';

    // Real-time monitoring every 2 seconds
    generationInterval = setInterval(async () => {
        try {
            // Check generator status
            const statusResponse = await fetch(`${GENERATOR_API_URL}/logs/status`);
            if (statusResponse.ok) {
                const statusResult = await statusResponse.json();

                if (!statusResult.data || !statusResult.data.active) {
                    // Generation stopped externally
                    stopLogGeneration();
                    showInfo('Log generation stopped externally');
                    return;
                }
            }

            // Get total logs count from LogParser
            const logsCountResponse = await fetch(`${API_BASE_URL}/logs/count`);
            if (logsCountResponse.ok) {
                const logsCountResult = await logsCountResponse.json();
                if (logsCountResult.status && logsCountResult.data) {
                    document.getElementById('total-logs-count').textContent = formatNumber(logsCountResult.data.count || 0);
                }
            }

            // Update uptime
            const elapsedSeconds = Math.floor((Date.now() - startTime) / 1000);
            const hours = Math.floor(elapsedSeconds / 3600);
            const minutes = Math.floor((elapsedSeconds % 3600) / 60);
            const seconds = elapsedSeconds % 60;

            let uptimeText;
            if (hours > 0) {
                uptimeText = `${hours}h ${minutes}m ${seconds}s`;
            } else if (minutes > 0) {
                uptimeText = `${minutes}m ${seconds}s`;
            } else {
                uptimeText = `${seconds}s`;
            }
            document.getElementById('generation-uptime').textContent = uptimeText;

        } catch (error) {
            // Suppress noisy errors while services (e.g., DB) are initializing
            // console.debug('Real-time monitoring fetch error:', error);
        }
    }, 2000);
}

async function stopLogGeneration() {
    try {
        // Call the stop API
        const response = await fetch(`${GENERATOR_API_URL}/logs/stop`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            }
        });

        if (response.ok) {
            const result = await response.json();
            if (result.status) {
                showSuccess('Log generation stopped successfully!');
            } else {
                showWarning('Stop request sent, but response was: ' + (result.message || 'Unknown'));
            }
        } else {
            showWarning(`Stop request sent (HTTP ${response.status})`);
        }
    } catch (error) {
        console.error('Error stopping log generation:', error);
        showWarning('Stop request sent, but network error occurred');
    }

    // Reset UI state
    if (generationInterval) {
        clearInterval(generationInterval);
        generationInterval = null;
    }

    isGenerating = false;
    updateGeneratorStatus();

    document.getElementById('start-generation').disabled = false;
    document.getElementById('stop-generation').disabled = true;
    document.getElementById('generation-status').textContent = 'Ready to generate logs';

    // Reset status alert styling
    const statusAlert = document.getElementById('generation-status-alert');
    statusAlert.className = 'alert alert-info';

    // Reset counters
    document.getElementById('generation-rate-display').textContent = '0/s';
    document.getElementById('generation-uptime').textContent = '0s';
}

async function checkInitialGeneratorStatus() {
    try {
        // Check generator status
        const response = await fetch(`${GENERATOR_API_URL}/logs/status`);
        if (response.ok) {
            const result = await response.json();
            if (result.status && result.data && result.data.active) {
                // Generator is already running
                isGenerating = true;
                document.getElementById('start-generation').disabled = true;
                document.getElementById('stop-generation').disabled = false;
                document.getElementById('generation-status').textContent = 'Generation is currently active (started externally)';

                // Update status alert
                const statusAlert = document.getElementById('generation-status-alert');
                statusAlert.className = 'alert alert-warning';

                // Show info message
                showInfo('Log generation is already running. You can stop it using the Stop button.');
            }
        }

        // Get initial total logs count
        const logsCountResponse = await fetch(`${API_BASE_URL}/logs/count`);
        if (logsCountResponse.ok) {
            const logsCountResult = await logsCountResponse.json();
            if (logsCountResult.status && logsCountResult.data) {
                document.getElementById('total-logs-count').textContent = formatNumber(logsCountResult.data.count || 0);
            }
        }
    } catch (error) {
        console.error('Error checking initial generator status:', error);
    }

    updateGeneratorStatus();
}

function updateGeneratorStatus() {
    const statusIndicator = document.getElementById('generator-status');
    const statusText = document.getElementById('generator-status-text');


    if (isGenerating) {
        statusIndicator.className = 'status-indicator status-running pulse';
        statusText.textContent = 'Running';
    } else {
        statusIndicator.className = 'status-indicator status-idle';
        statusText.textContent = 'Idle';
    }
}

// Log Viewer Functions
async function fetchLogs() {
    const limit = document.getElementById('logLimit').value;
    const statusFilter = document.getElementById('statusFilter').value;

    try {
        const url = new URL(`${API_BASE_URL}/logs`, window.location.origin);
        url.searchParams.set('limit', limit);
        if (statusFilter !== 'all') {
            url.searchParams.set('status', statusFilter);
        }
        // Apply cursor-based pagination if set
        if (currentCursorToken) {
            const parts = String(currentCursorToken).split('&id=');
            if (parts.length === 2) {
                url.searchParams.set('cursor', parts[0]);
                url.searchParams.set('id', parts[1]);
            } else {
                url.searchParams.set('cursor', currentCursorToken);
            }
        }

        const res = await fetch(url.toString());
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const payload = await res.json();
        const data = payload.data || {};
        const rawLogs = data.logs || [];
        const logs = rawLogs.map(l => {
            const req = (l.request || '').split(' ');
            return {
                timestamp: l.time_local,
                ip: l.remote_addr,
                method: req[0] || '',
                url: req[1] || l.request || '',
                status: l.status,
                size: l.body_bytes_sent,
                userAgent: l.http_user_agent
            };
        });
        updateLogsTable(logs);

        // Update pagination UI/state
        lastPaging = (data && data.paging) || null;
        updatePaginationUI(lastPaging);
    } catch (error) {
        showError('Error fetching logs: ' + error.message);
    }
}

function generateMockLogs(count, statusFilter) {
    const logs = [];
    const ips = ['192.168.1.1', '192.168.1.2', '10.0.0.1', '203.0.113.1', '198.51.100.1'];
    const methods = ['GET', 'POST', 'PUT', 'DELETE'];
    const urls = ['/api/users', '/api/orders', '/dashboard', '/login', '/api/products'];
    const statuses = statusFilter === 'all' ? [200, 301, 404, 500] : [parseInt(statusFilter)];
    const userAgents = [
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
        'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36'
    ];

    for (let i = 0; i < count; i++) {
        const timestamp = new Date(Date.now() - Math.random() * 86400000); // Last 24 hours
        logs.push({
            timestamp: timestamp.toISOString(),
            ip: ips[Math.floor(Math.random() * ips.length)],
            method: methods[Math.floor(Math.random() * methods.length)],
            url: urls[Math.floor(Math.random() * urls.length)],
            status: statuses[Math.floor(Math.random() * statuses.length)],
            size: Math.floor(Math.random() * 5000) + 200,
            userAgent: userAgents[Math.floor(Math.random() * userAgents.length)]
        });
    }

    return logs.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
}

function updateLogsTable(logs) {
    const tbody = document.getElementById('logs-table');
    tbody.innerHTML = '';

    if (logs.length === 0) {
        const row = tbody.insertRow();
        const cell = row.insertCell(0);
        cell.colSpan = 7;
        cell.className = 'text-center text-muted';
        cell.innerHTML = '<i class="fas fa-info-circle me-2"></i>No logs found';
        return;
    }

    logs.forEach(log => {
        const row = tbody.insertRow();
        row.insertCell(0).textContent = new Date(log.timestamp).toLocaleString();
        row.insertCell(1).textContent = log.ip;
        row.insertCell(2).innerHTML = `<span class="badge bg-primary">${log.method}</span>`;
        row.insertCell(3).textContent = log.url;

        const statusCell = row.insertCell(4);
        const statusClass = log.status < 300 ? 'success' : log.status < 400 ? 'warning' : 'danger';
        statusCell.innerHTML = `<span class="badge bg-${statusClass}">${log.status}</span>`;

        row.insertCell(5).textContent = formatNumber(log.size) + ' bytes';
        row.insertCell(6).textContent = log.userAgent.substring(0, 50) + '...';
    });

    document.getElementById('logs-count').textContent = logs.length;
}

// Tab change handler
function handleTabChange(event) {
    const targetTab = event.target.getAttribute('data-bs-target');

    if (targetTab === '#logs') {
        fetchLogs();
    } else if (targetTab === '#analytics') {
        // Refresh ML analytics when tab is shown
        if (previousData.total_logs) {
            updateMLAnalytics(previousData, [], []);
        }
    }
}

// ML Features Integration
function initializeMLFeatures() {
    // Initialize ML charts and features
    fetchMLInsights();

    // Set up periodic ML updates (every 30 seconds)
    setInterval(fetchMLInsights, 30 * 1000);
}

async function fetchMLInsights() {
    try {
        const response = await fetch(`${API_BASE_URL}/ml/insights`);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);

        const result = await response.json();
        if (result.status && result.data) {
            updateMLDashboard(result.data);
        }
    } catch (error) {
        console.warn('ML insights not available:', error.message);
        // Gracefully handle ML service unavailability
        updateMLDashboardWithMockData();
    }
}

function updateMLDashboard(insights) {
    // Update anomaly detection
    updateAnomalyChart(insights.anomalies || []);
    updateAnomalyStatus(insights.anomalies || []);

    // Update predictions
    updatePredictionChart(insights.predictions || []);
    updatePredictionDisplay(insights.predictions || []);

    // Update security threats
    updateSecurityThreats(insights.security_threats || []);

    // Update user clustering
    updateUserClusters(insights.clusters || []);

    // Update trend analysis
    updateTrendAnalysis(insights.trend_analysis || {});
}

function updateAnomalyChart(anomalies) {
    if (!anomalyChart) return;

    // Process anomalies for chart display
    const last24Hours = anomalies.filter(a => {
        const anomalyTime = new Date(a.timestamp);
        const now = new Date();
        return (now - anomalyTime) <= 24 * 60 * 60 * 1000;
    });

    const hourlyAnomalies = {};
    last24Hours.forEach(anomaly => {
        const hour = new Date(anomaly.timestamp).getHours();
        if (!hourlyAnomalies[hour]) hourlyAnomalies[hour] = 0;
        if (anomaly.is_anomaly) hourlyAnomalies[hour]++;
    });

    const labels = Array.from({length: 24}, (_, i) => `${i}:00`);
    const data = labels.map((_, i) => hourlyAnomalies[i] || 0);

    anomalyChart.data.labels = labels;
    anomalyChart.data.datasets[0].data = data;
    anomalyChart.update();
}

function updateAnomalyStatus(anomalies) {
    const recentAnomalies = anomalies.filter(a => {
        const anomalyTime = new Date(a.timestamp);
        const now = new Date();
        return (now - anomalyTime) <= 60 * 60 * 1000 && a.is_anomaly; // Last hour
    });

    const statusElement = document.getElementById('anomaly-status');
    if (statusElement) {
        if (recentAnomalies.length === 0) {
            statusElement.textContent = 'Normal - No anomalies detected';
            statusElement.parentElement.className = 'alert alert-success';
        } else {
            const highSeverity = recentAnomalies.filter(a => a.severity === 'high' || a.severity === 'critical');
            if (highSeverity.length > 0) {
                statusElement.textContent = `Critical - ${highSeverity.length} high-severity anomalies detected`;
                statusElement.parentElement.className = 'alert alert-danger';
            } else {
                statusElement.textContent = `Warning - ${recentAnomalies.length} anomalies detected`;
                statusElement.parentElement.className = 'alert alert-warning';
            }
        }
    }
}

function updatePredictionChart(predictions) {
    if (!predictionChart) return;

    const next12Hours = predictions.slice(0, 12);
    const labels = next12Hours.map(p => {
        const time = new Date(p.timestamp);
        return `${time.getHours()}:00`;
    });

    const predictedValues = next12Hours.map(p => p.predicted_value);
    const upperBounds = next12Hours.map(p => p.upper_bound);
    const lowerBounds = next12Hours.map(p => p.lower_bound);

    predictionChart.data.labels = labels;
    predictionChart.data.datasets[0].data = predictedValues;
    predictionChart.data.datasets[1].data = upperBounds;
    predictionChart.data.datasets[2].data = lowerBounds;
    predictionChart.update();
}

function updatePredictionDisplay(predictions) {
    const nextHourPrediction = predictions.find(p => {
        const predTime = new Date(p.timestamp);
        const nextHour = new Date();
        nextHour.setHours(nextHour.getHours() + 1);
        return Math.abs(predTime - nextHour) < 30 * 60 * 1000; // Within 30 minutes
    });

    if (nextHourPrediction) {
        const predictedLogsElement = document.getElementById('predicted-logs');
        const confidenceElement = document.getElementById('prediction-confidence');

        if (predictedLogsElement) {
            predictedLogsElement.textContent = `~${Math.round(nextHourPrediction.predicted_value)}`;
        }

        if (confidenceElement) {
            const confidence = Math.round(nextHourPrediction.confidence_level * 100);
            confidenceElement.textContent = `${confidence}% Confidence`;
            confidenceElement.style.width = `${confidence}%`;
        }
    }
}

function updateSecurityThreats(threats) {
    // Update security threat indicators
    const highThreats = threats.filter(t => t.severity === 'high' || t.severity === 'critical');
    const mediumThreats = threats.filter(t => t.severity === 'medium');

    // Update threat counters if elements exist
    const threatElements = {
        'high-threats': highThreats.length,
        'medium-threats': mediumThreats.length,
        'total-threats': threats.length
    };

    Object.entries(threatElements).forEach(([id, count]) => {
        const element = document.getElementById(id);
        if (element) element.textContent = count;
    });
}

function updateUserClusters(clusters) {
    // Group clusters by type
    const clusterGroups = {};
    clusters.forEach(cluster => {
        if (!clusterGroups[cluster.cluster_name]) {
            clusterGroups[cluster.cluster_name] = [];
        }
        clusterGroups[cluster.cluster_name].push(cluster);
    });

    // Update cluster displays
    const heavyUsers = clusterGroups['Heavy Users'] || [];
    const mediumUsers = clusterGroups['Medium Users'] || [];
    const lightUsers = clusterGroups['Light Users'] || [];

    const clusterElements = {
        'heavy-users': heavyUsers.length,
        'medium-users': mediumUsers.length,
        'light-users': lightUsers.length
    };

    Object.entries(clusterElements).forEach(([id, count]) => {
        const element = document.getElementById(id);
        if (element) element.textContent = count;
    });

    // Update cluster chart if it exists
    if (clusterChart) {
        const data = [lightUsers.length, mediumUsers.length, heavyUsers.length];
        clusterChart.data.datasets[0].data = data;
        clusterChart.update();
    }
}

function updateTrendAnalysis(trendAnalysis) {
    // Update trend indicators
    const trendElement = document.getElementById('trend-direction');
    if (trendElement && trendAnalysis.trend) {
        trendElement.textContent = trendAnalysis.trend.charAt(0).toUpperCase() + trendAnalysis.trend.slice(1);

        // Update trend icon/color based on direction
        trendElement.className = `trend-${trendAnalysis.trend}`;
    }
}

function updateMLDashboardWithMockData() {
    // Provide mock data when ML service is unavailable
    const mockInsights = {
        anomalies: [],
        predictions: Array.from({length: 12}, (_, i) => ({
            timestamp: new Date(Date.now() + i * 60 * 60 * 1000),
            predicted_value: 100 + Math.random() * 50,
            confidence_level: 0.75,
            upper_bound: 150 + Math.random() * 25,
            lower_bound: 50 + Math.random() * 25
        })),
        security_threats: [],
        clusters: [],
        trend_analysis: {
            trend: 'stable',
            slope: 0,
            correlation: 0.5,
            seasonality: false
        }
    };

    updateMLDashboard(mockInsights);
}
