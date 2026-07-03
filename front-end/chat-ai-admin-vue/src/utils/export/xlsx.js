const normalizeSheetName = (sheetName = 'Sheet1') =>
  String(sheetName || 'Sheet1')
    .replace(/[\\/*?:[\]]/g, '_')
    .slice(0, 31) || 'Sheet1'

export const exportAoAToXlsx = async ({ sheetName = 'Sheet1', fileName, rows = [], cols = [] }) => {
  const XLSX = await import('xlsx')
  const sheet = XLSX.utils.aoa_to_sheet(rows)

  if (Array.isArray(cols) && cols.length) {
    sheet['!cols'] = cols
  }

  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, sheet, normalizeSheetName(sheetName))
  XLSX.writeFile(workbook, fileName, { compression: true })
}
