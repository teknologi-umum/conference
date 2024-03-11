package nocodb_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"conf/nocodb"
	"conf/nocodb/nocodbmock"
	"github.com/rs/zerolog/log"
)

var client *nocodb.Client
var tableId string

func TestMain(m *testing.M) {
	baseUrl := os.Getenv("NOCODB_BASE_URL")
	apiToken := os.Getenv("NOCODB_API_KEY")
	tableId = os.Getenv("NOCODB_TABLE_ID")
	if tableId == "" {
		tableId = "aabbcc"
	}

	var err error = nil
	var mockServer *httptest.Server = nil
	if baseUrl != "" && apiToken != "" {
		client, err = nocodb.NewClient(nocodb.ClientOptions{
			ApiToken: apiToken,
			BaseUrl:  baseUrl,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("creating nocodb client")
			return
		}
	} else {
		mockServer, err = nocodbmock.NewNocoDBMockServer()
		if err != nil {
			mockServer.Close()
			log.Fatal().Err(err).Msg("creating mock server")
			return
		}

		client, err = nocodb.NewClient(nocodb.ClientOptions{
			ApiToken:   "testing",
			BaseUrl:    mockServer.URL,
			HttpClient: mockServer.Client(),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("creating nocodb client")
			return
		}
	}

	exitCode := m.Run()

	if mockServer != nil {
		mockServer.Close()
	}

	os.Exit(exitCode)
}

type testBody struct {
	Id         int64 `json:"Id,omitempty"`
	Title      string
	Age        int
	RandomText string
}

func TestIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	randomText := base64.StdEncoding.EncodeToString(randomBytes)
	payload := testBody{
		Title:      "John Doe",
		Age:        49,
		RandomText: randomText,
	}

	err := client.CreateTableRecords(ctx, tableId, []any{payload})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	var outPayload []testBody
	pageInfo, err := client.ListTableRecords(ctx, tableId, &outPayload, nocodb.ListTableRecordOptions{})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	found := false
	foundPayload := testBody{}
	for _, out := range outPayload {
		if out.RandomText == randomText {
			found = true
			foundPayload = out
		}
	}
	if !found {
		t.Errorf("expecting just inserted entry to be found, got not found")
	}
	if pageInfo.TotalRows <= 0 {
		t.Errorf("expecting pageInfo.TotalRows to be a positive number greater than one, got %d", pageInfo.TotalRows)
	}

	err = client.UpdateTableRecords(ctx, tableId, []any{map[string]any{"Id": foundPayload.Id, "Age": 320}})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	var anotherOutPayload testBody
	err = client.ReadTableRecords(ctx, tableId, strconv.FormatInt(foundPayload.Id, 10), &anotherOutPayload, nocodb.ReadTableRecordsOptions{})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if anotherOutPayload.Age != 320 {
		t.Errorf("expecting Age to be updated to 320, got %d", anotherOutPayload.Age)
	}
}
