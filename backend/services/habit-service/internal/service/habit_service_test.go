package service

import (
	"fmt"
	"habit-service/db"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestHabitService_GetHabitByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}

	defer mockDB.Close()

	habitID := uuid.New()
	userID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).
				AddRow(userID, userID, "habito 1", "descrição", "teste.png", time.Now()),
		)
	queries := db.New(mockDB)
	fmt.Println(queries)
	/* habitService := repository.NewHabitRepository(queries)
	routineService := repository.NewRoutineRepository(queries)
	service := NewHabitService(habitService, routineService) */
}
