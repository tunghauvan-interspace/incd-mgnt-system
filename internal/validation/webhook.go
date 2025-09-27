package validation

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// AlertmanagerWebhookSchema defines the JSON schema for Alertmanager webhooks
const AlertmanagerWebhookSchema = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"required": ["version", "status", "alerts"],
	"properties": {
		"version": {
			"type": "string"
		},
		"groupKey": {
			"type": "string"
		},
		"status": {
			"type": "string",
			"enum": ["firing", "resolved"]
		},
		"receiver": {
			"type": "string"
		},
		"groupLabels": {
			"type": "object",
			"additionalProperties": {
				"type": "string"
			}
		},
		"commonLabels": {
			"type": "object",
			"additionalProperties": {
				"type": "string"
			}
		},
		"commonAnnotations": {
			"type": "object",
			"additionalProperties": {
				"type": "string"
			}
		},
		"externalURL": {
			"type": "string"
		},
		"alerts": {
			"type": "array",
			"minItems": 1,
			"items": {
				"type": "object",
				"required": ["fingerprint", "status", "startsAt", "labels"],
				"properties": {
					"fingerprint": {
						"type": "string",
						"minLength": 1
					},
					"status": {
						"type": "string",
						"enum": ["firing", "resolved"]
					},
					"startsAt": {
						"type": "string",
						"format": "date-time"
					},
					"endsAt": {
						"type": "string",
						"format": "date-time"
					},
					"labels": {
						"type": "object",
						"additionalProperties": {
							"type": "string"
						}
					},
					"annotations": {
						"type": "object",
						"additionalProperties": {
							"type": "string"
						}
					}
				}
			}
		}
	}
}`

// WebhookValidator handles webhook validation
type WebhookValidator struct {
	alertmanagerSchema gojsonschema.JSONLoader
}

// NewWebhookValidator creates a new webhook validator
func NewWebhookValidator() *WebhookValidator {
	schemaLoader := gojsonschema.NewStringLoader(AlertmanagerWebhookSchema)
	return &WebhookValidator{
		alertmanagerSchema: schemaLoader,
	}
}

// ValidateAlertmanagerWebhook validates an Alertmanager webhook payload against the schema
func (v *WebhookValidator) ValidateAlertmanagerWebhook(payload []byte) error {
	// Parse JSON to check basic structure
	var jsonData interface{}
	if err := json.Unmarshal(payload, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate against schema
	documentLoader := gojsonschema.NewBytesLoader(payload)
	result, err := gojsonschema.Validate(v.alertmanagerSchema, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	if !result.Valid() {
		var errMsg string
		for i, err := range result.Errors() {
			if i > 0 {
				errMsg += ", "
			}
			errMsg += err.String()
		}
		return fmt.Errorf("webhook validation failed: %s", errMsg)
	}

	return nil
}