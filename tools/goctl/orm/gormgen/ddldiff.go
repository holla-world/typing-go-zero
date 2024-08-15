// CODE BY CHATGPT
// CODE BY CHATGPT
// CODE BY CHATGPT
package gormgen

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Column represents a database column with its name and definition.
type Column struct {
	Name       string
	Definition string
	Position   int
}

// Table represents a database table with its name and columns.
type Table struct {
	Name    string
	Columns map[string]Column
}

// ParseSQLFile parses the given SQL file content and returns a map of table names to their structures.
func ParseSQLFile(content string) (map[string]Table, error) {
	tables := make(map[string]Table)
	var currentTable *Table
	scanner := bufio.NewScanner(strings.NewReader(content))
	position := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "CREATE TABLE") {
			// 忽略包含 LIKE 的 CREATE TABLE 语句
			if strings.Contains(line, "LIKE") {
				continue
			}
			tableName := extractTableName(line)
			currentTable = &Table{Name: tableName, Columns: make(map[string]Column)}
			tables[tableName] = *currentTable
			position = 0
		} else if currentTable != nil && strings.HasPrefix(line, ")") {
			currentTable = nil
		} else if currentTable != nil {
			columnParts := strings.Fields(line)
			if len(columnParts) >= 2 {
				columnName := strings.Trim(columnParts[0], "`") // 去掉列名两侧的反引号
				columnDefinition := strings.Join(columnParts[1:], " ")
				columnDefinition = strings.TrimRight(columnDefinition, ",") // 移除行末的逗号
				position++
				currentTable.Columns[columnName] = Column{
					Name:       columnName,
					Definition: columnDefinition,
					Position:   position,
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// extractTableName extracts the table name from a CREATE TABLE statement.
func extractTableName(line string) string {
	parts := strings.Fields(line)
	if len(parts) >= 3 {
		return strings.Trim(parts[2], "`") // 去掉表名两侧的反引号
	}
	return ""
}

// CompareTables compares two sets of tables and generates combined DDL statements for differences.
func CompareTables(currentTables, stableTables map[string]Table) map[string]string {
	ddlMap := make(map[string][]string)

	for tableName, currentTable := range currentTables {
		if stableTable, exists := stableTables[tableName]; exists {
			// Compare columns in current vs stable
			for columnName, currentColumn := range currentTable.Columns {
				if stableColumn, colExists := stableTable.Columns[columnName]; colExists {
					if currentColumn.Definition != stableColumn.Definition {
						ddl := fmt.Sprintf("MODIFY COLUMN `%s` %s", columnName, currentColumn.Definition)
						if pos := findColumnPosition(currentColumn.Position, currentTable); pos != "" {
							ddl += fmt.Sprintf(" AFTER `%s`", pos)
						}
						ddlMap[tableName] = append(ddlMap[tableName], ddl)
					}
				} else {
					ddl := fmt.Sprintf("ADD COLUMN `%s` %s", columnName, currentColumn.Definition)
					if pos := findColumnPosition(currentColumn.Position, currentTable); pos != "" {
						ddl += fmt.Sprintf(" AFTER `%s`", pos)
					}
					ddlMap[tableName] = append(ddlMap[tableName], ddl)
				}
			}
		} else {
			// If a table exists in current but not in stable, create it
			ddl := fmt.Sprintf("CREATE TABLE `%s` (...);", tableName)
			ddlMap[tableName] = append(ddlMap[tableName], ddl)
		}
	}

	var ddls = map[string]string{}
	for tableName, operations := range ddlMap {
		ddl := fmt.Sprintf("ALTER TABLE `%s` %s;", tableName, strings.Join(operations, ", "))
		ddls[tableName] = ddl
	}

	return ddls
}

// findColumnPosition finds the name of the column that precedes the given position.
func findColumnPosition(position int, table Table) string {
	for _, column := range table.Columns {
		if column.Position == position-1 {
			return column.Name
		}
	}
	return ""
}

// GetSQLContentFromGit retrieves the content of a SQL file from a specific branch.
func GetSQLContentFromGit(branch, filePath string) (string, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", branch, filePath))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func doDDlCompare(filePath string, cfg GenerateSpec) {
	// Define the file path and branches
	stableBranch := "stable"

	// Get the SQL content from the stable branch
	stableSQLContent, err := GetSQLContentFromGit(stableBranch, filePath)
	if err != nil {
		fmt.Println("Error retrieving stable branch SQL content:", err)
		return
	}

	// Get the SQL content from the current branch
	currentSQLContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error retrieving current branch SQL content:", err)
		return
	}

	// Parse the SQL files
	stableTables, err := ParseSQLFile(stableSQLContent)
	if err != nil {
		fmt.Println("Error parsing stable branch SQL:", err)
		return
	}

	currentTables, err := ParseSQLFile(string(currentSQLContent))
	if err != nil {
		fmt.Println("Error parsing current branch SQL:", err)
		return
	}

	// Compare the tables and generate DDLs
	ddls := CompareTables(currentTables, stableTables)

	for _, v := range cfg.TableSpec {
		ddl, ok := ddls[v.TableName]
		if !ok {
			continue
		}
		if v.Shards <= 1 {
			fmt.Printf("上线前请执行DDL：%s\n", ddl)
			return
		}
		for i := 0; i < v.Shards; i++ {
			s := strings.ReplaceAll(ddl, v.TableName, fmt.Sprintf("%s_%d", v.TableName, i))
			fmt.Printf("上线前请执行DDL：%s\n", s)
		}
	}
}
