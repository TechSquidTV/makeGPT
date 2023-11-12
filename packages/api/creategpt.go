package api

import (
	"encoding/json"
	"strings"

	"github.com/TechSquidTV/makeGPT/packages/config"
	"github.com/charmbracelet/log"
)

type CreateGPTDisplayInfo struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	WelcomeMessage string `json:"welcome_message"`
}

type CreateGPTTool struct {
	Type GPTTool `json:"type"`
}

type CreateGPTPayload struct {
	Instructions string               `json:"instructions"`
	Display      CreateGPTDisplayInfo `json:"display"`
	Tools        []CreateGPTTool      `json:"tools"`
	Files        []any                `json:"files"` // remains the same, adjust as needed
}

func CreateGPT(config config.GPTConfig) (ResponseGizmo, error) {
	configJSON, err := json.Marshal(config.Config)
	if err != nil {
		log.Error("While attempting to preview the config for a debug statement, we were unable to marshal the config as valid JSON")
		log.Fatalf("Failed to marshal config: \n%v\n", err)
	}

	log.Debugf("Creating GPT with config: \n%v\n", string(configJSON))
	payload := CreateGPTPayload{
		Instructions: config.Config.Instructions,
		Display: CreateGPTDisplayInfo{
			Name:           config.Config.Name,
			Description:    config.Config.Description,
			WelcomeMessage: "Hello",
		},
		Tools: []CreateGPTTool{
			{
				Type: dalle,
			},
			{
				Type: browser,
			},
		},
		Files: []any{},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal payload: \n%v\n", err)
		return ResponseGizmo{}, err
	}
	log.Debugf("Posting payload to backend-api/gizmos: \n%v\n", string(payloadJSON))
	req, err := NewAuthorizedRequest("POST", "backend-api/gizmos", strings.NewReader(string(payloadJSON)))
	if err != nil {
		log.Warnf("Request:  \n%v\n", req)
		log.Fatalf("Failed to create GPT: \n%v\n", err)
		return ResponseGizmo{}, err
	}
	resp, err := Client().Do(req)
	if err != nil {
		log.Warnf("Response: \n%v\n", resp)
		log.Fatalf("Failed to create GPT: \n%v\n", err)
		return ResponseGizmo{}, err
	}

	defer resp.Body.Close()

	var reponseGizmo ResponseGizmo

	if resp.StatusCode != 200 {
		log.Errorf("Status code: \n%v\n", resp.StatusCode)
		log.Warnf("Request Body: \n%v\n", payload)
		log.Warnf("Request Headers: \n%v\n", req.Header)
		log.Errorf("Response: \n%v\n", resp)
		log.Fatalf("Failed to create GPT: \n%v\n", err)
		return ResponseGizmo{}, err
	}

	err = json.NewDecoder(resp.Body).Decode(&reponseGizmo)
	if err != nil {
		log.Warnf("Request Headers: \n%v\n", req.Header)
		log.Warnf("Request Body: \n%v\n", payload)
		log.Warnf("Response: \n%v\n", reponseGizmo)
		log.Fatal("Failed to decode GPT response", err)
		return ResponseGizmo{}, err
	}

	log.Info("Created GPT " + config.Config.Name + " at " + reponseGizmo.Gizmo.ShortURL)
	return reponseGizmo, nil
}
