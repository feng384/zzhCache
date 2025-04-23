# 清理旧进程和文件
Remove-Item -Path "server.exe" -ErrorAction SilentlyContinue

# 编译并运行服务器
go build -o server.exe
Start-Process -NoNewWindow -FilePath ".\server.exe" -ArgumentList "-port=8001"
Start-Process -NoNewWindow -FilePath ".\server.exe" -ArgumentList "-port=8002"
Start-Process -NoNewWindow -FilePath ".\server.exe" -ArgumentList "-port=8003 -api=1"

# 等待服务器启动
Start-Sleep -Seconds 2

# 发送请求
Invoke-WebRequest -Uri "http://localhost:9999/api?key=Tom" | Select-Object -ExpandProperty Content
Invoke-WebRequest -Uri "http://localhost:9999/api?key=Tom" | Select-Object -ExpandProperty Content
Invoke-WebRequest -Uri "http://localhost:9999/api?key=Tom" | Select-Object -ExpandProperty Content

# 结束所有 server.exe 进程
Get-Process -Name "server" -ErrorAction SilentlyContinue | Stop-Process