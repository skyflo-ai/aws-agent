package realtime

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/DavisAndn/go-aws-crawler/internal/backend"
	"github.com/DavisAndn/go-aws-crawler/internal/config"
)

// CloudTrailEvent represents a minimal view of a CloudTrail event.
type CloudTrailEvent struct {
	DetailType string `json:"detail-type"`
	Source     string `json:"source"`
	Detail     struct {
		EventName         string                 `json:"eventName"`
		RequestParameters map[string]interface{} `json:"requestParameters"`
		ResponseElements  map[string]interface{} `json:"responseElements"`
		// Additional fields can be added as needed.
	} `json:"detail"`
}

// EventDelta is the simplified payload that will be sent to the backend.
type EventDelta struct {
	Action       string                 `json:"action"`       // "create", "update", or "delete"
	ResourceType string                 `json:"resourceType"` // e.g., "EC2Instance", "VPC", etc.
	ResourceID   string                 `json:"resourceId"`
	NewState     map[string]interface{} `json:"newState,omitempty"`
}

// mustMarshal is a helper that marshals data or returns an empty JSON object on error.
func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return b
}

// mapEventToDelta converts a CloudTrailEvent to an EventDelta by examining the source and event name.
// It extracts a resource identifier and sets the appropriate action.
func mapEventToDelta(cte CloudTrailEvent) EventDelta {
	var delta EventDelta
	eventName := strings.ToLower(cte.Detail.EventName)

	switch cte.Source {
	// ----- EC2 -----
	case "aws.ec2":
		delta.ResourceType = "EC2Instance"
		if eventName == "runinstances" {
			delta.Action = "create"
			// Extract the instance ID from responseElements.
			if resp, ok := cte.Detail.ResponseElements["instancesSet"].(map[string]interface{}); ok {
				if items, ok := resp["items"].([]interface{}); ok && len(items) > 0 {
					if first, ok := items[0].(map[string]interface{}); ok {
						if id, ok := first["instanceId"].(string); ok {
							delta.ResourceID = id
						}
					}
				}
			}
		} else if eventName == "terminateinstances" {
			delta.Action = "delete"
			// Extract instance ID from requestParameters.
			if req, ok := cte.Detail.RequestParameters["instanceId"]; ok {
				switch v := req.(type) {
				case []interface{}:
					if len(v) > 0 {
						if id, ok := v[0].(string); ok {
							delta.ResourceID = id
						}
					}
				case string:
					delta.ResourceID = v
				}
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["instanceId"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- VPC -----
	case "aws.vpc":
		delta.ResourceType = "VPC"
		if eventName == "createvpc" {
			delta.Action = "create"
			if resp, ok := cte.Detail.ResponseElements["vpc"].(map[string]interface{}); ok {
				if id, ok := resp["vpcId"].(string); ok {
					delta.ResourceID = id
				}
			}
		} else if eventName == "deletevpc" {
			delta.Action = "delete"
			if req, ok := cte.Detail.RequestParameters["vpcId"].(string); ok {
				delta.ResourceID = req
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["vpcId"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- IAM -----
	case "aws.iam":
		delta.ResourceType = "IAMUser"
		if eventName == "createuser" {
			delta.Action = "create"
			if resp, ok := cte.Detail.ResponseElements["user"].(map[string]interface{}); ok {
				if id, ok := resp["userName"].(string); ok {
					delta.ResourceID = id
				}
			}
		} else if eventName == "deleteuser" {
			delta.Action = "delete"
			if req, ok := cte.Detail.RequestParameters["userName"].(string); ok {
				delta.ResourceID = req
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["userName"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- AutoScaling -----
	case "aws.autoscaling":
		delta.ResourceType = "AutoScalingGroup"
		if eventName == "createautoscalinggroup" {
			delta.Action = "create"
			if resp, ok := cte.Detail.ResponseElements["autoScalingGroupName"].(string); ok {
				delta.ResourceID = resp
			} else if req, ok := cte.Detail.RequestParameters["autoScalingGroupName"].(string); ok {
				delta.ResourceID = req
			}
		} else if eventName == "deleteautoscalinggroup" {
			delta.Action = "delete"
			if req, ok := cte.Detail.RequestParameters["autoScalingGroupName"].(string); ok {
				delta.ResourceID = req
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["autoScalingGroupName"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- ELB -----
	case "aws.elb":
		delta.ResourceType = "LoadBalancer"
		if eventName == "createloadbalancer" {
			delta.Action = "create"
			if resp, ok := cte.Detail.ResponseElements["loadBalancerName"].(string); ok {
				delta.ResourceID = resp
			} else if req, ok := cte.Detail.RequestParameters["loadBalancerName"].(string); ok {
				delta.ResourceID = req
			}
		} else if eventName == "deleteloadbalancer" {
			delta.Action = "delete"
			if req, ok := cte.Detail.RequestParameters["loadBalancerName"].(string); ok {
				delta.ResourceID = req
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["loadBalancerName"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- EKS -----
	case "aws.eks":
		delta.ResourceType = "EKSCluster"
		if eventName == "createcluster" {
			delta.Action = "create"
			if resp, ok := cte.Detail.ResponseElements["cluster"].(map[string]interface{}); ok {
				if id, ok := resp["name"].(string); ok {
					delta.ResourceID = id
				}
			}
		} else if eventName == "deletecluster" {
			delta.Action = "delete"
			if req, ok := cte.Detail.RequestParameters["name"].(string); ok {
				delta.ResourceID = req
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["name"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- ElastiCache -----
	case "aws.elasticache":
		delta.ResourceType = "ElastiCache"
		if eventName == "createcachecluster" {
			delta.Action = "create"
			if resp, ok := cte.Detail.ResponseElements["cacheClusterId"].(string); ok {
				delta.ResourceID = resp
			} else if req, ok := cte.Detail.RequestParameters["cacheClusterId"].(string); ok {
				delta.ResourceID = req
			}
		} else if eventName == "deletecachecluster" {
			delta.Action = "delete"
			if req, ok := cte.Detail.RequestParameters["cacheClusterId"].(string); ok {
				delta.ResourceID = req
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["cacheClusterId"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- Route53 -----
	case "aws.route53":
		delta.ResourceType = "Route53HostedZone"
		if eventName == "createhostedzone" {
			delta.Action = "create"
			if resp, ok := cte.Detail.ResponseElements["hostedZone"].(map[string]interface{}); ok {
				if id, ok := resp["id"].(string); ok {
					delta.ResourceID = id
				}
			}
		} else if eventName == "deletehostedzone" {
			delta.Action = "delete"
			if req, ok := cte.Detail.RequestParameters["id"].(string); ok {
				delta.ResourceID = req
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["id"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- S3 -----
	case "aws.s3":
		delta.ResourceType = "S3Bucket"
		if eventName == "createbucket" {
			delta.Action = "create"
			if resp, ok := cte.Detail.ResponseElements["bucketName"].(string); ok {
				delta.ResourceID = resp
			} else if req, ok := cte.Detail.RequestParameters["bucket"].(string); ok {
				delta.ResourceID = req
			}
		} else if eventName == "deletebucket" {
			delta.Action = "delete"
			if req, ok := cte.Detail.RequestParameters["bucket"].(string); ok {
				delta.ResourceID = req
			}
		} else {
			delta.Action = "update"
			if req, ok := cte.Detail.RequestParameters["bucket"].(string); ok {
				delta.ResourceID = req
			}
		}
	// ----- Fallback -----
	default:
		delta.Action = "update"
		delta.ResourceType = "Unknown"
	}

	// Always include the full detail as NewState.
	var detailMap map[string]interface{}
	if err := json.Unmarshal(mustMarshal(cte.Detail), &detailMap); err == nil {
		delta.NewState = map[string]interface{}{"detail": detailMap}
	}

	return delta
}

// eventHandler is the HTTP handler for incoming CloudTrail events.
func eventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	var cte CloudTrailEvent
	if err := json.Unmarshal(body, &cte); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	log.Printf("Received CloudTrail event: %+v", cte)

	delta := mapEventToDelta(cte)
	if delta.ResourceID == "" {
		log.Printf("No resource ID extracted; ignoring event: %+v", cte)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No resource ID; event ignored"))
		return
	}

	deltaJSON, err := json.Marshal(delta)
	if err != nil {
		http.Error(w, "failed to marshal delta", http.StatusInternalServerError)
		return
	}

	log.Printf("Forwarding delta: %s", string(deltaJSON))

	// Load configuration to get the backend endpoint.
	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "config error", http.StatusInternalServerError)
		return
	}

	// Forward the delta update to the backend.
	err = backend.SendInitialCrawlResults(cfg.BackendEndpoint, deltaJSON)
	if err != nil {
		log.Printf("Error forwarding delta: %v", err)
		http.Error(w, "failed to forward event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event delta forwarded successfully"))
}

// StartEventServer starts an HTTP server that listens for CloudTrail events on the /events endpoint.
func StartEventServer(addr string) {
	http.HandleFunc("/events", eventHandler)
	log.Printf("Starting event server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
