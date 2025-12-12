# Script para probar el login y ver qué devuelve
$body = @{
    usernameOrEmail = "admin"
    password = "admin123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost/api/users/login" -Method Post -ContentType "application/json" -Body $body

Write-Host "=== Respuesta del Login ===" -ForegroundColor Green
$response | ConvertTo-Json -Depth 10

Write-Host "`n=== Usuario en la respuesta ===" -ForegroundColor Green
$response.user | ConvertTo-Json -Depth 10

Write-Host "`n=== Role del usuario ===" -ForegroundColor Green
Write-Host "Role: $($response.user.role)" -ForegroundColor Yellow

