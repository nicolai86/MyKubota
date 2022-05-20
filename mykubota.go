// Package mykubota implements an API SDK matching the MyKubota app
package mykubota

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// taken from MyKubota app on iOS
var AppEndpoint = "https://app.mykubota.com"
var AppClientID = "1e74fe67-9753-4f65-b6e4-dd65a8132ea2"
var AppClientSecret = "TCDx0qg5kFQhIdCxW0t1iFlESodtWfaR49vy4JdbYjc"

type Session struct {
	client *http.Client
	token  *oauth2.Token
}

func New(ctx context.Context, username, password string) (*Session, error) {
	cfg := oauth2.Config{
		ClientID:     AppClientID,
		ClientSecret: AppClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "",
			TokenURL: fmt.Sprintf("%s/oauth/token", AppEndpoint),
		},
		Scopes: []string{"read"},
	}
	token, err := cfg.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed oauth2: %v", err)
	}
	s := Session{
		client: cfg.Client(ctx, token),
		token:  token,
	}
	return &s, nil
}

type User struct {
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	EmailVerified bool   `json:"email_verified"`
	MfaEnabled    bool   `json:"mfa_enabled"`
}

func (s *Session) User(ctx context.Context) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth/user", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}

	res := User{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type EquipmentLocation struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	AltitudeMeters       float64 `json:"altitudeMeters"`
	PositionHeadingAngle float64 `json:"positionHeadingAngle"`
}

type EquipmentRestartInhibitStatus struct {
	CanModify       bool   `json:"canModify"`
	CommandStatus   string `json:"commandStatus"`
	EquipmentStatus string `json:"equipmentStatus"`
}

type EquipmentTelematics struct {
	LocationTime             time.Time                     `json:"locationTime"`
	CumulativeOperatingHours float64                       `json:"cumulativeOperatingHours"`
	Location                 EquipmentLocation             `json:"location"`
	EngineRunning            bool                          `json:"engineRunning"`
	FuelTempCelsius          int                           `json:"fuelTempCelsius"`
	FuelRemainingPercent     int                           `json:"fuelRemainingPercent"`
	DEFTempCelsius           int                           `json:"defTempCelsius"`
	DEFQualityPercent        float64                       `json:"defQualityPercent"`
	DEFRemainingPercent      float64                       `json:"defRemainingPercent"`
	DEFPressureKPascal       float64                       `json:"defPressureKPascal"`
	EngineRPM                int                           `json:"engineRPM"`
	CoolantTempCelsius       int                           `json:"coolantTempCelsius"`
	HydraulicTempCelsius     int                           `json:"hydraulicTempCelsius"`
	ExtPowerVolts            float64                       `json:"extPowerVolts"`
	AirInletTempCelsius      float64                       `json:"airInletTempCelsius"`
	AmbientAirTempCelsius    float64                       `json:"ambientAirTempCelsius"`
	RunNumber                int                           `json:"runNumber"`
	MotionState              string                        `json:"motionState"`
	FaultCodes               []interface{}                 `json:"faultCodes"`
	RestartInhibitStatus     EquipmentRestartInhibitStatus `json:"restartInhibitStatus"`
	InsideGeofences          []interface{}                 `json:"insideGeofences"`
}

type ManualEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VideoEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Equipment struct {
	ID                      string  `json:"id"`
	Model                   string  `json:"model"`
	CategoryID              int     `json:"categoryId"`
	Category                string  `json:"category"`
	SubCategory             string  `json:"subcategory"`
	IdentifierType          string  `json:"identifierType"`
	Type                    string  `json:"type"`
	PinOrSerial             string  `json:"pinOrSerial"`
	Pin                     string  `json:"pin"`
	Serial                  string  `json:"serial"`
	Nickname                string  `json:"nickName"`
	UserEnteredEngineHours  float64 `json:"userEnteredEngineHours"`
	HasTelematics           bool    `json:"hasTelematics"`
	HasFaultCodes           bool    `json:"hasFaultCodes"`
	HasMaintenanceSchedules bool    `json:"hasMaintenanceSchedules"`

	// TODO - no use for these today, but they exist
	// "modelHeroUrl"
	// "modelFullUrl"
	// "modelIconUrl"
	// "warrantyUrl"
	// "guideUrl"
	ManualEntries []ManualEntry `json:"manualEntries"`
	VideoEntries  []VideoEntry  `json:"videoEntries"`

	Telematics EquipmentTelematics `json:"telematics"`
}

