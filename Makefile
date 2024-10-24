# Переменные
MOCKGEN = $(GOPATH)/bin/mockgen
MOCKS_DIR = mocks
DEPS_FILE=deps.go

# Установка mockgen
install_mockgen:
	@command -v $(MOCKGEN) >/dev/null 2>&1 || { \
		echo "Установка mockgen..."; \
		go install go.uber.org/mock/mockgen@v0.4.0; \
	}

# Функция для генерации моков для одного пакета
generate_mocks: install_mockgen
	@echo "Генерация моков для всех директорий, содержащих $(DEPS_FILE)..."
	@for dir in $$(find . -type f -name $(DEPS_FILE) -exec dirname {} \; | grep -v /vendor/); do \
		package_name=$$(basename $$dir); \
		if [ $$package_name = "cmd" ]; then \
			package_name="main"; \
		fi; \
		echo "Обработка директории $$dir с пакетом $$package_name"; \
		$(MOCKGEN) -source=$$dir/$(DEPS_FILE) -destination=$$dir/$$(basename $$dir)_mocks.go -package=$$package_name || { \
			echo "Ошибка генерации моков для директории $$dir"; \
			exit 1; \
		}; \
	done


mocks: generate_mocks

# Удаление сгенерированных моков
clean:
	rm -rf $(MOCKS_DIR)
