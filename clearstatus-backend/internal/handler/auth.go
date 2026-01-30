package handler

import (
	"net/http"
	"time"

	"github.com/clearstatus/backend/internal/auth"
	"github.com/clearstatus/backend/internal/config"
	"github.com/clearstatus/backend/internal/model"
	"github.com/clearstatus/backend/internal/store"
	"github.com/google/uuid"
)

type AuthHandler struct {
	cfg  *config.Config
	user *store.UserStore
}

func NewAuthHandler(cfg *config.Config, user *store.UserStore) *AuthHandler {
	return &AuthHandler{cfg: cfg, user: user}
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

type AuthResponse struct {
	Token string           `json:"token"`
	User  model.UserPublic `json:"user"`
}

func (h *AuthHandler) Google(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req GoogleAuthRequest
	if err := decodeJSON(r, &req); err != nil || req.IDToken == "" {
		Error(w, http.StatusBadRequest, "id_token is required")
		return
	}
	info, err := auth.VerifyGoogleIDToken(r.Context(), req.IDToken, h.cfg.GoogleClientID)
	if err != nil {
		Error(w, http.StatusUnauthorized, "invalid google token: "+err.Error())
		return
	}
	// Find by google_id or email; create if not exists
	u, err := h.user.GetByGoogleID(r.Context(), info.Sub)
	if err != nil {
		u, err = h.user.GetByEmail(r.Context(), info.Email)
		if err != nil {
			// Create new user
			u = &model.User{
				ID:        uuid.New().String(),
				Email:     info.Email,
				Name:      info.Name,
				AvatarURL: info.Picture,
				GoogleID:  &info.Sub,
			}
			if err := h.user.Create(r.Context(), u); err != nil {
				Error(w, http.StatusInternalServerError, "failed to create user")
				return
			}
		} else {
			// Link google_id if not set
			if u.GoogleID == nil {
				u.GoogleID = &info.Sub
				u.Name = info.Name
				u.AvatarURL = info.Picture
				_ = h.user.Update(r.Context(), u)
			}
		}
	} else {
		// Update name/picture if changed
		if u.Name != info.Name || u.AvatarURL != info.Picture {
			u.Name = info.Name
			u.AvatarURL = info.Picture
			_ = h.user.Update(r.Context(), u)
		}
	}
	token, err := auth.IssueJWT(h.cfg.JWTSecret, u.ID, u.Email, 7*24*time.Hour) // 7 days
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to issue token")
		return
	}
	JSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  userToPublic(u),
	})
}

func userToPublic(u *model.User) model.UserPublic {
	return model.UserPublic{ID: u.ID, Email: u.Email, Name: u.Name, AvatarURL: u.AvatarURL}
}
