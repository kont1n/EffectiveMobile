package models

type Song struct {
	ID      string `json:"id"`
	Name    string `json:"song"`
	Artist  string `json:"group"`
	Release string `json:"release_date"`
	Text    string `json:"text"`
	Link    string `json:"link"`
}

type Topic struct {
	Namespace string `json:"namespace"`
	Topic     string `json:"topic"`
	Private   bool   `json:"private"`
	ViewCount int64  `json:"viewCount"`
}

// TopicsPostRequest is the request to POST /topics.
type TopicsPostRequest struct {
	Topic
}

type TopicsPostResponse struct {
	ID string `json:"id"`
}

// TopicsGetResponse is the response to GET /topics.
type TopicsGetResponse struct {
	Topics []TopicRecord `json:"topics"`
}

type TopicRecord struct {
	// ID of the topic record.
	ID string `json:"id"`
	Topic
}
