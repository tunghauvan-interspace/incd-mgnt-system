let statusChart = null;
let severityChart = null;

async function loadDashboard() {
    try {
        // Load metrics
        const metricsResponse = await fetch('/api/metrics');
        const metrics = await metricsResponse.json();

        // Update metric cards
        document.getElementById('total-incidents').textContent = metrics.total_incidents || 0;
        document.getElementById('open-incidents').textContent = metrics.open_incidents || 0;
        document.getElementById('mtta').textContent = formatDuration(metrics.mtta);
        document.getElementById('mttr').textContent = formatDuration(metrics.mttr);

        // Update charts
        updateStatusChart(metrics.incidents_by_status || {});
        updateSeverityChart(metrics.incidents_by_severity || {});

        // Load recent incidents
        await loadRecentIncidents();

    } catch (error) {
        console.error('Error loading dashboard:', error);
    }
}

function updateStatusChart(data) {
    const ctx = document.getElementById('status-chart').getContext('2d');
    
    if (statusChart) {
        statusChart.destroy();
    }

    statusChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ['Open', 'Acknowledged', 'Resolved'],
            datasets: [{
                data: [
                    data.open || 0,
                    data.acknowledged || 0,
                    data.resolved || 0
                ],
                backgroundColor: [
                    '#dc2626',
                    '#d97706',
                    '#16a34a'
                ],
                borderWidth: 0
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom'
                }
            }
        }
    });
}

function updateSeverityChart(data) {
    const ctx = document.getElementById('severity-chart').getContext('2d');
    
    if (severityChart) {
        severityChart.destroy();
    }

    severityChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ['Critical', 'High', 'Medium', 'Low'],
            datasets: [{
                data: [
                    data.critical || 0,
                    data.high || 0,
                    data.medium || 0,
                    data.low || 0
                ],
                backgroundColor: [
                    '#dc2626',
                    '#ea580c',
                    '#d97706',
                    '#0369a1'
                ],
                borderWidth: 0
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: 1
                    }
                }
            }
        }
    });
}

async function loadRecentIncidents() {
    try {
        const response = await fetch('/api/incidents');
        const incidents = await response.json();

        // Sort by creation date (newest first) and take first 5
        const recentIncidents = incidents
            .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
            .slice(0, 5);

        const tableContainer = document.getElementById('recent-incidents-table');
        
        if (recentIncidents.length === 0) {
            tableContainer.innerHTML = '<div class="loading">No incidents found</div>';
            return;
        }

        let tableHTML = `
            <table class="incidents-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Title</th>
                        <th>Status</th>
                        <th>Severity</th>
                        <th>Created</th>
                    </tr>
                </thead>
                <tbody>
        `;

        recentIncidents.forEach(incident => {
            tableHTML += `
                <tr onclick="window.location.href='/incidents'">
                    <td>${incident.id.substring(0, 8)}...</td>
                    <td>${escapeHtml(incident.title)}</td>
                    <td><span class="status-badge status-${incident.status}">${incident.status}</span></td>
                    <td><span class="severity-badge severity-${incident.severity}">${incident.severity}</span></td>
                    <td>${formatDate(incident.created_at)}</td>
                </tr>
            `;
        });

        tableHTML += '</tbody></table>';
        tableContainer.innerHTML = tableHTML;

    } catch (error) {
        console.error('Error loading recent incidents:', error);
        document.getElementById('recent-incidents-table').innerHTML = '<div class="loading">Error loading incidents</div>';
    }
}

function formatDuration(nanoseconds) {
    if (!nanoseconds || nanoseconds === 0) {
        return '-';
    }

    const seconds = nanoseconds / 1000000000;
    
    if (seconds < 60) {
        return `${Math.round(seconds)}s`;
    } else if (seconds < 3600) {
        const minutes = Math.round(seconds / 60);
        return `${minutes}m`;
    } else if (seconds < 86400) {
        const hours = Math.round(seconds / 3600);
        return `${hours}h`;
    } else {
        const days = Math.round(seconds / 86400);
        return `${days}d`;
    }
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function refreshDashboard() {
    loadDashboard();
}

// Auto-refresh every 30 seconds
setInterval(loadDashboard, 30000);

// Load dashboard on page load
document.addEventListener('DOMContentLoaded', loadDashboard);