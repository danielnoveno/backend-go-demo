package user_test

import (
	"deeply/service/user"
)

import (
	"bytes"
	"deeply/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T){
	userStore := newMockUserStore()
	handler := user.NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T){
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName: "123",
			Email: "123",
			Password: "password",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest{
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName: "123",
			Email: "123@gmail.com",
			Password: "password",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated{
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}
type mockUserStore struct {
	users map[string]*types.User
}

func newMockUserStore() *mockUserStore {
	return &mockUserStore{
		users: make(map[string]*types.User),
	}
}

type User struct {
	Email string
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID (id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	m.users[user.Email] = &user
	return nil
}