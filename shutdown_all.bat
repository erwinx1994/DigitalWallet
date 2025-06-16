taskkill /f /fi "IMAGENAME eq api_gateway.exe" /im *
taskkill /f /fi "IMAGENAME eq balance_service.exe" /im *
taskkill /f /fi "IMAGENAME eq deposit_service.exe" /im *
taskkill /f /fi "IMAGENAME eq transaction_history_service.exe" /im *
taskkill /f /fi "IMAGENAME eq transfer_service.exe" /im *
taskkill /f /fi "IMAGENAME eq withdraw_service.exe" /im *
