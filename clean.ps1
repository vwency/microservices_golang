# Получаем путь к текущей директории
$rootPath = Get-Location

# Функция для удаления строк с комментариями
function Remove-LineComments {
    param (
        [string]$filePath
    )

    # Читаем все строки из файла
    $lines = Get-Content -Path $filePath

    # Фильтруем строки, удаляя те, которые являются комментариями
    $cleanedLines = $lines | ForEach-Object {
        if ($_ -match '^\s*//') {
            # Если строка начинается с //, то удаляем её
            return $null
        }
        return $_
    }

    # Записываем очищенные строки обратно в файл
    Set-Content -Path $filePath -Value $cleanedLines
}

# Рекурсивно проходим по всем файлам в текущей директории и её подкаталогах
Get-ChildItem -Path $rootPath -Recurse -File | ForEach-Object {
    $filePath = $_.FullName
    # Проверим, что это файл с расширением, в котором могут быть комментарии (например, .go, .js, .cpp)
    if ($filePath -match '\.(go|js|cpp|java|ts)$') {
        Write-Host "Processing file: $filePath"
        Remove-LineComments -filePath $filePath
    }
}
