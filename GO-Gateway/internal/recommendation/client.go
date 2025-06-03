package recommendation

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RESTRecommendationClient struct {
	baseURL string
	app     *fiber.App
}

func NewRESTRecommendationClient(baseURL string) *RESTRecommendationClient {
	return &RESTRecommendationClient{
		baseURL: baseURL,
		app:     fiber.New(),
	}
}

func (c *RESTRecommendationClient) SendUserAction(action UserAction) error {
	url := fmt.Sprintf("%s/api/actions", c.baseURL)
	jsonData, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal user action: %v", err)
	}

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.SetRequestURI(url)
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(jsonData)

	if err := agent.Parse(); err != nil {
		return fmt.Errorf("failed to parse request: %v", err)
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("failed to send user action: %v", errs[0])
	}

	if code != fiber.StatusOK {
		return fmt.Errorf("failed to send user action: status %d, body: %s", code, string(body))
	}

	return nil
}

func (c *RESTRecommendationClient) GetRecommendations(userID string, limit int) ([]Boardgame, error) {
	url := fmt.Sprintf("%s/api/recommendations", c.baseURL)
	reqBody := map[string]interface{}{
		"user_id": userID,
		"limit":   limit,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.SetRequestURI(url)
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(jsonData)

	if err := agent.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get recommendations: %v", errs[0])
	}

	if code != fiber.StatusOK {
		return nil, fmt.Errorf("failed to get recommendations: status %d, body: %s", code, string(body))
	}

	var response struct {
		Boardgames []Boardgame `json:"boardgames"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Boardgames, nil
}

func (c *RESTRecommendationClient) SendAllBoardgames(boardgames []Boardgame) error {
	url := fmt.Sprintf("%s/api/boardgames", c.baseURL)
	reqBody := map[string][]Boardgame{
		"boardgames": boardgames,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal boardgames: %v", err)
	}

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.SetRequestURI(url)
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(jsonData)

	if err := agent.Parse(); err != nil {
		return fmt.Errorf("failed to parse request: %v", err)
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("failed to send boardgames: %v", errs[0])
	}

	if code != fiber.StatusOK {
		return fmt.Errorf("failed to send boardgames: status %d, body: %s", code, string(body))
	}

	return nil
}

func (c *RESTRecommendationClient) GetAllBoardgames() ([]Boardgame, error) {
	url := fmt.Sprintf("%s/api/boardgames", c.baseURL)

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.SetRequestURI(url)
	req.Header.SetMethod(fiber.MethodGet)

	if err := agent.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get boardgames: %v", errs[0])
	}

	if code != fiber.StatusOK {
		return nil, fmt.Errorf("failed to get boardgames: status %d, body: %s", code, string(body))
	}

	var response struct {
		Boardgames []Boardgame `json:"boardgames"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Boardgames, nil
}

func (c *RESTRecommendationClient) GetPopularBoardgames(limit int) ([]Boardgame, error) {
	url := fmt.Sprintf("%s/api/boardgames/popular?limit=%d", c.baseURL, limit)

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.SetRequestURI(url)
	req.Header.SetMethod(fiber.MethodGet)

	if err := agent.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get popular boardgames: %v", errs[0])
	}

	if code != fiber.StatusOK {
		return nil, fmt.Errorf("failed to get popular boardgames: status %d, body: %s", code, string(body))
	}

	var response struct {
		Boardgames []Boardgame `json:"boardgames"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Boardgames, nil
}

func (c *RESTRecommendationClient) GetUserActions(userID string) ([]UserAction, error) {
	url := fmt.Sprintf("%s/api/actions/user/%s", c.baseURL, userID)

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.SetRequestURI(url)
	req.Header.SetMethod(fiber.MethodGet)

	if err := agent.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get user actions: %v", errs[0])
	}

	if code != fiber.StatusOK {
		return nil, fmt.Errorf("failed to get user actions: status %d, body: %s", code, string(body))
	}

	var result struct {
		Actions []UserAction `json:"actions"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Actions, nil
}

func (c *RESTRecommendationClient) GetBoardgameActions(boardgameID string) ([]UserAction, error) {
	url := fmt.Sprintf("%s/api/actions/boardgame/%s", c.baseURL, boardgameID)

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.SetRequestURI(url)
	req.Header.SetMethod(fiber.MethodGet)

	if err := agent.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get boardgame actions: %v", errs[0])
	}

	if code != fiber.StatusOK {
		return nil, fmt.Errorf("failed to get boardgame actions: status %d, body: %s", code, string(body))
	}

	var result struct {
		Actions []UserAction `json:"actions"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Actions, nil
}
