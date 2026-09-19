package model

import "time"

type Target struct { ID string `json:"id"`; Name string `json:"name"`; Origins []string `json:"origins"`; Description string `json:"description,omitempty"` }
type Lease struct { ID string `json:"id"`; TargetID string `json:"target_id"`; ExpiresAt time.Time `json:"expires_at"`; Budget int `json:"budget"`; Used int `json:"used"`; Revoked bool `json:"revoked"` }
type Event struct { ID string `json:"id"`; LeaseID string `json:"lease_id"`; Type string `json:"type"`; Method string `json:"method,omitempty"`; URL string `json:"url,omitempty"`; Status int `json:"status,omitempty"`; Payload string `json:"payload,omitempty"`; CreatedAt time.Time `json:"created_at"` }
type RedactRequest struct { Text string `json:"text"`; Rules map[string]string `json:"rules,omitempty"`; Stable bool `json:"stable,omitempty"` }
type RedactResponse struct { Text string `json:"text"`; Tokens map[string]string `json:"tokens"` }
