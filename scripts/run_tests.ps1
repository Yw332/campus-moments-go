$server = Start-Process -FilePath "go" -ArgumentList "run cmd/api/main.go" -PassThru -NoNewWindow
Write-Host "Server started with PID $($server.Id). Waiting 15 seconds for compilation and startup..."
Start-Sleep -Seconds 15

Write-Host "Running tests..."
go run cmd/test_client/test_client.go

Write-Host "Running interaction tests..."
go run cmd/test_interaction/main.go

Write-Host "Stopping server..."
Stop-Process -Id $server.Id -Force
