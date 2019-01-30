package main

type healthCheckResponse struct {
	Status string `json:"status"`
}

type snowflakeResponse struct {
	ID       snowflakeID `json:"id"`
	IDString string      `json:"idStr"`
}
