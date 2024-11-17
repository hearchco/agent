package result

type WebOutput struct {
	webOutputJSON
}

type webOutputJSON struct {
	Web

	FqdnHash          string `json:"fqdn_hash,omitempty"`
	FqdnHashTimestamp string `json:"fqdn_hash_timestamp,omitempty"`
}
