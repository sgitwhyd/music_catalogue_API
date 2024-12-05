package repositorys

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sgitwhyd/music-catalogue/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_userRepository_Upsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	

	type args struct {
		model models.User
	}

	tests := []struct {
		name    string
		r       *userRepository
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		// TODO: Add test cases.
		{
			name: "should fail to insert user due to duplicate username",
			args: args{
				model: models.User{
					Email:    "developer@testing.com",
					Password: "password",
					Username: "developer",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					args.model.Email,
					args.model.Username,
					args.model.Password,
				).WillReturnError(fmt.Errorf("UNIQUE constraint failed: users.username"))
				mock.ExpectRollback()
			},
		},
		{
			name: "should fail to insert user due to duplicate email",
			args: args{
				model: models.User{
					Email:    "developer@testing.com",
					Password: "password",
					Username: "developer",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					args.model.Email,
					args.model.Username,
					args.model.Password,
				).WillReturnError(fmt.Errorf("UNIQUE constraint failed: users.email"))
				mock.ExpectRollback()
			},
		},
		{
			name: "should insert new user to database",
			args: args{
				model: models.User{
					Email:    "developer@testing.com",
					Password: "password",
					Username: "developer",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					args.model.Email,
					args.model.Username,
					args.model.Password,
				).WillReturnRows(
					sqlmock.NewRows([]string{"id"}).AddRow(1),
				)
				mock.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &userRepository{
				db: gormDB,
			}
			if err := r.Upsert(tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_userRepository_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()

	type args struct {
		email    string
		username string
		id       uint
	}

	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "should successfully find by username",
			args: args{
				email:    "",
				username: "developer",
				id:       0,
			},
			want: &models.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Email:    "developer@gmail.com",
				Username: "developer",
				Password: "password",
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" .*`).
					WithArgs(args.email, args.username, args.id, 1).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "created_at", "updated_at", "email", "username", "password"},
					).AddRow(1, now, now, "developer@gmail.com", "developer", "password"))
			},
		},
		{
			name: "should successfully find by id",
			args: args{
				email:    "",
				username: "",
				id:       1,
			},
			want: &models.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Email:    "developer@gmail.com",
				Username: "developer",
				Password: "password",
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" .*`).
					WithArgs(args.email, args.username, args.id, 1).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "created_at", "updated_at", "email", "username", "password"},
					).AddRow(1, now, now, "developer@gmail.com", "developer", "password"))
			},
		},
		{
			name: "should successfully find by email",
			args: args{
				email:    "developer@gmail.com",
				username: "",
				id:       0,
			},
			want: &models.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Email:    "developer@gmail.com",
				Username: "developer",
				Password: "password",
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" .*`).
					WithArgs(args.email, args.username, args.id, 1).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "created_at", "updated_at", "email", "username", "password"},
					).AddRow(1, now, now, "developer@gmail.com", "developer", "password"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			r := &userRepository{
				db: gormDB,
			}
			got, err := r.Find(tt.args.email, tt.args.username, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.Find() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

