package metrics

import (
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	dbQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of DB queries",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query"},
	)

	dbTableRowCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_table_row_count",
			Help: "Number of rows in DB tables",
		},
		[]string{"table"},
	)

	deletedItemsCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deleted_items_count",
			Help: "Number of deleted items",
		},
		[]string{"table"},
	)
)

func init() {
	prometheus.MustRegister(dbQueryDuration)
	prometheus.MustRegister(dbTableRowCount)
	prometheus.MustRegister(deletedItemsCount)
}

func ObserveDBQuery(queryName string, execFunc func() error) error {
	start := time.Now()
	err := execFunc()
	duration := time.Since(start)

	dbQueryDuration.WithLabelValues(queryName).Observe(duration.Seconds())
	return err
}

func UpdateTableRowCount(db *sql.DB, tableName string) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	if err == nil {
		dbTableRowCount.WithLabelValues(tableName).Set(float64(count))
	}
}

func IncrementDeletedItems(table string, amount int) {
	deletedItemsCount.WithLabelValues(table).Add(float64(amount))
}
