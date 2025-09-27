package validation

import (
	"testing"
)

func TestWebhookValidator_ValidateAlertmanagerWebhook(t *testing.T) {
	validator := NewWebhookValidator()

	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name: "valid webhook payload",
			payload: `{
				"version": "4",
				"status": "firing",
				"alerts": [
					{
						"fingerprint": "test123",
						"status": "firing",
						"startsAt": "2023-10-01T12:00:00Z",
						"labels": {"alertname": "test"},
						"annotations": {}
					}
				]
			}`,
			wantErr: false,
		},
		{
			name: "invalid JSON",
			payload: `{
				"version": "4",
				"status": "firing"
				"alerts": []
			}`,
			wantErr: true,
		},
		{
			name: "missing required field",
			payload: `{
				"version": "4",
				"status": "firing"
			}`,
			wantErr: true,
		},
		{
			name: "empty alerts array",
			payload: `{
				"version": "4",
				"status": "firing",
				"alerts": []
			}`,
			wantErr: true,
		},
		{
			name: "invalid status",
			payload: `{
				"version": "4",
				"status": "invalid_status",
				"alerts": [
					{
						"fingerprint": "test123",
						"status": "firing",
						"startsAt": "2023-10-01T12:00:00Z",
						"labels": {},
						"annotations": {}
					}
				]
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAlertmanagerWebhook([]byte(tt.payload))
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAlertmanagerWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}