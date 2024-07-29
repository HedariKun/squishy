package main

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func excelizing(dbPath string, output string) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	excelFile := excelize.NewFile()
	if err != nil {
		log.Fatal(err)
	}

	tables := getTables(db)

	for _, table := range tables {
		_, err := excelFile.NewSheet(table)
		if err != nil {
			log.Fatal(err)
		}

		columns := getTableColumns(db, table)
		for index, column := range columns {
			cell, err := excelize.CoordinatesToCellName(index+1, 1)
			if err != nil {
				log.Fatal(err)
			}
			excelFile.SetCellValue(table, cell, column)
		}

		rows := getTableRows(db, table)
		for index, row := range rows {
			for columnIndex, column := range columns {
				value := row[column]
				cell, err := excelize.CoordinatesToCellName(columnIndex+1, index+2)
				if err != nil {
					log.Fatal(err)
				}
				excelFile.SetCellValue(table, cell, value)
			}

		}

	}

	excelFile.DeleteSheet(excelFile.GetSheetName(0))

	excelFile.SaveAs(output + ".xlsx")

}

func getTables(db *gorm.DB) []string {
	var tables []string
	err := db.Raw("SELECT name FROM sqlite_master WHERE type = 'table'").Scan(&tables).Error

	if err != nil {
		log.Fatal(err)
	}

	return tables
}

func getTableColumns(db *gorm.DB, table string) []string {
	rows, err := db.Raw(fmt.Sprintf("SELECT * FROM %s", table)).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	return columnNames
}

func getTableRows(db *gorm.DB, table string) []map[string]interface{} {
	rows, err := db.Raw(fmt.Sprintf("SELECT * FROM %s", table)).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	results := []map[string]interface{}{}

	for rows.Next() {
		values := make([]interface{}, len(columnNames))
		valuePtrs := make([]interface{}, len(columnNames))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			panic(err)
		}

		result := map[string]interface{}{}
		for i, col := range columnNames {
			result[col] = values[i]
		}

		results = append(results, result)
	}

	return results
}
