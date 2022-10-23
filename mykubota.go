// Package mykubota implements an API SDK matching the MyKubota app
package mykubota

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// taken from MyKubota app on iOS
var (
	AppEndpoint     = "https://app.mykubota.com"
	AppClientID     = "1e74fe67-9753-4f65-b6e4-dd65a8132ea2"
	AppClientSecret = "TCDx0qg5kFQhIdCxW0t1iFlESodtWfaR49vy4JdbYjc"
	oauthConfig     = oauth2.Config{
		ClientID:     AppClientID,
		ClientSecret: AppClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "",
			TokenURL: fmt.Sprintf("%s/oauth/token", AppEndpoint),
		},
		Scopes: []string{"read"},
	}
)

// Client allows location specific access to public content from the MyKubota app
type Client struct {
	client *http.Client
	locale string
	debug  bool
}

// New creates a new MyKubota client for public content in the region specified by the locale
// locale must be of format `{ISO 639-1}-{ISO 3166}`, ie en-US or en-CA
func New(locale string) *Client {
	return &Client{
		client: &http.Client{},
		locale: locale,
		debug:  os.Getenv("DEBUG") != "",
	}
}

// Maintenance contains required maintenance including intervals for a specific model
type Maintenance struct {
	ID                  string `json:"id"`
	CheckPoint          string `json:"checkPoint"`
	Measures            string `json:"measures"`
	FirstCheckValue     int    `json:"firstCheckValue"`
	DisplayIntervalType string `json:"displayIntervalType"`
	IntervalTyp         string `json:"intervalType"`
	IntervalValue       int    `json:"intervalValue"`
	SortOrder           int    `json:"sortOrder"`
}

func (c *Client) MaintenanceSchedule(model string) ([]Maintenance, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/maintenanceSchedule/%s", AppEndpoint, model), nil)
	if err != nil {
		return nil, err
	}
	res := []Maintenance{}
	if err := c.do(req, []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Client) do(req *http.Request, acceptableHTTPCodes []int, res any) error {
	req.Header.Set("version", "2022_R03")
	// locale is used by the backend to filter results for different countries. Ensure it's set to the country you're located in
	req.Header.Set("Accept-Language", s.locale)

	if s.debug {
		bs, _ := httputil.DumpRequest(req, true)
		log.Printf("> %s\n", string(bs))
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if s.debug {
		bs, _ := httputil.DumpResponse(resp, true)
		log.Printf("< %s\n", string(bs))
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

// Session allows location specific access to authenticated content
type Session struct {
	client *http.Client
	Token  *oauth2.Token
	locale string
	debug  bool
}

// SessionFromToken restores a session from an existing token
func (c *Client) SessionFromToken(ctx context.Context, t *oauth2.Token) (*Session, error) {
	return &Session{
		client: oauthConfig.Client(ctx, t),
		Token:  t,
		locale: c.locale,
		debug:  c.debug,
	}, nil
}

// Authenticate performs a password authentication with the MyKubota oauth API
func (c *Client) Authenticate(ctx context.Context, username, password string) (*Session, error) {
	token, err := oauthConfig.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed oauth2: %v", err)
	}
	s := Session{
		client: oauthConfig.Client(ctx, token),
		Token:  token,
		locale: c.locale,
	}
	return &s, nil
}

// User contains basic informations about your MyKubota registration
type User struct {
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	EmailVerified bool   `json:"email_verified"`
	MfaEnabled    bool   `json:"mfa_enabled"`
}

// User fetches the authenticated user for the current session
func (s *Session) User(ctx context.Context) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth/user", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}

	res := User{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, jsonDecodeProcessor(&res)); err != nil {
		return nil, err
	}
	return &res, nil
}

// EquipmentLocation contains location information for telematics enabled equipment
type EquipmentLocation struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	AltitudeMeters       float64 `json:"altitudeMeters"`
	PositionHeadingAngle float64 `json:"positionHeadingAngle"`
}

// EquipmentRestartInhibitStatus contains restart inhibit information for telematics enabled equipment
type EquipmentRestartInhibitStatus struct {
	CanModify       bool   `json:"canModify"`
	CommandStatus   string `json:"commandStatus"`
	EquipmentStatus string `json:"equipmentStatus"`
}

// EquipmentTelematics contains basic telematics data for telematics enabled equipment
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

type ManualEntries []ManualEntry
type VideoEntries []VideoEntry

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
	ManualEntries ManualEntries `json:"manualEntries"`
	VideoEntries  VideoEntries  `json:"videoEntries"`

	Telematics EquipmentTelematics `json:"telematics"`
}

