package sqlman_test

import (
	"log"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/client/sqlman"
	"github.com/boxgo/box/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
)

func TestDemo(t *testing.T) {
	db := sqlman.StdConfig("default").Build()

	db.DB().SetConnMaxLifetime(time.Minute * 3)
	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(10)

	rows, err := db.DB().Query("SHOW DATABASES")
	if err != nil {
		t.Fatal(err)
	}

	columnTypes, _ := rows.ColumnTypes()

	for _, ct := range columnTypes {
		log.Println(ct, ct.ScanType().Kind(), ct.ScanType().String())
	}

	for rows != nil && rows.Next() {
		str := ""
		rows.Scan(&str)

		log.Println(str)
	}
}

func TestSetFunc(t *testing.T) {
	db := sqlman.StdConfig("default").Build()

	taskCh := make(chan int)
	stopCh := make(chan struct{})
	runningCh := make(chan struct{}, 10)

	go func() {
		for id := 0; id < 100; id++ {
			taskCh <- id
		}
	}()

	go func() {
		for range taskCh {
			runningCh <- struct{}{}

			go func() {
				db.DB().Exec("SHOW DATABASES")
				<-runningCh
			}()
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)

			if i == 3 {
				db.DB().SetMaxOpenConns(2)
			} else if i == 6 {
				db.DB().SetMaxOpenConns(50)
			} else if i == 9 {
				db.DB().SetMaxOpenConns(10)
			}

			logger.Debugf("%+v", db.DB().Stats())
		}

		stopCh <- struct{}{}
	}()

	<-stopCh
}
