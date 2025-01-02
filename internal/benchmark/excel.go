package benchmark

import (
	"fmt"
	"log"
	"time"

	"github.com/xuri/excelize/v2"
)

func SaveResultsToExcel(servers []string, results []BenchmarkResult) string {
	f := excelize.NewFile()
	sheet := "DNS测试结果"
	f.SetSheetName("Sheet1", sheet)

	// 计算每个服务器结果需要的列数
	// 每个服务器占用4列（域名、响应时间、重试次数、错误信息）
	columnsPerServer := 6

	// 写入标题行
	for i, result := range results {
		baseCol := i * columnsPerServer // 每个服务器的起始列

		// 写入服务器标题
		serverCol := getColumnName(baseCol)
		f.SetCellValue(sheet, fmt.Sprintf("%s1", serverCol),
			fmt.Sprintf("DNS服务器 #%d: %s", i+1, result.Server))

		// 写入汇总信息表头
		f.SetCellValue(sheet, fmt.Sprintf("%s2", serverCol), "平均响应时间(ms)")
		f.SetCellValue(sheet, fmt.Sprintf("%s2", getColumnName(baseCol+1)), "成功率(%)")
		f.SetCellValue(sheet, fmt.Sprintf("%s2", getColumnName(baseCol+2)), "重试次数")

		// 写入汇总数据
		f.SetCellValue(sheet, fmt.Sprintf("%s3", serverCol),
			float64(result.AvgResponseTime.Milliseconds()))
		f.SetCellValue(sheet, fmt.Sprintf("%s3", getColumnName(baseCol+1)),
			result.SuccessRate*100)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", getColumnName(baseCol+2)),
			result.TotalRetries)

		// 写入详情表头
		f.SetCellValue(sheet, fmt.Sprintf("%s5", serverCol), "域名")
		f.SetCellValue(sheet, fmt.Sprintf("%s5", getColumnName(baseCol+1)), "响应时间(ms)")
		f.SetCellValue(sheet, fmt.Sprintf("%s5", getColumnName(baseCol+2)), "重试次数")
		f.SetCellValue(sheet, fmt.Sprintf("%s5", getColumnName(baseCol+3)), "错误信息")
		f.SetCellValue(sheet, fmt.Sprintf("%s5", getColumnName(baseCol+4)), "解析结果")
		f.SetCellValue(sheet, fmt.Sprintf("%s5", getColumnName(baseCol+5)), "平均延迟(ms)")

		// 写入域名测试详情
		for rowIdx, domain := range result.DomainResults {
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", serverCol, rowIdx+6),
				domain.Domain)
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", getColumnName(baseCol+1), rowIdx+6),
				float64(domain.ResponseTime.Milliseconds()))
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", getColumnName(baseCol+2), rowIdx+6),
				domain.RetryCount)
			if domain.Error != nil {
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", getColumnName(baseCol+3), rowIdx+6),
					domain.Error.Error())
			}
			ips := make([]string, len(domain.DnsPingResults.PingResults))
			for i, result := range domain.DnsPingResults.PingResults {
				ips[i] = result.IP
			}
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", getColumnName(baseCol+4), rowIdx+6),
				ips)
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", getColumnName(baseCol+5), rowIdx+6),
				float64(domain.DnsPingResults.AvgRTT.Milliseconds()))
		}

		// 设置列宽
		for j := 0; j < columnsPerServer; j++ {
			col := getColumnName(baseCol + j)
			width := 15.0
			if j == 0 { // 域名列
				width = 40.0
			}
			f.SetColWidth(sheet, col, col, width)
		}
	}

	// 保存文件
	filename := fmt.Sprintf("dns_benchmark_%s.xlsx", time.Now().Format("20060102_150405"))
	if err := f.SaveAs(filename); err != nil {
		log.Printf("保存Excel文件失败: %v\n", err)
		return ""
	}

	return filename
}

// 将列索引转换为Excel列名（A, B, C, ..., Z, AA, AB, ...）
func getColumnName(index int) string {
	name := ""
	for index >= 0 {
		name = string(rune('A'+(index%26))) + name
		index = index/26 - 1
	}
	return name
}
