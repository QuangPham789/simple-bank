package api

//
//import (
//	"database/sql"
//	"fmt"
//	mockdb "github.com/QuangPham789/simple-bank/db/mock"
//	db "github.com/QuangPham789/simple-bank/db/sqlc"
//	"github.com/QuangPham789/simple-bank/util"
//	"github.com/golang/mock/gomock"
//	"github.com/jackc/pgx/v5/pgtype"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//)
//
//func TestGetAccount(t *testing.T) {
//	testCases := []struct {
//		name         string
//		accountID    int64
//		buildStubs   func(store *mockdb.MockStore)
//		expectStatus int
//	}{
//		{
//			name:      "OK",
//			accountID: 1,
//			buildStubs: func(store *mockdb.MockStore) {
//				const id int64 = 1
//				store.EXPECT().
//					GetAccount(gomock.Any(), gomock.Eq(id)).
//					Return(randomAccount(id), nil).
//					Times(1)
//			},
//			expectStatus: http.StatusOK,
//		},
//		{
//			name:      "NotFound",
//			accountID: 2,
//			buildStubs: func(store *mockdb.MockStore) {
//				const id int64 = 2
//				store.EXPECT().
//					GetAccount(gomock.Any(), gomock.Eq(id)).
//					Return(db.Account{}, sql.ErrNoRows).
//					Times(1)
//			},
//			expectStatus: http.StatusNotFound,
//		},
//		{
//			name:      "InternalError",
//			accountID: 3,
//			buildStubs: func(store *mockdb.MockStore) {
//				const id int64 = 3
//				store.EXPECT().
//					GetAccount(gomock.Any(), gomock.Eq(id)).
//					Return(db.Account{}, sql.ErrConnDone).
//					Times(1)
//			},
//			expectStatus: http.StatusInternalServerError,
//		},
//		{
//			name:      "BadRequest",
//			accountID: 0,
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					GetAccount(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			expectStatus: http.StatusBadRequest,
//		},
//	}
//
//	for i := range testCases {
//		tc := testCases[i]
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			store := mockdb.NewMockStore(ctrl)
//			tc.buildStubs(store)
//
//			config := util.Config{
//				TokenSymmetricKey:   util.RandomString(32),
//				AccessTokenDuration: time.Minute,
//			}
//
//			server, _ := NewServer(config, store)
//			recorder := httptest.NewRecorder()
//
//			url := fmt.Sprintf("/accounts/%d", tc.accountID)
//			request, err := http.NewRequest(http.MethodGet, url, nil)
//			require.NoError(t, err)
//
//			server.router.ServeHTTP(recorder, request)
//			require.Equal(t, tc.expectStatus, recorder.Code)
//		})
//	}
//}
//
//func randomAccount(id int64) db.Account {
//	now := time.Now()
//
//	// Convert time.Time to pgtype.Timestamptz
//	timestamp := pgtype.Timestamptz{
//		Time:  now,
//		Valid: true, // Ensures the timestamp is considered valid
//	}
//	return db.Account{
//		ID:        id,
//		Owner:     util.RandomOwner(),
//		Balance:   util.RandomMoney(),
//		Currency:  util.RandomCurrency(),
//		CreatedAt: timestamp,
//	}
//}