func jsonDecodeProcessor(v any) func(*http.Response) error {
	return func(resp *http.Response) error {
		return json.NewDecoder(resp.Body).Decode(v)
	}
}

func noopProcessor(response *http.Response) error {
	return nil
}

func (s *Session) do(req *http.Request, acceptableHTTPCodes []int, responseProcessor func(*http.Response) error) error {
	req.Header.Set("version", "2022_R03")
	// locale is used by the backend to filter results for different countries. Ensure it's set to the country you're located in
	req.Header.Set("Accept-Language", s.locale)

	if s.debug {
		bs, _ := httputil.DumpRequest(req, true)
		log.Printf("> %s\n", string(bs))
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if s.debug {
		bs, _ := httputil.DumpResponse(resp, true)
		log.Printf("< %s\n", string(bs))
	}
	defer resp.Body.Close()
	isSuccessful := false
	for _, code := range acceptableHTTPCodes {
		isSuccessful = isSuccessful || code == resp.StatusCode
	}
	if !isSuccessful {
		return fmt.Errorf("response code %d didn't match any expected http status codes %v", resp.StatusCode, acceptableHTTPCodes)
	}
	return responseProcessor(resp)
}

// ListEquipment retrieves all equipment registered with the MyKubota app
func (s *Session) ListEquipment(ctx context.Context) ([]Equipment, error) {
	// TODO does the app support pagination? not that I can tell
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/equipment", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}
	var res = []Equipment{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, jsonDecodeProcessor(&res)); err != nil {
		return nil, err
	}
	return res, nil
}

type Settings struct {
	MeasurementUnit string `json:"measurementUnit"`
	// subscribedToAlerts, subscribedToMessages, subscribedToNotifications
}

// Settings loads user settings made in the MyKubota app
func (s *Session) Settings(ctx context.Context) (*Settings, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/settings", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}
	type settingsResponse struct {
		Settings Settings `json:"settings"`
	}
	res := settingsResponse{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, jsonDecodeProcessor(&res)); err != nil {
		return nil, err
	}
	return &res.Settings, nil
}

// GetEquipment fetches a particular equipment by its ID
func (s *Session) GetEquipment(ctx context.Context, id string) (*Equipment, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/equipment/%s", AppEndpoint, id), nil)
	if err != nil {
		return nil, err
	}
	var res = Equipment{}
	if err := s.do(req.WithContext(ctx), []int{http.StatusOK}, jsonDecodeProcessor(&res)); err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteEquipment removes equipment associations for the current user
func (s *Session) DeleteEquipment(ctx context.Context, id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/user/equipment/%s", AppEndpoint, id), nil)
	if err != nil {
		return err
	}
	return s.do(req.WithContext(ctx), []int{http.StatusOK}, noopProcessor)
}

type AddEquipmentRequest struct {
	Model       *Model
	PinOrSerial string
}

type UpdateEquipmentRequest struct {
	EquipmentID string  `json:"id"`
	EngineHours float64 `json:"engineHours"`
	NickName    string  `json:"nickName"`
}

func (s *Session) UpdateEquipment(ctx context.Context, req UpdateEquipmentRequest) (*Equipment, error) {
	bs := bytes.Buffer{}
	if err := json.NewEncoder(&bs).Encode(req); err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/user/equipment/update", AppEndpoint), &bs)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	res := []Equipment{}
	err = s.do(httpReq, []int{http.StatusOK}, jsonDecodeProcessor(&res))
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

// AddEquipment adds equipment associations for the current user
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
	return s.do(req.WithContext(ctx), []int{http.StatusOK}, noopProcessor)
}

type CategoryNode struct {
	ID       int
	Name     string
	ParentID *int

	SubCategories []*CategoryNode
	Models        []Model
}

// GetModelTree returns a category tree containing all machines/ attachments currently offered by Kubota
func (c *Client) GetModelTree(ctx context.Context) ([]*CategoryNode, error) {
	cs, ms, err := c.loadCategoriesAndModels(ctx)
	if err != nil {
		return nil, err
	}

	categoryModels := map[int][]Model{}
	for _, m := range ms {
		vs, ok := categoryModels[m.CategoryID]
		if !ok {
			vs = []Model{}
		}
		categoryModels[m.CategoryID] = append(vs, m)
	}

	roots := []*CategoryNode{}
	categories := map[int]*CategoryNode{}
	for _, c := range cs {
		node := &CategoryNode{
			ID:            c.ID,
			Name:          c.Name,
			SubCategories: []*CategoryNode{},
			Models:        categoryModels[c.ID],
			ParentID:      c.ParentID,
		}
		if c.ParentID == nil {
			roots = append(roots, node)
		}
		categories[c.ID] = node
	}
	for _, node := range categories {
		if node.ParentID == nil {
			continue
		}
		parent := categories[*node.ParentID]
		vs := parent.SubCategories
		parent.SubCategories = append(vs, node)
	}

	return roots, nil
}

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parentId"`
	// heroUrl, fullUrl, iconUrl
}

func (c *Client) loadCategoriesAndModels(ctx context.Context) ([]Category, []Model, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/models", AppEndpoint), nil)
	if err != nil {
		return nil, nil, err
	}
	type modelsResponse struct {
		Categories []Category `json:"categories"`
		Models     []Model    `json:"models"`
	}
	var res = modelsResponse{}
	if err := c.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, nil, err
	}
	return res.Categories, res.Models, nil
}

