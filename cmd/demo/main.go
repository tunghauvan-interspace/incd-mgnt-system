package main

// Simple test to demonstrate enhanced incident features without authentication
// This is for testing purposes only

import (
	"fmt"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

func main() {
	fmt.Println("🚀 Enhanced Incident Features Demo")
	fmt.Println("==================================")

	// Initialize services
	store, err := storage.NewMemoryStore()
	if err != nil {
		fmt.Printf("Failed to create memory store: %v\n", err)
		return
	}
	defer store.Close()

	metricsService := services.NewMetricsService()
	incidentService := services.NewIncidentService(store, metricsService)

	// 1. Create a test incident
	fmt.Println("📝 1. Creating a test incident...")
	incident, err := incidentService.CreateIncident(
		"Database Connection Issues",
		"Users are experiencing timeout errors when accessing the application",
		models.SeverityHigh,
		[]string{},
	)
	if err != nil {
		fmt.Printf("Failed to create incident: %v\n", err)
		return
	}
	fmt.Printf("✅ Created incident: %s (ID: %s)\n", incident.Title, incident.ID)
	fmt.Println()

	// 2. Add comments to demonstrate timeline tracking
	fmt.Println("💬 2. Adding comments to track investigation...")
	comment1, err := incidentService.AddComment(
		incident.ID,
		"engineer-alice",
		"Initial investigation started. Checking database connections and query performance.",
		models.CommentTypeComment,
		nil,
	)
	if err != nil {
		fmt.Printf("Failed to add comment: %v\n", err)
		return
	}
	fmt.Printf("✅ Added comment: %s\n", comment1.Content[:50]+"...")

	comment2, err := incidentService.AddComment(
		incident.ID,
		"engineer-alice", 
		"Found high CPU usage on database server. Investigating potential queries causing the load.",
		models.CommentTypeComment,
		nil,
	)
	if err != nil {
		fmt.Printf("Failed to add comment: %v\n", err)
		return
	}
	fmt.Printf("✅ Added comment: %s\n", comment2.Content[:50]+"...")
	fmt.Println()

	// 3. Add tags to categorize the incident
	fmt.Println("🏷️  3. Adding tags to categorize the incident...")
	tags := []models.TemplateTag{
		{Name: "component", Value: "database", Color: "#007bff"},
		{Name: "priority", Value: "high", Color: "#dc3545"},
		{Name: "environment", Value: "production", Color: "#28a745"},
		{Name: "team", Value: "backend", Color: "#6f42c1"},
	}

	err = incidentService.AddTags(incident.ID, "engineer-alice", tags)
	if err != nil {
		fmt.Printf("Failed to add tags: %v\n", err)
		return
	}
	fmt.Printf("✅ Added %d tags to categorize the incident\n", len(tags))
	fmt.Println()

	// 4. Create an incident template
	fmt.Println("📋 4. Creating incident template for future use...")
	template := &models.IncidentTemplate{
		Name:                "Database Performance Issue",
		Description:         "Template for database-related performance incidents",
		TitleTemplate:       "Database Performance: {{service_name}}",
		DescriptionTemplate: "Performance issues detected in {{service_name}} database.\n\nSymptoms:\n- {{symptoms}}\n\nAffected Services:\n- {{affected_services}}\n\nInvestigation Steps:\n1. Check database metrics\n2. Review slow query log\n3. Monitor connection pool\n4. Check for blocking queries",
		Severity:            models.SeverityHigh,
		DefaultTags: []models.TemplateTag{
			{Name: "component", Value: "database", Color: "#007bff"},
			{Name: "type", Value: "performance", Color: "#ffc107"},
		},
	}

	err = incidentService.CreateTemplate(template)
	if err != nil {
		fmt.Printf("Failed to create template: %v\n", err)
		return
	}
	fmt.Printf("✅ Created template: %s (ID: %s)\n", template.Name, template.ID)
	fmt.Println()

	// 5. Create incident from template
	fmt.Println("🔄 5. Creating incident from template...")
	templateReq := &models.CreateIncidentFromTemplateRequest{
		TemplateID: template.ID,
		Variables: map[string]string{
			"service_name":      "user-authentication-db",
			"symptoms":          "Slow login responses, connection timeouts",
			"affected_services": "Login Service, User Profile Service, Session Management",
		},
		AdditionalTags: []models.TemplateTag{
			{Name: "severity_impact", Value: "user_facing", Color: "#dc3545"},
		},
	}

	templateIncident, err := incidentService.UseTemplate(templateReq, "engineer-bob")
	if err != nil {
		fmt.Printf("Failed to create incident from template: %v\n", err)
		return
	}
	fmt.Printf("✅ Created incident from template: %s (ID: %s)\n", templateIncident.Title, templateIncident.ID)
	fmt.Println()

	// 6. Demonstrate search functionality
	fmt.Println("🔍 6. Searching for database-related incidents...")
	searchReq := &models.IncidentSearchRequest{
		Query:    "database",
		Status:   []models.IncidentStatus{models.IncidentStatusOpen},
		Severity: []models.IncidentSeverity{models.SeverityHigh},
		Tags:     []string{"database"},
		Page:     1,
		Limit:    10,
	}

	searchResp, err := incidentService.SearchIncidents(searchReq)
	if err != nil {
		fmt.Printf("Failed to search incidents: %v\n", err)
		return
	}
	fmt.Printf("✅ Found %d incidents matching search criteria\n", searchResp.Total)
	for i, inc := range searchResp.Incidents {
		fmt.Printf("   %d. %s (Severity: %s)\n", i+1, inc.Title, inc.Severity)
	}
	fmt.Println()

	// 7. Demonstrate assignment workflow
	fmt.Println("👤 7. Assigning incident to specialist...")
	err = incidentService.AssignIncident(incident.ID, "database-specialist-carol", "manager-dave")
	if err != nil {
		fmt.Printf("Failed to assign incident: %v\n", err)
		return
	}
	fmt.Printf("✅ Assigned incident to database specialist\n")
	fmt.Println()

	// 8. Show timeline with all events
	fmt.Println("📅 8. Incident timeline (comments + system events)...")
	timeline, err := incidentService.GetTimeline(incident.ID)
	if err != nil {
		fmt.Printf("Failed to get timeline: %v\n", err)
		return
	}
	
	fmt.Printf("Timeline has %d events:\n", len(timeline))
	for i, event := range timeline {
		fmt.Printf("   %d. [%s] %s: %s\n", 
			i+1, 
			event.CommentType, 
			event.CreatedAt.Format("15:04:05"),
			event.Content[:min(60, len(event.Content))])
	}
	fmt.Println()

	// 9. Show tags
	fmt.Println("🏷️  9. Current incident tags...")
	incidentTags, err := incidentService.GetTags(incident.ID)
	if err != nil {
		fmt.Printf("Failed to get tags: %v\n", err)
		return
	}

	fmt.Printf("Incident has %d tags:\n", len(incidentTags))
	for _, tag := range incidentTags {
		value := ""
		if tag.TagValue != nil {
			value = *tag.TagValue
		}
		fmt.Printf("   • %s: %s (color: %s)\n", tag.TagName, value, tag.Color)
	}
	fmt.Println()

	// 10. Demonstrate bulk operations
	fmt.Println("📊 10. Bulk operations on multiple incidents...")
	
	// Create a few more incidents for bulk demo
	incident2, _ := incidentService.CreateIncident("API Response Timeout", "API slow response times", models.SeverityMedium, []string{})
	incident3, _ := incidentService.CreateIncident("Memory Usage High", "High memory usage on servers", models.SeverityHigh, []string{})
	
	// Bulk acknowledge
	bulkResp, err := incidentService.BulkAcknowledge(
		[]string{incident.ID, incident2.ID, incident3.ID},
		"oncall-engineer",
		"manager-dave",
	)
	if err != nil {
		fmt.Printf("Failed to bulk acknowledge: %v\n", err)
		return
	}
	fmt.Printf("✅ Bulk acknowledged %d incidents (%d failed)\n", bulkResp.ProcessedCount, bulkResp.FailedCount)
	fmt.Println()

	// 11. List all templates
	fmt.Println("📋 11. Available incident templates...")
	templates, err := incidentService.ListTemplates()
	if err != nil {
		fmt.Printf("Failed to list templates: %v\n", err)
		return
	}
	fmt.Printf("Available templates (%d):\n", len(templates))
	for i, tmpl := range templates {
		fmt.Printf("   %d. %s - %s (Severity: %s)\n", i+1, tmpl.Name, tmpl.Description, tmpl.Severity)
	}
	fmt.Println()

	fmt.Println("🎉 Enhanced Features Demo Completed Successfully!")
	fmt.Println("===============================================")
	fmt.Println("✅ Features demonstrated:")
	fmt.Println("   • Comments and timeline tracking")
	fmt.Println("   • Flexible tagging system")
	fmt.Println("   • Template-based incident creation")
	fmt.Println("   • Advanced search and filtering")
	fmt.Println("   • Assignment workflow")
	fmt.Println("   • Bulk operations")
	fmt.Println("   • Rich metadata and categorization")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}