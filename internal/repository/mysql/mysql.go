package mysql

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lazylex/watch-store-store/internal/config"
	"github.com/lazylex/watch-store-store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store-store/internal/dto"
	"github.com/lazylex/watch-store-store/internal/helpers/constants/prefixes"
	"github.com/lazylex/watch-store-store/internal/logger"
	"github.com/lazylex/watch-store-store/internal/ports/repository"
	"github.com/lazylex/watch-store-store/internal/service"
	"log/slog"
	"os"
	"strings"
	"time"
)

const txIsolationLevel = sql.LevelSerializable

type Repository struct {
	db *sql.DB
}

func mysqlErr(text string) error {
	return errors.New(prefixes.MySQLPrefix + text)
}

var (
	ErrNilConfigPointer = mysqlErr("nil config pointer")
)

// DB возвращает структуру DB репозитория.
func (r *Repository) DB() *sql.DB {
	return r.db
}

// Close закрывает пул подключений к БД.
func (r *Repository) Close() error {
	log := slog.With(slog.String(logger.OPLabel, "mysql.Close"))
	err := r.db.Close()
	if err != nil {
		log.Error("error close repository")
	}
	log.Info("close repository")
	return r.ConvertToCommonErr(err)
}

// createDSN создает строку подключения к БД из параметров, переданных в конфигурации.
func createDSN(cfg *config.Storage) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&interpolateParams=true",
		cfg.DatabaseLogin,
		cfg.DatabasePassword,
		cfg.DatabaseAddress,
		cfg.DatabaseName)
}

// generateTransactionNumber генерирует строку с номером транзакции. Код позаимствовал из функции, генерирующей номер
// http запроса в роутере chi.
func generateTransactionNumber() string {
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		_, err := rand.Read(buf[:])
		if err != nil {
			return ""
		}
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	return b64
}

// WithRepository служит для инициализации репозитория и внедрение его в сервис, используя паттерн Options.
func WithRepository(cfg *config.Storage) service.Option {
	log := slog.With(logger.OPLabel, "repository.mysql.WithRepository")
	if cfg == nil {
		log.Error(ErrNilConfigPointer.Error())
		os.Exit(1)
	}

	return func(s *service.Service) {
		db, err := sql.Open("mysql", createDSN(cfg))
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		db.SetMaxOpenConns(cfg.DatabaseMaxOpenConnections)

		if err = db.Ping(); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		log.Info("successfully ping db")
		repo := &Repository{db: db}
		s.Repository = repo
		s.SQLRepository = repo
	}
}

type txKey struct{}

// queryExecutorInterface интерфейс предоставления доступных методов для выполнения как одиночных, так и транзакционных
// запросов. Используется для объявления переменных, которые в дальнейшем будут инициализироваться методом
// Repository.executor.
type queryExecutorInterface interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

// extractTx если транзакция имеется в контексте, то возвращает объект транзакции *sql.Tx и true. В противном случае
// возвращает nil и false.
func (r *Repository) extractTx(ctx context.Context) (*sql.Tx, bool) {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx, true
	}
	return nil, false
}

// executor в зависимости от наличия в передаваемом контексте транзакции по ключу txKey{} либо возвращает эту
// транзакцию, либо пул соединений для выполнения отдельного запроса к базе данных. В контекст транзакционный объект
// внедряется в функции WithinTransaction.
func (r *Repository) executor(ctx context.Context) queryExecutorInterface {
	var inOuterTX bool
	var executor queryExecutorInterface
	if executor, inOuterTX = r.extractTx(ctx); inOuterTX == true {
		return executor
	}
	executor = r.db

	return executor
}

