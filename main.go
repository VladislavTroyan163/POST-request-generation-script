package main

import (
	"log"

	configpack "github.com/VladislavTroyan163/POST-request-generation-script/config"
	loggerpack "github.com/VladislavTroyan163/POST-request-generation-script/logger"
	registriespack "github.com/VladislavTroyan163/POST-request-generation-script/registries"
	authorizationpack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/authorization"
	changesupplierspack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/changeSupplier"
	createregistryofidspack "github.com/VladislavTroyan163/POST-request-generation-script/requestsModules/createRegistryOfIDs"
	utilspack "github.com/VladislavTroyan163/POST-request-generation-script/utils"
)

func main() {
	//Создание логирования main
	logger, logFile, err := loggerpack.CreateLogger("main")
	if err != nil {
		log.Fatal("main: ошибка создания логера: \n", err)
	}
	defer logFile.Close()
	//----------------------------------------------------------------------------------

	//Чтение конфига приложения
	configApp, err := configpack.ReadConfigApp()
	if err != nil {
		logger.Println("main: ошибка чтения конфигурации приложения: \n", err)
		return
	}
	logger.Println("main: конфигурация прочитана успешно")
	//----------------------------------------------------------------------------------

	//Авторизация: конфигурация, формирование запроса, выполнение запроса.
	tokenAuthorization, err := authorizationpack.GetTokenAuthorization()
	if err != nil {
		logger.Println("main: ошибка при выполнении запроса авторизации: \n", err)
		return
	}
	logger.Println("main: авторизация прошла успешно")
	//----------------------------------------------------------------------------------

	//Получние от пользователя настроек загружаемого файла.
	fileName, sheetNumInt, columnName, err := utilspack.ReadConsoleSettingsFile()
	if err != nil {
		logger.Println("main:\nReadConsoleSettingsFile: \n", err)
		return
	}
	logger.Println("main: настройки загружаемого файла получены")
	//----------------------------------------------------------------------------------

	//Чтение реестра лицевых счетов
	registryNumber, err := registriespack.ReadRegistry(configApp.Mode, fileName, sheetNumInt, columnName)
	if err != nil {
		logger.Println("main: ошибка чтения реестра: \n", err)
		return
	}
	logger.Println("main: чтение реестра прошло успешно")
	//----------------------------------------------------------------------------------

	//Установка количества потоков
	workersCount := utilspack.GetWorkersCount(configApp.WorkersCount, len(registryNumber))
	if workersCount == 0 {
		logger.Println("main: число потоков не может быть равно нулю!")
		return
	}
	logger.Println("main: число потоков установлено:", workersCount)
	//----------------------------------------------------------------------------------

	//Преобразование объектов реестра к виду .json списка
	err = createregistryofidspack.CreateWorkersNumRegistryOfIDs(
		registryNumber,
		workersCount,
		tokenAuthorization,
	)
	if err != nil {
		logger.Println("main: ошибка преобразования реестра: \n", err)
		return
	}
	logger.Println("main: номера успешно переведены в row_id и сохранены")
	//----------------------------------------------------------------------------------

	//Преобразование .json списка к виду объекта
	rowIDs, err := createregistryofidspack.ReadFileId()
	if err != nil {
		logger.Println("main: ошибка преобразования .json к виду объекта: \n", err)
		return
	}
	logger.Println("main: .json успешно приведен к виду объекта")
	//----------------------------------------------------------------------------------

	//Отправка запросов на проставление поставщиков
	err = changesupplierspack.CreateWorkersChangeSupplier(
		workersCount,
		tokenAuthorization,
		rowIDs,
	)
	if err != nil {
		logger.Println("main: ошибка при проставлении поставщиков: \n", err)
		return
	}
	logger.Println("main: поставщики проставлены успешно")
	//----------------------------------------------------------------------------------
}
