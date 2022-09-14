package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const handleSetProfilePic = "handleSetProfilePic"

type setProfilePicReq struct {
	ProfilePic string `json:"profile_pic"`
}

// handleSetProfilePic handles the request to set a user's profile picture.
// Since I don't have time to implement blob storage for this demo
// instead I'm going to accept a hypothetical URL and move on
func (app *App) handleSetProfilePic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars[urlVarUserID]

		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var profPicData setProfilePicReq
		if err := json.Unmarshal(reqBodyBytes, &profPicData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		// set the profile pic url
		if dErr := app.service.SetUserProfilePic(r.Context(), userID, profPicData.ProfilePic); dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleSetProfilePic).Infof("user %s profile pic set", userID)
		w.WriteHeader(http.StatusNoContent)
	}
}
