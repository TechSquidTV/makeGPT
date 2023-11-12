package api

import "time"

type ResponseGizmo struct {
	Gizmo struct {
		ID             string `json:"id"`
		OrganizationID string `json:"organization_id"`
		ShortURL       string `json:"short_url"`
		Author         struct {
			UserID          string `json:"user_id"`
			DisplayName     string `json:"display_name"`
			LinkTo          string `json:"link_to"`
			SelectedDisplay string `json:"selected_display"`
			IsVerified      bool   `json:"is_verified"`
		} `json:"author"`
		Voice struct {
			ID string `json:"id"`
		} `json:"voice"`
		WorkspaceID  any    `json:"workspace_id"`
		Model        any    `json:"model"`
		Instructions string `json:"instructions"`
		Settings     struct {
		} `json:"settings"`
		Display struct {
			Name              string `json:"name"`
			Description       string `json:"description"`
			WelcomeMessage    string `json:"welcome_message"`
			PromptStarters    any    `json:"prompt_starters"`
			ProfilePictureURL any    `json:"profile_picture_url"`
			Categories        []any  `json:"categories"`
		} `json:"display"`
		ShareRecipient           string    `json:"share_recipient"`
		UpdatedAt                time.Time `json:"updated_at"`
		LastInteractedAt         any       `json:"last_interacted_at"`
		Tags                     []string  `json:"tags"`
		Version                  int       `json:"version"`
		LiveVersion              int       `json:"live_version"`
		TrainingDisabled         any       `json:"training_disabled"`
		AllowedSharingRecipients []string  `json:"allowed_sharing_recipients"`
		ReviewInfo               any       `json:"review_info"`
		AppealInfo               any       `json:"appeal_info"`
		VanityMetrics            struct {
			NumConversations       int `json:"num_conversations"`
			NumPins                int `json:"num_pins"`
			NumUsersInteractedWith int `json:"num_users_interacted_with"`
		} `json:"vanity_metrics"`
	} `json:"gizmo"`
	Tools           []any `json:"tools"`
	Files           []any `json:"files"`
	ProductFeatures struct {
		Attachments struct {
			Type                  string   `json:"type"`
			AcceptedMimeTypes     []string `json:"accepted_mime_types"`
			ImageMimeTypes        []string `json:"image_mime_types"`
			CanAcceptAllMimeTypes bool     `json:"can_accept_all_mime_types"`
		} `json:"attachments"`
	} `json:"product_features"`
}

type GPTTool string

const (
	dalle   GPTTool = "dalle"
	browser GPTTool = "browser"
)
