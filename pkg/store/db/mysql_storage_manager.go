package db

import (
	"context"
	"errors"

	"github.com/siyu-nuro/go-example/pkg/entity"
	"github.com/siyu-nuro/go-example/pkg/store"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

var (
	errPreparedStatementCreation              = errors.New("failed to create prepared statement")
	errIdTranslation                          = errors.New("id translation failed")
	errPreparedStatementCompletion            = errors.New("prepared statement failed to execute and/or scan")
)

// storageManager implements the store.StorageManager interface
type storageManager struct {
	db                                            *sqlx.DB
	logger                                        *zap.SugaredLogger
	createOrder                                   *sqlx.NamedStmt
}

// NewStorageManager constructs a StorageManager for MySQL
func NewMySQLStorageManager(db *sqlx.DB, logger *zap.SugaredLogger) (store.StorageManager, error) {
	createOrder, err := db.PrepareNamed(
		"INSERT INTO `order` (pickup_address_street_1, pickup_address_street_2, pickup_address_city, pickup_address_state, pickup_address_zip, pickup_location, " +
			"dropoff_address_street_1, dropoff_address_street_2, dropoff_address_city, dropoff_address_state, dropoff_address_zip, dropoff_location, " +
			"earliest_pickup_utc, earliest_dropoff_utc, latest_dropoff_utc, dispatcher_notes, merchant_notes, last_modified_at_utc, store_order_id, store_hmi_access_code, consumer_hmi_access_code, store_id, " +
			"consumer_first_name, consumer_last_name, consumer_email, consumer_notification_type, consumer_phone, store_notification_type, store_phone, tracking_url, created_at_utc, " +
			"grocery_tote_count, large_item_count) " +
			"VALUES (:pickup_address_street_1, :pickup_address_street_2, :pickup_address_city, :pickup_address_state, :pickup_address_zip, ST_GeomFromWKB(:pickup_location), " +
			":dropoff_address_street_1, :dropoff_address_street_2, :dropoff_address_city, :dropoff_address_state, :dropoff_address_zip, ST_GeomFromWKB(:dropoff_location), " +
			":earliest_pickup_utc, :earliest_dropoff_utc, :latest_dropoff_utc, :dispatcher_notes, :merchant_notes, UTC_TIMESTAMP, :store_order_id, :store_hmi_access_code, :consumer_hmi_access_code, :store_id, " +
			":consumer_first_name, :consumer_last_name, :consumer_email, :consumer_notification_type, :consumer_phone, :store_notification_type, :store_phone, :tracking_url, UTC_TIMESTAMP, " +
			":grocery_tote_count, :large_item_count)")
	if err != nil {
		logger.Error("Unable to prepare statement for createOrder: ", err)
		return nil, errPreparedStatementCreation
	}

	return &storageManager{
		db:                                            db,
		logger:                                        logger,
		createOrder:                                   createOrder,
	}, nil
}

// CreateOrder creates an order with given order info and returns new orderID if succeeds
func (s *storageManager) CreateOrder(ctx context.Context, inputOrderDetail *entity.InputOrderDetail) (*entity.OrderID, error) {
    // implement

	entityOrderID := entity.OrderID("newOrderID")
	return &entityOrderID, nil
}