// Listcategories returns all product categories offered by Kubota
func (s *Client) ListCategories(ctx context.Context) ([]Category, error) {
	cs, _, err := s.loadCategoriesAndModels(ctx)
	return cs, err
}

type Model struct {
	CategoryID            int      `json:"categoryId"`
	Type                  string   `json:"type"`
	CompatibleAttachments []string `json:"compatibleAttachments"`
	// categoryFullUrl, categoryHeroUrl, categoryIconUrl, guideUrl string
	HasFaultCodes           bool          `json:"hasFaultCodes"`
	HasMaintenanceSchedules bool          `json:"hasMaintenanceSchedules"`
	ManualEntries           []ManualEntry `json:"manualEntries"`
	VideoEntries            []VideoEntry  `json:"videoEntries"`
	Model                   string        `json:"model"`
	// modelFullUrl, modelHeroUrl, modelIconUrl string
	// subcategoryFullUrl, subcategoryHeroUrl, subcategoryIconUrl string
	// warrantyUrl string
}

// ListModels returns all machines/ attachments offered by Kubota
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	_, ms, err := c.loadCategoriesAndModels(ctx)
	return ms, err
}

type SearchMachineRequest struct {
	PartialModel string
	Serial       string
}

// SearchMachine performs a location aware search in Kubotas registry for a matching model/ serial combination
func (c *Client) SearchMachine(ctx context.Context, request SearchMachineRequest) (*Model, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/models", AppEndpoint), nil)
	if err != nil {
		return nil, err
	}
	values := req.URL.Query()
	values.Set("partialModel", request.PartialModel)
	values.Set("serial", request.Serial)
	req.URL.RawQuery = values.Encode()

	type modelsResponse struct {
		Models []Model `json:"models"`
	}
	var res = modelsResponse{}
	if err := c.do(req.WithContext(ctx), []int{http.StatusOK}, &res); err != nil {
		return nil, err
	}
	if len(res.Models) < 1 {
		return nil, fmt.Errorf("didn't find a matching model")
	}
	return &res.Models[0], nil
}

type MaintenanceHistory struct {
	ID                   string    `json:"id"`
	IntervalType         string    `json:"intervalType"`
	IntervalValue        int       `json:"intervalValue"`
	CompletedEngineHours float32   `json:"completedEngineHours"`
	Notes                string    `json:"notes"`
	UpdatedDate          time.Time `json:"updatedDate"`
	// map of MaintenanceSchedule id to performed <Y/N>
	MaintenanceCheckList map[string]bool `json:"maintenanceCheckList"`
}

func (s *Session) MaintenanceHistory(equipmentID string) ([]MaintenanceHistory, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/equipment/%s/maintenanceHistory", AppEndpoint, equipmentID), nil)
	if err != nil {
		return nil, err
	}
	res := []MaintenanceHistory{}
	if err := s.do(req, []int{http.StatusOK}, jsonDecodeProcessor(&res)); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Session) RecordMaintenance(equipmentID string, entry MaintenanceHistory) error {
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	// TODO maybe sanity check the interval and checks align?
	payload := bytes.Buffer{}
	if err := json.NewEncoder(&payload).Encode(entry); err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/user/equipment/%s/maintenanceHistory", AppEndpoint, equipmentID), &payload)
	if err != nil {
		return err
	}
	return s.do(req, []int{http.StatusOK}, nil)
}