// WithinTransaction запускает функцию tFunc с контекстом, содержащим внутри транзакционный объект. Транзакция
// завершается, если функция завершается без ошибок. Взял идею такой работы с транзакциями в этой статье:
// https://habr.com/ru/articles/651799/. Модифицировал предложенную в статье идею, добавив возможность внутри функции
// tFunc использовать вызов WithinTransaction. При этом вложенный вызов WithinTransaction будет использовать ту же
// транзакцию, что и внешний.
func (r *Repository) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	var tx *sql.Tx
	var err error
	var internalCall bool
	var log *slog.Logger

	// если это внешний (первый) вызов функции
	if tx, internalCall = r.extractTx(ctx); !internalCall {
		// начинаем транзакцию
		ctx = context.WithValue(ctx, logger.TxId, generateTransactionNumber())
		log = logger.LogWithCtxData(ctx, slog.Default())
		log.Info("start transaction")

		if tx, err = r.db.BeginTx(ctx, &sql.TxOptions{Isolation: txIsolationLevel}); err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		log = slog.With(logger.TxLabel, ctx.Value(logger.TxId))
	}

	// запускаем callback
	err = tFunc(context.WithValue(ctx, txKey{}, tx))

	// при ошибке делаем откат
	if err != nil {
		// во внутреннем вызове не делаем откат, чтобы не вызвать ошибку при внешнем откате
		if internalCall {
			return err
		}
		log.Error(err.Error())

		if errRollback := tx.Rollback(); errRollback != nil {
			log.Error("error rollback transaction")
			return errRollback
		}
		log.Info("transaction rolled back")
		return err
	}

	// Если нет ошибок и мы не во вложенной транзакции, то выполняем коммит
	if !internalCall {
		if err = tx.Commit(); err != nil {
			log.Error(err.Error())
			return err
		}
		log.Info("transaction commit")
	}

	return nil
}

// ConvertToCommonErr замещает ошибки конкретной реализации на подобные по смыслу ошибки из доступных в абстрактном
// репозитории. К примеру, меняет sql.ErrNoRows на repository.ErrNoRecord.
func (r *Repository) ConvertToCommonErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case strings.HasPrefix(err.Error(), "Error 1062"):
		return repository.ErrDuplicate
	case errors.Is(err, sql.ErrNoRows):
		return repository.ErrNoRecord
	case errors.Is(err, context.DeadlineExceeded):
		return repository.ErrTimeout
	default:
		return err
	}
}

////////////////////////////
// Ниже по коду идёт CRUD //
////////////////////////////

// CreateStock сохраняет в БД запись о товаре.
func (r *Repository) CreateStock(ctx context.Context, data *dto.NamedProductDTO) error {
	var err error
	stmt := `INSERT INTO stock (article, name, price, amount) VALUES (?,?,?,?)`

	_, err = r.executor(ctx).ExecContext(ctx, stmt, data.Article, data.Name, data.Price, data.Amount)

	return r.ConvertToCommonErr(err)
}

// ReadStock возвращает запись из БД о товаре, находящемся в продаже в виде dto.NamedProductDTO.
func (r *Repository) ReadStock(ctx context.Context, data *dto.ArticleDTO) (dto.NamedProductDTO, error) {
	var result dto.NamedProductDTO
	var err error
	stmt := `SELECT article, name, price, amount FROM stock WHERE article = ?`

	row := r.executor(ctx).QueryRowContext(ctx, stmt, data.Article)
	err = row.Scan(&result.Article, &result.Name, &result.Price, &result.Amount)

	return result, r.ConvertToCommonErr(err)
}

// ReadStockAmount возвращает количество товара с артикулом, переданным в dto.ArticleDTO из находящегося в продаже.
func (r *Repository) ReadStockAmount(ctx context.Context, data *dto.ArticleDTO) (uint, error) {
	var amount uint
	stmt := `SELECT amount FROM stock WHERE article = ?`

	row := r.executor(ctx).QueryRowContext(ctx, stmt, data.Article)
	if err := row.Scan(&amount); err != nil {
		return 0, r.ConvertToCommonErr(err)
	}

	return amount, nil
}

// ReadStockPrice возвращает цену товара с артикулом, переданным в dto.ArticleDTO, из находящегося в продаже.
func (r *Repository) ReadStockPrice(ctx context.Context, data *dto.ArticleDTO) (float64, error) {
	var price float64
	stmt := `SELECT amount FROM stock WHERE article = ?`

	row := r.executor(ctx).QueryRowContext(ctx, stmt, data.Article)
	if err := row.Scan(&price); err != nil {
		return 0, r.ConvertToCommonErr(err)
	}

	return price, nil
}

