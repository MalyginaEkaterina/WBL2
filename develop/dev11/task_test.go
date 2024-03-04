package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "Positive test with correct parameters",
			body: "user_id=34&name=action&date=2024-03-04",
			want: want{statusCode: 200},
		},
		{
			name: "Negative test with empty name parameter",
			body: "user_id=34&date=2024-03-04",
			want: want{statusCode: 400},
		},
		{
			name: "Negative test with empty user_id parameter",
			body: "name=action&date=2024-03-04",
			want: want{statusCode: 400},
		},
		{
			name: "Negative test with empty date parameter",
			body: "user_id=34&name=action",
			want: want{statusCode: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := getHandler()
			ts := httptest.NewServer(handler)
			defer ts.Close()

			resp := makePostRequest(ts, handler, "/create_event/", tt.body)

			assert.Equal(t, tt.want.statusCode, resp.Code)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tt.want.statusCode == 200 {
				var respOK Response
				err := json.Unmarshal(respBody, &respOK)
				require.NoError(t, err)
				result := respOK.Result.(map[string]interface{})
				require.EqualValues(t, Created, result["status"])
			} else {
				var respErr ErrorResponse
				err := json.Unmarshal(respBody, &respErr)
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name       string
		body       string
		doNotAddID bool
		want       want
	}{
		{
			name: "Positive test with correct parameters",
			body: "user_id=34&name=action2&date=2024-03-05",
			want: want{statusCode: 200},
		},
		{
			name:       "Negative test with empty id parameter",
			body:       "user_id=34&name=action2&date=2024-03-05",
			doNotAddID: true,
			want:       want{statusCode: 400},
		},
		{
			name: "Negative test with empty name parameter",
			body: "user_id=34&date=2024-03-05",
			want: want{statusCode: 400},
		},
		{
			name: "Negative test with empty user_id parameter",
			body: "name=action2&date=2024-03-05",
			want: want{statusCode: 400},
		},
		{
			name: "Negative test with wrong user_id parameter",
			body: "user_id=33&name=action2&date=2024-03-05",
			want: want{statusCode: 400},
		},
		{
			name:       "Negative test with wrong id parameter",
			body:       "user_id=34&name=action2&date=2024-03-05&id=dfghjkl",
			doNotAddID: true,
			want:       want{statusCode: 400},
		},
		{
			name: "Negative test with empty date parameter",
			body: "user_id=34&name=action2",
			want: want{statusCode: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := getHandler()
			ts := httptest.NewServer(handler)
			defer ts.Close()

			id, err := createEventAndGetID(ts, handler, "user_id=34&name=action&date=2024-03-04")
			require.NoError(t, err)

			body := tt.body
			if !tt.doNotAddID {
				body = body + "&id=" + id
			}

			resp := makePostRequest(ts, handler, "/update_event/", body)

			assert.Equal(t, tt.want.statusCode, resp.Code)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tt.want.statusCode == 200 {
				var respOK Response
				err := json.Unmarshal(respBody, &respOK)
				require.NoError(t, err)
				result := respOK.Result.(map[string]interface{})
				require.EqualValues(t, Updated, result["status"])
			} else {
				var respErr ErrorResponse
				err := json.Unmarshal(respBody, &respErr)
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name       string
		body       string
		doNotAddID bool
		want       want
	}{
		{
			name: "Positive test with correct parameters",
			body: "user_id=34",
			want: want{statusCode: 200},
		},
		{
			name:       "Negative test with empty id parameter",
			body:       "user_id=34",
			doNotAddID: true,
			want:       want{statusCode: 400},
		},
		{
			name: "Negative test with empty user_id parameter",
			body: "",
			want: want{statusCode: 400},
		},
		{
			name: "Negative test with wrong user_id parameter",
			body: "user_id=33",
			want: want{statusCode: 400},
		},
		{
			name:       "Negative test with wrong id parameter",
			body:       "user_id=34&id=dfghjkl",
			doNotAddID: true,
			want:       want{statusCode: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := getHandler()
			ts := httptest.NewServer(handler)
			defer ts.Close()

			id, err := createEventAndGetID(ts, handler, "user_id=34&name=action&date=2024-03-04")
			require.NoError(t, err)

			body := tt.body
			if !tt.doNotAddID {
				body = body + "&id=" + id
			}

			resp := makePostRequest(ts, handler, "/delete_event/", body)

			assert.Equal(t, tt.want.statusCode, resp.Code)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tt.want.statusCode == 200 {
				var respOK Response
				err := json.Unmarshal(respBody, &respOK)
				require.NoError(t, err)
				result := respOK.Result.(map[string]interface{})
				require.EqualValues(t, Deleted, result["status"])
			} else {
				var respErr ErrorResponse
				err := json.Unmarshal(respBody, &respErr)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEventPerDay(t *testing.T) {
	type want struct {
		statusCode int
		lenRes     int
	}
	tests := []struct {
		name  string
		query string
		want  want
	}{
		{
			name:  "Positive test with correct parameters",
			query: "user_id=34&date=2024-03-04",
			want:  want{statusCode: 200, lenRes: 2},
		},
		{
			name:  "Positive test with empty result",
			query: "user_id=34&date=2024-03-05",
			want:  want{statusCode: 200, lenRes: 0},
		},
		{
			name:  "Negative test with empty user_id parameter",
			query: "date=2024-03-04",
			want:  want{statusCode: 400},
		},
		{
			name:  "Negative test with empty date parameter",
			query: "user_id=34",
			want:  want{statusCode: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := getHandler()
			ts := httptest.NewServer(handler)
			defer ts.Close()

			_, err := createEventAndGetID(ts, handler, "user_id=34&name=action&date=2024-03-04")
			require.NoError(t, err)
			_, err = createEventAndGetID(ts, handler, "user_id=34&name=action2&date=2024-03-04")
			require.NoError(t, err)
			_, err = createEventAndGetID(ts, handler, "user_id=34&name=action3&date=2024-03-03")
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodGet, ts.URL+"/events_for_day/?"+tt.query, nil)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, request)

			assert.Equal(t, tt.want.statusCode, resp.Code)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tt.want.statusCode == 200 {
				var respOK Response
				err := json.Unmarshal(respBody, &respOK)
				require.NoError(t, err)
				result := respOK.Result.([]interface{})
				require.EqualValues(t, tt.want.lenRes, len(result))
			} else {
				var respErr ErrorResponse
				err := json.Unmarshal(respBody, &respErr)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEventPerWeek(t *testing.T) {
	type want struct {
		statusCode int
		lenRes     int
	}
	tests := []struct {
		name  string
		query string
		want  want
	}{
		{
			name:  "Positive test with correct parameters",
			query: "user_id=34&date=2024-03-04",
			want:  want{statusCode: 200, lenRes: 2},
		},
		{
			name:  "Positive test with empty result",
			query: "user_id=34&date=2024-02-26",
			want:  want{statusCode: 200, lenRes: 0},
		},
		{
			name:  "Negative test with empty user_id parameter",
			query: "date=2024-03-04",
			want:  want{statusCode: 400},
		},
		{
			name:  "Negative test with empty date parameter",
			query: "user_id=34",
			want:  want{statusCode: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := getHandler()
			ts := httptest.NewServer(handler)
			defer ts.Close()

			_, err := createEventAndGetID(ts, handler, "user_id=34&name=action&date=2024-03-04")
			require.NoError(t, err)
			_, err = createEventAndGetID(ts, handler, "user_id=34&name=action2&date=2024-03-10")
			require.NoError(t, err)
			_, err = createEventAndGetID(ts, handler, "user_id=34&name=action3&date=2024-03-11")
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodGet, ts.URL+"/events_for_week/?"+tt.query, nil)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, request)

			assert.Equal(t, tt.want.statusCode, resp.Code)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tt.want.statusCode == 200 {
				var respOK Response
				err := json.Unmarshal(respBody, &respOK)
				require.NoError(t, err)
				result := respOK.Result.([]interface{})
				require.EqualValues(t, tt.want.lenRes, len(result))
			} else {
				var respErr ErrorResponse
				err := json.Unmarshal(respBody, &respErr)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEventPerMonth(t *testing.T) {
	type want struct {
		statusCode int
		lenRes     int
	}
	tests := []struct {
		name  string
		query string
		want  want
	}{
		{
			name:  "Positive test with correct parameters",
			query: "user_id=34&year=2024&month=3",
			want:  want{statusCode: 200, lenRes: 2},
		},
		{
			name:  "Positive test with empty result",
			query: "user_id=34&year=2024&month=2",
			want:  want{statusCode: 200, lenRes: 0},
		},
		{
			name:  "Negative test with empty user_id parameter",
			query: "year=2024&month=3",
			want:  want{statusCode: 400},
		},
		{
			name:  "Negative test with empty year parameter",
			query: "user_id=34&month=3",
			want:  want{statusCode: 400},
		},
		{
			name:  "Negative test with empty month parameter",
			query: "user_id=34&year=2024",
			want:  want{statusCode: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := getHandler()
			ts := httptest.NewServer(handler)
			defer ts.Close()

			_, err := createEventAndGetID(ts, handler, "user_id=34&name=action&date=2024-03-01")
			require.NoError(t, err)
			_, err = createEventAndGetID(ts, handler, "user_id=34&name=action2&date=2024-03-31")
			require.NoError(t, err)
			_, err = createEventAndGetID(ts, handler, "user_id=34&name=action3&date=2024-04-11")
			require.NoError(t, err)
			_, err = createEventAndGetID(ts, handler, "user_id=34&name=action3&date=2023-03-11")
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodGet, ts.URL+"/events_for_month/?"+tt.query, nil)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, request)

			assert.Equal(t, tt.want.statusCode, resp.Code)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tt.want.statusCode == 200 {
				var respOK Response
				err := json.Unmarshal(respBody, &respOK)
				require.NoError(t, err)
				result := respOK.Result.([]interface{})
				require.EqualValues(t, tt.want.lenRes, len(result))
			} else {
				var respErr ErrorResponse
				err := json.Unmarshal(respBody, &respErr)
				require.NoError(t, err)
			}
		})
	}
}

func createEventAndGetID(ts *httptest.Server, handler http.Handler, body string) (string, error) {
	resp := makePostRequest(ts, handler, "/create_event/", body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var createResp Response
	err = json.Unmarshal(respBody, &createResp)
	if err != nil {
		return "", err
	}
	id := createResp.Result.(map[string]interface{})["id"].(string)
	return id, nil
}

func makePostRequest(ts *httptest.Server, handler http.Handler, path string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, ts.URL+path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, request)
	return resp
}
