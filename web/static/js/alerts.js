let allAlerts = [];
let filteredAlerts = [];

async function loadAlerts() {
    try {
        const response = await fetch('/api/alerts');
        allAlerts = await response.json();
        
        // Sort by creation date (newest first)
        allAlerts.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
        
        filteredAlerts = [...allAlerts];
        renderAlertsTable();
    } catch (error) {
        console.error('Error loading alerts:', error);
        document.getElementById('alerts-table-body').innerHTML = 
            '<tr><td colspan="7" class="loading">Error loading alerts</td></tr>';
    }
}

function renderAlertsTable() {
    const tbody = document.getElementById('alerts-table-body');
    
    if (filteredAlerts.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="loading">No alerts found</td></tr>';
        return;
    }

    tbody.innerHTML = filteredAlerts.map(alert => {
        const alertName = alert.labels.alertname || 'Unknown';
        const instance = alert.labels.instance || '-';
        const incidentLink = alert.incident_id ? 
            `<a href="/incidents" onclick="event.stopPropagation()">${alert.incident_id.substring(0, 8)}...</a>` : 
            '-';
        
        return `
            <tr onclick="showAlertDetails('${alert.id}')" style="cursor: pointer;">
                <td>${alert.id.substring(0, 8)}...</td>
                <td>${escapeHtml(alertName)}</td>
                <td><span class="status-badge ${alert.status === 'firing' ? 'status-open' : 'status-resolved'}">${alert.status}</span></td>
                <td>${escapeHtml(instance)}</td>
                <td>${formatDate(alert.starts_at)}</td>
                <td>${incidentLink}</td>
                <td onclick="event.stopPropagation()">
                    <button class="btn btn-primary" onclick="showAlertDetails('${alert.id}')">Details</button>
                </td>
            </tr>
        `;
    }).join('');
}

async function showAlertDetails(alertId) {
    try {
        const response = await fetch(`/api/alerts`);
        const alerts = await response.json();
        const alert = alerts.find(a => a.id === alertId);
        
        if (!alert) {
            alert('Alert not found');
            return;
        }
        
        const modal = document.getElementById('alert-modal');
        const detailsDiv = document.getElementById('alert-details');
        
        // Format labels and annotations
        const labelsHtml = Object.entries(alert.labels || {})
            .map(([key, value]) => `<div><strong>${escapeHtml(key)}:</strong> ${escapeHtml(value)}</div>`)
            .join('');
        
        const annotationsHtml = Object.entries(alert.annotations || {})
            .map(([key, value]) => `<div><strong>${escapeHtml(key)}:</strong> ${escapeHtml(value)}</div>`)
            .join('');
        
        detailsDiv.innerHTML = `
            <h2>Alert Details</h2>
            <div style="margin: 20px 0;">
                <strong>ID:</strong> ${alert.id}<br>
                <strong>Fingerprint:</strong> ${alert.fingerprint}<br>
                <strong>Status:</strong> <span class="status-badge ${alert.status === 'firing' ? 'status-open' : 'status-resolved'}">${alert.status}</span><br>
                <strong>Started:</strong> ${formatDate(alert.starts_at)}<br>
                ${alert.ends_at && alert.ends_at !== '0001-01-01T00:00:00Z' ? 
                    `<strong>Ended:</strong> ${formatDate(alert.ends_at)}<br>` : 
                    ''
                }
                <strong>Created:</strong> ${formatDate(alert.created_at)}<br>
                ${alert.incident_id ? `<strong>Incident:</strong> ${alert.incident_id}<br>` : ''}
            </div>
            
            ${labelsHtml ? `
                <div style="margin: 20px 0;">
                    <h3>Labels</h3>
                    <div style="background: #f8f9fa; padding: 10px; border-radius: 4px; margin-top: 5px;">
                        ${labelsHtml}
                    </div>
                </div>
            ` : ''}
            
            ${annotationsHtml ? `
                <div style="margin: 20px 0;">
                    <h3>Annotations</h3>
                    <div style="background: #f8f9fa; padding: 10px; border-radius: 4px; margin-top: 5px;">
                        ${annotationsHtml}
                    </div>
                </div>
            ` : ''}
        `;
        
        modal.style.display = 'block';
    } catch (error) {
        console.error('Error loading alert details:', error);
        alert('Error loading alert details');
    }
}

function filterAlerts() {
    const statusFilter = document.getElementById('status-filter').value;
    
    filteredAlerts = allAlerts.filter(alert => {
        const statusMatch = !statusFilter || alert.status === statusFilter;
        return statusMatch;
    });
    
    renderAlertsTable();
}

function closeModal() {
    document.getElementById('alert-modal').style.display = 'none';
}

function refreshAlerts() {
    loadAlerts();
}

function formatDate(dateString) {
    if (!dateString || dateString === '0001-01-01T00:00:00Z') {
        return '-';
    }
    const date = new Date(dateString);
    return date.toLocaleString();
}

function escapeHtml(text) {
    if (typeof text !== 'string') {
        text = String(text);
    }
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Close modal when clicking outside of it
window.onclick = function(event) {
    const modal = document.getElementById('alert-modal');
    if (event.target === modal) {
        modal.style.display = 'none';
    }
}

// Auto-refresh every 30 seconds
setInterval(loadAlerts, 30000);

// Load alerts on page load
document.addEventListener('DOMContentLoaded', loadAlerts);