// UpdateStock обновляет запись о товаре в БД, в соответствии с переданными в dto.NamedProductDTO данными.
func (r *Repository) UpdateStock(ctx context.Context, data *dto.NamedProductDTO) error {
	var err error
	stmt := `UPDATE stock SET name, price, amount = (?,?,?) WHERE article = ?`

	_, err = r.executor(ctx).ExecContext(ctx, stmt, data.Name, data.Price, data.Amount, data.Article)

	return r.ConvertToCommonErr(err)
}

// UpdateStockAmount обновляет количество доступного для продажи товара в соответствии с переданными в
// dto.ArticleWithAmountDTO данными.
func (r *Repository) UpdateStockAmount(ctx context.Context, data *dto.ArticleWithAmountDTO) error {
	var err error
	stmt := `UPDATE stock SET amount = ? WHERE article = ?`

	_, err = r.executor(ctx).ExecContext(ctx, stmt, data.Amount, data.Article)

	return r.ConvertToCommonErr(err)
}

// UpdateStockPrice обновляет цену доступного для продажи товара в соответствии с переданными в
// dto.ArticleWithPriceDTO данными.
func (r *Repository) UpdateStockPrice(ctx context.Context, data *dto.ArticleWithPriceDTO) error {
	var err error
	stmt := `UPDATE stock SET price = ? WHERE article = ?`

	_, err = r.executor(ctx).ExecContext(ctx, stmt, data.Price, data.Article)

	return r.ConvertToCommonErr(err)
}

// CreateReservation выполняет резервирование товаров в таблицу on_processing. Если в передаваемом контексте уже
// содержится транзакция, то запросы к БД выполняются в этой внешней транзакции. В противном случае создается новая
// транзакция и запросы выполняются в ней.
func (r *Repository) CreateReservation(ctx context.Context, data *dto.ReservationDTO) error {
	stmt := `INSERT
			 INTO on_processing 
			 (article, price, amount, date_of_reservation, updated_at, order_number, status) 
			 values (?,?,?,?,?,?,?)`

	currentDate := time.Now()

	f := func(txCtx context.Context) error {
		for _, p := range data.Products {
			_, err := r.executor(txCtx).ExecContext(ctx, stmt,
				p.Article, p.Price, p.Amount, currentDate, currentDate, data.OrderNumber, data.State)
			if err != nil {
				return r.ConvertToCommonErr(err)
			}
		}
		return nil
	}

	return r.WithinTransaction(ctx, f)
}

// ReadReservation возвращает в виде dto.ReservationDTO  данные о бронировании товаров с номером заказа, переданным в
// dto.OrderNumberDTO.
func (r *Repository) ReadReservation(ctx context.Context, data *dto.OrderNumberDTO) (dto.ReservationDTO, error) {
	stmt := `SELECT article, price, amount, date_of_reservation, order_number, status
    		 FROM on_processing 
    		 WHERE order_number = ?`

	rows, err := r.executor(ctx).QueryContext(ctx, stmt, data.OrderNumber)
	if err != nil {
		return dto.ReservationDTO{}, err
	}

	var state uint
	var date time.Time
	var orderNumber reservation.OrderNumber
	var products []dto.ProductDTO

	for rows.Next() {
		var product dto.ProductDTO
		err = rows.Scan(&product.Article, &product.Price, &product.Amount, &date, &orderNumber, &state)
		if err != nil {
			return dto.ReservationDTO{}, r.ConvertToCommonErr(err)
		}
		products = append(products, product)
	}

	return dto.ReservationDTO{OrderNumber: orderNumber, Date: date, State: state, Products: products}, nil
}

// UpdateReservation обновляет в БД записи о бронировании, в соответствии с переданными в dto.ReservationDTO данными
// (кроме идентификатора записи в таблицы и времени бронирования).
func (r *Repository) UpdateReservation(ctx context.Context, data *dto.ReservationDTO) error {
	var err error
	stmt := `UPDATE on_processing 
			 SET article = ?, price = ?, amount = ?, status= ? , updated_at= ? 
			 WHERE order_number = ? AND article = ? AND price = ? AND amount = ?`

	for _, p := range data.Products {
		_, err = r.executor(ctx).ExecContext(ctx, stmt,
			p.Article, p.Price, p.Amount, data.State, data.Date, data.OrderNumber, p.Article, p.Price, p.Amount)
		if err != nil {
			return r.ConvertToCommonErr(err)
		}
	}

	return nil
}

