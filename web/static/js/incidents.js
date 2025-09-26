let allIncidents = [];
let filteredIncidents = [];

async function loadIncidents() {
    try {
        const response = await fetch('/api/incidents');
        allIncidents = await response.json();
        
        // Sort by creation date (newest first)
        allIncidents.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
        
        filteredIncidents = [...allIncidents];
        renderIncidentsTable();
    } catch (error) {
        console.error('Error loading incidents:', error);
        document.getElementById('incidents-table-body').innerHTML = 
            '<tr><td colspan="7" class="loading">Error loading incidents</td></tr>';
    }
}

function renderIncidentsTable() {
    const tbody = document.getElementById('incidents-table-body');
    
    if (filteredIncidents.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="loading">No incidents found</td></tr>';
        return;
    }

    tbody.innerHTML = filteredIncidents.map(incident => {
        const duration = calculateDuration(incident);
        const actions = generateActionButtons(incident);
        
        return `
            <tr onclick="showIncidentDetails('${incident.id}')" style="cursor: pointer;">
                <td>${incident.id.substring(0, 8)}...</td>
                <td>${escapeHtml(incident.title)}</td>
                <td><span class="status-badge status-${incident.status}">${incident.status}</span></td>
                <td><span class="severity-badge severity-${incident.severity}">${incident.severity}</span></td>
                <td>${formatDate(incident.created_at)}</td>
                <td>${duration}</td>
                <td onclick="event.stopPropagation()">${actions}</td>
            </tr>
        `;
    }).join('');
}

function generateActionButtons(incident) {
    let buttons = '';
    
    if (incident.status === 'open') {
        buttons += `<button class="btn btn-warning" onclick="acknowledgeIncident('${incident.id}')">Acknowledge</button> `;
    }
    
    if (incident.status !== 'resolved') {
        buttons += `<button class="btn btn-success" onclick="resolveIncident('${incident.id}')">Resolve</button>`;
    }
    
    return buttons;
}

function calculateDuration(incident) {
    const start = new Date(incident.created_at);
    let end;
    
    if (incident.resolved_at) {
        end = new Date(incident.resolved_at);
    } else {
        end = new Date();
    }
    
    const diff = end - start;
    return formatDuration(diff * 1000000); // Convert to nanoseconds for consistency
}

async function showIncidentDetails(incidentId) {
    try {
        const response = await fetch(`/api/incidents/${incidentId}`);
        const incident = await response.json();
        
        const modal = document.getElementById('incident-modal');
        const detailsDiv = document.getElementById('incident-details');
        
        detailsDiv.innerHTML = `
            <h2>${escapeHtml(incident.title)}</h2>
            <div style="margin: 20px 0;">
                <strong>ID:</strong> ${incident.id}<br>
                <strong>Status:</strong> <span class="status-badge status-${incident.status}">${incident.status}</span><br>
                <strong>Severity:</strong> <span class="severity-badge severity-${incident.severity}">${incident.severity}</span><br>
                <strong>Created:</strong> ${formatDate(incident.created_at)}<br>
                ${incident.acked_at ? `<strong>Acknowledged:</strong> ${formatDate(incident.acked_at)}<br>` : ''}
                ${incident.resolved_at ? `<strong>Resolved:</strong> ${formatDate(incident.resolved_at)}<br>` : ''}
                ${incident.assignee_id ? `<strong>Assignee:</strong> ${incident.assignee_id}<br>` : ''}
            </div>
            <div style="margin: 20px 0;">
                <strong>Description:</strong><br>
                <div style="background: #f8f9fa; padding: 10px; border-radius: 4px; margin-top: 5px;">
                    ${escapeHtml(incident.description).replace(/\n/g, '<br>')}
                </div>
            </div>
            ${incident.alert_ids && incident.alert_ids.length > 0 ? `
                <div style="margin: 20px 0;">
                    <strong>Related Alerts:</strong> ${incident.alert_ids.length}
                </div>
            ` : ''}
            <div style="margin-top: 20px;">
                ${generateActionButtons(incident)}
            </div>
        `;
        
        modal.style.display = 'block';
    } catch (error) {
        console.error('Error loading incident details:', error);
        alert('Error loading incident details');
    }
}

async function acknowledgeIncident(incidentId) {
    const assigneeId = prompt('Enter your user ID:');
    if (!assigneeId) return;
    
    try {
        const response = await fetch(`/api/incidents/${incidentId}/acknowledge`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ assignee_id: assigneeId })
        });
        
        if (response.ok) {
            loadIncidents();
            alert('Incident acknowledged successfully');
        } else {
            alert('Failed to acknowledge incident');
        }
    } catch (error) {
        console.error('Error acknowledging incident:', error);
        alert('Error acknowledging incident');
    }
}

async function resolveIncident(incidentId) {
    if (!confirm('Are you sure you want to resolve this incident?')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/incidents/${incidentId}/resolve`, {
            method: 'POST'
        });
        
        if (response.ok) {
            loadIncidents();
            alert('Incident resolved successfully');
        } else {
            alert('Failed to resolve incident');
        }
    } catch (error) {
        console.error('Error resolving incident:', error);
        alert('Error resolving incident');
    }
}

function filterIncidents() {
    const statusFilter = document.getElementById('status-filter').value;
    const severityFilter = document.getElementById('severity-filter').value;
    
    filteredIncidents = allIncidents.filter(incident => {
        const statusMatch = !statusFilter || incident.status === statusFilter;
        const severityMatch = !severityFilter || incident.severity === severityFilter;
        return statusMatch && severityMatch;
    });
    
    renderIncidentsTable();
}

function closeModal() {
    document.getElementById('incident-modal').style.display = 'none';
}

function refreshIncidents() {
    loadIncidents();
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
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

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Close modal when clicking outside of it
window.onclick = function(event) {
    const modal = document.getElementById('incident-modal');
    if (event.target === modal) {
        modal.style.display = 'none';
    }
}

// Auto-refresh every 30 seconds
setInterval(loadIncidents, 30000);

// Load incidents on page load
document.addEventListener('DOMContentLoaded', loadIncidents);