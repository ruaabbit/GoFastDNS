package benchmark

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"time"
)

func SaveResultsToExcel(servers []string, results []BenchmarkResult, formatResults []CSVResult, serverStats map[string]*ServerStats) string {
	// 创建新的Excel文件
	f := excelize.NewFile()
	sheet := "DNS测试结果"
	f.SetSheetName("Sheet1", sheet)

	// 写入表头
	headers := []string{"Domain"}
	for _, result := range results {
		headers = append(headers, result.Server)
	}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// 写入域名数据
	for rowID, res := range formatResults {
		row := []interface{}{res.Domain}
		for _, server := range servers {
			stats := res.ServerStats[server]
			if stats.Success {
				row = append(row, fmt.Sprintf("%.2f ms", stats.ResponseTime))
			} else {
				row = append(row, fmt.Sprintf("ERROR: %s", stats.Error))
			}
		}
		// 写入一行数据
		for colID, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colID+1, rowID+2)
			f.SetCellValue(sheet, cell, value)
		}
	}

	// 写入统计信息，从最后一行数据往下3行开始
	startRow := len(formatResults) + 5
	for serverIdx, server := range servers {
		stats := serverStats[server]
		colLetter, _ := excelize.CoordinatesToCellName(serverIdx+2, startRow)
		titleCell, _ := excelize.CoordinatesToCellName(serverIdx+2, startRow-1)

		// 写入服务器标题
		f.SetCellValue(sheet, titleCell, fmt.Sprintf("%s 统计信息", server))

		// 写入统计数据
		f.SetCellValue(sheet, colLetter, fmt.Sprintf("平均: %.2f ms", stats.AvgTime))
		f.SetCellValue(sheet, incrementRow(colLetter, 1), fmt.Sprintf("成功率: %.1f%%", stats.SuccessRate))
		f.SetCellValue(sheet, incrementRow(colLetter, 2), fmt.Sprintf("最快: %.2f ms", stats.MinTime))
		f.SetCellValue(sheet, incrementRow(colLetter, 3), fmt.Sprintf("最慢: %.2f ms", stats.MaxTime))
		f.SetCellValue(sheet, incrementRow(colLetter, 4), fmt.Sprintf("总查询: %d", stats.TotalQueries))
		f.SetCellValue(sheet, incrementRow(colLetter, 5), fmt.Sprintf("成功: %d", stats.SuccessQueries))
		f.SetCellValue(sheet, incrementRow(colLetter, 6), fmt.Sprintf("失败: %d", stats.FailedQueries))
	}

	// 设置样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#CCE5FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheet, 1, 1, headerStyle)

	// 自动调整列宽
	for i := 1; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)

		// 获取该列所有单元格的内容，找出最长的内容
		maxWidth := 10.0 // 设置最小宽度
		for row := 1; row <= len(formatResults)+1; row++ {
			cell, _ := excelize.CoordinatesToCellName(i, row)
			if cellValue, _ := f.GetCellValue(sheet, cell); len(cellValue) > 0 {
				// 根据内容长度计算所需宽度（假设每个字符约需要1.2个单位宽度）
				width := float64(len(cellValue)) * 1.2
				if width > maxWidth {
					maxWidth = width
				}
			}
		}

		// 设置列宽
		f.SetColWidth(sheet, col, col, maxWidth)

		// 保持自动筛选功能
		f.AutoFilter(sheet, fmt.Sprintf("%s1:%s%d", col, col, len(formatResults)+1), nil)
	}

	// 保存文件
	filename := fmt.Sprintf("dns_benchmark_%s.xlsx", time.Now().Format("20060102_150405"))
	if err := f.SaveAs(filename); err != nil {
		log.Printf("保存Excel文件失败: %v\n", err)
		return ""
	}

	return filename
}

// 辅助函数：递增单元格行号
func incrementRow(cell string, increment int) string {
	col, row, _ := excelize.CellNameToCoordinates(cell)
	newCell, _ := excelize.CoordinatesToCellName(col, row+increment)
	return newCell
}