// DeleteReservation удаляет из БД записи с номером заказа, переданным в dto.OrderNumberDTO.
func (r *Repository) DeleteReservation(ctx context.Context, data *dto.OrderNumberDTO) error {
	stmt := `DELETE FROM on_processing WHERE order_number = ?`

	_, err := r.executor(ctx).ExecContext(ctx, stmt, data.OrderNumber)

	return r.ConvertToCommonErr(err)
}

// CreateSoldRecord сохраняет в БД запись об проданном товаре.
func (r *Repository) CreateSoldRecord(ctx context.Context, data *dto.SoldDTO) error {
	var err error
	stmt := `INSERT INTO sold (article, price, amount, date_of_sale) VALUES (?,?,?,?)`

	_, err = r.executor(ctx).ExecContext(ctx, stmt, data.Article, data.Price, data.Amount, data.Date)

	return r.ConvertToCommonErr(err)
}

// ReadSoldRecords возвращает все записи о продажах товара с переданным в dto.ArticleDTO артикулом.
func (r *Repository) ReadSoldRecords(ctx context.Context, data *dto.ArticleDTO) ([]dto.SoldDTO, error) {
	var result []dto.SoldDTO
	stmt := `SELECT article, price, amount, date_of_sale FROM stock WHERE article = ?`

	rows, err := r.executor(ctx).QueryContext(ctx, stmt, data.Article)
	if rows == nil || err != nil {
		return result, r.ConvertToCommonErr(err)
	}

	for rows.Next() {
		var record dto.SoldDTO
		if err = rows.Scan(&record.Article, &record.Price, &record.Amount, &record.Date); err != nil {
			return result, r.ConvertToCommonErr(err)
		}
		result = append(result, record)
	}

	return result, r.ConvertToCommonErr(err)
}

// ReadSoldAmount возвращает количество проданного товара с переданным в *dto.ArticleDTO артикулом (за весь период).
func (r *Repository) ReadSoldAmount(ctx context.Context, data *dto.ArticleDTO) (uint, error) {
	var result sql.NullInt64

	stmt := `SELECT SUM(amount) FROM sold WHERE article = ?`

	row := r.executor(ctx).QueryRowContext(ctx, stmt, data.Article)
	if err := row.Scan(&result); err != nil {
		return 0, r.ConvertToCommonErr(err)
	}
	if result.Valid {
		return uint(result.Int64), nil
	}

	return 0, nil
}

// ReadSoldRecordsInPeriod возвращает все записи о продажах товара с переданным в dto.ArticleWithPeriodDTO артикулом
// в период между датами From и To включительно.
func (r *Repository) ReadSoldRecordsInPeriod(ctx context.Context, data *dto.ArticleWithPeriodDTO) ([]dto.SoldDTO, error) {
	var result []dto.SoldDTO
	stmt := `SELECT article, price, amount, date_of_sale 
			 FROM stock 
			 WHERE article = ? AND date_of_sale >= ? AND date_of_sale <= ?`

	rows, err := r.executor(ctx).QueryContext(ctx, stmt, data.Article, data.From, data.To)
	if rows == nil || err != nil {
		return result, r.ConvertToCommonErr(err)
	}

	for rows.Next() {
		var record dto.SoldDTO
		if err = rows.Scan(&record.Article, &record.Price, &record.Amount, &record.Date); err != nil {
			return result, r.ConvertToCommonErr(err)
		}
		result = append(result, record)
	}

	return result, r.ConvertToCommonErr(err)
}

// ReadSoldAmountInPeriod возвращает количество проданного товара за определенный период.
func (r *Repository) ReadSoldAmountInPeriod(ctx context.Context, data *dto.ArticleWithPeriodDTO) (uint, error) {
	var result sql.NullInt64

	stmt := `SELECT SUM(amount) FROM sold WHERE article = ? AND date_of_sale >= ? AND date_of_sale <= ?`

	row := r.executor(ctx).QueryRowContext(ctx, stmt, data.Article, data.From, data.To)
	if err := row.Scan(&result); err != nil {
		return 0, r.ConvertToCommonErr(err)
	}
	if result.Valid {
		return uint(result.Int64), nil
	}

	return 0, nil
}