func (s *Session) do(req *http.Request, acceptableHTTPCodes []int, res any) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	isSuccessful := false
	for _, code := range acceptableHTTPCodes {
		isSuccessful = isSuccessful || code == resp.StatusCode
	}
	if !isSuccessful {
		return fmt.Errorf("response code %d didn't match any expected http status codes %v", resp.StatusCode, acceptableHTTPCodes)
	}
	return json.NewDecoder(resp.Body).Decode(res)
}

func (s *Session) ListEquipment(ctx context.Context) ([]Equipment, error) {
	// TODO does the app support pagination? not that I can tell
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/equipment", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}
	var res = []Equipment{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

type Settings struct {
	MeasurementUnit string `json:"measurementUnit"`
	// subscribedToAlerts, subscribedToMessages, subscribedToNotifications
}

func (s *Session) Settings(ctx context.Context) (*Settings, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/settings", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}
	type settingsResponse struct {
		Settings Settings `json:"settings"`
	}
	res := settingsResponse{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	return &res.Settings, nil
}

func (s *Session) GetEquipment(ctx context.Context, id string) (*Equipment, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/equipment/%s", AppEndpoint, id), nil)
	if err != nil {
		return nil, err
	}
	var res = Equipment{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Session) DeleteEquipment(ctx context.Context, id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/user/equipment/%s", AppEndpoint, id), nil)
	if err != nil {
		return err
	}
	return s.do(req.WithContext(ctx), []int{http.StatusOK}, struct{}{})
}

type AddEquipmentRequest struct {
	Model       *Model
	PinOrSerial string
}

func (s *Session) AddEquipment(ctx context.Context, request AddEquipmentRequest) error {
	type addMachineRequest struct {
		Model       string `json:"model"`
		PinOrSerial string `json:"pinOrSerial"`
		Identifier  string `json:"identifierType"`
		EngineHours int    `json:"engineHours"`
		Type        string `json:"type"`
	}
	bs := bytes.Buffer{}
	json.NewEncoder(&bs).Encode(addMachineRequest{
		Model:       request.Model.Model,
		PinOrSerial: request.PinOrSerial,
		Identifier:  "Serial",
		Type:        "machine",
	})

	req, err := http.NewRequest("POST", "/api/user/equipment/addFromScan", bytes.NewReader(bs.Bytes()))
	if err != nil {
		return err
	}
	return s.do(req.WithContext(ctx), []int{http.StatusOK}, struct{}{})
}

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parentId"`
	// heroUrl, fullUrl, iconUrl
}

func (s *Session) Categories(ctx context.Context) ([]Category, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/models", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("version", "2021_R06")
	type modelsResponse struct {
		Categories []Category `json:"categories"`
	}
	var res = modelsResponse{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	return res.Categories, nil
}

type Model struct {
	Category    string `json:"category"`
	SubCategory string `json:"subcategory"`
	// categoryFullUrl, categoryHeroUrl, categoryIconUrl, guideUrl string
	HasFaultCodes           bool          `json:"hasFaultCodes"`
	HasMaintenanceSchedules bool          `json:"hasMaintenanceSchedules"`
	ManualEntries           []ManualEntry `json:"manualEntries"`
	Model                   string        `json:"model"`
	// modelFullUrl, modelHeroUrl, modelIconUrl string
	// subcategoryFullUrl, subcategoryHeroUrl, subcategoryIconUrl string
	VideoEntries []VideoEntry `json:"videoEntries"`
	// warrantyUrl string
}

func (s *Session) Models(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/models", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("version", "2021_R06")
	type modelsResponse struct {
		Models []Model `json:"models"`
	}
	var res = modelsResponse{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	return res.Models, nil
}

type SearchMachineRequest struct {
	PartialModel string
	Serial       string
	Locale       string
}

func (s *Session) SearchMachine(ctx context.Context, request SearchMachineRequest) (*Model, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/models", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("version", "2021_R06")
	req.Header.Set("Accept-Language", request.Locale)

	values := req.URL.Query()
	values.Set("partialModel", request.PartialModel)
	values.Set("serial", request.Serial)
	req.URL.RawQuery = values.Encode()

	type modelsResponse struct {
		Models []Model `json:"models"`
	}
	var res = modelsResponse{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	if len(res.Models) < 1 {
		return nil, fmt.Errorf("didn't find a matching model")
	}
	return &res.Models[0], nil
}
