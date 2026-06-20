# check-ports.ps1
# Verifica que los puertos requeridos por docker-compose esten disponibles

$puertos = @(
    [PSCustomObject]@{ Puerto = 5433; Servicio = "PostgreSQL (Docker)"; Detalle = "Puerto externo no estandar (sin conflicto con postgres local)" }
    [PSCustomObject]@{ Puerto = 8080; Servicio = "Go Backend API";      Detalle = "API REST" }
    [PSCustomObject]@{ Puerto = 3000; Servicio = "Vue Frontend";        Detalle = "Interfaz web (nginx interno en 80)" }
)

$linea = "  " + ("-" * 55)

Write-Host ""
Write-Host "  TaskGo -- Verificacion de puertos" -ForegroundColor Cyan
Write-Host $linea -ForegroundColor DarkGray
Write-Host ""

$ocupados = 0

foreach ($p in $puertos) {
    $conexiones = Get-NetTCPConnection -LocalPort $p.Puerto -State Listen -ErrorAction SilentlyContinue

    if ($conexiones) {
        $ocupados++
        $primera  = $conexiones | Select-Object -First 1
        $proceso  = Get-Process -Id $primera.OwningProcess -ErrorAction SilentlyContinue
        $procInfo = if ($proceso) { "$($proceso.Name)  (PID $($proceso.Id))" } else { "PID $($primera.OwningProcess)" }

        Write-Host "  [OCUPADO]  :$($p.Puerto)  $($p.Servicio)" -ForegroundColor Red
        Write-Host "             Proceso: $procInfo" -ForegroundColor DarkGray
        Write-Host "             Para liberar: Stop-Process -Id $($primera.OwningProcess) -Force" -ForegroundColor DarkYellow
        Write-Host ""
    } else {
        Write-Host "  [  OK   ]  :$($p.Puerto)  $($p.Servicio)" -ForegroundColor Green
        Write-Host "             $($p.Detalle)" -ForegroundColor DarkGray
        Write-Host ""
    }
}

Write-Host $linea -ForegroundColor DarkGray
Write-Host ""

if ($ocupados -eq 0) {
    Write-Host "  Todo OK. Puedes levantar la app:" -ForegroundColor Green
    Write-Host ""
    Write-Host "    docker-compose up --build" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  URLs una vez levantado:" -ForegroundColor DarkGray
    Write-Host "    Frontend   ->  http://localhost:3000" -ForegroundColor White
    Write-Host "    API        ->  http://localhost:8080/api/v1" -ForegroundColor White
    Write-Host "    PostgreSQL ->  localhost:5433  (usuario: postgres)" -ForegroundColor White
    Write-Host ""
    exit 0
} else {
    Write-Host "  $ocupados puerto(s) ocupado(s). Libera los procesos antes de ejecutar Docker." -ForegroundColor Red
    Write-Host ""
    Write-Host "  Comandos utiles:" -ForegroundColor DarkGray
    Write-Host "    Ver que usa un puerto:  Get-NetTCPConnection -LocalPort XXXX -State Listen" -ForegroundColor DarkYellow
    Write-Host "    Matar por PID:          Stop-Process -Id PID -Force" -ForegroundColor DarkYellow
    Write-Host "    Matar por nombre:       Stop-Process -Name nombre -Force" -ForegroundColor DarkYellow
    Write-Host ""
    exit 1
}
