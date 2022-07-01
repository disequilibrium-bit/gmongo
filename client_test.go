package gmongo

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_NewDatabase(t *testing.T) {
	client := NewClient()
	database, _ := client.NewDatabase("", "", SetMaxPoolSize(10), SetMaxConnIdleTime(5*time.Second), SetMinPoolSize(1))
	assert.NotNil(t, database)
}
