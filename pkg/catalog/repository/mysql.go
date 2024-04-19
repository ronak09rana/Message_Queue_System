package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"message_queue_system/domain/entity"
	"message_queue_system/domain/interfaces/repository"
)

type ProductRepo struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) repository.IProductRepo {
	return ProductRepo{
		DB: db,
	}
}

func (pr ProductRepo) Upsert(ctx context.Context, product entity.Product) (int, error) {
	conn, err := pr.DB.Conn(ctx)
	if err != nil {
		log.Printf("Error: %v\n, unable_to_db_connect\n\n", err.Error())
		return 0, errors.New("unable to db connect")
	}
	defer conn.Close()

	sqlQuery := "INSERT INTO product(product_name, product_description, product_price, product_images)" +
		" VALUES(?, ?, ?, ?) ON DUPLICATE KEY UPDATE" +
		" product_description=values(product_description), product_price=values(product_price), product_images=values(product_images)"

	productImageBytes, err := json.Marshal(product.Images)
	if err != nil {
		log.Printf("Error: %v\n, unable_to_marshal_array_to_json\n\n", err.Error())
		return 0, errors.New("unable to marshal array to json")
	}

	args := []interface{}{product.Name, product.Description, product.Price, string(productImageBytes)}
	result, err := conn.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		log.Printf("Error: %v\n, unable_to_execute_sql_query\n\n", err.Error())
		return 0, errors.New("unable to execute sql query")
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error: %v\n, unable_to_fetch_last_insertedId\n\n", err.Error())
		return 0, errors.New("unable to fetch last insertedId")
	}
	return int(lastInsertedId), nil
}

func (pr ProductRepo) Get(ctx context.Context, productId int) ([]string, error) {
	conn, err := pr.DB.Conn(ctx)
	if err != nil {
		log.Printf("Error: %v\n, unable_to_db_connect\n\n", err.Error())
		return nil, errors.New("unable to db connect")
	}
	defer conn.Close()

	var jsonString string
	sqlQuery := "SELECT JSON_EXTRACT(product_images, '$') AS prod_img_array FROM product WHERE product_id = ?"
	args := []interface{}{productId}
	err = conn.QueryRowContext(ctx, sqlQuery, args...).Scan(&jsonString)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_execute_query", err.Error())
		return nil, errors.New("unable to execute query")
	}

	var productImages []string
	err = json.Unmarshal([]byte(jsonString), &productImages)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_unmarshal_data", err.Error())
		return nil, errors.New("unable to unmarshal data")
	}
	return productImages, nil
}

func (pr ProductRepo) Save(ctx context.Context, productId int, imagesArr []string) error {
	conn, err := pr.DB.Conn(ctx)
	if err != nil {
		log.Printf("Error: %v\n, unable_to_db_connect\n\n", err.Error())
		return errors.New("unable to db connect")
	}
	defer conn.Close()

	sqlQuery := "UPDATE product SET compressed_product_images = ? WHERE product_id = ?"

	productImageBytes, err := json.Marshal(imagesArr)
	if err != nil {
		log.Printf("Error: %v\n, unable_to_marshal_array_to_json\n\n", err.Error())
		return errors.New("unable to marshal array to json")
	}

	args := []interface{}{string(productImageBytes), productId}
	_, err = conn.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		log.Printf("Error: %v\n, failed_to_execute_sql_query\n\n", err.Error())
		return errors.New("unable to execute sql query")
	}
	return nil
}
