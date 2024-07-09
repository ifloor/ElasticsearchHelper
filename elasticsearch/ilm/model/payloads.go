package model

type IlmResponse struct {
	Indices map[string]IlmIndex `json:"indices"`
}

type IlmIndex struct {
	IndexName               string         `json:"index"`
	Managed                 bool           `json:"managed"`
	Policy                  string         `json:"policy"`
	IndexCreationDateMillis int64          `json:"index_creation_date_millis"`
	TimeSinceIndexCreation  string         `json:"time_since_index_creation"`
	LifecycleDateMillis     int64          `json:"lifecycle_date_millis"`
	Age                     string         `json:"age"`
	Phase                   Phase          `json:"phase"`
	PhaseTimeMillis         int64          `json:"phase_time_millis"`
	Action                  string         `json:"action"`
	ActionTimeMillis        int64          `json:"action_time_millis"`
	Step                    string         `json:"step"`
	PhaseExecution          PhaseExecution `json:"phase_execution"`
}

type PhaseExecution struct {
	Policy               string          `json:"policy"`
	PhaseDefinition      PhaseDefinition `json:"phase_definition"`
	Version              int64           `json:"version"`
	ModifiedDateInMillis int64           `json:"modified_date_in_millis"`
}

type PhaseDefinition struct {
	MinAge  string                           `json:"min_age"`
	Actions map[string]PhaseDefinitionAction `json:"actions"`
}

type PhaseDefinitionAction struct {
	MaxSize  *string `json:"max_size,omitempty"`
	MaxAge   *string `json:"max_age,omitempty"`
	Priority *int64  `json:"priority,omitempty"`
}

type Phase string

const (
	Hot  Phase = "hot"
	Warm Phase = "warm"
	Cold Phase = "cold"
)
