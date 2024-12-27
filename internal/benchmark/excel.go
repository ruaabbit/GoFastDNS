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

	currentRow := 1

	for i, result := range results {
		// 写入服务器汇总信息
		f.SetCellValue(sheet, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("DNS服务器 #%d: %s", i+1, result.Server))
		currentRow++

		// 汇总信息表头
		summaryHeaders := []string{"平均响应时间(ms)", "成功率(%)", "重试次数"}
		for j, header := range summaryHeaders {
			f.SetCellValue(sheet, fmt.Sprintf("%c%d", 'B'+j, currentRow), header)
		}
		currentRow++

		// 写入汇总数据
		f.SetCellValue(sheet, fmt.Sprintf("B%d", currentRow), float64(result.AvgResponseTime.Milliseconds()))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", currentRow), result.SuccessRate*100)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", currentRow), result.TotalRetries)
		currentRow += 2

		// 域名详情表头
		detailHeaders := []string{"域名", "响应时间(ms)", "重试次数", "错误信息"}
		for j, header := range detailHeaders {
			f.SetCellValue(sheet, fmt.Sprintf("%c%d", 'A'+j, currentRow), header)
		}
		currentRow++

		// 写入域名测试详情
		for _, domain := range result.DomainResults {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", currentRow), domain.Domain)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", currentRow), float64(domain.ResponseTime.Milliseconds()))
			f.SetCellValue(sheet, fmt.Sprintf("C%d", currentRow), domain.RetryCount)
			if domain.Error != nil {
				f.SetCellValue(sheet, fmt.Sprintf("D%d", currentRow), domain.Error.Error())
			}
			currentRow++
		}

		// 添加空行分隔不同服务器的数据
		currentRow += 2
	}

	// 设置列宽
	f.SetColWidth(sheet, "A", "A", 40)
	f.SetColWidth(sheet, "B", "D", 15)

	// 保存文件
	filename := fmt.Sprintf("dns_benchmark_%s.xlsx", time.Now().Format("20060102_150405"))
	if err := f.SaveAs(filename); err != nil {
		log.Printf("保存Excel文件失败: %v\n", err)
		return ""
	}

	return filename
}
